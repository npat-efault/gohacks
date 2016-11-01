// Copyright (c) 2016, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE.txt file.

package task_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/npat-efault/gohacks/task"
)

var ErrCanceled = errors.New("Canceled")

type FooServer struct {
	*task.SrvCtl
	// ... other server fields ...
}

func NewFooServer( /* ... params ...*/ ) *FooServer {
	fs := &FooServer{}
	fs.SrvCtl = task.NewSrvCtlCtx(fs.serve, fs.init)
	// ... foo server instance init ...
	return fs
}

func (fs *FooServer) serve(ctx context.Context) error {
	fmt.Println("Started foo server")
	// ... server processing here ...
	select {
	case <-ctx.Done():
		fmt.Println("Foo server canceled")
		return ErrCanceled
	case <-time.After(10 * time.Second):
		fmt.Println("Foo server timed-out, exiting")

	}
	return nil
}

func (fs *FooServer) init() {
	fmt.Println("Foo server init called")
	// ... other pre-start initilization ...
}

func ExampleSrvCtl_embedded() {
	fs := NewFooServer( /* ... params ... */ )
	// Start the server
	fs.Start()
	// Wait a bit
	time.Sleep(1 * time.Second)
	// Re-start it
	fs.Start()
	// Kill it and wait to terminate
	err := fs.Kill().Wait()
	fmt.Println(err)
	// Output:
	// Foo server init called
	// Started foo server
	// Foo server canceled
	// Foo server init called
	// Started foo server
	// Foo server canceled
	// Canceled
}
