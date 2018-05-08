package onering

import (
	"sync/atomic"
)

type SPSC struct {
	_ [8]int64
	wc int64
	ring
	_ [4]byte
}

func (r *SPSC) Get(i interface{}) bool {
	var rc = r.rc
	if wp := atomic.LoadInt64(&r.wp); rc >= wp {
		if rc > r.rp {
			atomic.StoreInt64(&r.rp, rc)
		}
		for ; rc >= wp; wp = atomic.LoadInt64(&r.wc) {
			if atomic.LoadInt32(&r.done) > 0 {
				return false
			}
			r.wait()
		}
	}

	inject(i, r.data[rc&r.mask])
	rc++
	r.rc = rc
	if r.rc-r.rp > r.maxbatch {
		atomic.StoreInt64(&r.rp, rc)
	}
	return true
}

func (r *SPSC) Consume(i interface{}) {
	var (
		fn       = extractfn(i)
		maxbatch = int(r.maxbatch)
		it       iter
	)
	for keep := true; keep; {
		var rc, wp = r.rc, atomic.LoadInt64(&r.wp)
		for rc >= wp {
			if atomic.LoadInt32(&r.done) > 0 {
				return
			}
			r.wait()
			wp = atomic.LoadInt64(&r.wc)
		}

		for i := 0; rc < wp && keep; it.inc() {
			if i++; i&maxbatch == 0 {
				r.rc = rc
				atomic.StoreInt64(&r.rp, rc)
			}
			fn(&it, r.data[rc&r.mask])
			rc++
			keep = !it.stop
		}
		r.rc = rc
		atomic.StoreInt64(&r.rp, rc)
	}
}

func (r *SPSC) Put(i interface{}) {
	var wc = r.wc
	if diff, rp := wc-r.mask, atomic.LoadInt64(&r.rp); diff >= rp {
		if wc > r.wp {
			atomic.StoreInt64(&r.wp, wc)
		}
		for ; diff >= rp; rp = atomic.LoadInt64(&r.rp) {
			r.wait()
		}
	}
	r.data[wc&r.mask] = extractptr(i)
	wc = atomic.AddInt64(&r.wc, 1)
	if wc-r.wp > r.maxbatch {
		atomic.StoreInt64(&r.wp, wc)
	}
}
