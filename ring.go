package onering

import (
	"math/bits"
	"sync/atomic"
)

const MaxBatch = 255
const spin = 512 - 1 // not used at the moment

type ring struct {
	wp   int64
	_    [7]int64
	rp   int64
	_    [7]int64
	data []int64
	mask int64
	done int32
}

func (r *ring) Init(size uint) {
	r.data = make([]int64, 1<<uint(64-bits.LeadingZeros(size-1)))
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

func (c *multi) Init(size uint) {
	c.ring.Init(size)
	c.seq = make([]int64, len(c.data))
}

// empty sync.Locker for conditionals
type NoLock struct{}

func (NoLock) Lock()   {}
func (NoLock) Unlock() {}
