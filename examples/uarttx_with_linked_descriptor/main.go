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

	from := make([]byte, 4)
	from2 := make([]byte, len(from))
	from3 := make([]byte, len(from))
	for i := range from {
		from[i] = byte(i)
		from2[i] = byte(i) + 0x40
		from3[i] = byte(i) + 0x80
	}

	desc3 := dma.NewDescriptor()
	desc3.UpdateDescriptor(dma.DescriptorConfig{
		SRC:    unsafe.Pointer(&from3[0]),
		DST:    unsafe.Pointer(&uart.Bus.DATA.Reg),
		SRCINC: dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC: dma.DMAC_SRAM_BTCTRL_DSTINC_DISABLE,
		SIZE:   uint32(len(from3)), // Total size of DMA transfer
	})

	desc2 := dma.NewDescriptor()
	desc2.UpdateDescriptor(dma.DescriptorConfig{
		SRC:    unsafe.Pointer(&from2[0]),
		DST:    unsafe.Pointer(&uart.Bus.DATA.Reg),
		SRCINC: dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC: dma.DMAC_SRAM_BTCTRL_DSTINC_DISABLE,
		SIZE:   uint32(len(from2)), // Total size of DMA transfer
		DESC:   unsafe.Pointer(desc3),
	})

	dmatx := dma.NewDMA(func(d *dma.DMA) {
		dbg5.Toggle()
		return
	})
	dmatx.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM5_TX)
	dmatx.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	desc := dmatx.GetDescriptor()
	desc.UpdateDescriptor(dma.DescriptorConfig{
		SRC:    unsafe.Pointer(&from[0]),
		DST:    unsafe.Pointer(&uart.Bus.DATA.Reg),
		SRCINC: dma.DMAC_SRAM_BTCTRL_SRCINC_ENABLE,
		DSTINC: dma.DMAC_SRAM_BTCTRL_DSTINC_DISABLE,
		SIZE:   uint32(len(from)), // Total size of DMA transfer
		DESC:   unsafe.Pointer(desc2),
	})

	for {
		dbg5.Toggle()
		s := time.Now()
		dmatx.Start()
		dbg6.Toggle()
		dmatx.Wait()
		e := time.Now()

		fmt.Printf("tx : %#v %#v %#v\r\n", from, from2, from3)
		fmt.Printf("%d us\r\n", e.Sub(s).Microseconds())

		for i := range from {
			from[i] = from[(i+1)%len(from)]
		}

		led.Toggle()
		time.Sleep(1000 * time.Millisecond)
	}
}
