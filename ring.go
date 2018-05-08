package onering

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type ring struct {
	_        [7]int64
	wp       int64
	_        [7]int64
	rp       int64
	_        [7]int64
	rc       int64 // reader cache
	_        [7]int64
	data     []unsafe.Pointer
	mask     int64
	size     int64
	maxbatch int64
	done     int32
}

func (r *ring) init(n *New) {
	r.data = make([]unsafe.Pointer, roundUp2(n.Size))
	r.mask = int64(len(r.data) - 1)

	var bs = n.BatchSize
	if bs == 0 {
		bs = DefaultMaxBatch
	}
	r.maxbatch = int64(roundUp2(bs) - 1)
}

func (r *ring) Close() {
	atomic.AddInt32(&r.done, 1)
}

func (r *ring) Done() bool {
	return atomic.LoadInt64(&r.wp) <= atomic.LoadInt64(&r.rp) && atomic.LoadInt32(&r.done) > 0
}

func (r *ring) wait() {
	runtime.Gosched()
}

func (r *ring) waitForEq(data *int64, val int64) (keep bool) {
	for keep = true; keep && atomic.LoadInt64(data) != val; runtime.Gosched() {
		keep = atomic.LoadInt64(&r.wp) > atomic.LoadInt64(&r.rp) || atomic.LoadInt32(&r.done) == 0
	}
	return
}

type multi struct {
	_ int64
	ring
	_ [42]byte
	seq []int64
}

func (c *multi) init(n *New) {
	c.ring.init(n)
	c.size = int64(len(c.data))
	c.seq = make([]int64, len(c.data))
	for i := range c.seq {
		c.seq[i] = int64(i)
	}
	c.wp = 1 // just to avoid 0-awkwardness with seq
	c.rp = 1
	c.rc = c.rp
}

func (c *multi) next(p *int64) int64 {
	return atomic.AddInt64(p, 1) - 1
}

func (c *multi) frame(p int64) (data *unsafe.Pointer, seq *int64) {
	var pos = c.mask & p
	return &c.data[pos], &c.seq[pos]
}

// empty sync.Locker for conditionals
type nolock struct{}

func (nolock) Lock()   {}
func (nolock) Unlock() {}
