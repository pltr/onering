package onering

import (
	"sync/atomic"
	"runtime"
	"unsafe"
)

type SPMC struct {
	commit
}


func (r *SPMC) Get(i *int64) bool {
	var (
		rp = atomic.AddInt64(&r.rp, 1) - 1
		pos = rp & r.mask
	)
	for atomic.LoadInt32(&r.log[pos]) == 0 {
		if atomic.LoadInt32(&r.done) > 0 {
			return false
		}
		runtime.Gosched()
	}
	*i = r.data[pos]
	atomic.StoreInt32(&r.log[pos], 0)
	return true
}

func (r *SPMC) Close() {
	atomic.StoreInt32(&r.done, 1)
}

func (r *SPMC) Open() bool {
	return atomic.LoadInt32(&r.done) == 0 || atomic.LoadInt64(&r.wp) - atomic.LoadInt64(&r.rp) > 0
}

func (r *SPMC) Put(i int64) {
	var pos = r.wp & r.mask
	for atomic.LoadInt32(&r.log[pos]) != 0 {
		runtime.Gosched()
	}
	r.data[pos] = i
	r.wp++
	atomic.StoreInt32(&r.log[pos], 1)
}

//
//type NoLock struct{}
//
//func(*NoLock) Lock() {}
//func(*NoLock) Unlock() {}
//

type IHeader struct {
	T, D unsafe.Pointer
}