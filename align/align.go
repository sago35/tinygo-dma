package align

import (
	"unsafe"
)

func Make(alignBytes, length int) []byte {
	slice := make([]byte, length+alignBytes-1)
	slice = Align(slice, alignBytes)
	return slice[:length]
}

func Align(in []byte, alignBytes int) []byte {
	ptr := uintptr(unsafe.Pointer(&in[0]))
	slice := in[(alignBytes-1)-((int(ptr)+alignBytes-1)%alignBytes):]
	return slice
}
