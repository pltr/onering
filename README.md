## One ring to queue them all

Well, no. It's not really one ring, but it is a collection of lock free rings for different scenarios

Microbenchmarks are *everything*

    BenchmarkSPSC_Get-8     	100000000	        61.5 ns/op
    BenchmarkSPSC_Batch-8   	300000000	        12.8 ns/op
    BenchmarkSPMC-8         	100000000	        50.1 ns/op
    BenchmarkMPSC_Batch-8   	100000000	        34.5 ns/op
    BenchmarkChanSPSC-8     	100000000	        55.0 ns/op
    BenchmarkChanMPSC/MPSC-64-8         	20000000	       263 ns/op