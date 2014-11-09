// Package cirq provides a circular double-ended queue
// implementation. The implementation is based on slices and supports
// the typical PushFront / PushBack, PopFront / PopBack, PeekFront /
// PeekBack operations. The queue can grow dynamically and stores
// elements of type interface{}.
//
// You can generate queue implementations specialized to specific
// element data-types using the cirq_gen.sh script. See Makefile for
// an example.
//
package cirq
