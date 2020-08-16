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
	dbg6 = machine.D4
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
		d.Trigger()
		return
	})
	d.SetTrigger(0) // Only software/event triggers
	d.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_TRANSACTION)

	// Only software/event triggers
	//d.SetTrigger(0)

	// やっぱりいらないものが多いので、 DMAConfig{} を作って処理するようにしたい
	// func (c *DMAConfig) Add(a DMAConfig) も出来るようにしたい
	// 最終的に d.AddDescriptor(c) して link も出来るようにする
	// といいつつ、 chan 毎の最初の 1 個目の descriptor はアドレスがある程度固定なので
	//   c := DMADescriptor{}
	//   c.Add(DMADescriptor{})
	//   d.SetDescriptor(c)
	// みたいなのが良いかも
	// 最終、決められたアドレスの場所の descriptor に書く必要があるのでコピーしつつ
	// channel 毎に enable にしていく必要がありそう
	// なんでも設定できる I/F を定義するのはかなり難しい感じなので、
	// 用途を絞って SPI tx に対してはこれ、みたいなのを定義するのが楽そう
	d.AddDescriptor(unsafe.Pointer(&source[0]), unsafe.Pointer(&destination[0]),
		1,       // Beat Size
		true,    // Source Address Increment
		true,    // Destination Address Increment
		1,       // Step Size
		false,   // Step Selection (true: source, false: destination)
		bufSize, // Total size of DMA transfer
	)

	desc2 := d.NewDescriptor(unsafe.Pointer(&source[0]), unsafe.Pointer(&destination[9]),
		1,         // Beat Size
		true,      // Source Address Increment
		true,      // Destination Address Increment
		1,         // Step Size
		false,     // Step Selection (true: source, false: destination)
		bufSize-9, // Total size of DMA transfer
	)

	for i := range source {
		source[i] = 0xFF - byte(i)
	}

	if true {
		fmt.Printf("new %08X\r\n", unsafe.Pointer(&desc2))
		dma.DmaDescriptorSection[d.Channel].Descaddr = uint32(uintptr(unsafe.Pointer(desc2)))
		fmt.Printf("%#v\r\n", dma.DmaDescriptorSection[0])
		fmt.Printf("%#v\r\n", dma.DmaDescriptorSection[1])
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
