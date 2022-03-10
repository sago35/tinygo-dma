//go:build atsamd51 || atsame5x
// +build atsamd51 atsame5x

package dma

import (
	"unsafe"

	"github.com/sago35/tinygo-dma/align"
)

type DMADescriptor struct {
	Btctrl   uint16
	Btcnt    uint16
	Srcaddr  uint32 //unsafe.Pointer
	Dstaddr  uint32 //unsafe.Pointer
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
	DESC     unsafe.Pointer
}

func NewDescriptor() *DMADescriptor {
	// Descriptors must live in SRAM and must be aligned on a 16-byte boundary.
	slice := align.Make(16, 16)

	return (*DMADescriptor)(unsafe.Pointer(&slice[0]))
}

func (d *DMADescriptor) UpdateDescriptor(cfg DescriptorConfig) {
	d.Btctrl = (1 << 0) | // VALID: Descriptor Valid
		cfg.EVOSEL | // EVOSEL=DISABLE: Event Output Selection
		cfg.BLOCKACT | // BLOCKACT=NOACT: Block Action
		cfg.BEATSIZE | // BEATSIZE: Beat Size
		cfg.SRCINC | // SRCINC: Source Address Increment Enable
		cfg.DSTINC | // DSTINC: Destination Address Increment Enable
		cfg.STEPSEL | // STEPSEL: Step Selection
		cfg.STEPSIZE // STEPSIZE: Address Increment Step Size
	d.Btcnt = uint16(cfg.SIZE >> (uint16(cfg.BEATSIZE) >> DMAC_SRAM_BTCTRL_BEATSIZE_Pos))
	d.Descaddr = 0

	if cfg.STEPSEL == (DMAC_SRAM_BTCTRL_STEPSEL_SRC >> DMAC_SRAM_BTCTRL_STEPSEL_Pos) {
		// STEPSEL == SRC
		if cfg.SRCINC == DMAC_SRAM_BTCTRL_SRCINC_ENABLE {
			d.Srcaddr = uint32(uintptr(cfg.SRC) + uintptr((cfg.SIZE)<<(uint32(cfg.STEPSIZE)>>DMAC_SRAM_BTCTRL_STEPSIZE_Pos)))
		} else {
			d.Srcaddr = uint32(uintptr(cfg.SRC))
		}
		if cfg.DSTINC == DMAC_SRAM_BTCTRL_DSTINC_ENABLE {
			d.Dstaddr = uint32(uintptr(cfg.DST) + uintptr(cfg.SIZE))
		} else {
			d.Dstaddr = uint32(uintptr(cfg.DST))
		}
	} else {
		// STEPSEL == DST
		if cfg.SRCINC == DMAC_SRAM_BTCTRL_SRCINC_ENABLE {
			d.Srcaddr = uint32(uintptr(cfg.SRC) + uintptr(cfg.SIZE))
		} else {
			d.Srcaddr = uint32(uintptr(cfg.SRC))
		}
		if cfg.DSTINC == DMAC_SRAM_BTCTRL_DSTINC_ENABLE {
			d.Dstaddr = uint32(uintptr(cfg.DST) + uintptr((cfg.SIZE)<<(uint32(cfg.STEPSIZE)>>DMAC_SRAM_BTCTRL_STEPSIZE_Pos)))
		} else {
			d.Dstaddr = uint32(uintptr(cfg.DST))
		}
	}

	d.Descaddr = uint32(uintptr(cfg.DESC))
}

func (d *DMADescriptor) AddDescriptor(next *DMADescriptor) {
	d.Descaddr = uint32(uintptr(unsafe.Pointer(next)))
}
