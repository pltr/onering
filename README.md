## One ring to queue them all

Well, no. It's not really one ring, but it is a collection of lock free rings for different scenarios

Microbenchmarks are *everything*

    BenchmarkChanSPSC-8     	100000000	        55.0 ns/op	       0 B/op	       0 allocs/op
    BenchmarkChanMPSC/MPSC-64-8         	20000000	       329 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSPSC_Get-8                 	100000000	        70.8 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSPSC_Batch-8               	300000000	        12.9 ns/op	       0 B/op	       0 allocs/op
    BenchmarkSPMC-8                     	100000000	        50.5 ns/op	       0 B/op	       0 allocs/op
    BenchmarkMPSC_Batch-8               	100000000	        32.2 ns/op	       0 B/op	       0 allocs/op