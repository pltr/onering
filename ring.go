package onering

import (
	"math/bits"
	"sync/atomic"
	"unsafe"
)

const MaxBatch = 255
const spin = 512 - 1 // not used at the moment

type ring struct {
	wp   int64
	_    [7]int64
	rp   int64
	_    [7]int64
	data []unsafe.Pointer
	typ unsafe.Pointer
	mask int64
	done int32

}

func (r *ring) Init(typ interface{}, size uint32) {
	r.data = make([]unsafe.Pointer, 1<<uint(32-bits.LeadingZeros32(size-1)))
	r.mask = int64(len(r.data) - 1)
	r.done = 0
}

func (r *ring) Close() {
	atomic.StoreInt32(&r.done, 1)
}

func (r *ring) Done() bool {
	return atomic.LoadInt32(&r.done) > 0 && atomic.LoadInt64(&r.wp) <= atomic.LoadInt64(&r.rp)
}

type multi struct {
	ring
	seq []int64
	_   [5]int64
}

func (c *multi) Init(size uint32) {
	c.ring.Init(nil, size)
	c.seq = make([]int64, len(c.data))
}

// empty sync.Locker for conditionals
type NoLock struct{}

func (NoLock) Lock()   {}
func (NoLock) Unlock() {}
