// Copyright (c) 2016, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE.txt file.

package task_test

import (
	"context"
	"fmt"
	"time"

	"github.com/npat-efault/gohacks/task"
)

type BarServer struct {
	// ... various server fields ...
}

func NewBarServer( /* ... params ...*/ ) *BarServer {
	bs := &BarServer{}
	// ... bar sever instance init ...
	return bs
}

func (bs *BarServer) serve(ctx context.Context) error {
	fmt.Println("Started bar server")
	// ... server processing here ...
	select {
	case <-ctx.Done():
		fmt.Println("Bar server canceled")
		return ErrCanceled
	case <-time.After(10 * time.Second):
		fmt.Println("Bar server timed-out, exiting")

	}
	return nil
}

func (bs *BarServer) init() {
	fmt.Println("Bar server init called")
	// ... other pre-start initilization ...
}

type CtledBarServer struct {
	*task.SrvCtl
	*BarServer
}

func NewCtledBarServer( /* ... params ...*/ ) *CtledBarServer {
	cfs := &CtledBarServer{}
	cfs.BarServer = NewBarServer( /* ... params ...*/ )
	cfs.SrvCtl = task.NewSrvCtlCtx(cfs.BarServer.serve, cfs.BarServer.init)
	return cfs
}

func ExampleSrvCtl_wrapped() {
	bs := NewCtledBarServer( /* ... params ... */ )
	// Start the server
	bs.Start()
	// Wait a bit
	time.Sleep(1 * time.Second)
	// Re-start it
	bs.Start()
	// Kill it and wait to terminate
	err := bs.Kill().Wait()
	fmt.Println(err)
	// Output:
	// Bar server init called
	// Started bar server
	// Bar server canceled
	// Bar server init called
	// Started bar server
	// Bar server canceled
	// Canceled
}
