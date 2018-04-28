package onering

import (
	"runtime"
	"sync/atomic"
)


type MPSC struct {
	commit
}

func (r *MPSC) Consume(fn func(int64)) {
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
			var pos = p & r.mask
			for atomic.LoadInt32(&r.log[pos]) == 0 {
				runtime.Gosched()
			}
			fn(r.data[pos])
			r.log[pos] = 0
		}
		atomic.StoreInt64(&r.rp, end)
	}
}

func (r *MPSC) Put(i int64) {
	var wp = atomic.AddInt64(&r.wp, 1) - 1
	var pos = wp & r.mask
	for wp-atomic.LoadInt64(&r.rp) >= r.mask {
		runtime.Gosched()
	}
	r.data[pos] = i
	atomic.StoreInt32(&r.log[pos], 1)
}
