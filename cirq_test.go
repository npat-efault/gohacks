package cirq

import (
	"strings"
	"testing"
)

type pushFn func(interface{}) bool
type popFn func() (interface{}, bool)

func testLIFO(t *testing.T, q *CQ, maxSz int, push pushFn, pop popFn) {
	if !q.Empty() || q.Full() || q.Len() != 0 || q.MaxCap() != maxSz {
		t.Fatalf("Initially: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.MaxCap())
	}
	for i := 0; i < maxSz; i++ {
		ok := push(i)
		if !ok {
			t.Fatalf("Cannot push %d", i)
		}
		if q.Empty() || q.Len() != i+1 || q.MaxCap() != maxSz {
			t.Fatalf("After %d Pushes: E=%v, F=%v, L=%d, C=%d",
				i+1, q.Empty(), q.Full(), q.Len(), q.MaxCap())
		}
		if i < maxSz-1 && q.Full() {
			t.Fatal("Bad queue-full inication")
		}
	}
	if !q.Full() {
		t.Fatal("Bad queue-not-full inication")
	}

	for i := 0; i < maxSz; i++ {
		ei, ok := pop()
		if !ok {
			t.Fatalf("Cannot pop %d", maxSz-i-1)
		}
		if q.Full() || q.Len() != maxSz-i-1 || q.MaxCap() != maxSz {
			t.Fatalf("Bad queue after %d Pops", i)
		}
		if i < maxSz-1 && q.Empty() {
			t.Fatal("Bad queue-empty inication")
		}
		if ei.(int) != maxSz-i-1 {
			t.Fatalf("Bad element %d: %d", i, ei.(int))
		}
	}
	if !q.Empty() || q.Full() || q.Len() != 0 || q.MaxCap() != maxSz {
		t.Fatalf("Finally: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.MaxCap())
	}
}

func testFIFO(t *testing.T, q *CQ, maxSz int, push pushFn, pop popFn) {
	if q.Full() || !q.Empty() || q.Len() != 0 || q.MaxCap() != maxSz {
		t.Fatalf("Initially: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.MaxCap())
	}
	for i := 0; i < maxSz; i++ {
		ok := push(i)
		if !ok {
			t.Fatalf("Cannot push %d", i)
		}
	}
	for i := 0; i < maxSz*5; i++ {
		ei, ok := pop()
		if !ok {
			t.Fatalf("Cannot pop %d", i)
		}
		if ei.(int) != i%maxSz {
			t.Fatalf("Bad element: %d != %d", ei.(int), i%maxSz)
		}
		ok = push(i % maxSz)
		if !ok {
			t.Fatalf("Cannot push %d", i)
		}
		if !q.Full() || q.Empty() ||
			q.Len() != maxSz || q.MaxCap() != maxSz {
			t.Fatalf("After %d ops: E=%v, F=%v, L=%d, C=%d",
				i+1, q.Empty(), q.Full(), q.Len(), q.MaxCap())
		}
	}
	for i := 0; i < maxSz; i++ {
		ei, ok := pop()
		if !ok {
			t.Fatalf("Cannot pop %d", i)
		}
		if ei.(int) != i {
			t.Fatalf("Bad element: %d != %d", ei.(int), i)
		}
	}
	if q.Full() || !q.Empty() || q.Len() != 0 || q.MaxCap() != maxSz {
		t.Fatalf("Finally: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.MaxCap())
	}
}

func testResize(t *testing.T, q *CQ, maxSz int) {
	if q.Full() || !q.Empty() || q.Len() != 0 || q.MaxCap() != maxSz {
		t.Fatalf("Initially: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.MaxCap())
	}
	q.Compact(1)
	for i := 1; i <= maxSz; i <<= 1 {
		for j := i >> 1; j < i; j++ {
			ok := q.PushBack(j)
			if !ok {
				t.Fatalf("Cannot push %d", j)
			}
			if q.Cap() != i {
				t.Fatalf("Queue cap %d != %d", q.Cap(), i)
			}
		}
	}
	for i := maxSz; i > 1; i >>= 1 {
		for j := i; j > i>>1; j-- {
			_, ok := q.PopFront()
			if !ok {
				t.Fatalf("Cannot pop %d", j)
			}
		}
		q.Compact(1)
		if q.Cap() != i>>1 {
			t.Fatalf("Queue cap %d != %d", q.Cap(), i>>1)
		}
	}
	if q.Full() || q.Empty() || q.Len() != 1 || q.MaxCap() != maxSz {
		t.Fatalf("Finally: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.MaxCap())
	}
}

func TestLIFO(t *testing.T) {
	maxSz := 2048
	q := New(1, maxSz)
	testLIFO(t, q, maxSz, q.PushFront, q.PopFront)
	testLIFO(t, q, maxSz, q.PushBack, q.PopBack)
}

func TestFIFO(t *testing.T) {
	maxSz := 2048
	q := New(1, maxSz)
	testFIFO(t, q, maxSz, q.PushBack, q.PopFront)
	testFIFO(t, q, maxSz, q.PushFront, q.PopBack)
}

func TestResize(t *testing.T) {
	maxSz := 2048
	q := New(1, maxSz)
	testResize(t, q, maxSz)
}

func TestPeek(t *testing.T) {
	maxSz := 2048
	q := New(1, maxSz)
	for i := 0; i < maxSz; i++ {
		q.PushBack(i)
	}
	for i := 0; i < maxSz; i++ {
		e, ok := q.PeekFront()
		if !ok {
			t.Fatalf("Failed to PeekFront %d", i)
		}
		if e.(int) != i {
			t.Fatalf("PeekFront %d != %d", e.(int), i)
		}
		q.PopFront()
	}
	for i := 0; i < maxSz; i++ {
		q.PushFront(i)
	}
	for i := 0; i < maxSz; i++ {
		e, ok := q.PeekBack()
		if !ok {
			t.Fatalf("Failed to PeekBack %d", i)
		}
		if e.(int) != i {
			t.Fatalf("PeekBack %d != %d", e.(int), i)
		}
		q.PopBack()
	}
}

func TestMustPeek(t *testing.T) {
	maxSz := 2048
	q := New(1, maxSz)
	for i := 0; i < maxSz; i++ {
		q.PushBack(i)
	}
	for i := 0; i < maxSz; i++ {
		e := q.MustPeekFront()
		if e.(int) != i {
			t.Fatalf("PeekFront %d != %d", e.(int), i)
		}
		q.PopFront()
	}
	for i := 0; i < maxSz; i++ {
		q.PushFront(i)
	}
	for i := 0; i < maxSz; i++ {
		e := q.MustPeekBack()
		if e.(int) != i {
			t.Fatalf("PeekBack %d != %d", e.(int), i)
		}
		q.PopBack()
	}
	func() {
		defer func() {
			x := recover()
			if x == nil {
				t.Fatal("No panic when peeking empty Q")
			}
			if !strings.HasPrefix(x.(string), "MustPeek") {
				panic(x)
			}
		}()
		// q is empty, so these should panic. The panic will
		// be caught by the recover, above.
		q.MustPeekBack()
		q.MustPeekFront()
	}()
}
