package main

import (
	"device/sam"
	"machine"
	"math"
	"time"
	"unsafe"

	dma "github.com/sago35/tinygo-dma"
)

var (
	dbg5 = machine.D5
	dbg6 = machine.D4
)

func initDbg() {
	dbg5.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dbg6.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func main() {
	initDbg()
	time.Sleep(2 * time.Second)

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	from := make([]uint16, 128)
	for i := 0; i < len(from); i++ {
		from[i] = uint16(math.Sin(float64(i)*math.Pi*2/float64(len(from)))*1000 + 0x0800)
	}

	dmadac := dma.NewDMA(func(d *dma.DMA) {
		dbg6.Toggle()
		return
	})
	dmadac.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_DISABLE)
	dmadac.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	desc := dmadac.GetDescriptor()
	desc.UpdateDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&from[0]),
		DST:      unsafe.Pointer(&sam.DAC.DATA[0].Reg),
		SRCINC:   dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC:   dma.DMAC_SRAM_BTCTRL_DSTINC_DISABLE,
		STEPSEL:  dma.DMAC_SRAM_BTCTRL_STEPSEL_SRC,
		BEATSIZE: dma.DMAC_SRAM_BTCTRL_BEATSIZE_HWORD,
		SIZE:     uint32(len(from)) * 2,
	})
	desc.AddDescriptor(desc)

	a0 := machine.A0
	a0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.DAC0.Configure(machine.DACConfig{})

	dmadac.Start()
	for {
		led.Toggle()
		dbg5.Toggle()
		dmadac.Trigger()
		time.Sleep(1 * time.Millisecond)
	}
}
