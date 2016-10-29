package cirq_test

import (
	"fmt"

	"github.com/npat-efault/gohacks/cirq"
)

func Example() {
	// Create a queue with initial size 4 and capacity (maximum
	// allowed size) 16
	q := cirq.New(4, 16)

	// Fill the queue up to capacity.
	for i := 0; !q.Full(); i++ {
		q.PushBack(i)
	}

	// Test that the queue is full.
	if !q.Full() {
		fmt.Println("Queue is not full!")
	}

	// Remove and print the first 4 elements.
	for i := 0; i < 4; i++ {
		el, ok := q.PopFront()
		if !ok {
			// Can't happen in this example.
			break
		}
		fmt.Println(el.(int))
	}

	// Remove the remaining elements from the queue.
	for !q.Empty() {
		q.PopFront()
	}

	// Print queue lenght (=0) and capacity (=16)
	fmt.Println(q.Len(), q.Cap())

	// Output:
	// 0
	// 1
	// 2
	// 3
	// 0 16
}
