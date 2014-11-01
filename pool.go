// Package pool implements simple object recycling pools. When a
// program creates and destroys a large number of objects of the same
// type, it may be advantageous (performance-wise) to recycle (some
// of) these objects instead of collecting them and subseqently
// allocating new ones. Package pool can be used for this.
//
// First you create a pool to be used for objects of a specific type
// (e.g. Header):
//
//     p = pool.New(1024, func() { return &Header{} })
//
// Then, whenever you need a new object, you request it from the pool:
//
//     o := p.Get().(Header)
//
// When you have finished using it, you put it back in the pool (you
// recycle it):
//
//     p.Put(o)
//
// Every time you request an object from the pool, you will either get
// one of the objects stored in it or (if the pool is empty) a new
// one allocated the the supplied function.
//
// It is safe to use the same pool concurrently from multiple
// goroutines.
//
package pool

// Pool is an object recycling pool.
type Pool struct {
	queue chan interface{}
	alloc func() interface{}
}

// New creates and returns an object-recycling pool with a capacity of
// "n" objects. Function "alloc" will be called when an object is
// requested (by Pool.Get) and the pool is empty. It is ok to pass nil
// for alloc; in this case, if the pool is empty, Pool.Get will return
// nil.
func New(n int, alloc func() interface{}) *Pool {
	p := &Pool{}
	p.alloc = alloc
	p.queue = make(chan interface{}, n)
	return p
}

// Put recycles (stores, returns) an object to the pool. If the pool
// is filled to capacity, the oject is dropped.
func (p Pool) Put(i interface{}) {
	select {
	case p.queue <- i:
	default:
	}
}

// Get retrieves and returns an object from the pool. If the pool is
// empty, a new object is allocated and returned---provided that an
// "alloc" function was given when the pool was created (see
// Pool.New). If no "alloc" function was given, and the pool is empty,
// Pool.Get returns nil.
func (p Pool) Get() interface{} {
	var i interface{}
	select {
	case i = <-p.queue:
	default:
	}
	if i == nil && p.alloc != nil {
		i = p.alloc()
	}
	return i
}

// Empty removes all objects from the pool.
func (p Pool) Empty() {
	for {
		select {
		case <-p.queue:
		default:
			return
		}
	}
}

// ByteSlicePool is a specialized pool for byte-slices.
type ByteSlicePool struct {
	queue chan []byte
	alloc func() []byte
}

// NewByteSlice creates and returns a recycling pool specialized for
// byte-slices. See function New for more.
func NewByteSlice(n int, alloc func() []byte) *ByteSlicePool {
	p := &ByteSlicePool{}
	p.alloc = alloc
	p.queue = make(chan []byte, n)
	return p
}

// Put recycles (stores, returns) a byte-slice to the pool. If the pool
// is filled to capacity, the oject is dropped.
func (p ByteSlicePool) Put(s []byte) {
	select {
	case p.queue <- s:
	default:
	}
}

// Get retrieves and returns a byte-slice from the pool. See Pool.Get
// for more.
func (p ByteSlicePool) Get() []byte {
	var s []byte
	select {
	case s = <-p.queue:
	default:
	}
	if s == nil && p.alloc != nil {
		s = p.alloc()
	}
	return s
}

// Empty removes all objects from the pool.
func (p ByteSlicePool) Empty() {
	for {
		select {
		case <-p.queue:
		default:
			return
		}
	}
}
