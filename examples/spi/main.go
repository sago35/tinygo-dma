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
	d := dma.NewDMA(func(d *dma.DMA) {
		dbg5.Toggle()
		//fmt.Printf("done\r\n")
		return
	})
	d.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM1_TX)
	d.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	desc := dma.NewDescriptor(dma.DescriptorConfig{
		SRC:     unsafe.Pointer(&from[0]),
		DST:     unsafe.Pointer(&spi0.Bus.DATA.Reg),
		SRCINC:  true,
		DSTINC:  false,
		STEPSEL: true,
		SIZE:    uint32(len(from)), // Total size of DMA transfer
		//BLOCKACT: 1,
	})
	d.SetDescriptor(desc)

	to := make([]byte, len(from))
	d2 := dma.NewDMA(func(d *dma.DMA) {
		dbg6.Toggle()
		return
	})
	d2.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM1_RX)
	d2.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	desc2 := dma.NewDescriptor(dma.DescriptorConfig{
		SRC:    unsafe.Pointer(&spi0.Bus.DATA.Reg),
		DST:    unsafe.Pointer(&to[0]),
		SRCINC: false,
		DSTINC: true,
		SIZE:   uint32(len(to)), // Total size of DMA transfer
		//BLOCKACT: 1,
	})
	d2.SetDescriptor(desc2)

	sam.DMAC.CHANNEL[d2.Channel].CHPRILVL.Set(1)

	for {
		dbg5.Toggle()
		dbg6.Toggle()
		d2.Start()
		d.Start()
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
