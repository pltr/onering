package onering

import (
	"runtime"
	"sync/atomic"
)

type SPMC struct {
	multi
}

func (r *SPMC) Get(i *int64) bool {
	var (
		rp  = atomic.AddInt64(&r.rp, 1) - 1
		pos = rp & r.mask
		seq = &r.seq[pos]
	)
	for next := rp + 1; atomic.LoadInt64(seq) != next; runtime.Gosched() {
		if r.Done() {
			return false
		}
	}
	*i = r.data[pos]
	atomic.StoreInt64(seq, 0)
	return true
}

func (r *SPMC) Put(i int64) {
	var (
		wp  = r.wp
		pos = wp & r.mask
		seq = &r.seq[pos]
	)
	for atomic.LoadInt64(seq) != 0 {
		runtime.Gosched()
	}
	r.data[pos] = i
	r.wp++
	atomic.StoreInt64(seq, r.wp)
}
