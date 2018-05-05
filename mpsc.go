package onering

import (
	"sync/atomic"
)

type MPSC struct {
	multi
}

func (r *MPSC) init(size uint32) {
	r.multi.init(size)
	r.rc = 1
}

func (r *MPSC) Get(i interface{}) bool {
	var (
		rc        = r.rc
		data, seq = r.frame(rc)
	)
	if sv := atomic.LoadInt64(seq); rc > sv {
		if rc > r.rp {
			atomic.StoreInt64(&r.rp, rc)
		}
		for ; rc > sv; sv = atomic.LoadInt64(seq) {
			if r.Done() {
				return false
			}
			r.wait()
		}
	}

	inject(i, *data)
	*seq = -rc
	rc++
	r.rc = rc
	if r.rc-r.rp > r.maxbatch {
		atomic.StoreInt64(&r.rp, rc)
	}
	return true
}

func (r *MPSC) Consume(i interface{}) {
	var (
		fn       = extractfn(i)
		maxbatch = int(r.maxbatch)
		it       iter
	)
	for keep := true; keep; {
		var rc, wp = r.rc, atomic.LoadInt64(&r.wp)
		for ; rc >= wp; wp = atomic.LoadInt64(&r.wp) {
			if atomic.LoadInt32(&r.done) > 0 {
				return
			}
			r.wait()
		}

		for i := 0; rc < wp && keep; it.inc() {
			var data, seq = r.frame(rc)
			if i++; atomic.LoadInt64(seq) <= 0 || i&maxbatch == 0 {
				r.rc = rc
				atomic.StoreInt64(&r.rp, rc)
				for atomic.LoadInt64(seq) <= 0 {
					r.wait()
				}
			}
			fn(&it, *data)
			*seq = -rc
			keep = !it.stop
			rc++
		}
		r.rc = rc
		atomic.StoreInt64(&r.rp, r.rc)
	}
}

func (r *MPSC) Put(i interface{}) {
	var wp = r.next(&r.wp)
	for diff := wp - r.mask; diff >= atomic.LoadInt64(&r.rp); {
		r.wait()
	}
	var pos = wp & r.mask
	r.data[pos] = extractptr(i)
	atomic.StoreInt64(&r.seq[pos], wp)
}
