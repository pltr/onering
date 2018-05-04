//+build debug

package onering

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

const DEBUG = true

func PrintSizes() {
	fmt.Printf("ring size: %v bytes\n", unsafe.Sizeof(ring{}))
	fmt.Printf("multi size: %v bytes\n", unsafe.Sizeof(multi{}))
	fmt.Printf("SPSC size: %v bytes\n", unsafe.Sizeof(SPSC{}))
	fmt.Printf("SPMC size: %v bytes\n", unsafe.Sizeof(SPMC{}))
	fmt.Printf("MPSC size: %v bytes\n", unsafe.Sizeof(MPSC{}))
	fmt.Printf("MPMC size: %v bytes\n", unsafe.Sizeof(MPMC{}))
}

var visited = struct {
	m  map[string]int
	mu sync.Mutex
}{
	m: map[string]int{},
}

func checkVisited(path string) bool {
	visited.mu.Lock()
	defer visited.mu.Unlock()
	visited.m[path]++
	return visited.m[path] > 1
}

func checkGetType(i interface{}, path string) {
	if checkVisited(path) {
		return
	}
	t1 := reflect.TypeOf(i)
	if t1.Kind() != reflect.Ptr {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: %s calls queue.Get(**%[1]v) argument type; Expected **%[1]v found: %v\n", path, t1))
		os.Exit(1)
	}
	t2 := t1.Elem()
	if kind := t2.Kind(); kind != reflect.Ptr && kind != reflect.UnsafePointer {
		os.Stderr.WriteString(fmt.Sprintf("ERROR: %[1]s calls queue.Get(*%[2]v) with an illegal argument type %[2]v\n", path, t1))
		os.Exit(1)
	}
}

func checkPutType(i interface{}, path string) {
	if checkVisited(path) {
		return
	}
	t1 := reflect.TypeOf(i)
	if t1.Kind() != reflect.Ptr {
		fmt.Printf("WARNING: %s calls .Put() with a non pointer type: <%s:%v>. ", path, t1, i)
		fmt.Println("This will cause a memory allocation")
	}
}

func checkFuncType(i interface{}, path string) {
	if checkVisited(path) {
		return
	}
	var errors = 0
	t1 := reflect.TypeOf(i)
	if t1.Kind() != reflect.Func || t1.NumIn() != 2 {
		fmt.Printf("ERROR: %s calls .Consume(func(Iter, *T)) with an illegal argument type: %v\n", path, t1)
		os.Exit(1)
	}
	a1, a2 := t1.In(0), t1.In(1)
	if a1.String() != "onering.Iter" {
		fmt.Printf("ERROR: %s calls .Consume(func(Iter, *T)) with an illegal first argument type: %v\n", path, a1)
		errors++
	}

	if a2.Kind() != reflect.Ptr {
		fmt.Printf("ERROR: %s calls .Consume(func(Iter, *T)) with a non-pointer second argument type : %v\n", path, a2)
		errors++
	}

	if t1.NumOut() > 0 {
		fmt.Printf("WARNING: %s calls .Consume(func(Iter, *T)) with a function that return values\n", path)
	}
	if errors > 0 {
		os.Exit(1)
	}

}

func getCallerPath() (path string) {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		path = fmt.Sprintf("%s:%d", file, line)
	} else {
		path = "<unknown>"
	}
	return
}
