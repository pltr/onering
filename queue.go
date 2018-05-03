package onering

type Consumer interface {
	Get(interface{}) bool
	Consume(interface{})
}

type Producer interface {
	Put(interface{})
	Close()
}

type Queue interface {
	Producer
	Consumer
}

type New struct {
	Size     uint32
	MaxBatch int32
}

func (n New) SPSC() Queue {
	var spsc = new(SPSC)
	spsc.init(n.Size)
	spsc.maxbatch = n.BatchSize()
	return spsc
}

func (n New) MPSC() Queue {
	var mpsc = new(MPSC)
	mpsc.init(n.Size)
	mpsc.maxbatch = n.BatchSize()
	return mpsc
}

func (n New) SPMC() Queue {
	var spmc = new(SPMC)
	spmc.init(n.Size)
	return spmc
}

func (n New) MPMC() Queue {
	var mpmc = new(MPMC)
	mpmc.init(n.Size)
	return mpmc
}

func (n *New) BatchSize() int64 {
	if n.MaxBatch > 0 {
		return int64(n.MaxBatch)
	}
	return DefaultMaxBatch
}

//type Waiter interface {
//	Wait()
//	Signal()
//	Broadcast()
//}
//
//type Ticket interface {
//	Try(interface{}) bool
//	Use(interface{}) bool
//}

type Iter interface {
	Stop()
	Count() int
}

type iter struct {
	count int
	stop  bool
}

func (i *iter) Stop() {
	i.stop = true
}

func (i *iter) Count() int {
	return i.count
}

func (i *iter) inc() {
	i.count++
}
