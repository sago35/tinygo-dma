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
		6<<sam.TC_COUNT32_CTRLA_PRESCALER_Pos |
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
	//sam.TC0_COUNT32.CC[0].Set(48)

	// 1000 ms
	sam.TC0_COUNT32.CC[0].Set(48 * 1000 * 1000)
	//sam.TC0_COUNT32.CC[0].Set(0xFFFFFFFF)

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

func evsysInit() {
	// Turn on clock to EVSYS (7)
	sam.MCLK.APBBMASK.SetBits(sam.MCLK_APBBMASK_EVSYS_)

	// 48Mhz
	sam.GCLK.PCHCTRL[11].Set((sam.GCLK_PCHCTRL_GEN_GCLK1 << sam.GCLK_PCHCTRL_GEN_Pos) |
		sam.GCLK_PCHCTRL_CHEN)

	// event_system_init();

	// hri_evsys_write_CHANNEL_reg(EVSYS, i, channel_confs[i]);
	// #define CONF_EVGEN_0 12
	// 12 : RTC period7
	// 0x49 : TC0 Overflow
	//sam.EVSYS.CHANNEL[12].CHANNEL.SetBits(sam.EVSYS_CHANNEL_CHANNEL_RUNSTDBY | sam.EVSYS_CHANNEL_CHANNEL_PATH_ASYNCHRONOUS<<sam.EVSYS_CHANNEL_CHANNEL_PATH_Pos | (0x49)<<sam.EVSYS_CHANNEL_CHANNEL_EVGEN_Pos)
	sam.EVSYS.CHANNEL[0].CHANNEL.SetBits(sam.EVSYS_CHANNEL_CHANNEL_RUNSTDBY | sam.EVSYS_CHANNEL_CHANNEL_PATH_ASYNCHRONOUS<<sam.EVSYS_CHANNEL_CHANNEL_PATH_Pos | (0x49)<<sam.EVSYS_CHANNEL_CHANNEL_EVGEN_Pos)

	// hri_evsys_write_USER_reg(EVSYS, i, user_mux_confs[i]);
	// #define CONF_CHANNEL_57 1
	// 55 : ADC0 START
	// 57 : ADC1 START
	// 1 : Channel 0 selected
	//sam.EVSYS.USER[55].Set(13 << sam.EVSYS_USER_CHANNEL_Pos)
	sam.EVSYS.USER[55].Set(1 << sam.EVSYS_USER_CHANNEL_Pos)

	//sam.EVSYS.USER[55].Set(2 << sam.EVSYS_USER_CHANNEL_Pos)
	//sam.EVSYS.USER[56].Set(1 << sam.EVSYS_USER_CHANNEL_Pos)
	//sam.EVSYS.USER[57].Set(1 << sam.EVSYS_USER_CHANNEL_Pos)
	//sam.EVSYS.USER[58].Set(1 << sam.EVSYS_USER_CHANNEL_Pos)
}

func main() {
	initDbg()

	machine.InitADC()
	time.Sleep(2 * time.Second)

	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})

	evsysInit()

	dmaadc := dma.NewDMA(func(d *dma.DMA) {
		dbg6.Toggle()
		return
	})
	dmaadc.SetTrigger(dma.DMAC_CHANNEL_CHCTRLA_TRIGSRC_ADC0_RESRDY)
	dmaadc.SetTriggerAction(sam.DMAC_CHANNEL_CHCTRLA_TRIGACT_BURST)

	to := uint16(0)
	//to := [1024]uint16{}

	desc := dmaadc.GetDescriptor()
	desc.UpdateDescriptor(dma.DescriptorConfig{
		SRC:      unsafe.Pointer(&sam.ADC0.RESULT.Reg),
		DST:      unsafe.Pointer(&to),
		SRCINC:   dma.DMAC_SRAM_BTCTRL_SRCINC_DISABLE,
		DSTINC:   dma.DMAC_SRAM_BTCTRL_DSTINC_ENABLE,
		STEPSEL:  dma.DMAC_SRAM_BTCTRL_STEPSEL_DST,
		BEATSIZE: dma.DMAC_SRAM_BTCTRL_BEATSIZE_HWORD,
		//BLOCKACT: dma.DMAC_SRAM_BTCTRL_BLOCKACT_INT,
		SIZE: 2,
		//SIZE: 2 * 8,
	})
	desc.AddDescriptor(desc)

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

		//bus.CTRLA.SetBits(sam.ADC_CTRLA_RUNSTDBY)
		bus.EVCTRL.SetBits(sam.ADC_EVCTRL_STARTEI)

		// Enable ADC
		bus.CTRLA.SetBits(sam.ADC_CTRLA_ENABLE)
		for bus.SYNCBUSY.HasBits(sam.ADC_SYNCBUSY_ENABLE) {
		}

		fmt.Printf("ADC.CTRLA : %08X\r\n", bus.CTRLA.Get())
		fmt.Printf("EVSYS.CHINTFLAG[0] : %02X\r\n", sam.EVSYS.CHANNEL[0].CHINTFLAG.Get())
		fmt.Printf("MCLCK.APB : %08X %08X %08X %08X\r\n", sam.MCLK.APBAMASK.Get(), sam.MCLK.APBBMASK.Get(), sam.MCLK.APBCMASK.Get(), sam.MCLK.APBDMASK.Get())
	}

	dmaadc.Start()
	timerInit()

	for {
		if false {
			val := sensor.Get()
			fmt.Printf("%04X\r\n", val)
		}
		// Start conversion
		//to = bus.RESULT.Get()
		//sam.EVSYS.SWEVT.Set(0x00001000)
		//bus.SWTRIG.SetBits(sam.ADC_SWTRIG_START)
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%04X\r\n", to)
		//fmt.Printf("%04X %04X %04X %04X %04X %04X %04X %04X\r\n", to[0], to[1], to[2], to[3], to[4], to[5], to[6], to[7])
		//fmt.Printf("%04X %02X %08X %08X\r\n", to, bus.STATUS.Get(), bus.SYNCBUSY.Get(),
		//	bus.DSEQSTAT.Get(),
		//)
		//fmt.Printf("EVSYS.CHINTFLAG[0] : %02X, CHSTATUS[0] : %02X \r\n", sam.EVSYS.CHANNEL[0].CHINTFLAG.Get(), sam.EVSYS.CHANNEL[0].CHSTATUS.Get())
		//fmt.Printf("ADC0.INTFLAG[0] : %02X\r\n", sam.ADC0.INTFLAG.Get())
		led.Toggle()
		//time.Sleep(time.Millisecond * 5000)
	}
}
