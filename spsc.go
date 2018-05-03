package onering

import (
	"sync/atomic"
)

type SPSC struct {
	ring
	_ [8]byte
}

func (r *SPSC) Get(i interface{}) bool {
	var rp = r.rp
	for ; rp >= atomic.LoadInt64(&r.wp); r.wait() {
		if r.Done() {
			return false
		}
	}
	inject(i, r.data[rp&r.mask])
	atomic.StoreInt64(&r.rp, rp+1)
	return true
}

func (r *SPSC) Consume(i interface{}) {
	var (
		fn       = extractfn(i)
		maxbatch = int(r.maxbatch)
		it       iter
	)
	for keep := true; keep; {
		var rp, wp = r.rp, atomic.LoadInt64(&r.wp)
		for ; rp >= wp; r.wait() {
			if r.Done() {
				return
			}
			wp = atomic.LoadInt64(&r.wp)
		}
		for i := 0; rp < wp && keep; it.inc() {
			if i++; i&maxbatch == 0 {
				atomic.StoreInt64(&r.rp, rp)
			}
			fn(&it, r.data[rp&r.mask])
			rp++
			keep = !it.stop
		}
		atomic.StoreInt64(&r.rp, rp)
	}
}

func (r *SPSC) Put(i interface{}) {
	var wp = r.wp
	for diff := wp - r.mask; diff >= atomic.LoadInt64(&r.rp); {
		r.wait()
	}
	r.data[wp&r.mask] = extractptr(i)
	atomic.StoreInt64(&r.wp, wp+1)
}
