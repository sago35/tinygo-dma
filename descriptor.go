// +build atsamd51 atsame5x

package dma

import (
	"unsafe"

	"github.com/sago35/tinygo-dma/align"
)

type DMADescriptor struct {
	btctrl   uint16
	btcnt    uint16
	srcaddr  uint32 //unsafe.Pointer
	dstaddr  uint32 //unsafe.Pointer
	Descaddr uint32 //unsafe.Pointer
}

type DescriptorConfig struct {
	EVOSEL   uint16 // DISABLE, BLOCK, BEAT
	BLOCKACT uint16 // 0:NOACT, 1:INT, 2:SUSPEND, 3:BOTH
	BEATSIZE uint16 // 1, 2, 4
	SRCINC   uint16
	DSTINC   uint16
	STEPSEL  uint16 // false: DST, true: SRC
	STEPSIZE uint16 // 1, 2, 4, 8, 16, 32, 64, 128
	SIZE     uint32
	SRC      unsafe.Pointer
	DST      unsafe.Pointer
}

func NewDescriptor(cfg DescriptorConfig) *DMADescriptor {
	// Descriptors must live in SRAM and must be aligned on a 16-byte boundary.
	slice := align.Make(16, 16)

	ret := (*DMADescriptor)(unsafe.Pointer(&slice[0]))
	ret.UpdateDescriptor(cfg)
	return ret
}

func (d *DMADescriptor) UpdateDescriptor(cfg DescriptorConfig) {
	d.btctrl = (1 << 0) | // VALID: Descriptor Valid
		cfg.EVOSEL | // EVOSEL=DISABLE: Event Output Selection
		cfg.BLOCKACT | // BLOCKACT=NOACT: Block Action
		cfg.BEATSIZE | // BEATSIZE: Beat Size
		cfg.SRCINC | // SRCINC: Source Address Increment Enable
		cfg.DSTINC | // DSTINC: Destination Address Increment Enable
		cfg.STEPSEL | // STEPSEL: Step Selection
		cfg.STEPSIZE // STEPSIZE: Address Increment Step Size
	d.btcnt = uint16(cfg.SIZE >> (uint16(cfg.BEATSIZE) >> DMAC_SRAM_BTCTRL_BEATSIZE_Pos))
	d.Descaddr = 0

	if cfg.STEPSEL == (DMAC_SRAM_BTCTRL_STEPSEL_SRC >> DMAC_SRAM_BTCTRL_STEPSEL_Pos) {
		// STEPSEL == SRC
		if cfg.SRCINC == DMAC_SRAM_BTCTRL_SRCINC_ENABLE {
			d.srcaddr = uint32(uintptr(cfg.SRC) + uintptr((cfg.SIZE)<<(uint32(cfg.STEPSIZE)>>DMAC_SRAM_BTCTRL_STEPSIZE_Pos)))
		} else {
			d.srcaddr = uint32(uintptr(cfg.SRC))
		}
		if cfg.DSTINC == DMAC_SRAM_BTCTRL_DSTINC_ENABLE {
			d.dstaddr = uint32(uintptr(cfg.DST) + uintptr(cfg.SIZE))
		} else {
			d.dstaddr = uint32(uintptr(cfg.DST))
		}
	} else {
		// STEPSEL == DST
		if cfg.SRCINC == DMAC_SRAM_BTCTRL_SRCINC_ENABLE {
			d.srcaddr = uint32(uintptr(cfg.SRC) + uintptr(cfg.SIZE))
		} else {
			d.srcaddr = uint32(uintptr(cfg.SRC))
		}
		if cfg.DSTINC == DMAC_SRAM_BTCTRL_DSTINC_ENABLE {
			d.dstaddr = uint32(uintptr(cfg.DST) + uintptr((cfg.SIZE)<<(uint32(cfg.STEPSIZE)>>DMAC_SRAM_BTCTRL_STEPSIZE_Pos)))
		} else {
			d.dstaddr = uint32(uintptr(cfg.DST))
		}
	}
}

func (d *DMADescriptor) AddDescriptor(next *DMADescriptor) {
	d.Descaddr = uint32(uintptr(unsafe.Pointer(next)))
}
