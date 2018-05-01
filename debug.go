package onering

import (
	"fmt"
	"unsafe"
)

func PrintSizes() {
	fmt.Printf("ring size: %v bytes\n", unsafe.Sizeof(ring{}))
	fmt.Printf("multi size: %v bytes\n", unsafe.Sizeof(multi{}))
	fmt.Printf("SPSC size: %v bytes\n", unsafe.Sizeof(SPSC{}))
	fmt.Printf("SPMC size: %v bytes\n", unsafe.Sizeof(SPMC{}))
	fmt.Printf("MPSC size: %v bytes\n", unsafe.Sizeof(MPSC{}))
	fmt.Printf("MPMC size: %v bytes\n", unsafe.Sizeof(MPMC{}))
}
