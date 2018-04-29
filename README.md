## One ring to queue them all

Well, no. It's not really one ring, but it is a collection of lock free rings for different scenarios

Microbenchmarks are *everything*

    BenchmarkSPSC_Get-8     	20000000	        59.9 ns/op
    BenchmarkSPSC_Batch-8   	100000000	        12.7 ns/op
    BenchmarkSPMC-8         	30000000	        41.4 ns/op
    BenchmarkMPSC_Get-8     	30000000	        60.2 ns/op
    BenchmarkMPSC_Batch-8   	50000000	        35.4 ns/op
    BenchmarkChanSPSC-8     	30000000	        54.8 ns/op
    BenchmarkChanMPSC/MPSC-64-8         	 3000000	       436 ns/op