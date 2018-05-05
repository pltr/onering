package onering

import (
	"math/bits"
	"runtime"
	"sync/atomic"
	"unsafe"
)

type ring struct {
	_        [8]int64
	wp       int64
	_        [7]int64
	wc       int64 // writer cache
	_        [7]int64
	rc       int64 // reader cache
	_        [7]int64
	rp       int64
	_        [7]int64
	data     []unsafe.Pointer
	mask     int64
	size     int64
	maxbatch int64
	done     int32
}

func (r *ring) init(size uint32) {
	r.data = make([]unsafe.Pointer, 1<<uint(32-bits.LeadingZeros32(size-1)))
	r.mask = int64(len(r.data) - 1)
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
	ring
	_   [42]byte
	seq []int64
}

func (c *multi) init(size uint32) {
	c.ring.init(size)
	c.size = int64(len(c.data))
	c.seq = make([]int64, len(c.data))
	c.wp = 1 // just to avoid 0-awkwardness with seq
	c.wc = c.wp
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
