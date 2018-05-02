package onering

import "unsafe"


const BatchExp = 8
const MaxBatch = (1 << BatchExp) - 1
const spin = 512 - 1 // not used at the moment

type Injector func(dst unsafe.Pointer, src unsafe.Pointer)

var (
	Pointer Injector = func(dst unsafe.Pointer, src unsafe.Pointer) {
		*(*unsafe.Pointer)(dst) = src
	}
	Interface Injector = func(dst unsafe.Pointer, src unsafe.Pointer) {
		(*IHeader)(dst).D = src
	}
)


type Batch interface {
	Next(interface{}) bool
}