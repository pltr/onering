package onering

import (
	"sync/atomic"
)

type MPSC struct {
	multi
}

func (r *MPSC) Get(i interface{}) bool {
	var (
		rp        = r.rp
		data, seq = r.frame(rp)
	)
	for rp > atomic.LoadInt64(seq) {
		if r.Done() {
			return false
		}
		r.wait()
	}
	inject(i, *data)
	*seq = -rp
	atomic.StoreInt64(&r.rp, rp+1)
	return true
}

func (r *MPSC) Consume(i interface{}) {
	var fn = extractfn(i)
	for {
		var rp, wp = r.rp, atomic.LoadInt64(&r.wp)
		for ; rp >= wp; r.wait() {
			if r.Done() {
				return
			}
			wp = atomic.LoadInt64(&r.wp)
		}

		for i := 0; rp < wp; rp++ {
			var data, seq = r.frame(rp)
			if i++; atomic.LoadInt64(seq) <= 0 || i&MaxBatch == 0 {
				atomic.StoreInt64(&r.rp, rp)
				for atomic.LoadInt64(seq) <= 0 {
					r.wait()
				}
			}
			fn(*data)
			*seq = -rp
		}
		atomic.StoreInt64(&r.rp, wp)
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
