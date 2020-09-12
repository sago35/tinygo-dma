package align

import (
	"testing"
	"unsafe"
)

func TestMake(t *testing.T) {
	for i := 1; i < 32; i++ {
		got := Make(16, i)
		addr := uintptr(unsafe.Pointer(&got[0]))

		if (addr & 0xF) != 0 {
			t.Errorf("align failed %d: got %08X", i, addr)
		}

		if len(got) != i {
			t.Errorf("length failed %d: got %d", i, len(got))
		}
	}
}

func TestAlign16(t *testing.T) {
	input64 := make([]uint64, 64)
	input := *(*[]byte)(unsafe.Pointer(&input64))

	for i := 0; i < 32-14; i++ {
		ptr := uintptr(unsafe.Pointer(&input[i]))

		got := Align(input[i:], 16)
		addr := uintptr(unsafe.Pointer(&got[0]))

		if (addr & 0xF) != 0 {
			t.Errorf("align failed %d: input %08X, got %08X", i, ptr, addr)
		}
	}
}

func TestAlign4(t *testing.T) {
	input := make([]byte, 32)

	for i := 0; i < 32-4; i++ {
		ptr := uintptr(unsafe.Pointer(&input[i]))

		got := Align(input[i:], 4)
		addr := uintptr(unsafe.Pointer(&got[0]))

		if (addr % 4) != 0 {
			t.Errorf("align failed %d: input %08X, got %08X", i, ptr, addr)
		}
	}
}
