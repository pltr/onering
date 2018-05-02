package onering

type Queue interface {
	Get(interface{}) bool
	Put(interface{})
	ReadTicket() Ticket
	WriteTicket() Ticket
}
type New struct {
	Type Injector
	Size uint32
}

func (nq New) SPSC() *SPSC {
	var q SPSC
	q.init(nq.Type, nq)
	return &q
}

type Waiter interface {
	Wait()
	Signal()
	Broadcast()
}

type Ticket interface {
	Try(interface{}) bool
	Use(interface{}) bool
}