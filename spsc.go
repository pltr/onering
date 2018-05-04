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
	atomic.StoreInt64(&r.wc, wc+1)
	if wc-wp > r.maxbatch {
		atomic.StoreInt64(&r.wp, wc+1)
	}
}
