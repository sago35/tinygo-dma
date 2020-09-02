// +build atsamd51

package dma

import (
	"unsafe"
)

type DescriptorConfig struct {
	EVOSEL   byte // DISABLE, BLOCK, BEAT
	BLOCKACT byte // 0:NOACT, 1:INT, 2:SUSPEND, 3:BOTH
	BEATSIZE byte // 1, 2, 4
	SRCINC   bool
	DSTINC   bool
	STEPSEL  bool // false: DST, true: SRC
	STEPSIZE byte // 1, 2, 4, 8, 16, 32, 64, 128
	SIZE     uint32
	SRC      unsafe.Pointer
	DST      unsafe.Pointer
}

func NewDescriptor(cfg DescriptorConfig) *DMADescriptor {

	si := uint16(0)
	if cfg.SRCINC {
		si = 1
	}

	di := uint16(0)
	if cfg.DSTINC {
		di = 1
	}

	ssel := uint16(0)
	if cfg.STEPSEL {
		ssel = 1
	}

	ss := uint16(0)
	switch cfg.STEPSIZE {
	case 0, 1:
		ss = 0
	case 2:
		ss = 1
	case 4:
		ss = 2
	case 8:
		ss = 3
	case 16:
		ss = 4
	case 32:
		ss = 5
	case 64:
		ss = 6
	case 128:
		ss = 7
	default:
		// TODO: error
		ss = 0
	}

	bs := uint16(0)
	switch cfg.BEATSIZE {
	case 0, 1:
		bs = 0
	case 2:
		bs = 1
	case 4:
		bs = 2
	default:
		// TODO: error
		bs = 0
	}

	//go:align 16
	var ret DMADescriptor

	ret.btctrl = (1 << 0) | // VALID: Descriptor Valid
		(uint16(cfg.EVOSEL) << 1) | // EVOSEL=DISABLE: Event Output Selection
		(uint16(cfg.BLOCKACT) << 3) | // BLOCKACT=NOACT: Block Action
		(bs << 8) | // BEATSIZE: Beat Size
		(si << 10) | // SRCINC: Source Address Increment Enable
		(di << 11) | // DSTINC: Destination Address Increment Enable
		(ssel << 12) | // STEPSEL: Step Selection
		(ss << 13) // STEPSIZE: Address Increment Step Size
	ret.btcnt = uint16(cfg.SIZE >> bs)
	ret.Descaddr = 0

	if cfg.STEPSEL {
		// STEPSEL == SRC
		ret.srcaddr = uint32(uintptr(cfg.SRC) + uintptr((cfg.SIZE)<<ss))
		ret.dstaddr = uint32(uintptr(cfg.DST) + uintptr(cfg.SIZE))
	} else {
		// STEPSEL == DST
		ret.srcaddr = uint32(uintptr(cfg.SRC) + uintptr(cfg.SIZE))
		ret.dstaddr = uint32(uintptr(cfg.DST) + uintptr((cfg.SIZE)<<ss))
	}

	return &ret
}

func (d *DMADescriptor) AddDescriptor(next *DMADescriptor) {
	d.Descaddr = uint32(uintptr(unsafe.Pointer(next)))
}
