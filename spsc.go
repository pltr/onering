package onering


import (
	"sync/atomic"
	"runtime"
)

// WARNING: this will ONLY work in SPSC situations

type SPSC struct {
	ring
	_ [5]int64
}

func (r *SPSC) Get() (i int64) {
	var rp = r.rp
	for rp >= atomic.LoadInt64(&r.wp) {
		runtime.Gosched()
	}
	i = r.data[rp&r.mask]
	atomic.AddInt64(&r.rp, 1)
	return
}

func (r *SPSC) Consume(fn func(int64)) {
	for {
		var (
			rp  = r.rp
			end int64
		)
		for {
			if end = atomic.LoadInt64(&r.wp); end > rp {
				break
			} else if atomic.LoadInt32(&r.done) > 0 {
				return
			}
			runtime.Gosched()
		}
		if end-rp > MaxBatch {
			end = rp + MaxBatch
		}

		for p := rp; p < end; p++ {
			fn(r.data[p & r.mask])
		}
		atomic.StoreInt64(&r.rp, end)
	}
}

func (r *SPSC) Put(i int64) {
	var wp = r.wp
	for wp - atomic.LoadInt64(&r.rp) >= r.mask {
		 runtime.Gosched()
	}
	r.data[wp&r.mask] = i
	atomic.AddInt64(&r.wp, 1)
}
