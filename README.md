# One Ring to Queue Them All

Well, no, it's not really just one ring, but a collection of lock-free (starvation/wait free even) ring buffers for different scenarios, so it's even better!
These queues don't use CAS operations to make them suitable for low latency/real-time environments and as a side effect of that,
they preserve total order of messages. As a reward for finding flaws/bugs in this, I offer 64bit of random numbers for each.

A couple of things in it were inspired by the very cool LMAX Disruptor, so thanks @mjpt777!
It's not anywhere near as intrusive and opinionated as the Disruptor though. It's not a framework and its main goal is to be (very) simple.

## Description

The package contains 4 related but different implementations
1. SPSC - Single Producer/Single Consumer - For cases when you just need to send messages from one thread/goroutine to another
2. MPSC - Multi-Producer/Single Consumer - When you need to send messages from many threads/goroutines into a single receiver
3. SPMC - Single Producer/Multi-Consumer - When you need to distribute messages from a single thread to many others
4. MPMC - Multi-Producer/Multi-Consumer - Many-to-Many


At the moment, all queues only support sending pointers (of any type). You can send non pointer types, but it will cause heap allocation. But you *can not* receive anything but pointers, don't even try, it will blow up.

If you build it with `-tags debug`, then all functions will be instrumented to check types at runtime.

There are 2 tests in the package that intentionally use value types to demonstrate it.

## How to use it

### Common interface
    var queeue = onering.New{Size: N}.QueueType()
    queue.Put(*T)
    queue.Get(**T)
    queue.Consume(fn(onering.Iter, *T))
    queue.Close()

### Simplest case
```go
   import "github.com/pltr/onering"
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
```
### Single consumer batching case
Batching consumption is strongly recommended in all single consumer cases, it's expected to have both higher throughput and lower latency

```go
    import "github.com/pltr/onering"
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
```
You can run both examples by `go run cmd/examples.go`


### Warnings
Currently this is highly experimental, so be careful. It also uses some dirty tricks to get around go's typesystem.
If you have a type mismatch between your sender and receiver or try to receive something unexpected, it will likely blow up.

Build it with `-tags debug` to ensure it's not the case.

### FAQ

* **Why four different implementations instead of just one (MPMC)?**
    _There are optimizations to be made in each case. They can have significant effect on performance._

* **Which one should I use?**
    _If you're not sure, MPMC will likely to be the safest choice._

* **I think I found a bug/something doesn't work as expectd**
    _Feel free to open an issue_

* **How fast is it?**
    _I haven't seen any faster, especially when it comes to latency and its distribution (see below)_

* **Did someone actually ask those questions above?**
    _No_

### Some benchmarks

`GOGC=off go test -bench=. -benchtime=3s -run=none`
Macbook pro 2.9 GHz Intel Core i7 (2017)

Rings:

    BenchmarkRingSPSC_GetPinned-8      	300000000	        12.9 ns/op
    BenchmarkRingSPSC_GetNoPin-8       	300000000	        15.8 ns/op
    BenchmarkRingSPSC_Consume-8        	300000000	        12.5 ns/op
    BenchmarkRingMPSC_GetPinned-8      	100000000	        30.2 ns/op
    BenchmarkRingMPSC_GetNoPin1CPU-8   	200000000	        21.1 ns/op
    BenchmarkRingMPSC_Consume-8        	200000000	        28.2 ns/op
    BenchmarkRingSPMC_Pinned-8         	100000000	        43.0 ns/op
    BenchmarkRingSPMC_NoPin1CPU-8      	200000000	        24.9 ns/op
    BenchmarkRingSPMC_Consume-8        	100000000	        44.0 ns/op
    BenchmarkRingMPMC/100P100C-8       	100000000	        46.0 ns/op
    BenchmarkRingMPMC/4P4C_Pinned-8    	100000000	        43.5 ns/op
    BenchmarkRingMPMC/4P4C_1CPU-8      	100000000	        37.9 ns/op


Go channels:

    BenchmarkChanMPMC_100P100C-8       	50000000	        83.9 ns/op
    BenchmarkChan/SPSC_Pinned-8        	100000000	        54.9 ns/op
    BenchmarkChan/SPSC_1CPU-8          	100000000	        45.8 ns/op
    BenchmarkChan/SPMC_Pinned100C-8    	20000000	       352 ns/op
    BenchmarkChan/SPSC_1CPU#01-8       	100000000	        46.9 ns/op
    BenchmarkChan/SPMC_Pinned100C#01-8 	10000000	       452 ns/op
    BenchmarkChan/SPMC_1CPU-8          	100000000	        46.3 ns/op

You can generally expect a 2-10x increase in performance, especially if you use a multicore setup.
Do note that batching methods in them *do not* increase latency but, in fact, do the opposite.

Here's some (however flawed - it's hard to measure it precisely, so had to sample) latency distribution (run with `-tags histogram`):
`GOGC=off go test -tags histogram -bench=. -benchtime=3s -run=none`

    BenchmarkResponseTimesRing-8
    [Sample size: 4096 messages] 50: 25ns	75: 25ns	90: 25ns	99: 25ns	99.9: 25ns	99.99: 25ns	99.999: 25ns	99.9999: 25ns
    [Sample size: 4096 messages] 50: 13ns	75: 13ns	90: 21ns	99: 31ns	99.9: 31ns	99.99: 31ns	99.999: 31ns	99.9999: 31ns
    [Sample size: 4096 messages] 50: 13ns	75: 14ns	90: 14ns	99: 28ns	99.9: 36ns	99.99: 39ns	99.999: 40ns	99.9999: 40ns
    [Sample size: 4096 messages] 50: 13ns	75: 14ns	90: 14ns	99: 28ns	99.9: 37ns	99.99: 43ns	99.999: 50ns	99.9999: 55ns

    BenchmarkResponseTimesChannel-8
    [Sample size: 4096 messages] 50: 86ns	75: 104ns	90: 104ns	99: 104ns	99.9: 104ns	99.99: 104ns	99.999: 104ns	99.9999: 104ns
    [Sample size: 4096 messages] 50: 92ns	75: 119ns	90: 144ns	99: 222ns	99.9: 244ns	99.99: 244ns	99.999: 244ns	99.9999: 244ns
    [Sample size: 4096 messages] 50: 101ns	75: 130ns	90: 154ns	99: 179ns	99.9: 216ns	99.99: 255ns	99.999: 276ns	99.9999: 276ns

This is WIP, so the API is unstable at the moment - there are no guarantees about anything

Also: https://github.com/kellabyte/go-benchmarks/tree/master/queues
Special thanks to @kellabyte and @egonelbre
