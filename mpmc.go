package onering

import (
	"runtime"
	"sync/atomic"
)

// my precious

// I haven't actually seen this implementation of WFPO MPMC before, would appreciate if someone knows something about it

type MPMC struct {
	multi
	size int64
}

func (r *MPMC) Init(size uint32) {
	r.multi.Init(size)
	r.size = int64(len(r.data))
	for i := range r.seq {
		r.seq[i] = -int64(i)
	}
	r.wp = r.size
	r.rp = r.size
}

func (r *MPMC) Get(i *int64) bool {
	var (
		rp  = atomic.AddInt64(&r.rp, 1) - 1
		pos = rp & r.mask
		seq = &r.seq[pos]
	)
	for ; atomic.LoadInt64(seq) != rp; runtime.Gosched() {
		if r.Done() {
			return false
		}
	}
	*i = r.data[pos]
	atomic.StoreInt64(seq, -rp)
	return true
}

func (r *MPMC) Put(i int64) {
	var (
		wp = atomic.AddInt64(&r.wp, 1) - 1
		pos = wp & r.mask
		seq = &r.seq[pos]
	)

	for pread := r.size - wp; atomic.LoadInt64(seq) != pread; {
		runtime.Gosched()
	}

	r.data[pos] = i
	atomic.StoreInt64(seq, wp)
}

func (r *MPMC) wait() {
	runtime.Gosched()
}