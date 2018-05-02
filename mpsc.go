package onering

import (
	"sync/atomic"
)

type MPSC struct {
	multi
}

func (r *MPSC) Get(i *int64) bool {
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
	*i = *data
	*seq = -rp
	atomic.StoreInt64(&r.rp, rp+1)
	return true
}

func (r *MPSC) Consume(fn func(int64)) {
	for {
		var rp, wp = r.rp, atomic.LoadInt64(&r.wp)
		for ; rp >= wp; r.wait() {
			if r.Done() {
				return
			}
			wp = atomic.LoadInt64(&r.wp)
		}

		for p, i := rp, 0; p < wp; p++ {
			var data, seq = r.frame(p)
			if i++; atomic.LoadInt64(seq) <= 0 || i&MaxBatch == 0 {
				atomic.StoreInt64(&r.rp, p)
				for atomic.LoadInt64(seq) == 0 {
					r.wait()
				}
			}

			fn(*data)
			*seq = -p
		}
		atomic.StoreInt64(&r.rp, wp)
	}
}

func (r *MPSC) Put(i int64) {
	var wp = r.next(&r.wp)
	for diff := wp - r.mask; diff >= atomic.LoadInt64(&r.rp); {
		r.wait()
	}
	var pos = wp & r.mask
	r.data[pos] = i
	atomic.StoreInt64(&r.seq[pos], wp)
}
