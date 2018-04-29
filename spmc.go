package onering

import (
	"sync/atomic"
	"runtime"
)

type SPMC struct {
	commit
}

func (r *SPMC) Get(i *int64) bool {
	var (
		next = atomic.AddInt64(&r.rp, 1)
		rp =  next - 1
		pos = rp & r.mask
		seq = &r.log[pos]
	)
	for atomic.LoadInt64(seq) != next {
		if !r.Opened() {
			return false
		}
		runtime.Gosched()
	}
	*i = r.data[pos]
	atomic.StoreInt64(seq, 0)
	return true
}


func (r *SPMC) Put(i int64) {
	var (
		wp = r.wp
		pos = wp & r.mask
		seq = &r.log[pos]
	)
	for atomic.LoadInt64(seq) != 0 {
		runtime.Gosched()
	}
	r.data[pos] = i
	r.wp++
	atomic.StoreInt64(seq, r.wp)
}

