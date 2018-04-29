## One ring to queue them all

Well, no, it's not really just one ring, but a collection of lock free ring buffers for different scenarios, so it's even better!
These queues don't use CAS operations to make them suitable for low latency environments and as a side effect of that,
they observe total order of messages. As a reward for finding flaws/bugs in this, I offer 8 bytes of random numbers for each.

Microbenchmarks are *everything*, the most important thing in the universe.

Rings:
    BenchmarkRingSPSC_Get-8     	100000000	        59.2 ns/op
    BenchmarkRingSPSC_Batch-8   	300000000	        12.7 ns/op
    BenchmarkRingSPMC-8         	100000000	        41.0 ns/op
    BenchmarkRingMPSC_Get-8     	100000000	        59.7 ns/op
    BenchmarkRingMPSC_Batch-8   	100000000	        33.0 ns/op

Go channels
    BenchmarkChanSPSC-8         	100000000	        54.8 ns/op
    BenchmarkChanMPSC/MPSC-64-8 	20000000	       333 ns/op

Generally a 4-10x increase in performance if you take advantage of batching.
Do note that batching methods in them *do not* increase latency but, in fact, do the opposite.

TODO: Add producer batching