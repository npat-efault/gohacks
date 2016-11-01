// Copyright (c) 2016, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE.txt file.

// Package task provides types and functions for managing tasks. Task
// is a process performed by a goroutine, or a group of related
// goroutines, that can be Kill'ed (stopped) and Wait'ed-for. In
// addition, a helper type is provided for controlling the life-cycle
// of servers (starting, stopping, waiting for, and restarting them).
package task

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotStarted = errors.New("Server not started")
)

// Task is a process performed by a goroutine, or a group of related
// goroutines, that can be Kill'ed (stopped) and Wait'ed-for. Starting
// the task is beyond the scope of this interface. Specific
// implementations (like the ones included in this package), provide
// methods to start the respective tasks. All Task interface methods
// can be safely called concurently from multiple goroutines.
type Task interface {
	// Kill requests that the task terminates as soon as
	// possible. Kill returns immediately and does not wait for
	// the task to terminate. It is ok to call Kill multiple
	// times. After the first, subsequent calls do nothing.
	Kill() Task
	// Wait waits for the task to terminate, and returns its exit
	// status. It is ok to call Wait multiple times. If called
	// after the task has terminated, it returns its exit status
	// immediately.
	Wait() error
	// WaitChan returns a channel that will be closed when the
	// task terminates. WaitChan does not wait for the task to
	// terminate. Receiving from the returned channel blocks until
	// it does. Once the receive from the returned channel
	// succeeds, Wait can be called to retrieve the task's exit
	// status.
	WaitChan() <-chan struct{}
}

// StartFunc is an entry-point function that can be run as a task. The
// function must take a context.Context as an argument, which may be
// canceled to request the task's termination.
type StartFunc func(context.Context) error

// Single is used to start a signle goroutine as a task (a goroutine
// that can be killed and waited-for). All Single methods can be
// called concurently.
type Single struct {
	cancel func()
	end    chan struct{}
	err    error
}

// Go starts a task using the given function as entry point.
func Go(f StartFunc) *Single {
	return GoWithContext(context.Background(), f)
}

// GoWithContext is similar to Go, but uses ctx as the parent of
// the context that will be used for the task's cancelation.
func GoWithContext(ctx context.Context, f StartFunc) *Single {
	s := &Single{}
	ctx, s.cancel = context.WithCancel(ctx)
	s.end = make(chan struct{})
	go func() {
		s.err = f(ctx)
		s.cancel()
		close(s.end)
	}()
	return s
}

// Kill requests task s to quit. Kill returns imediatelly (does not
// wait for the task to terminate).
func (s *Single) Kill() Task {
	s.cancel()
	return s
}

// Wait waits for task s to terminate and returns its exit status
// (i.e. the return value of its entry-point StartFunc).
func (s *Single) Wait() error {
	<-s.end
	return s.err
}

// WaitChan returns a channel that will be closed when the task
// terminates. After the receive from the returned channel suceeds,
// Wait can be called to retrieve the task's exit status. Usefull for
// select statements.
func (s *Single) WaitChan() <-chan struct{} {
	return s.end
}

// Grp is used to start a set (a group) of goroutines as a single
// task. All goroutines share the same cancelation context and thus
// can be killed and waited-for collectively.
type Grp struct {
	sync.Mutex
	ctx         context.Context
	cancel      func()
	wg          sync.WaitGroup
	killOnError bool
	end         chan struct{}
	err         error
}

// NewGrp creates and returns a new group. Grp implements the Task
// interface. Initially the group is empty (no running
// goroutines). Grp.Go must be subsequently called to start goroutines
// in the group.
func NewGrp() *Grp {
	g := &Grp{}
	g.ctx, g.cancel = context.WithCancel(context.Background())
	return g
}

// NewGrpWithContext is similar to NewGrp, but uses ctx as the parent
// of the context that will be used for the task's cancelation.
func NewGrpWithContext(ctx context.Context) *Grp {
	g := &Grp{}
	g.ctx, g.cancel = context.WithCancel(ctx)
	return g
}

// KillOnError enables the kill-on-error behavior for the group (by
// default disabled). If enabled, the task is Kill'ed if one of it's
// goroutines returns a non-nil error. Usually KillOnError is called
// before starting the group's goroutines.
func (g *Grp) KillOnError() *Grp {
	g.Lock()
	g.killOnError = true
	g.Unlock()
	return g
}

// Go starts a goroutine in the group using the StartFunc f as an
// entry point. Go can be called multiple times to start multiple
// goroutines.
func (g *Grp) Go(f StartFunc) *Grp {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(g.ctx); err != nil {
			g.Lock()
			if g.err == nil {
				// keep first non-nil error
				g.err = err
			}
			if g.killOnError {
				g.cancel()
			}
			g.Unlock()
		}
	}()
	return g
}

// Kill requests task g to quit (signals all its goroutines to
// terminate). Kill returns imediatelly (does not wait for the task to
// terminate).
func (g *Grp) Kill() Task {
	g.cancel()
	return g
}

// Wait waits for task g to terminate and returns its exit-status. The
// task terminates when all it's goroutines exit. As exit status of
// the task is considered the first non-nil value returned by the
// StartFunc of one of it's goroutines.
func (g *Grp) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

// WaitChan returns a channel that will be closed when task g
// terminates.
func (g *Grp) WaitChan() <-chan struct{} {
	g.Lock()
	defer g.Unlock()
	if g.end == nil {
		g.end = make(chan struct{})
		go func() {
			g.Wait()
			close(g.end)
		}()
	}
	return g.end
}

// SrvCtl is a server controller. It is a helper type that povides
// methods for starting, stopping, waiting, and re-starting
// server-instances in a convenient, race-free manner. A pointer to
// SrvCtl can be embedded (either anonymoysly or not) in the
// controlled server type. SrvCtl is initialized (via NewSrvCtl) by
// two functions (most likely method-values): One is spawned as the
// sever-task (goroutine) when Start is called; the other (which may
// be nil) is called immediately before. It is safe to call all SrvCtl
// methods concurently. See example for details.
//
type SrvCtl struct {
	m       sync.Mutex
	t       *Single
	fnStart StartFunc // func(context.Context) error
	fnPre   func()
}

// NewSrvCtl returns (a pointer to) an initialized server
// controller. The returned pointer can be embedded (either
// anonymously, or as a named field) in the controlled server type /
// structure. It is intialized with two functions (most likely
// method-values): fnStart is the function that will be spawned as the
// server task (goroutine). fnPre is a function to be called before
// spawning the server goroutine. fnPre may be nil.
func NewSrvCtl(fnStart StartFunc, fnPre func()) *SrvCtl {
	return &SrvCtl{fnStart: fnStart, fnPre: fnPre}
}

// task is a helper that reads and returns a pointer to the
// task-structure atomically
func (sc *SrvCtl) task() *Single {
	sc.m.Lock()
	t := sc.t
	sc.m.Unlock()
	return t
}

// Start starts the controlled server instance as a task (i.e in its
// own goroutine). If already running, it stops it, waits for it to
// terminate, and re-starts it.
func (sc *SrvCtl) Start() Task {
	t := sc.task()
	if t != nil {
		t.Kill().Wait()
	}
	sc.m.Lock()
	defer sc.m.Unlock()
	if sc.t != t {
		// Concurrent re-start.
		// Another racer won.
		return sc
	}
	if sc.fnPre != nil {
		sc.fnPre()
	}
	sc.t = Go(sc.fnStart)
	return sc
}

// Kill requests the termination of the controlled server instance. It
// does not wait for the server to terminate.
func (sc *SrvCtl) Kill() Task {
	t := sc.task()
	if t == nil {
		return sc
	}
	t.Kill()
	return sc
}

// Wait waits for the controlled server to terminate and returns it's
// exit code (the return-value of the fnStart function, see
// NewSrvCtl). If the managed server has never been started, returns
// ErrNotStarted.
func (sc *SrvCtl) Wait() error {
	t := sc.task()
	if t == nil {
		return ErrNotStarted
	}
	return t.Wait()
}

// WaitChan returns a channel that will be closed when the
// controlled-server terminates. After the receive from the returned
// channel suceeds, Wait can be called to retrieve the server's exit
// status. Usefull for select statements.
func (sc *SrvCtl) WaitChan() <-chan struct{} {
	t := sc.task()
	if t == nil {
		cerr := make(chan struct{})
		close(cerr)
		return cerr
	}
	return t.WaitChan()
}
