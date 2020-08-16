package main

import (
	"device/sam"
	"fmt"
	"machine"
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

	spi0 := machine.SPI0
	spi0.Configure(machine.SPIConfig{
		SCK: machine.SPI0_SCK_PIN,
		SDO: machine.SPI0_SDO_PIN,
		SDI: machine.SPI0_SDI_PIN,
	})

	//	from := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	//	d := dma.NewDMA(func(d dma.DMA) {
	//		dbg5.Toggle()
	//		//fmt.Printf("done\r\n")
	//		return
	//	})
	//	d.SetTrigger(0x07) // SERCOM1 TX trigger
	//	d.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)
	//
	//	d.AddDescriptor(unsafe.Pointer(&from[0]), unsafe.Pointer(&spi0.Bus.DATA.Reg),
	//		1,                 // Beat Size
	//		true,              // Source Address Increment
	//		false,             // Destination Address Increment
	//		1,                 // Step Size
	//		true,              // Step Selection (true: source, false: destination)
	//		uint16(len(from)), // Total size of DMA transfer
	//	)
	//
	//	if false {
	//		for {
	//			d.Start()
	//			//d.Trigger()
	//			//d.Wait()
	//			dbg6.Toggle()
	//
	//			got := byte(spi0.Bus.DATA.Get())
	//			fmt.Printf("tx %02X - rx %02X\r\n", from[0], got)
	//			fmt.Printf("rx : %#v\r\n", to)
	//			from[0]++
	//			led.Toggle()
	//			time.Sleep(1500 * time.Millisecond)
	//			if false {
	//				break
	//			}
	//		}
	//	}

	//	to := []byte{0, 0, 0, 0, 0, 0, 0}
	//	d2 := dma.NewDMA(func(d dma.DMA) {
	//		dbg5.Toggle()
	//		return
	//	})
	//	d2.SetTrigger(0x06) // SERCOM1 RX trigger
	//	d2.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)
	//
	//	d2.AddDescriptor(unsafe.Pointer(&spi0.Bus.DATA.Reg), unsafe.Pointer(&to[0]),
	//		1,     // Beat Size
	//		false, // Source Address Increment
	//		true,  // Destination Address Increment
	//		1,     // Step Size
	//		false, // Step Selection (true: source, false: destination)
	//		1,     // Total size of DMA transfer
	//		//uint16(len(to)), // Total size of DMA transfer
	//	)
	//	//sam.DMAC.CHANNEL[d2.Channel].CHPRILVL.Set(1)
	//
	//	txData := byte(1)
	//	for {
	//		d2.Start()
	//		spi0.Bus.DATA.Set(uint32(txData))
	//		time.Sleep(500 * time.Millisecond)
	//		fmt.Printf("tx : %02X rx : %02X\r\n", txData, to[0])
	//		txData++
	//		led.Toggle()
	//		time.Sleep(1000 * time.Millisecond)
	//	}

	from := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	d := dma.NewDMA(func(d *dma.DMA) {
		dbg5.Toggle()
		//fmt.Printf("done\r\n")
		return
	})
	d.SetTrigger(0x07) // SERCOM1 TX trigger
	d.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	d.AddDescriptor(unsafe.Pointer(&from[0]), unsafe.Pointer(&spi0.Bus.DATA.Reg),
		1,                 // Beat Size
		true,              // Source Address Increment
		false,             // Destination Address Increment
		1,                 // Step Size
		true,              // Step Selection (true: source, false: destination)
		uint16(len(from)), // Total size of DMA transfer
	)

	to := make([]byte, 8)
	d2 := dma.NewDMA(func(d *dma.DMA) {
		dbg6.Toggle()
		return
	})
	d2.SetTrigger(0x06) // SERCOM1 RX trigger
	d2.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	d2.AddDescriptor(unsafe.Pointer(&spi0.Bus.DATA.Reg), unsafe.Pointer(&to[0]),
		1,               // Beat Size
		false,           // Source Address Increment
		true,            // Destination Address Increment
		1,               // Step Size
		false,           // Step Selection (true: source, false: destination)
		uint16(len(to)), // Total size of DMA transfer
	)
	sam.DMAC.CHANNEL[d2.Channel].CHPRILVL.Set(1)

	for {
		dbg5.Toggle()
		d.Start()
		dbg6.Toggle()
		d2.Start()
		d.Wait()
		d2.Wait()
		time.Sleep(500 * time.Millisecond)

		//got := byte(spi0.Bus.DATA.Get())
		//fmt.Printf("tx %02X - rx %02X\r\n", from[0], got)
		fmt.Printf("tx : %#v\r\n", from)
		fmt.Printf("rx : %#v\r\n", to)
		from[0]++
		led.Toggle()
		time.Sleep(1000 * time.Millisecond)
		if false {
			break
		}
	}
}
