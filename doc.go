// Copyright (c) 2014, Nick Patavalis (npat@efault.net).
// All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.

// Package cirq provides a circular double-ended queue
// implementation. The implementation is based on slices and supports
// the typical PushFront / PushBack, PopFront / PopBack, PeekFront /
// PeekBack operations. The queue can grow dynamically and stores
// elements of type interface{}.
//
// You can generate queue implementations specialized to specific
// element data-types using the cirq_gen.sh script. See "doc.go" in
// the package sources for an example.
//
package cirq

//go:generate ./cirq_gen.sh cirq.go cirq CQ New "interface{}"
