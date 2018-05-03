package onering

import "unsafe"

type iface struct {
	t, d unsafe.Pointer
}

func extractptr(i interface{}) unsafe.Pointer {
	return (*iface)(unsafe.Pointer(&i)).d
}

func extractfn(i interface{}) func(unsafe.Pointer) {
	var ptr = (*iface)(unsafe.Pointer(&i)).d
	return *(*func(unsafe.Pointer))(unsafe.Pointer(&ptr))
}

func inject(i interface{}, ptr unsafe.Pointer) {
	var v = (*unsafe.Pointer)((*iface)(unsafe.Pointer(&i)).d)
	*v = ptr
}
