package onering

import (
	"sync/atomic"
)

type SPSC struct {
	ring
	_ [4]byte
}

func (r *SPSC) Get(i interface{}) bool {
	var rc = r.rc
	for wp := atomic.LoadInt64(&r.wp); rc >= wp; r.wait() {
		wp = atomic.LoadInt64(&r.wc)
		if r.rp < rc {
			atomic.StoreInt64(&r.rp, rc)
		} else if atomic.LoadInt32(&r.done) > 0 && rc >= wp {
			return false
		}
	}
	inject(i, r.data[rc&r.mask])
	r.rc++
	if r.rc-r.rp > r.maxbatch {
		atomic.StoreInt64(&r.rp, r.rc)
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
		for ; rc >= wp; r.wait() {
			wp = atomic.LoadInt64(&r.wc)
			if rc > r.rp {
				r.rc = rc
				atomic.StoreInt64(&r.rp, rc)
			} else if atomic.LoadInt32(&r.done) > 0 && rc >= wp {
				return
			}
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
	for diff := wc - r.mask; diff >= atomic.LoadInt64(&r.rp); {
		if wc > r.wp {
			atomic.StoreInt64(&r.wp, wc)
		}
		r.wait()
	}
	r.data[wc&r.mask] = extractptr(i)
	wc = atomic.AddInt64(&r.wc, 1)
	if wc-r.wp > r.maxbatch {
		atomic.StoreInt64(&r.wp, wc)
	}
}
