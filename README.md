## One ring to queue them all

Well, no, it's not really just one ring, but a collection of lock free ring buffers for different scenarios, so it's even better!
These queues don't use CAS operations to make them suitable for low latency environments and as a side effect of that,
they observe total order of messages. As a reword for finding flaws/bugs in this, I offer 8 bytes of random numbers for each.

Microbenchmarks are *everything*, the most important thing in the universe.

    BenchmarkSPSC_Get-8     	20000000	        59.9 ns/op
    BenchmarkSPSC_Batch-8   	100000000	        12.7 ns/op
    BenchmarkSPMC-8         	30000000	        41.4 ns/op
    BenchmarkMPSC_Get-8     	30000000	        60.2 ns/op
    BenchmarkMPSC_Batch-8   	50000000	        35.4 ns/op
    BenchmarkChanSPSC-8     	30000000	        54.8 ns/op
    BenchmarkChanMPSC/MPSC-64-8         	 3000000	       436 ns/op