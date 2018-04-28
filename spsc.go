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

func (r *SPSC) Get(i *int64) bool {
	var rp = r.rp
	for rp >= atomic.LoadInt64(&r.wp) {
		if atomic.LoadInt32(&r.done) > 0 {
			return false
		}
		runtime.Gosched()
	}
	*i = r.data[rp&r.mask]
	atomic.AddInt64(&r.rp, 1)
	return true
}

func (r *SPSC) Consume(fn func(int64)) {
	for {
		var (
			rp  = r.rp
			wp int64
		)
		for {
			if wp = atomic.LoadInt64(&r.wp); wp > rp {
				break
			} else if atomic.LoadInt32(&r.done) > 0 {
				return
			}
			runtime.Gosched()
		}
		var i = 0
		for p := rp; p < wp; p++ {
			fn(r.data[p & r.mask])
			if i++; i & MaxBatch == 0 {
				atomic.StoreInt64(&r.rp, p)
			}
		}
		atomic.StoreInt64(&r.rp, wp)
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
