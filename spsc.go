package onering

import (
	"sync/atomic"
	"unsafe"
)

// WARNING: this will ONLY work in SPSC situations

type SPSC struct {
	ring
	_ [24]byte
}

func (r *SPSC) Get(i interface{}) bool {
	var rp = r.rp
	for rp >= atomic.LoadInt64(&r.wp) {
		if r.Done() {
			return false
		}
		r.wait()
	}
	r.inject(extractptr(i), r.data[rp&r.mask])
	atomic.StoreInt64(&r.rp, rp+1)
	return true
}

func (r *SPSC) Consume(fn func(pointer unsafe.Pointer)) {
	for {
		var rp, wp = r.rp, atomic.LoadInt64(&r.wp)
		for ; rp >= wp; r.wait() {
			if r.Done() {
				return
			}
			wp = atomic.LoadInt64(&r.wp)
		}
		for i := 0; rp < wp; rp++ {
			if i++; i&MaxBatch == 0 {
				atomic.StoreInt64(&r.rp, rp)
			}
			fn(r.data[rp&r.mask])
		}
		atomic.StoreInt64(&r.rp, wp)
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


type spscbatch struct {
	ring *SPSC
	offset int
	rp int64
	wp int64
}


func splitbatches(rp, wp int64) (int64, int64) {
	var total = wp - rp
	var batches = total >> BatchExp
	var rest = total & MaxBatch
	return batches, rest
}