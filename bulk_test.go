package onering
//
//import (
//	"fmt"
//	"github.com/codahale/hdrhistogram"
//	"runtime"
//	"sync"
//	"testing"
//	"time"
//)
//
//// courtesy of Kelly Sommers aka @kellabyte
//
//const (
//	sampleE     = 11
//	sampleTimes = 1 << sampleE
//	sampleMask  = sampleTimes - 1
//)
//
//func BenchmarkResponseTimesRing(b *testing.B) {
//	var ring SPSC
//	ring.Init(8192)
//	var wg sync.WaitGroup
//	wg.Add(2)
//	var diffs = make([]int64, (b.N/sampleTimes)+1)
//	b.ResetTimer()
//	go func(n int) {
//		runtime.LockOSThread()
//		var t1 = time.Now().UnixNano()
//		for i := 1; i < n; i++ {
//			var v int64 = 0
//			if i&sampleMask == 0 {
//				v = t1
//				t1 = time.Now().UnixNano()
//			}
//			ring.Put(v)
//		}
//		wg.Done()
//	}(b.N + 1)
//	go func(n int) {
//		runtime.LockOSThread()
//		var i int = 0
//		ring.Consume(func(v int64) {
//			if v != 0 {
//				diffs[i] = (time.Now().UnixNano() - v) / sampleTimes
//				i++
//			}
//			n--
//			if n <= 0 {
//				ring.Close()
//			}
//		})
//		wg.Done()
//	}(b.N)
//
//	wg.Wait()
//	recordLatencyDistribution("BenchmarkResponseTimesRing", diffs)
//}
//
//func BenchmarkResponseTimesChannel(b *testing.B) {
//	var ch = make(chan int64, 8192)
//	var wg sync.WaitGroup
//	wg.Add(2)
//	var diffs = make([]int64, (b.N/sampleTimes)+1)
//	b.ResetTimer()
//	go func(n int) {
//		runtime.LockOSThread()
//		var t1 = time.Now().UnixNano()
//		for i := 1; i < n; i++ {
//			var v int64 = 0
//			if i&sampleMask == 0 {
//				v = t1
//				t1 = time.Now().UnixNano()
//			}
//			ch <- v
//		}
//		close(ch)
//		wg.Done()
//	}(b.N + 1)
//	go func(n int) {
//		runtime.LockOSThread()
//		var i = 0
//		for v := range ch {
//			if v != 0 {
//				diffs[i] = (time.Now().UnixNano() - v) / sampleTimes
//				i++
//			}
//		}
//		wg.Done()
//	}(b.N)
//
//	wg.Wait()
//	recordLatencyDistribution("BenchmarkResponseTimesChannel", diffs)
//}
//
//func recordLatencyDistribution(name string, diffs []int64) {
//	fmt.Printf("[Sample size: %v messages] ", sampleTimes)
//	histogram := hdrhistogram.New(1, 1000000, 5)
//	for _, d := range diffs {
//		if d != 0 {
//			histogram.RecordValue(d)
//		}
//	}
//
//	fmt.Printf("50: %dns\t75: %dns\t90: %dns\t99: %dns\t99.9: %dns\t99.99: %dns\t99.999: %dns\t99.9999: %dns\n",
//		histogram.ValueAtQuantile(50),
//		histogram.ValueAtQuantile(75),
//		histogram.ValueAtQuantile(90),
//		histogram.ValueAtQuantile(99),
//		histogram.ValueAtQuantile(99.9),
//		histogram.ValueAtQuantile(99.99),
//		histogram.ValueAtQuantile(99.999),
//		histogram.ValueAtQuantile(99.9999))
//
//	//histwriter.WriteDistributionFile(histogram, histwriter.Percentiles{50, 75, 90, 99, 99.9, 99.99, 99.999, 99.9999}, 1.0, name+".histogram")
//}
