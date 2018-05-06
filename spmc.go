package onering

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type SPMC struct {
	multi
}

func (r *SPMC) Get(i interface{}) bool {
	var (
		rp        = r.next(&r.rp)
		data, seq = r.frame(rp)
	)

	for pread := -rp; atomic.LoadInt64(seq) != pread; runtime.Gosched() {
		if atomic.LoadInt32(&r.done) > 0 && atomic.LoadInt64(&r.wp) <= rp {
			return false
		}
	}
	inject(i, *data)
	atomic.StoreInt64(seq, rp+r.size)
	return true
}

func (r *SPMC) Consume(i interface{}) {
	var (
		fn  = extractfn(i)
		it  iter
		ptr unsafe.Pointer
	)
	for !it.stop && r.Get(&ptr) {
		fn(&it, ptr)
	}
}

func (r *SPMC) Put(i interface{}) {
	var (
		wp        = r.wp
		data, seq = r.frame(wp)
	)
	for atomic.LoadInt64(seq) < 0 {
		runtime.Gosched()
	}
	*data = extractptr(i)
	r.wp++
	atomic.StoreInt64(seq, -wp)
}
