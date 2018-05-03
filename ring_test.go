package onering

import (
	"runtime"
	"sync"
	"testing"
	"fmt"
)

const MULTI = 100

func mknumslice(n int) []int {
	var s = make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func BenchmarkRingSPSC_Get(b *testing.B) {
	var ring = New{8192}.SPSC()
	var wg sync.WaitGroup
	wg.Add(2)
	type T struct {
		i int
	}
	var v = T{5}
	b.ResetTimer()
	go func(n int) {
		runtime.LockOSThread()
		for i := 0; i < b.N; i++ {
			ring.Put(&v)
		}
		ring.Close()
		wg.Done()
	}(b.N)

	go func(n int) {
		runtime.LockOSThread()
		var v *T
		for i := 0; ring.Get(&v); i++ {
			_ = *v
		}
		wg.Done()
	}(b.N)

	wg.Wait()

}

func BenchmarkRingSPSC_Consume(b *testing.B) {
	var numbers = mknumslice(b.N)
	var ring = New{Size: 8192}.SPSC()
	var wg sync.WaitGroup
	wg.Add(2)
	//pp := runtime.GOMAXPROCS(8)
	b.ResetTimer()
	go func(n int) {
		runtime.LockOSThread()
		for i := range numbers {
			ring.Put(&numbers[i])
		}
		ring.Close()
		wg.Done()
	}(b.N)
	go func(n int) {
		runtime.LockOSThread()
		var i int
		ring.Consume(func(v *int) {
			if *v != i {
				b.Fatalf("Expected %d got %d", i, v)
			}
			i++
		})
		wg.Done()
	}(b.N)

	wg.Wait()
	//runtime.GOMAXPROCS(pp)
}

func BenchmarkRingSPMC(b *testing.B) {
	var numbers = mknumslice(b.N)
	var ring = New{Size:8192}.SPMC()
	var wg sync.WaitGroup
	var readers = MULTI
	wg.Add(readers + 1)
	//pp := runtime.GOMAXPROCS(8)
	for c := 0; c < readers; c++ {
		go func(c int) {
			var v *int64
			for ring.Get(&v) {
				_ = *v
			}
			wg.Done()
		}(c)
	}
	go func(n int) {
		runtime.LockOSThread()
		for i := range numbers {
			ring.Put(&numbers[i])
		}
		ring.Close()
		wg.Done()
	}(b.N)
	wg.Wait()
	//runtime.GOMAXPROCS(pp)
}

func BenchmarkRingMPSC_Get(b *testing.B) {
	var ring = New{Size:8192}.MPSC()
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = MULTI
	wg.Add(producers + 1)
	for p := 0; p < producers; p++ {
		go func(p int) {
			var total = b.N/producers + 1
			var numbers = mknumslice(total)
			for i := range numbers {
				ring.Put(&numbers[i])
			}
			wg.Done()
		}(p)
	}
	go func(n int) {
		runtime.LockOSThread()
		var v *int
		for ring.Get(&v) {
			n--
			if n <= 0 {
				ring.Close()
			}
		}
		wg.Done()
	}(b.N)

	wg.Wait()
	//runtime.GOMAXPROCS(pp)
}
//
func BenchmarkRingMPSC_Batch(b *testing.B) {
	var ring = New{Size:8192}.MPSC()
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = MULTI
	wg.Add(producers + 1)
	for p := 0; p < producers; p++ {
		go func(p int) {
			var total = b.N/producers + 1
			var numbers = mknumslice(total)
			for i := range numbers {
				ring.Put(&numbers[i])
			}
			wg.Done()
		}(p)
	}
	go func(n int) {
		runtime.LockOSThread()
		ring.Consume(func(i *int) {
			n--
			if n <= 0 {
				ring.Close()
			}
		})
		wg.Done()
	}(b.N)

	wg.Wait()
	//runtime.GOMAXPROCS(pp)
}

func BenchmarkRingMPMC_Get(b *testing.B) {
	var ring = New{Size:8192}.MPMC()
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = 50
	wg.Add(producers * 2)
	//var total = int32(b.N)
	for p := 0; p < producers; p++ {
		go func(p int) {
			var size = b.N/producers + 1
			var numbers = mknumslice(size)
			for i := range numbers {
				ring.Put(&numbers[i])
			}
			wg.Done()
		}(p)
	}

	for p := 0; p < producers; p++ {
		go func(c int) {
			var v int64
			var total = b.N/producers + 1
			for i := 0; ring.Get(&v); {
				if i++; i == total {
					break
				}
			}
			wg.Done()
		}(p)
	}
	wg.Wait()
	//runtime.GOMAXPROCS(pp)
}

func BenchmarkChanMPMC(b *testing.B) {
	var ch = make(chan int64, 8192)
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = 50
	wg.Add(producers * 2)
	for p := 0; p < producers; p++ {
		go func(p int) {
			var size = b.N/producers + 1
			for i := 0; i < size; i++ {
				ch <- int64(i)
			}
			wg.Done()
		}(p)
	}

	for p := 0; p < producers; p++ {
		go func(c int) {
			for n := b.N/producers + 1; n > 0; n-- {
				v := <-ch
				_ = v
			}

			wg.Done()
		}(p)
	}
	wg.Wait()
	//runtime.GOMAXPROCS(pp)
}

func BenchmarkChan(b *testing.B) {
	b.Run("SPSC", func(b *testing.B) {
		q := make(chan int64, 8192)

		var wg sync.WaitGroup
		wg.Add(2)

		go func(n int) {
			runtime.LockOSThread()
			for i := 0; i < n; i++ {
				q <- int64(i)
			}
			wg.Done()
		}(b.N)

		go func(n int) {
			runtime.LockOSThread()
			for i := 0; i < n; i++ {
				<-q
			}
			wg.Done()
		}(b.N)

		wg.Wait()
	})

	for i := 64; i <= 64; i <<= 1 {
		producers := i
		b.Run(fmt.Sprintf("SPMC-%d", producers), func(b *testing.B) {
			single := b.N/producers + 1
			total := single * producers
			q := make(chan int64, 8192)
			var wg sync.WaitGroup
			wg.Add(producers + 1)
			for p := 0; p < producers; p++ {
				go func(n int) {
					for i := 0; i < single; i++ {
						<-q
					}
					wg.Done()
				}(b.N)
			}
			go func(n int) {
				runtime.LockOSThread()
				for i := 0; i < total; i++ {
					q <- int64(i)
				}
				wg.Done()
			}(b.N)
			wg.Wait()
		})
	}

	for i := 64; i <= 64; i <<= 1 {
		producers := i
		b.Run(fmt.Sprintf("MPSC-%d", producers), func(b *testing.B) {
			single := b.N/producers + 1
			total := single * producers
			q := make(chan int64, 8192)
			var wg sync.WaitGroup
			wg.Add(producers + 1)
			for p := 0; p < producers; p++ {
				go func(n int) {
					for i := 0; i < single; i++ {
						q <- int64(i)
					}
					wg.Done()
				}(b.N)
			}
			go func(n int) {
				runtime.LockOSThread()
				for i := 0; i < total; i++ {
					<-q
				}
				wg.Done()
			}(b.N)
			wg.Wait()
		})
	}
}

//// courtesy or Egon Elbre
func TestXOneringSPMC(t *testing.T) {
	const P = 4
	const N = 100
	var q = New{4}.SPMC()

	var wg sync.WaitGroup
	wg.Add(P + 1)
	go func() {
		defer wg.Done()
		for i := 0; i < N*P; i++ {
			q.Put(int64(i + 1))
		}
	}()

	errs := make(chan error)
	go func() {
		wg.Wait()
		close(errs)
	}()

	for i := 0; i < P; i++ {
		go func(p int) {
			defer wg.Done()
			var lastSeen int64
			for i := 0; i < N; i++ {
				var v *int64
				if !q.Get(&v) {
					errs <- fmt.Errorf("failed get")
				}
				//fmt.Println(p, v)
				if *v <= lastSeen {
					errs <- fmt.Errorf("got %v last seen %v on producer %v", v, lastSeen, p)
				}
				lastSeen = *v
			}
		}(i)
	}

	for err := range errs {
		t.Fatal(err)
	}
}

func TestXOneringMPSCBatch(t *testing.T) {
	var q = New{2}.MPSC()
	const P = 4
	const C = 2
	var wg sync.WaitGroup
	wg.Add(P + 1)
	for id := 0; id < P; id++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < C; i++ {
				q.Put(int64(id<<32 | i))
			}
		}(id)
	}

	go func() {
		defer wg.Done()
		total := C * P
		q.Consume(func(val *int64) {
			total--
			if total == 0 {
				q.Close()
			} else if total < 0 {
				t.Fatal("invalid value")
				q.Close()
				return
			}
		})
	}()
	wg.Wait()
}
//
func TestRingMPMC_Get(t *testing.T) {
	var numbers = mknumslice(1000)
	var ring = New{Size:8192}.MPMC()
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = 4
	wg.Add(producers * 2)
	var N = len(numbers)
	for p := 0; p < producers; p++ {
		go func(p int) {
			runtime.LockOSThread()
			for i := 0; i < N; i++ {
				ring.Put(&numbers[i])
			}
			wg.Done()
		}(p)
	}
	var ch = make(chan int, N*producers)
	for p := 0; p < producers; p++ {
		go func(c int) {
			runtime.LockOSThread()
			var v *int
			var i int
			for ring.Get(&v) {
				i++
				ch <- *v
				if i == N {
					break
				}
			}
			wg.Done()
		}(p)
	}

	var m = map[int]int{}
	for i := 0; i < N*4; i++ {
		v := <-ch
		m[v]++
	}

	for k, v := range m {
		if v != producers {
			t.Fatalf("%v(%v) != 4: %v", k, v, m)
		}
	}
	fmt.Println("waiting")
	wg.Wait()

	//runtime.GOMAXPROCS(pp)
}


func TestRingMPSC_Get(t *testing.T) {
	var numbers = mknumslice(1000)
	var ring = New{Size:8192}.MPSC()
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = 50
	wg.Add(producers+1)
	var N = len(numbers) * producers
	for p := 0; p < producers; p++ {
		go func(p int) {
			for i := range numbers {
				ring.Put(&numbers[i])
			}
			wg.Done()
		}(p)
	}
	var ch = make(chan int, N)
	go func() {
		runtime.LockOSThread()
		var v *int
		var i int
		for ring.Get(&v) {
			ch <- *v
			if i++; i == N {
				break
			}
		}
		wg.Done()
	}()

	var m = map[int]int{}
	for i := 0; i < N; i++ {
		v := <-ch
		m[v]++
	}

	for k, v := range m {
		if v != producers {
			t.Fatalf("%v(%v) != 4: %v", k, v, m)
		}
	}
	wg.Wait()

	//runtime.GOMAXPROCS(pp)
}