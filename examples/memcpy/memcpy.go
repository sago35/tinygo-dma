package main

import (
	"device/sam"
	"fmt"
	"machine"
	"time"
	"unsafe"

	dma "github.com/sago35/tinygo-dma"
)

const (
	bufSize = 1024
)

var (
	source      [bufSize]byte
	destination [bufSize]byte
)

var (
	dbg5 = machine.D5
	dbg6 = machine.D6
)

func initDbg() {
	dbg5.Configure(machine.PinConfig{Mode: machine.PinOutput})
	dbg6.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func main() {
	initDbg()

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// wait for USB-CDC
	time.Sleep(5 * time.Second)

	fmt.Printf("- %b\r\n", sam.DMAC.CHANNEL[0].CHSTATUS.Get())
	d := dma.NewDMA(func(d *dma.DMA) {
		dbg6.Toggle()
		//d.Trigger()
		return
	})
	d.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_DISABLE) // Only software/event triggers
	d.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_TRANSACTION)

	d1 := d.GetDescriptor()
	d1.UpdateDescriptor(dma.DescriptorConfig{
		SRC:    unsafe.Pointer(&source[0]),
		DST:    unsafe.Pointer(&destination[0]),
		SRCINC: dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC: dma.DMAC_SRAM_BTCTRL_DSTINC_ENABLE,
		SIZE:   3,
	})

	d2 := dma.NewDescriptor()
	d2.UpdateDescriptor(dma.DescriptorConfig{
		SRC:    unsafe.Pointer(&source[0]),
		DST:    unsafe.Pointer(&destination[8]),
		SRCINC: dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC: dma.DMAC_SRAM_BTCTRL_DSTINC_ENABLE,
		SIZE:   bufSize - 8,
	})

	d1.AddDescriptor(d2)

	for i := range source {
		source[i] = 0xFF - byte(i)
	}

	fmt.Printf("%b\r\n", sam.DMAC.CHANNEL[0].CHSTATUS.Get())
	fmt.Printf("INTSTATUS %08X\r\n", sam.DMAC.INTSTATUS.Get())
	d.Start()

	fmt.Printf("--\r\n")
	fmt.Printf("%b\r\n", sam.DMAC.CHANNEL[0].CHSTATUS.Get())
	fmt.Printf("pend: %b\r\n", sam.DMAC.PENDCH.Get())
	fmt.Printf("source %v %v\r\n", source[:8], source[8:16])
	fmt.Printf("before %v %v\r\n", destination[:8], destination[8:16])
	fmt.Printf("INTSTATUS %08X\r\n", sam.DMAC.INTSTATUS.Get())
	fmt.Printf("INTPEND %04X\r\n", sam.DMAC.INTPEND.Get())
	dbg5.High()
	dbg6.Toggle()
	d.Trigger()
	//sam.DMAC.CHANNEL[0].CHCTRLB.Set(0x02)
	//time.Sleep(1 * time.Second)
	d.Wait()
	//d.Trigger()
	//time.Sleep(1 * time.Second)
	dbg5.Low()
	fmt.Printf("%b\r\n", sam.DMAC.CHANNEL[0].CHSTATUS.Get())
	//d.Wait()
	//time.Sleep(5 * time.Second)
	fmt.Printf("after  %v %v\r\n", destination[:8], destination[8:16])
	fmt.Printf("%b\r\n", sam.DMAC.CHANNEL[0].CHSTATUS.Get())

	//for i := range source {
	//	destination[i] = 0
	//}

	dbg5.High()
	for i := range source {
		destination[i] = source[i]
	}
	dbg5.Low()
	fmt.Printf("after2 %v %v\r\n", destination[:8], destination[8:16])

	for {
		led.Toggle()
		time.Sleep(100 * time.Millisecond)
	}
}
