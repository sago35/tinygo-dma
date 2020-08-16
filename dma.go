// +build atsamd51

package dma

import (
	"device/sam"
	"fmt"
	"machine"
	"runtime/interrupt"
	"unsafe"
)

type DMA struct {
	Channel       uint8
	triggerSource uint8
	triggerAction uint8
	wait          chan bool
}

type DMADescriptor struct {
	btctrl   uint16
	btcnt    uint16
	srcaddr  uint32 //unsafe.Pointer
	dstaddr  uint32 //unsafe.Pointer
	Descaddr uint32 //unsafe.Pointer
}

const (
	dmaDescriptors       = 32
	maxDMAchannels       = 32
	maxDMATriggerSources = 32
)

//go:align 16
var DmaDescriptorSection [dmaDescriptors]DMADescriptor

//go:align 16
var dmaDescriptorWritebackSection [dmaDescriptors]DMADescriptor

var (
	nextDMAIndex uint8
	DMAChannels  [maxDMAchannels]DMA
	dmaCallbacks [maxDMAchannels]func(*DMA)
)

var (
	dbg5 = machine.D5
	dbg6 = machine.D4
)

func init() {
	for i := range DMAChannels {
		DMAChannels[i].wait = make(chan bool, 1)
	}
}

func handleDMACInterrupt(intr interrupt.Interrupt) {
	//fmt.Printf("interrupt %#v %04X\r\n", intr, sam.DMAC.INTPEND.Get())
	//susp := sam.DMAC.INTPEND.HasBits(sam.DMAC_INTPEND_SUSP_Pos)
	//tcmpl := sam.DMAC.INTPEND.HasBits(sam.DMAC_INTPEND_TCMPL_Pos)
	//terr := sam.DMAC.INTPEND.HasBits(sam.DMAC_INTPEND_TERR_Pos)
	//channel := (sam.DMAC.INTPEND.Get() >> sam.DMAC_INTPEND_ID_Pos) & sam.DMAC_INTPEND_ID_Msk
	channel := (sam.DMAC.INTPEND.Get() >> sam.DMAC_INTPEND_ID_Pos) & sam.DMAC_INTPEND_ID_Msk

	select {
	case DMAChannels[channel].wait <- true:
	}

	if sam.DMAC.CHANNEL[channel].CHINTFLAG.HasBits(sam.DMAC_CHANNEL_CHINTFLAG_SUSP) {
		sam.DMAC.CHANNEL[channel].CHINTFLAG.Set(sam.DMAC_CHANNEL_CHINTFLAG_SUSP)
	}
	if sam.DMAC.CHANNEL[channel].CHINTFLAG.HasBits(sam.DMAC_CHANNEL_CHINTFLAG_TCMPL) {
		sam.DMAC.CHANNEL[channel].CHINTFLAG.Set(sam.DMAC_CHANNEL_CHINTFLAG_TCMPL)
		dmaCallbacks[channel](&DMAChannels[channel])
	}
	if sam.DMAC.CHANNEL[channel].CHINTFLAG.HasBits(sam.DMAC_CHANNEL_CHINTFLAG_TERR) {
		sam.DMAC.CHANNEL[channel].CHINTFLAG.Set(sam.DMAC_CHANNEL_CHINTFLAG_TERR)
	}
}

func NewDMA(callback func(*DMA)) *DMA {
	if maxDMAchannels <= nextDMAIndex {
		return nil
	}

	{
		// DMAC peripheral has not yet been initialized. Initialize it now.
		// Init DMAC.
		// First configure the clocks, then configure the DMA descriptors. Those
		// descriptors must live in SRAM and must be aligned on a 16-byte boundary.
		// http://www.lucadavidian.com/2018/03/08/wifi-controlled-neo-pixels-strips/
		// https://svn.larosterna.com/oss/trunk/arduino/zerotimer/zerodma.cpp
		sam.MCLK.AHBMASK.SetBits(sam.MCLK_AHBMASK_DMAC_)

		sam.DMAC.CTRL.ClearBits(sam.DMAC_CTRL_DMAENABLE)
		sam.DMAC.CTRL.SetBits(sam.DMAC_CTRL_SWRST)
		for sam.DMAC.CTRL.HasBits(sam.DMAC_CTRL_SWRST) {
		}

		sam.DMAC.BASEADDR.Set(uint32(uintptr(unsafe.Pointer(&DmaDescriptorSection))))
		sam.DMAC.WRBADDR.Set(uint32(uintptr(unsafe.Pointer(&dmaDescriptorWritebackSection))))
		//fmt.Printf("BASE %08X\r\n", uint32(uintptr(unsafe.Pointer(&DmaDescriptorSection))))
		//fmt.Printf("WRBA %08X\r\n", uint32(uintptr(unsafe.Pointer(&dmaDescriptorWritebackSection))))

		sam.DMAC.CTRL.SetBits(sam.DMAC_CTRL_LVLEN0 | sam.DMAC_CTRL_LVLEN1 | sam.DMAC_CTRL_LVLEN2 | sam.DMAC_CTRL_LVLEN3)

		interrupt.New(sam.IRQ_DMAC_0, handleDMACInterrupt).Enable()
		interrupt.New(sam.IRQ_DMAC_1, handleDMACInterrupt).Enable()
		interrupt.New(sam.IRQ_DMAC_2, handleDMACInterrupt).Enable()
		interrupt.New(sam.IRQ_DMAC_3, handleDMACInterrupt).Enable()
		interrupt.New(sam.IRQ_DMAC_OTHER, handleDMACInterrupt).Enable()

	}

	dma := &DMAChannels[nextDMAIndex]
	dma.Channel = nextDMAIndex
	dmaCallbacks[nextDMAIndex] = callback

	// Enable peripheral with all priorities.
	sam.DMAC.CTRL.SetBits(sam.DMAC_CTRL_DMAENABLE)

	nextDMAIndex++
	return dma
}

// SetTrigger sets trigger source of Channel Control A register
func (dma *DMA) SetTrigger(triggerSource uint8) error {
	if maxDMATriggerSources <= triggerSource {
		return fmt.Errorf("trigger source must be smaller than 32")
	}
	dma.triggerSource = triggerSource
	return nil
}

// SetTriggerAction sets trigger action of Channel Control A register.
//   sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BLOCK
//     0x0 : BLOCK       : One trigger required for each block transfer
//   sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST
//     0x2 : BURST       : One trigger required for each burst transfer
//   sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_TRANSACTION
//     0x3 : TRANSACTION : One trigger required for each transaction
func (dma *DMA) SetTriggerAction(triggerAction uint8) error {
	if 0x03 <= triggerAction {
		return fmt.Errorf("trigger action must be smaller than 4")
	}
	//
	dma.triggerAction = triggerAction
	return nil
}

func (dma *DMA) Start() {
	// Reset channel.
	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.ClearBits(sam.DMAC_CHANNEL_CHCTRLA_ENABLE)
	for sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.HasBits(sam.DMAC_CHANNEL_CHCTRLA_ENABLE) {
	}

	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.SetBits(sam.DMAC_CHANNEL_CHCTRLA_SWRST)
	for sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.HasBits(sam.DMAC_CHANNEL_CHCTRLA_SWRST) {
	}

	// Configure channel.
	sam.DMAC.CHANNEL[dma.Channel].CHINTENSET.SetBits(sam.DMAC_CHANNEL_CHINTENSET_SUSP | sam.DMAC_CHANNEL_CHINTENSET_TCMPL | sam.DMAC_CHANNEL_CHINTENSET_TERR)
	//sam.DMAC.CHANNEL[dma.Channel].CHINTENSET.SetBits(sam.DMAC_CHANNEL_CHINTENSET_TCMPL | sam.DMAC_CHANNEL_CHINTENSET_TERR)

	sam.DMAC.CHANNEL[dma.Channel].CHPRILVL.Set(0)

	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.Set((uint32(dma.triggerSource) << sam.DMAC_CHANNEL_CHCTRLA_TRIGSRC_Pos) |
		(uint32(dma.triggerAction) << sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_Pos) |
		(sam.DMAC_CHANNEL_CHCTRLA_BURSTLEN_SINGLE << sam.DMAC_CHANNEL_CHCTRLA_BURSTLEN_Pos))

	//sam.DMAC.CHANNEL[dma.Channel].CHCTRLB.Set(sam.DMAC_CHANNEL_CHCTRLB_CMD_RESUME << sam.DMAC_CHANNEL_CHCTRLB_CMD_Pos)

	//sam.DMAC.CHANNEL[dma.Channel].CHINTENSET.SetBits(sam.DMAC_CHANNEL_CHINTENSET_SUSP | sam.DMAC_CHANNEL_CHINTENSET_TCMPL | sam.DMAC_CHANNEL_CHINTENSET_TERR)
	//sam.DMAC.CHANNEL[dma.Channel].CHINTENSET.SetBits(sam.DMAC_CHANNEL_CHINTENSET_TCMPL)

	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.SetBits(sam.DMAC_CHANNEL_CHCTRLA_ENABLE)
}

func (dma *DMA) Trigger() {
	sam.DMAC.SWTRIGCTRL.SetBits(1 << dma.Channel)
}

func (dma *DMA) SetAction() {
	// block beat transaction ...
}

func (dma *DMA) SetDescriptor(desc DMADescriptor) {
	DmaDescriptorSection[dma.Channel] = desc
}

func (dma *DMA) Wait() {
	//for !sam.DMAC.CHANNEL[dma.Channel].CHINTFLAG.HasBits(sam.DMAC_CHANNEL_CHINTFLAG_TCMPL) {
	//}
	//sam.DMAC.CHANNEL[dma.Channel].CHINTFLAG.SetBits(sam.DMAC_CHANNEL_CHINTFLAG_TCMPL)
	<-DMAChannels[dma.Channel].wait
}

type DMADescriptorHelper struct {
	StepSize uint8          // Address Increment Step Size : 1, 2, 4, 8, 16, 32, 64, 128
	StepSel  uint8          // Step Selection : 0=DST, 1=SRC
	DstInc   bool           // Destination Address Increment Enable
	SrcInc   bool           // Source Address Increment Enable
	BeatSize uint8          // Beat Size : 1=BYTE, 2:HWORD, 4:WORD
	BlockAct uint8          // Block Action : 0=NOACT, 1=INT, 2=SUSPEND, 3=BOTH
	EvoSel   uint8          // Event Output Selection : 0=DISABLE, 1=BLOCK, 3:BEAT
	Valid    bool           // Descriptor Valid
	Length   uint32         // Length is `BTCNT * BeatSize` or `BTCNT * BeatSize * StepSize`
	SrcAddr  unsafe.Pointer // Block Transfer Source Address
	DstAddr  unsafe.Pointer // Block Transfer Destination Address
	DescAddr unsafe.Pointer // Next Descriptor Address (must be 128-bit aligned)
}

func GetDMADescriptorHelper() *DMADescriptorHelper {
	ret := &DMADescriptorHelper{
		StepSize: 1,
		StepSel:  0,
		DstInc:   false,
		SrcInc:   false,
		BeatSize: 1,
		BlockAct: 0,
		EvoSel:   0,
		Valid:    true,
		Length:   0,
		SrcAddr:  unsafe.Pointer(nil),
		DstAddr:  unsafe.Pointer(nil),
		DescAddr: unsafe.Pointer(nil),
	}

	return ret
}

func (dma *DMA) AddDescriptor(src unsafe.Pointer, dst unsafe.Pointer, beatSize uint8, srcInc, dstInc bool, stepSize uint8, stepSrc bool, size uint16) {
	si := uint16(0)
	if srcInc {
		si = 1
	}

	di := uint16(0)
	if dstInc {
		di = 1
	}

	ssel := uint16(0)
	if stepSrc {
		ssel = 1
	}

	ss := uint16(0)
	switch stepSize {
	case 1:
		ss = 0
	case 2:
		ss = 1
	case 4:
		ss = 2
	case 8:
		ss = 3
	case 16:
		ss = 4
	case 32:
		ss = 5
	case 64:
		ss = 6
	case 128:
		ss = 7
	default:
		// TODO: error
		ss = 0
	}

	bs := uint16(0)
	switch beatSize {
	case 1:
		bs = 0
	case 2:
		bs = 1
	case 4:
		bs = 2
	default:
		// TODO: error
		bs = 0
	}

	// 1
	DmaDescriptorSection[dma.Channel] = DMADescriptor{
		btctrl: (1 << 0) | // VALID: Descriptor Valid
			(0 << 1) | // EVOSEL=DISABLE: Event Output Selection
			(1 << 3) | // BLOCKACT=NOACT: Block Action
			(bs << 8) | // BEATSIZE: Beat Size
			(si << 10) | // SRCINC: Source Address Increment Enable
			(di << 11) | // DSTINC: Destination Address Increment Enable
			(ssel << 12) | // STEPSEL: Step Selection
			(ss << 13), // STEPSIZE: Address Increment Step Size
		btcnt:    size >> bs,
		Descaddr: 0,
	}

	if stepSrc {
		// STEPSEL == SRC
		DmaDescriptorSection[dma.Channel].srcaddr = uint32(uintptr(src) + uintptr((size)<<ss))
		if dstInc {
			DmaDescriptorSection[dma.Channel].dstaddr = uint32(uintptr(dst) + uintptr(size))
		} else {
			DmaDescriptorSection[dma.Channel].dstaddr = uint32(uintptr(dst))
		}
	} else {
		// STEPSEL == DST
		if srcInc {
			DmaDescriptorSection[dma.Channel].srcaddr = uint32(uintptr(src) + uintptr(size))
		} else {
			DmaDescriptorSection[dma.Channel].srcaddr = uint32(uintptr(src))
		}
		DmaDescriptorSection[dma.Channel].dstaddr = uint32(uintptr(dst) + uintptr((size)<<ss))
	}
}

func (dma *DMA) NewDescriptor(src unsafe.Pointer, dst unsafe.Pointer, beatSize uint8, srcInc, dstInc bool, stepSize uint8, stepSrc bool, size uint16) *DMADescriptor {
	si := uint16(0)
	if srcInc {
		si = 1
	}

	di := uint16(0)
	if dstInc {
		di = 1
	}

	ssel := uint16(0)
	if stepSrc {
		ssel = 1
	}

	ss := uint16(0)
	switch stepSize {
	case 1:
		ss = 0
	case 2:
		ss = 1
	case 4:
		ss = 2
	case 8:
		ss = 3
	case 16:
		ss = 4
	case 32:
		ss = 5
	case 64:
		ss = 6
	case 128:
		ss = 7
	default:
		// TODO: error
		ss = 0
	}

	bs := uint16(0)
	switch beatSize {
	case 1:
		bs = 0
	case 2:
		bs = 1
	case 4:
		bs = 2
	default:
		// TODO: error
		bs = 0
	}

	// 2
	//go:align 16
	DmaDescriptorSection[1] = DMADescriptor{
		btctrl: (1 << 0) | // VALID: Descriptor Valid
			(0 << 1) | // EVOSEL=DISABLE: Event Output Selection
			(0 << 3) | // BLOCKACT=NOACT: Block Action
			(bs << 8) | // BEATSIZE: Beat Size
			(si << 10) | // SRCINC: Source Address Increment Enable
			(di << 11) | // DSTINC: Destination Address Increment Enable
			(ssel << 12) | // STEPSEL: Step Selection
			(ss << 13), // STEPSIZE: Address Increment Step Size
		btcnt:    size >> bs,
		Descaddr: 0,
	}

	if stepSrc {
		// STEPSEL == SRC
		DmaDescriptorSection[1].srcaddr = uint32(uintptr(src) + uintptr((size)<<ss))
		DmaDescriptorSection[1].dstaddr = uint32(uintptr(dst) + uintptr(size))
	} else {
		// STEPSEL == DST
		DmaDescriptorSection[1].srcaddr = uint32(uintptr(src) + uintptr(size))
		DmaDescriptorSection[1].dstaddr = uint32(uintptr(dst) + uintptr((size)<<ss))
	}

	return &DmaDescriptorSection[1]
}
