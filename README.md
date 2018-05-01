## One Ring to Queue Them All

Well, no, it's not really just one ring, but a collection of lock-free (WFPO, even) ring buffers for different scenarios, so it's even better!
These queues don't use CAS operations to make them suitable for low latency/real-time environments and as a side effect of that,
they preserve total order of messages. As a reward for finding flaws/bugs in this, I offer 64bit of random numbers for each.

Microbenchmarks are *everything*, the most important thing in the universe.

Rings:

    BenchmarkRingSPSC_Get-8     	300000000	        59.0 ns/op
    BenchmarkRingSPSC_Batch-8   	1000000000	        12.7 ns/op
    BenchmarkRingSPMC-8         	300000000	        42.2 ns/op
    BenchmarkRingMPSC_Get-8     	200000000	        60.5 ns/op
    BenchmarkRingMPSC_Batch-8   	500000000	        32.6 ns/op
    BenchmarkRingMPMC_Get-8         100000000	        52.6 ns/op

Go channels:

    BenchmarkChan/SPSC-8         	300000000	        54.8 ns/op
    BenchmarkChan/SPMC-64-8      	50000000	       327 ns/op
    BenchmarkChan/MPSC-64-8      	100000000	       332 ns/op

Generally a 4-10x increase in performance if you take advantage of batching.
Do note that batching methods in them *do not* increase latency but, in fact, do the opposite.

    BenchmarkResponseTimesRing-8

    [Sample size: 2048 messages] 50: 18ns	75: 19ns	90: 22ns	99: 22ns	99.9: 22ns	99.99: 22ns	99.999: 22ns	99.9999: 22ns
    [Sample size: 2048 messages] 50: 17ns	75: 19ns	90: 21ns	99: 25ns	99.9: 33ns	99.99: 33ns	99.999: 33ns	99.9999: 33ns
    [Sample size: 2048 messages] 50: 15ns	75: 17ns	90: 18ns	99: 34ns	99.9: 46ns	99.99: 54ns	99.999: 77ns	99.9999: 77ns

    BenchmarkResponseTimesChannel-8
    [Sample size: 2048 messages] 50: 169ns	75: 170ns	90: 170ns	99: 170ns	99.9: 170ns	99.99: 170ns	99.999: 170ns	99.9999: 170ns
    [Sample size: 2048 messages] 50: 157ns	75: 205ns	90: 251ns	99: 352ns	99.9: 421ns	99.99: 421ns	99.999: 421ns	99.9999: 421ns
    [Sample size: 2048 messages] 50: 163ns	75: 222ns	90: 266ns	99: 317ns	99.9: 393ns	99.99: 448ns	99.999: 459ns	99.9999: 459ns

The API is unstable at the moment, there are no guarantees about anything

### TODO
 * Add producer batching

Also: https://github.com/kellabyte/go-benchmarks/tree/master/queues
SPSC Get (bounded by time.Now() call)
![chart](https://camo.githubusercontent.com/553d9f8936ed5f298e1b3c0de1724d71b5c57cea/68747470733a2f2f692e696d6775722e636f6d2f78547a397645432e706e67
 "Queue Benchmark")

Special thanks to @kellabyte and @egonelbre
