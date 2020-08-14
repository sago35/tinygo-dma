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

	from := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	d := dma.NewDMA(func(d dma.DMA) {
		dbg5.Toggle()
		//fmt.Printf("done\r\n")
		return
	})
	d.SetTrigger(0x07) // SERCOM1 TX trigger
	d.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)
	fmt.Printf("&from[0]: %08X\r\n", unsafe.Pointer(&from[0]))
	fmt.Printf("&DATA   : %08X\r\n", unsafe.Pointer(&spi0.Bus.DATA.Reg))
	d.AddDescriptor(unsafe.Pointer(&from[0]), unsafe.Pointer(&spi0.Bus.DATA.Reg),
		1,                 // Beat Size
		true,              // Source Address Increment
		false,             // Destination Address Increment
		1,                 // Step Size
		true,              // Step Selection (true: source, false: destination)
		uint16(len(from)), // Total size of DMA transfer
	)

	for {
		d.Start()
		//d.Trigger()
		//d.Wait()
		dbg6.Toggle()

		got := byte(spi0.Bus.DATA.Get())
		fmt.Printf("tx %02X - rx %02X\r\n", from[0], got)
		from[0]++
		led.Toggle()
		time.Sleep(1500 * time.Millisecond)
		if false {
			break
		}
	}

	txData := byte(0)
	for {
		rxData, err := spi0.Transfer(txData)
		if err != nil {
			fmt.Printf("err : %s\r\n", err.Error())
		}
		fmt.Printf("tx %02X - rx %02X\r\n", txData, rxData)
		txData++
		led.Toggle()
		time.Sleep(500 * time.Millisecond)
	}
}
