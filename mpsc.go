package onering

import (
	"runtime"
	"sync/atomic"
)

type MPSC struct {
	multi
}

func (r *MPSC) Get(i *int64) bool {
	var (
		rp  = r.rp
		pos = rp & r.mask
		seq = &r.seq[pos]
	)
	for rp >= atomic.LoadInt64(seq) {
		if !r.Opened() {
			return false
		}
		runtime.Gosched()
	}
	*i = r.data[pos]
	*seq = 0
	atomic.AddInt64(&r.rp, 1)
	return true
}

func (r *MPSC) Consume(fn func(int64)) {
	for {
		var (
			rp = r.rp
			wp int64
		)
		for {
			if wp = atomic.LoadInt64(&r.wp); wp > rp {
				break
			} else if atomic.LoadInt32(&r.done) > 0 {
				return
			}
			runtime.Gosched()
		}
		var i = 0
		for p := rp; p < wp; p++ {
			var (
				pos = p & r.mask
				seq = &r.seq[pos]
			)
			for atomic.LoadInt64(seq) == 0 {
				runtime.Gosched()
			}
			fn(r.data[pos])
			*seq = 0
			if i++; i&MaxBatch == 0 {
				atomic.StoreInt64(&r.rp, p)
			}
		}
		atomic.StoreInt64(&r.rp, wp)
	}
}

func (r *MPSC) Put(i int64) {
	var (
		next = atomic.AddInt64(&r.wp, 1)
		wp   = next - 1
		pos  = wp & r.mask
	)
	for wp-atomic.LoadInt64(&r.rp) >= r.mask {
		runtime.Gosched()
	}
	r.data[pos] = i
	atomic.StoreInt64(&r.seq[pos], next)
}
