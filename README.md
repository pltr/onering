## One ring to queue them all

Well, no. It's not really one ring, but it is a collection of lock free rings for different scenarios

Microbenchmarks are *everything*

    BenchmarkSPSC_Get-8     	100000000	        59.0 ns/op
    BenchmarkSPSC_Batch-8   	300000000	        12.8 ns/op
    BenchmarkSPMC-8         	100000000	        50.1 ns/op
    BenchmarkMPSC_Batch-8   	100000000	        32.5 ns/op
    BenchmarkChanSPSC-8     	100000000	        55.1 ns/op
    BenchmarkChanMPSC-8        	10000000	        342 ns/op