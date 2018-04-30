package onering

import (
	"runtime"
	"sync/atomic"
)

// WARNING: this will ONLY work in SPSC situations

type SPSC struct {
	ring
	_ [5]int64
}

func (r *SPSC) Get(i *int64) bool {
	var rp = r.rp
	for rp >= atomic.LoadInt64(&r.wp) {
		if r.Done() {
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
		var rp, wp = r.rp, atomic.LoadInt64(&r.wp)
		for ; rp >= wp; runtime.Gosched() {
			if r.Done() {
				return
			}
			wp = atomic.LoadInt64(&r.wp)
		}
		var i = 0
		for p := rp; p < wp; p++ {
			fn(r.data[p&r.mask])
			if i++; i&MaxBatch == 0 {
				atomic.StoreInt64(&r.rp, p)
			}
		}
		atomic.StoreInt64(&r.rp, wp)
	}
}

func (r *SPSC) Write(i interface{}) {
	var wp = r.wp
	for wp-atomic.LoadInt64(&r.rp) >= r.mask {
		runtime.Gosched()
	}
	r.data[wp&r.mask] = i
	atomic.AddInt64(&r.wp, 1)
}
