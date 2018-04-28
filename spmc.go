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
		rp = atomic.AddInt64(&r.rp, 1) - 1
		pos = rp & r.mask
	)
	for atomic.LoadInt32(&r.log[pos]) != int32(rp+1) {
		if atomic.LoadInt32(&r.done) > 0 {
			return false
		}
		runtime.Gosched()
	}
	*i = r.data[pos]
	atomic.StoreInt32(&r.log[pos], 0)
	return true
}


func (r *SPMC) Put(i int64) {
	var wp = r.wp
	var pos = wp & r.mask
	for atomic.LoadInt32(&r.log[pos]) != 0 {
		runtime.Gosched()
	}
	r.data[pos] = i
	r.wp++
	atomic.StoreInt32(&r.log[pos], int32(r.wp))
}

