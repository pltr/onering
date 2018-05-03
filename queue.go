package onering

type Queue interface {
	Get(interface{}) bool
	Put(interface{})
	ReadTicket() Ticket
	WriteTicket() Ticket
}
type New struct {
	Size uint32
}

func (n New) SPSC() (spsc *SPSC) {
	spsc = new(SPSC)
	spsc.init(n.Size)
	return
}

func (n New) MPSC() (mpsc *MPSC) {
	mpsc = new(MPSC)
	mpsc.init(n.Size)
	return
}

func (n New) SPMC() (spmc *SPMC) {
	spmc = new(SPMC)
	spmc.init(n.Size)
	return
}

func (n New) MPMC() (mpmc *MPMC) {
	mpmc = new(MPMC)
	mpmc.init(n.Size)
	return
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