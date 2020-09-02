// +build atsamd51

package dma

const (
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_DISABLE    = 0x00 // Only software/event triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_RTC        = 0x01 // TIMESTAMP DMA RTC timestamp trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_DSU_DCC0   = 0x02 // DMAC ID for DCC0 register
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_DSU_DCC1   = 0x03 // DMAC ID for DCC1 register
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM0_RX = 0x04 // Index of DMA RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM0_TX = 0x05 // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM1_RX = 0x06 // Index of DMA RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM1_TX = 0x07 // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM2_RX = 0x08 // Index of DMA RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM2_TX = 0x09 // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM3_RX = 0x0A // Index of DMA RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM3_TX = 0x0B // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM4_RX = 0x0C // Index of DMA RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM4_TX = 0x0D // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM5_RX = 0x0E // Index of DMA RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM5_TX = 0x0F // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM6_RX = 0x10 // Index of DMA RX trigger

	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM6_TX = 0x11        // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM7_RX = 0x12        // Index of DMA RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_SERCOM7_TX = 0x13        // Index of DMA TX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_CAN0_DEBUG = 0x14        // DMA CAN Debug Req
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_CAN1_DEBUG = 0x15        // DMA CAN Debug Req
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC0_OVF   = 0x16        // DMA overflow/underflow/retrigger trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC0_MC    = 0x1C - 0x17 // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC1_OVF   = 0x1D        // DMA overflow/underflow/retrigger trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC1_MC    = 0x21 - 0x1E // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC2_OVF   = 0x22        // DMA overflow/underflow/retrigger trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC2_MC    = 0x25 - 0x23 // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC3_OVF   = 0x26        // DMA overflow/underflow/retrigger trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC3_MC    = 0x28 - 0x27 // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC4_OVF   = 0x29        // DMA overflow/underflow/retrigger trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TCC4_MC    = 0x2B - 0x2A // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC0_OVF    = 0x2C        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC0_MC     = 0x2E - 0x2D // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC1_OVF    = 0x2F        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC1_MC     = 0x31 - 0x30 // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC2_OVF    = 0x32        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC2_MC     = 0x34 - 0x33 // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC3_OVF    = 0x35        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC3_MC     = 0x37 - 0x36 // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC4_OVF    = 0x38        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC4_MC     = 0x3A - 0x39 // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC5_OVF    = 0x3B        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC5_MC     = 0x3D - 0x3C // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC6_OVF    = 0x3E        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC6_MC     = 0x40 - 0x3F // Indexes of DMA Match/Compare triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC7_OVF    = 0x41        // Indexes of DMA Overflow trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_TC7_MC     = 0x43 - 0x41 // Indexes of DMA Match/Compare triggers

	DMAC_CHANNEL_CHCTRLA_TRIGSRC_ADC0_RESRDY = 0x44        // index of DMA RESRDY trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_ADC0_SEQ    = 0x45        // Index of DMA SEQ trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_ADC1_RESRDY = 0x46        // Index of DMA RESRDY trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_ADC1_SEQ    = 0x47        // Index of DMA SEQ trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_DAC_EMPTY0  = 0x48        // DMA DAC0 Empty Req
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_DAC_EMPTY1  = 0x49        // DMA DAC1 Empty Req
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_DAC_RESRDY0 = 0x4A        // DMA DAC0 Result Ready Req
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_DAC_RESRDY1 = 0x4B        // DMA DAC1 Result Ready Req
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_I2S_RX      = 0x4D - 0x4C // Indexes of DMA RX triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_I2S_TX      = 0x4F - 0x4E // Indexes of DMA TX triggers
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_PCC_RX      = 0x50        // Indexes of PCC RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_AES_WR      = 0x51        // DMA DATA Write trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_AES_RD      = 0x52        // DMA DATA Read trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_QSPI_RX     = 0x53        // Indexes of QSPI RX trigger
	DMAC_CHANNEL_CHCTRLA_TRIGSRC_QSPI_TX     = 0x54        // Indexes of QSPI TX trigger
)
