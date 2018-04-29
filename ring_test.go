package onering

import (
	"testing"
	"sync"
	"fmt"
	"runtime"
	"time"
)

func BenchmarkSPSC_Get(b *testing.B) {
	var ring SPSC
	ring.Init(8192)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(n int) {
		runtime.LockOSThread()
		for i := 0; i < n; i++ {
			ring.Put(int64(i))
		}
		ring.Close()
		wg.Done()
	}(b.N)

	go func(n int64) {
		runtime.LockOSThread()
		var i, v int64
		for ring.Get(&v) {
			if v != i {
				fmt.Printf("Expected %d got %d", i, v)
				panic(v)
			}
			i++
		}
		wg.Done()
	}(int64(b.N))

	wg.Wait()

}



func BenchmarkSPSC_Batch(b *testing.B) {
	var ring SPSC
	ring.Init(8192)
	var wg sync.WaitGroup
	wg.Add(2)
	go func(n int) {
		runtime.LockOSThread()
		for i := 0; i < n; i++ {
			ring.Put(int64(i))
		}
		wg.Done()
	}(b.N)
	go func(n int) {
		runtime.LockOSThread()
		ring.Consume(func(i int64) {
			n--
			if n <= 0 {
				ring.Close()
			}
		})
		wg.Done()
	}(b.N)

	wg.Wait()
}


func BenchmarkSPMC(b *testing.B) {
	var ring SPMC
	ring.Init(8192)
	var wg sync.WaitGroup
	var readers = 64
	wg.Add(readers+1)
	pp := runtime.GOMAXPROCS(8)
	for c := 0; c < readers; c++ {
		go func(c int) {
			var i int64
			for ring.Get(&i) {
				_ = i
			}
			wg.Done()
		}(c)
	}
	time.Sleep(1000)
	go func(n int) {
		runtime.LockOSThread()
		for i := 0; i < n; i++ {
			ring.Put(int64(i))
		}
		ring.Close()
		wg.Done()
	}(b.N)
	wg.Wait()
	runtime.GOMAXPROCS(pp)
}

func BenchmarkMPSC_Get(b *testing.B) {
	var ring MPSC
	ring.Init(8192)
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = 64
	wg.Add(producers+1)
	for p := 0; p < producers; p++ {
		go func(p int) {
			var total = b.N / producers + 1
			var start = p * total
			var end = start + total
			for i := start; i < end; i++ {
				ring.Put(int64(i))
			}
			wg.Done()
		}(p)
	}
	go func(n int) {
		runtime.LockOSThread()
		var v int64
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

func BenchmarkMPSC_Batch(b *testing.B) {
	var ring MPSC
	ring.Init(8192)
	var wg sync.WaitGroup
	//pp := runtime.GOMAXPROCS(8)
	var producers = 64
	wg.Add(producers+1)
	for p := 0; p < producers; p++ {
		go func(p int) {
			var total = b.N / producers + 1
			var start = p * total
			var end = start + total
			for i := start; i < end; i++ {
				ring.Put(int64(p))
				//time.Sleep(100 * time.Nanosecond)
			}
			wg.Done()
		}(p)
	}
	go func(n int) {
		runtime.LockOSThread()
		ring.Consume(func(i int64) {
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


func BenchmarkChanSPSC(b *testing.B) {
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
}

func BenchmarkChanMPSC(b *testing.B) {
	for i := 64; i <= 64; i <<= 1 {
		producers := i
		b.Run(fmt.Sprintf("MPSC-%d", producers), func(b *testing.B) {
			single := b.N / producers+1
			total := single * producers
			q := make(chan int64, 8192)
			var wg sync.WaitGroup
			wg.Add(producers+1)
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

// courtesy or Egon Elbre
func TestXOneringSPMC(t *testing.T) {
	const P = 4
	const N = 100
	var q SPMC
	q.Init(4)

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
				var v int64
				if !q.Get(&v) {
					errs <- fmt.Errorf("failed get")
				}
				//fmt.Println(p, v)
				if v <= lastSeen {
					errs <- fmt.Errorf("got %v last seen %v on producer %v", v, lastSeen, p)
				}
				lastSeen = v
			}
		}(i)
	}

	for err := range errs {
		t.Fatal(err)
	}
}