package main

import (
	"fmt"
	"github.com/pltr/onering"
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
	queue.Put(6) // WARNING: this will allocate memory on the heap and copy the value into it
	queue.Close()

	queue.Consume(func(it onering.Iter, dst *int64) {
		if *dst != src {
			panic("i don't know what's going on")
		}
		it.Stop()
	})
	// still one element left in the queue
	var dst *int64
	// Get will always expect a pointer to a pointer
	if !queue.Get(&dst) || *dst != 6 {
		panic("uh oh")
	}
	fmt.Println("Yay, batching works")
}

func main() {
	simple()
	batching()
}
