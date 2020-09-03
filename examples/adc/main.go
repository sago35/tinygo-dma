package main

import (
	"device/sam"
	"fmt"
	"machine"
	"time"
	"unsafe"

	dma "github.com/sago35/tinygo-dma"
)

// This example assumes that an analog sensor such as a rotary dial is connected to pin ADC0.
// When the dial is turned past the midway point, the built-in LED will light up.

func main() {
	machine.InitADC()

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	dmaadc := dma.NewDMA(func(d *dma.DMA) {
		return
	})
	dmaadc.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_ADC0_RESRDY)
	dmaadc.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	to := uint16(0)

	desc := dma.NewDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&sam.ADC0.RESULT.Reg),
		DST:      unsafe.Pointer(&to),
		SRCINC:   false,
		DSTINC:   true,
		SIZE:     2,
		STEPSEL:  false,
		BEATSIZE: 2,
	})
	desc.AddDescriptor(desc)
	dmaadc.SetDescriptor(desc)

	sensor := machine.ADC{Pin: machine.A0}
	sensor.Configure()

	bus := sam.ADC0
	ch := 0
	{

		for bus.SYNCBUSY.HasBits(sam.ADC_SYNCBUSY_INPUTCTRL) {
		}

		// Selection for the positive ADC input channel
		bus.INPUTCTRL.ClearBits(sam.ADC_INPUTCTRL_MUXPOS_Msk)
		for bus.SYNCBUSY.HasBits(sam.ADC_SYNCBUSY_ENABLE) {
		}
		bus.INPUTCTRL.SetBits((uint16(ch) & sam.ADC_INPUTCTRL_MUXPOS_Msk) << sam.ADC_INPUTCTRL_MUXPOS_Pos)
		for bus.SYNCBUSY.HasBits(sam.ADC_SYNCBUSY_ENABLE) {
		}

		// Enable ADC
		bus.CTRLA.SetBits(sam.ADC_CTRLA_ENABLE)
		for bus.SYNCBUSY.HasBits(sam.ADC_SYNCBUSY_ENABLE) {
		}

	}

	dmaadc.Start()
	for {
		if false {
			val := sensor.Get()
			fmt.Printf("%04X\r\n", val)
		}
		// Start conversion
		bus.SWTRIG.SetBits(sam.ADC_SWTRIG_START)
		time.Sleep(time.Millisecond * 100)
		fmt.Printf("%04X\r\n", to)
		led.Toggle()
	}
}
