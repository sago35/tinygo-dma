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

func timerInit() {
	// Turn on clock to TC0 / TC1
	sam.MCLK.APBAMASK.SetBits(sam.MCLK_APBAMASK_TC0_)

	// 48Mhz
	sam.GCLK.PCHCTRL[9].Set((sam.GCLK_PCHCTRL_GEN_GCLK1 << sam.GCLK_PCHCTRL_GEN_Pos) |
		sam.GCLK_PCHCTRL_CHEN)

	// init
	if !sam.TC0_COUNT32.SYNCBUSY.HasBits(sam.TC_COUNT32_SYNCBUSY_SWRST) {
		if sam.TC0_COUNT32.CTRLA.HasBits(sam.TC_COUNT32_CTRLA_ENABLE) {
			sam.TC0_COUNT32.CTRLA.ClearBits(sam.TC_COUNT32_CTRLA_ENABLE)
			for sam.TC0_COUNT32.SYNCBUSY.HasBits(sam.TC_COUNT32_SYNCBUSY_ENABLE) {
			}
		}
		sam.TC0_COUNT32.CTRLA.SetBits(sam.TC_COUNT32_CTRLA_SWRST)
	}
	for sam.TC0_COUNT32.SYNCBUSY.HasBits(sam.TC_COUNT32_SYNCBUSY_SWRST) {
	}

	//hri_tc_get_INTFLAG_reg(TC0, (TC_INTFLAG_OVF | TC_INTFLAG_ERR | TC_INTFLAG_MC0 | TC_INTFLAG_MC1));
	sam.TC0_COUNT32.INTFLAG.SetBits(sam.TC_COUNT32_INTFLAG_OVF |
		sam.TC_COUNT32_INTFLAG_ERR |
		sam.TC_COUNT32_INTFLAG_MC0 |
		sam.TC_COUNT32_INTFLAG_MC1)

	//hri_tc_write_CTRLA_reg(TC0,
	//                       0 << TC_CTRLA_CAPTMODE0_Pos       /* Capture mode Channel 0: 0 */
	//                           | 0 << TC_CTRLA_CAPTMODE1_Pos /* Capture mode Channel 1: 0 */
	//                           | 0 << TC_CTRLA_COPEN0_Pos    /* Capture Pin 0 Enable: disabled */
	//                           | 0 << TC_CTRLA_COPEN1_Pos    /* Capture Pin 1 Enable: disabled */
	//                           | 0 << TC_CTRLA_CAPTEN0_Pos   /* Capture Channel 0 Enable: disabled */
	//                           | 0 << TC_CTRLA_CAPTEN1_Pos   /* Capture Channel 1 Enable: disabled */
	//                           | 0 << TC_CTRLA_ALOCK_Pos     /* Auto Lock: disabled */
	//                           | 0 << TC_CTRLA_PRESCSYNC_Pos /* Prescaler and Counter Synchronization: 0 */
	//                           | 0 << TC_CTRLA_ONDEMAND_Pos  /* Clock On Demand: disabled */
	//                           | 0 << TC_CTRLA_RUNSTDBY_Pos  /* Run in Standby: disabled */
	//                           | 0 << TC_CTRLA_PRESCALER_Pos /* Setting: 0 */
	//                           | 0x2 << TC_CTRLA_MODE_Pos);  /* Operating Mode: 0x2 */
	sam.TC0_COUNT32.CTRLA.Set(0 |
		0<<sam.TC_COUNT32_CTRLA_CAPTMODE0_Pos |
		0<<sam.TC_COUNT32_CTRLA_CAPTMODE1_Pos |
		0<<sam.TC_COUNT32_CTRLA_COPEN0_Pos |
		0<<sam.TC_COUNT32_CTRLA_COPEN1_Pos |
		0<<sam.TC_COUNT32_CTRLA_CAPTEN0_Pos |
		0<<sam.TC_COUNT32_CTRLA_CAPTEN1_Pos |
		0<<sam.TC_COUNT32_CTRLA_ALOCK_Pos |
		0<<sam.TC_COUNT32_CTRLA_PRESCSYNC_Pos |
		0<<sam.TC_COUNT32_CTRLA_ONDEMAND_Pos |
		0<<sam.TC_COUNT32_CTRLA_RUNSTDBY_Pos |
		0<<sam.TC_COUNT32_CTRLA_PRESCALER_Pos |
		2<<sam.TC_COUNT32_CTRLA_MODE_Pos | // Time Counter Mode : Counter in 32-bit mode
		0)

	//hri_tc_write_CTRLB_reg(TC0,
	//                       0 << TC_CTRLBSET_CMD_Pos           /* Command: 0 */
	//                           | 0 << TC_CTRLBSET_ONESHOT_Pos /* One-Shot: disabled */
	//                           | 0 << TC_CTRLBCLR_LUPD_Pos    /* Setting: disabled */
	//                           | 0 << TC_CTRLBSET_DIR_Pos);   /* Counter Direction: disabled */
	sam.TC0_COUNT32.CTRLBSET.Set(0 |
		0<<sam.TC_COUNT32_CTRLBSET_CMD_Pos |
		0<<sam.TC_COUNT32_CTRLBSET_ONESHOT_Pos |
		0<<sam.TC_COUNT32_CTRLBSET_LUPD_Pos |
		0<<sam.TC_COUNT32_CTRLBSET_DIR_Pos |
		0)

	//hri_tc_write_WAVE_reg(TC0, 1); /* Waveform Generation Mode: 1 */
	sam.TC0_COUNT32.WAVE.Set(sam.TC_COUNT32_WAVE_WAVEGEN_MFRQ)

	//hri_tccount32_write_CC_reg(TC0, 0, 0xbb); /* Compare/Capture Value: 0xbb */
	//sam.TC0_COUNT32.CC[0].Set(48)
	// 1us
	sam.TC0_COUNT32.CC[0].Set(48)

	//hri_tc_write_EVCTRL_reg(
	//    TC0,
	//    0 << TC_EVCTRL_MCEO0_Pos       /* Match or Capture Channel 0 Event Output Enable: disabled */
	//        | 0 << TC_EVCTRL_MCEO1_Pos /* Match or Capture Channel 1 Event Output Enable: disabled */
	//        | 1 << TC_EVCTRL_OVFEO_Pos /* Overflow/Underflow Event Output Enable: enabled */
	//        | 0 << TC_EVCTRL_TCEI_Pos  /* TC Event Input: disabled */
	//        | 0 << TC_EVCTRL_TCINV_Pos /* TC Inverted Event Input: disabled */
	//        | 0);                      /* Event Action: 0 */
	sam.TC0_COUNT32.EVCTRL.Set(0 |
		0<<sam.TC_COUNT32_EVCTRL_MCEO0_Pos |
		0<<sam.TC_COUNT32_EVCTRL_MCEO1_Pos |
		1<<sam.TC_COUNT32_EVCTRL_OVFEO_Pos |
		0<<sam.TC_COUNT32_EVCTRL_TCEI_Pos |
		0<<sam.TC_COUNT32_EVCTRL_TCINV_Pos |
		0<<sam.TC_COUNT32_EVCTRL_EVACT_Pos |
		0)

	//hri_tc_write_CTRLA_ENABLE_bit(TC0, 1 << TC_CTRLA_ENABLE_Pos); /* Enable: enabled */
	sam.TC0_COUNT32.CTRLA.SetBits(sam.TC_COUNT32_CTRLA_ENABLE)
	for sam.TC0_COUNT32.SYNCBUSY.HasBits(sam.TC_COUNT32_SYNCBUSY_ENABLE) {
	}
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
	dmadac.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC0_OVF)
	dmadac.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	desc := dmadac.GetDescriptor()
	desc.UpdateDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&from[0]),
		DST:      unsafe.Pointer(&sam.DAC.DATA[0].Reg),
		SRCINC:   true,
		DSTINC:   false,
		STEPSEL:  true,
		BEATSIZE: 2,
		SIZE:     uint32(len(from)) * 2,
	})
	desc.AddDescriptor(desc)

	a0 := machine.A0
	a0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.DAC0.Configure(machine.DACConfig{})

	dmadac.Start()
	timerInit()

	for {
		led.Toggle()
		dbg5.Toggle()
		time.Sleep(1 * time.Millisecond)
	}
}
