package cirq

import "testing"

type pushFn func(interface{}) bool
type popFn func() (interface{}, bool)

func testLIFO(t *testing.T, q *CQ, maxSz int, push pushFn, pop popFn) {
	if !q.Empty() || q.Full() || q.Len() != 0 || q.Cap() != maxSz {
		t.Fatalf("Initially: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.Cap())
	}
	for i := 0; i < maxSz; i++ {
		ok := push(i)
		if !ok {
			t.Fatalf("Cannot push %d", i)
		}
		if q.Empty() || q.Len() != i+1 || q.Cap() != maxSz {
			t.Fatalf("After %d Pushes: E=%v, F=%v, L=%d, C=%d",
				i+1, q.Empty(), q.Full(), q.Len(), q.Cap())
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
		if q.Full() || q.Len() != maxSz-i-1 || q.Cap() != maxSz {
			t.Fatalf("Bad queue after %d Pops", i)
		}
		if i < maxSz-1 && q.Empty() {
			t.Fatal("Bad queue-empty inication")
		}
		if ei.(int) != maxSz-i-1 {
			t.Fatalf("Bad element %d: %d", i, ei.(int))
		}
	}
	if !q.Empty() || q.Full() || q.Len() != 0 || q.Cap() != maxSz {
		t.Fatalf("Finally: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.Cap())
	}
}

func testFIFO(t *testing.T, q *CQ, maxSz int, push pushFn, pop popFn) {
	if q.Full() || !q.Empty() || q.Len() != 0 || q.Cap() != maxSz {
		t.Fatalf("Initially: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.Cap())
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
			q.Len() != maxSz || q.Cap() != maxSz {
			t.Fatalf("After %d ops: E=%v, F=%v, L=%d, C=%d",
				i+1, q.Empty(), q.Full(), q.Len(), q.Cap())
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
	if q.Full() || !q.Empty() || q.Len() != 0 || q.Cap() != maxSz {
		t.Fatalf("Finally: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.Cap())
	}
}

func testResize(t *testing.T, q *CQ, maxSz int) {
	if q.Full() || !q.Empty() || q.Len() != 0 || q.Cap() != maxSz {
		t.Fatalf("Initially: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.Cap())
	}
	q.Compact(1)
	for i := 1; i <= maxSz; i <<= 1 {
		for j := i >> 1; j < i; j++ {
			ok := q.PushBack(j)
			if !ok {
				t.Fatalf("Cannot push %d", j)
			}
			if q.sz != uint32(i) {
				t.Fatalf("Queue sz %d != %d", q.sz, i)
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
		if q.sz != uint32(i>>1) {
			t.Fatalf("Queue sz %d != %d", q.sz, i>>1)
		}
	}
	if q.Full() || q.Empty() || q.Len() != 1 || q.Cap() != maxSz {
		t.Fatalf("Finally: E=%v, F=%v, L=%d, C=%d",
			q.Empty(), q.Full(), q.Len(), q.Cap())
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
