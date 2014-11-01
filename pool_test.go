package pool_test

import (
	"sync"
	"testing"

	"github.com/npat-efault/tst/pool"
)

type S struct {
	id  int
	buf [64]byte
}

func TestPoolNil(t *testing.T) {
	p := pool.New(10, nil)
	for i := 0; i < 12; i++ {
		p.Put(&S{id: i})
	}
	for i := 0; i < 10; i++ {
		if s := p.Get(); s == nil {
			t.Fatalf("Elem %d is nil", i)
		} else {
			ss := s.(*S)
			if ss.id != i {
				t.Fatalf("Elem id %d != %d", ss.id, i)
			}
		}
	}
	if s := p.Get(); s != nil {
		t.Fatal("Too many items in pool")
	}
}

func TestPoolAlloc(t *testing.T) {
	p := pool.New(10, func() interface{} { return &S{} })
	for i := 0; i < 12; i++ {
		p.Put(&S{id: i})
	}
	for i := 0; i < 10; i++ {
		if s := p.Get(); s == nil {
			t.Fatalf("Elem %d is nil", i)
		} else {
			ss := s.(*S)
			if ss.id != i {
				t.Fatalf("Elem id %d != %d", ss.id, i)
			}
		}
	}
	if s := p.Get(); s == nil {
		t.Fatal("No new allocation")
	} else if ss := s.(*S); ss.id != 0 {
		t.Fatal("Unitialized new allocation")
	}
}

func TestPoolEmpty(t *testing.T) {
	p := pool.New(10, nil)
	for i := 0; i < 12; i++ {
		p.Put(&S{id: i})
	}
	p.Empty()
	if s := p.Get(); s != nil {
		t.Fatal("Pool not empty")
	}
}

func TestBSPoolNil(t *testing.T) {
	p := pool.NewByteSlice(10, nil)
	for i := 0; i < 12; i++ {
		b := make([]byte, 64)
		b[63] = byte(i)
		p.Put(b)
	}
	for i := 0; i < 10; i++ {
		if b := p.Get(); b == nil {
			t.Fatalf("Elem %d is nil", i)
		} else if b[63] != byte(i) {
			t.Fatalf("Elem %d[0] != %d", b[63], i)
		}
	}
	if b := p.Get(); b != nil {
		t.Fatal("Too many items in pool")
	}
}

func TestBSPoolAlloc(t *testing.T) {
	p := pool.NewByteSlice(10, func() []byte { return make([]byte, 64) })
	for i := 0; i < 12; i++ {
		b := make([]byte, 64)
		b[63] = byte(i)
		p.Put(b)
	}
	for i := 0; i < 10; i++ {
		if b := p.Get(); b == nil {
			t.Fatalf("Elem %d is nil", i)
		} else if b[63] != byte(i) {
			t.Fatalf("Elem %d[0] != %d", b[63], i)
		}
	}
	if b := p.Get(); b == nil {
		t.Fatal("No new allocation")
	} else if b[63] != 0 {
		t.Fatal("Unitialized new allocation")
	}
}

func TestBSPoolEmpty(t *testing.T) {
	p := pool.NewByteSlice(10, nil)
	for i := 0; i < 12; i++ {
		b := make([]byte, 64)
		b[63] = byte(i)
		p.Put(b)
	}
	p.Empty()
	if b := p.Get(); b != nil {
		t.Fatal("Pool not empty")
	}
}

// See allocations due to conversions from []byte to
// interface{}. Avoided by specialized ByteSlice pool.
//
// run with:
//      go test -v -memprofile=mem-bsp.out -bench='ByteSlicePool$'
//   or go test -v -memprofile=mem-p.out -bench='Pool$'
//   or go test -v -memprofile=mem-sp.out -bench='SyncPool$'
//
// see results with
//
//      go tool pprof --alloc_objects ./pool.test mem-xxx.out -text
//      ... etc ...

func BenchmarkAllocByteSlicePool(b *testing.B) {
	s := make([]byte, 10)
	p := pool.NewByteSlice(1, nil)
	for i := 0; i < b.N; i++ {
		p.Put(s)
		s = p.Get()
	}
}

func BenchmarkAllocPool(b *testing.B) {
	s := make([]byte, 10)
	p := pool.New(1, nil)
	for i := 0; i < b.N; i++ {
		p.Put(s)
		s = p.Get().([]byte)
	}
}

func BenchmarkAllocSyncPool(b *testing.B) {
	var p sync.Pool
	s := make([]byte, 10)
	for i := 0; i < b.N; i++ {
		p.Put(s)
		s = p.Get().([]byte)
	}
}
