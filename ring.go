package onering

import (
	"math/bits"
	"sync/atomic"
)

const MaxBatch = 255

type ring struct {
	wp      int64
	_       [7]int64
	rp      int64
	_       [7]int64
	data    []int64
	mask    int64
	done    int32
}


func (r *ring) Init(size uint) {
	r.data = make([]int64, 1 << uint(64 - bits.LeadingZeros(size-1)))
	r.mask = int64(len(r.data) - 1)
	r.done = 0
}


func (r *ring) Close() {
	atomic.StoreInt32(&r.done, 1)
}


func (r *ring) Open() bool {
	return atomic.LoadInt32(&r.done) == 0 || atomic.LoadInt64(&r.wp) - atomic.LoadInt64(&r.rp) > 0
}

type commit struct {
	ring
	log []int64
	_ [4]int64
}

func (c *commit) Init(size uint) {
	c.ring.Init(size)
	c.log = make([]int64, len(c.data))
}