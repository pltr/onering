package onering

import (
	"runtime"
	"sync/atomic"
)

type MPSC struct {
	multi
}

func (r *MPSC) wait() {
	runtime.Gosched()
}

func (r *MPSC) Get(i *int64) bool {
	var (
		rp  = r.rp
		pos = rp & r.mask
		seq = &r.seq[pos]
	)
	for rp >= atomic.LoadInt64(seq) {
		if r.Done() {
			return false
		}
		r.wait()
	}
	*i = r.data[pos]
	*seq = 0
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
			var (
				pos = p & r.mask
				seq = &r.seq[pos]
			)
			if i++; atomic.LoadInt64(seq) == 0 || i&MaxBatch == 0 {
				atomic.StoreInt64(&r.rp, p)
				for atomic.LoadInt64(seq) == 0 {
					runtime.Gosched()
				}
			}

			fn(r.data[pos])
			*seq = 0
		}
		atomic.StoreInt64(&r.rp, wp)
	}
}

func (r *MPSC) Put(i int64) {
	var wp = atomic.AddInt64(&r.wp, 1) - 1
	for wp-atomic.LoadInt64(&r.rp) >= r.mask {
		runtime.Gosched()
	}
	var pos = wp & r.mask
	r.data[pos] = i
	atomic.StoreInt64(&r.seq[pos], wp+1)
}
