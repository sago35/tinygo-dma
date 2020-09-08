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
	dbg6 = machine.D6
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
		SCK:       machine.SPI0_SCK_PIN,
		SDO:       machine.SPI0_SDO_PIN,
		SDI:       machine.SPI0_SDI_PIN,
		Frequency: 48000000,
	})

	from := make([]byte, 8)
	for i := range from {
		from[i] = byte(i)
	}
	to := make([]byte, len(from))

	dmatx := dma.NewDMA(func(d *dma.DMA) {
		dbg5.Toggle()
		return
	})
	dmatx.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM1_TX)
	dmatx.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	dmatx.GetDescriptor().UpdateDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&from[0]),
		DST:      unsafe.Pointer(&spi0.Bus.DATA.Reg),
		SRCINC:   dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC:   dma.DMAC_SRAM_BTCTRL_DSTINC_DISABLE,
		SIZE:     uint32(len(from)), // Total size of DMA transfer
		BLOCKACT: 1,
	})

	dmarx := dma.NewDMA(func(d *dma.DMA) {
		dbg6.Toggle()
		return
	})
	dmarx.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM1_RX)
	dmarx.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	dmarx.GetDescriptor().UpdateDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&spi0.Bus.DATA.Reg),
		DST:      unsafe.Pointer(&to[0]),
		SRCINC:   dma.DMAC_SRAM_BTCTRL_SRCINC_DISABLE,
		DSTINC:   dma.DMAC_SRAM_BTCTRL_DSTINC_ENABLE,
		SIZE:     uint32(len(to)), // Total size of DMA transfer
		BLOCKACT: 1,
	})

	for {
		dbg5.Toggle()
		dbg6.Toggle()
		dmarx.Start()
		dmatx.Start()

		dmatx.Wait()
		dmarx.Wait()
		time.Sleep(500 * time.Millisecond)

		fmt.Printf("tx : %#v\r\n", from)
		fmt.Printf("rx : %#v\r\n", to)

		for i := range from {
			from[i] = from[(i+1)%len(from)]
		}

		led.Toggle()
		time.Sleep(1000 * time.Millisecond)
	}
}
