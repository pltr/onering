package onering

import (
	"sync/atomic"
)

type SPMC struct {
	multi
}

func (r *SPMC) Get(i *int64) bool {
	var (
		rp        = r.next(&r.rp)
		data, seq = r.frame(rp)
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

func (r *SPMC) Put(i int64) {
	var (
		wp        = r.wp
		data, seq = r.frame(wp)
	)
	for atomic.LoadInt64(seq) > 0 {
		r.wait()
	}
	*data = i
	r.wp++
	atomic.StoreInt64(seq, wp)
}
