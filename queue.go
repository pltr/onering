package onering

type Queue interface {
	Read(interface{}) bool
	Write(interface{})
}

type New struct {
	Type interface{}
	Capacity uint32
	BatchSize uint32
	Wait Waiter
}

func (nq New) SPSC() *SPSC {
	var q SPSC
	q.Init(nq.Type, nq.Capacity)
	return &q
}

type Waiter interface {
	Wait()
	Signal()
	Broadcast()
}