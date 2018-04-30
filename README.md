## One Ring to Queue Them All

Well, no, it's not really just one ring, but a collection of lock free ring buffers for different scenarios, so it's even better!
These queues don't use CAS operations to make them suitable for low latency environments and as a side effect of that,
they observe total order of messages. As a reward for finding flaws/bugs in this, I offer 8 bytes of random numbers for each.

Microbenchmarks are *everything*, the most important thing in the universe.

Rings:

    BenchmarkRingSPSC_Get-8     	300000000	        59.0 ns/op
    BenchmarkRingSPSC_Batch-8   	1000000000	        12.7 ns/op
    BenchmarkRingSPMC-8         	300000000	        42.2 ns/op
    BenchmarkRingMPSC_Get-8     	200000000	        60.5 ns/op
    BenchmarkRingMPSC_Batch-8   	500000000	        32.6 ns/op

Go channels:

    BenchmarkChan/SPSC-8         	300000000	        54.8 ns/op
    BenchmarkChan/SPMC-64-8      	50000000	       327 ns/op
    BenchmarkChan/MPSC-64-8      	100000000	       332 ns/op

Generally a 4-10x increase in performance if you take advantage of batching.
Do note that batching methods in them *do not* increase latency but, in fact, do the opposite.

The API is unstable at the moment, there are no guarantees about anything

### TODO
 * Add producer batching
