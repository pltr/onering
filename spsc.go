package onering

import (
	"sync/atomic"
)

type SPSC struct {
	ring
}

func (r *SPSC) Get(i interface{}) bool {
	var rp = r.rc
	for wp := atomic.LoadInt64(&r.wp); rp >= wp; r.wait() {
		if r.rp < rp {
			atomic.StoreInt64(&r.rp, rp)
		} else if atomic.LoadInt32(&r.done) > 0 {
			return false
		}
		wp = atomic.LoadInt64(&r.wc) // in case the writer is idle, start reading its cache
	}
	inject(i, r.data[rp&r.mask])
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
		var rp, wp = r.rc, atomic.LoadInt64(&r.wp)
		for ; rp >= wp; r.wait() {
			if rp > r.rp {
				r.rc = rp
				atomic.StoreInt64(&r.rp, r.rc)
			} else if r.Done() {
				return
			}
			wp = atomic.LoadInt64(&r.wc)
		}
		for i := 0; rp < wp && keep; it.inc() {
			if i++; i&maxbatch == 0 {
				r.rc = rp
				atomic.StoreInt64(&r.rp, r.rc)
			}
			fn(&it, r.data[rp&r.mask])
			rp++
			keep = !it.stop
		}
		r.rc = rp
		atomic.StoreInt64(&r.rp, r.rc)
	}
}

func (r *SPSC) Put(i interface{}) {
	var (
		wc = atomic.LoadInt64(&r.wc)
		wp = atomic.LoadInt64(&r.wp)
	)
	for diff := wc - r.mask; diff >= atomic.LoadInt64(&r.rp); {
		if wc > wp {
			atomic.StoreInt64(&r.wp, wc)
		}
		r.wait()
	}
	r.data[wc&r.mask] = extractptr(i)
	if wc-wp > r.maxbatch {
		r.wc++
		atomic.StoreInt64(&r.wp, r.wc)
	} else {
		atomic.StoreInt64(&r.wc, wp+1)
	}
}
