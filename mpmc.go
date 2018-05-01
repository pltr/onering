package onering

import (
	"sync/atomic"
)

// my precious

type MPMC struct {
	multi
}

func (r *MPMC) Init(size uint32) {
	r.multi.Init(size)
	for i := range r.seq {
		r.seq[i] = -int64(i)
	}
	r.wp = r.size
	r.rp = r.size
}

func (r *MPMC) Get(i *int64) bool {
	var (
		rp        = r.next(&r.rp)
		data, seq = r.contents(rp)
	)
	for ; atomic.LoadInt64(seq) != rp; r.wait() {
		if r.Done() {
			return false
		}
	}
	*i = *data
	atomic.StoreInt64(seq, -rp)
	return true
}

func (r *MPMC) Put(i int64) {
	var (
		wp        = r.next(&r.wp)
		data, seq = r.contents(wp)
	)
	for pread := r.size - wp; atomic.LoadInt64(seq) != pread; {
		r.wait()
	}

	*data = i
	atomic.StoreInt64(seq, wp)
}
