package onering

/*
Consumer represents consuming queue
*/
type Consumer interface {
	// Get(**T) expects a pointer to a pointer and returns a boolean value informing if the read was successful
	Get(interface{}) bool
	// Consume(func(onering.Iter, *T)) expects a function with 2 arguments, onering.Iter and *T
	Consume(interface{})
}

/*
Producer represents producing queue
*/
type Producer interface {
	// Put(T) will accept anything, but it's strongly recommended
	// to only call it with pointers to avoid heap allocation
	Put(interface{})
	// Close() closes the queue.
	// The actual consumption will only stop after all pending messages have been consumed.
	Close()
}

// Generic queue interface mathing all implementations
type Queue interface {
	Producer
	Consumer
}

// Iter is a generic loop interface
type Iter interface {
	// stops consuming function
	Stop()
	// returns the current iteration count
	Count() int
}

// New is a configuration structure for the queue constructor
type New struct {
	// Size (Capacity) of the queue
	Size     uint32
	// Maximum number of batched messages
	MaxBatch int32
}

// SPSC constructs a Single Producer/Single Consumer queue
func (n New) SPSC() Queue {
	var spsc = new(SPSC)
	spsc.init(n.Size)
	spsc.maxbatch = n.BatchSize()
	return spsc
}

// MPSC constructs a Multi-Producer/Single Consumer queue
func (n New) MPSC() Queue {
	var mpsc = new(MPSC)
	mpsc.init(n.Size)
	mpsc.maxbatch = n.BatchSize()
	return mpsc
}

// SPMC constructs a Single Producer/Multi-Consumer queue
func (n New) SPMC() Queue {
	var spmc = new(SPMC)
	spmc.init(n.Size)
	return spmc
}

// MPMC constructs a Multi-Producer/Multi-Consumer queue.
// This is the default and the most versatile/safest queue.
// However it will not implement many of the optimizations available to other queue types
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
