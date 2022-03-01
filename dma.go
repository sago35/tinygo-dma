// +build atsamd51 atsame5x

package dma

import (
	"device/sam"
	"fmt"
	"runtime"
	"runtime/interrupt"
	"unsafe"
)

type DMA struct {
	Channel       uint8
	triggerSource uint8
	triggerAction uint8
	cmdOnStart    uint8
	wait          chan bool
}

const (
	dmaDescriptors = 32
	maxDMAchannels = 32
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

func init() {
	for i := range DMAChannels {
		DMAChannels[i].wait = make(chan bool, 1)
	}
}

func handleDMACInterrupt(intr interrupt.Interrupt) {
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

	if nextDMAIndex == 0 {
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
	if DMAC_CHANNEL_CHCTRLA_TRIGSRC_Max < triggerSource {
		return fmt.Errorf("trigger source must be smaller than %d", DMAC_CHANNEL_CHCTRLA_TRIGSRC_Max)
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
	if sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_TRANSACTION < triggerAction {
		return fmt.Errorf("trigger action must be smaller than 4")
	}
	dma.triggerAction = triggerAction
	return nil
}

func (dma *DMA) SetCmdOnStart(cmdOnStart uint8) error {
	if sam.DMAC_CHANNEL_CHCTRLB_CMD_RESUME < cmdOnStart {
		return fmt.Errorf("software Command must be smaller than 3")
	}
	dma.cmdOnStart = cmdOnStart
	return nil
}

// Start starts DMA transfer.
func (dma *DMA) Start() {
	// Reset channel.
	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.ClearBits(sam.DMAC_CHANNEL_CHCTRLA_ENABLE)
	for sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.HasBits(sam.DMAC_CHANNEL_CHCTRLA_ENABLE) {
	}

	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.SetBits(sam.DMAC_CHANNEL_CHCTRLA_SWRST)
	for sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.HasBits(sam.DMAC_CHANNEL_CHCTRLA_SWRST) {
	}

	// Configure channel.
	sam.DMAC.CHANNEL[dma.Channel].CHINTENSET.SetBits(sam.DMAC_CHANNEL_CHINTENSET_TCMPL | sam.DMAC_CHANNEL_CHINTENSET_TERR)

	sam.DMAC.CHANNEL[dma.Channel].CHEVCTRL.Set(sam.DMAC_CHANNEL_CHEVCTRL_EVOE | sam.DMAC_CHANNEL_CHEVCTRL_EVIE | sam.DMAC_CHANNEL_CHEVCTRL_EVACT_RESUME)

	sam.DMAC.CHANNEL[dma.Channel].CHPRILVL.Set(0)

	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.Set((uint32(dma.triggerSource) << sam.DMAC_CHANNEL_CHCTRLA_TRIGSRC_Pos) |
		(uint32(dma.triggerAction) << sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_Pos) |
		(sam.DMAC_CHANNEL_CHCTRLA_BURSTLEN_SINGLE << sam.DMAC_CHANNEL_CHCTRLA_BURSTLEN_Pos))

	sam.DMAC.CHANNEL[dma.Channel].CHCTRLB.Set(dma.cmdOnStart)

	sam.DMAC.CHANNEL[dma.Channel].CHCTRLA.SetBits(sam.DMAC_CHANNEL_CHCTRLA_ENABLE)
}

// Trigger generates a DMA software trigger on correspond channel
func (dma *DMA) Trigger() {
	sam.DMAC.SWTRIGCTRL.SetBits(1 << dma.Channel)
}

// Resume generates a DMA Channel resume operation command on correspond channel
func (dma *DMA) Resume() {
	sam.DMAC.CHANNEL[dma.Channel].CHCTRLB.Set(sam.DMAC_CHANNEL_CHCTRLB_CMD_RESUME)
}

// Suspend generates a DMA Channel suspend operation command on correspond channel
func (dma *DMA) Suspend() {
	sam.DMAC.CHANNEL[dma.Channel].CHCTRLB.Set(sam.DMAC_CHANNEL_CHCTRLB_CMD_SUSPEND)
}

func (dma *DMA) SetAction() {
	// block beat transaction ...
}

func (dma *DMA) GetDescriptor() *DMADescriptor {
	return &DmaDescriptorSection[dma.Channel]
}

func (dma *DMA) Wait() {
	//for !sam.DMAC.CHANNEL[dma.Channel].CHINTFLAG.HasBits(sam.DMAC_CHANNEL_CHINTFLAG_TCMPL) {
	//}
	//sam.DMAC.CHANNEL[dma.Channel].CHINTFLAG.SetBits(sam.DMAC_CHANNEL_CHINTFLAG_TCMPL)
	for {
		select {
		case <-DMAChannels[dma.Channel].wait:
			return
		default:
			runtime.Gosched()
		}
	}
}
