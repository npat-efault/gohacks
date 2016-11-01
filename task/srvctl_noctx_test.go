// Copyright (c) 2016, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE.txt file.

package task_test

import (
	"fmt"
	"time"

	"github.com/npat-efault/gohacks/task"
)

type QuxServer struct {
	*task.SrvCtl
	quit chan struct{}
	// ... various other server fields ...
}

func NewQuxServer( /* ... params ...*/ ) *QuxServer {
	qs := &QuxServer{}
	qs.SrvCtl = task.NewSrvCtl(qs.serve, qs.reset, qs.kill)
	qs.quit = make(chan struct{})
	// ... other qux sever instance init ...
	return qs
}

func (qs *QuxServer) serve() error {
	fmt.Println("Started qux server")
	// ... server processing here ...
	select {
	case <-qs.quit:
		fmt.Println("Qux server quiting")
		return ErrCanceled
	case <-time.After(10 * time.Second):
		fmt.Println("Qux server timed-out, exiting")
	}
	return nil
}

func (qs *QuxServer) reset() {
	fmt.Println("Qux server reset called")
	qs.quit = make(chan struct{})
	// ... other pre-start initilization ...
}

func (qs *QuxServer) kill() {
	fmt.Println("Qux server kill called")
	close(qs.quit)
}

func ExampleSrvCtlNoCtx() {
	bs := NewQuxServer( /* ... params ... */ )
	// Start the server
	bs.Start()
	// Wait a bit
	time.Sleep(1 * time.Second)
	// Re-start it
	bs.Start()
	// Wait a bit
	time.Sleep(1 * time.Second)
	// Kill it and wait to terminate
	err := bs.Kill().Wait()
	fmt.Println(err)
	// Output:
	// Qux server reset called
	// Started qux server
	// Qux server kill called
	// Qux server quiting
	// Qux server reset called
	// Started qux server
	// Qux server kill called
	// Qux server quiting
	// Canceled
}
