package onering

import "unsafe"

type IHeader struct {
	T, D unsafe.Pointer
}

func extractptr(i interface{}) unsafe.Pointer {
	return (*IHeader)(unsafe.Pointer(&i)).D
}