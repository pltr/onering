package onering

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

// As it turns out, it's a Chris' Tomasson's modification of Vyukov's queue
// https://groups.google.com/forum/#!searchin/lock-free/thomasson/lock-free/acjQ3-89abE/a6-Di0GZsyEJ
// http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue

type MPMC struct {
	multi
}

func (r *MPMC) init(n *New) {
	r.multi.init(n)
	for i := range r.seq {
		r.seq[i] = int64(i)
	}
	r.seq[0] = r.size
}

func (r *MPMC) Get(i interface{}) bool {
	var (
		rp        = r.next(&r.rp)
		data, seq = r.frame(rp)
	)

	for pread := -rp; atomic.LoadInt64(seq) != pread; {
		if atomic.LoadInt32(&r.done) > 0 && atomic.LoadInt64(&r.wp) <= rp {
			return false
		}
		runtime.Gosched()
	}

	inject(i, *data)
	atomic.StoreInt64(seq, rp+r.size)
	return true
}

func (r *MPMC) Consume(i interface{}) {
	var (
		fn  = extractfn(i)
		it  iter
		ptr unsafe.Pointer
	)
	for !it.stop && r.Get(&ptr) {
		fn(&it, ptr)
	}
}

func (r *MPMC) Put(i interface{}) {
	var (
		wp        = r.next(&r.wp)
		data, seq = r.frame(wp)
	)

	for atomic.LoadInt64(seq) != wp {
		runtime.Gosched()
	}

	*data = extractptr(i)
	atomic.StoreInt64(seq, -wp)
}
