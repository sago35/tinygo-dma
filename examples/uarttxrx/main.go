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

	uart := machine.UART1
	uart.Configure(machine.UARTConfig{
		BaudRate: 115200,
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
	dmatx.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM5_TX)
	dmatx.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	dmatx.GetDescriptor().UpdateDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&from[0]),
		DST:      unsafe.Pointer(&uart.Bus.DATA.Reg),
		SRCINC:   dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC:   dma.DMAC_SRAM_BTCTRL_DSTINC_DISABLE,
		SIZE:     uint32(len(from)), // Total size of DMA transfer
		BLOCKACT: dma.DMAC_SRAM_BTCTRL_BLOCKACT_NOACT,
	})

	dmarx := dma.NewDMA(func(d *dma.DMA) {
		dbg6.Toggle()
		return
	})
	dmarx.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM5_RX)
	dmarx.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	dmarx.GetDescriptor().UpdateDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&uart.Bus.DATA.Reg),
		DST:      unsafe.Pointer(&to[0]),
		SRCINC:   dma.DMAC_SRAM_BTCTRL_SRCINC_DISABLE,
		DSTINC:   dma.DMAC_SRAM_BTCTRL_DSTINC_ENABLE,
		SIZE:     uint32(len(to)), // Total size of DMA transfer
		BLOCKACT: dma.DMAC_SRAM_BTCTRL_BLOCKACT_NOACT,
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
