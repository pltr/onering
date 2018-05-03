package onering

import (
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
	for ; atomic.LoadInt64(seq) != rp; r.wait() {
		if r.Done() {
			return false
		}
	}
	inject(i, *data)
	atomic.StoreInt64(seq, -rp)
	return true
}

func (r *SPMC) Consume(i interface{}) {
	var (
		fn  = extractfn(i)
		ptr unsafe.Pointer
		it  iter
	)
	for ; !it.stop && r.Get(&ptr); it.inc() {
		fn(&it, ptr)
	}
}

func (r *SPMC) Put(i interface{}) {
	var (
		wp        = r.wp
		data, seq = r.frame(wp)
	)
	for atomic.LoadInt64(seq) > 0 {
		r.wait()
	}
	*data = extractptr(i)
	r.wp++
	atomic.StoreInt64(seq, wp)
}
