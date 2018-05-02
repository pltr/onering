package onering

import (
	"sync/atomic"
)

// my precious

type MPMC struct {
	multi
}

func (r *MPMC) init(inj Injector, size uint32) {
	r.multi.init(inj, size)
	for i := range r.seq {
		r.seq[i] = -int64(i)
	}
	r.wp = r.size
	r.rp = r.size
}

func (r *MPMC) Get(i interface{}) bool {
	var (
		rp        = r.next(&r.rp)
		data, seq = r.frame(rp)
	)
	for ; atomic.LoadInt64(seq) != rp; r.wait() {
		if r.Done() {
			return false
		}
	}
	r.inject(extractptr(i), data)
	atomic.StoreInt64(seq, -rp)
	return true
}

func (r *MPMC) Put(i interface{}) {
	var (
		wp        = r.next(&r.wp)
		data, seq = r.frame(wp)
	)
	for pread := r.size - wp; atomic.LoadInt64(seq) != pread; {
		r.wait()
	}

	*data = extractptr(i)
	atomic.StoreInt64(seq, wp)
}
