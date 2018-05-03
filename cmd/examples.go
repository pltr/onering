package main

import (
	"github.com/pltr/onering"
	"fmt"
)

func simple() {
	var queue = onering.New{Size: 8192}.MPMC()

	var src = int64(5)
	queue.Put(&src)
	queue.Close()
	var dst *int64
	// .Get expects a pointer to a pointer
	for queue.Get(&dst) {
		if *dst != src {
			panic("i don't know what's going on")
		}
	}
	fmt.Println("Yay, get works")
}

func batching() {
	var queue = onering.New{Size: 8192}.SPSC()

	var src = int64(5)
	queue.Put(&src)
	queue.Close()

	queue.Consume(func(dst *int64) {
		if *dst != src {
			panic("i don't know what's going on")
		}
	})
	fmt.Println("Yay, batching works")
}


func main() {
	simple()
	batching()
}
