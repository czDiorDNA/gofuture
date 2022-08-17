// Copyright (c) 2022 DiorDNA. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package gofuture

import (
	"context"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	FUTURE_IDLE = iota
	FUTURE_WAIT
	FUTURE_READY
	FUTURE_ERROR
)

// Future is an object that can retrieve a value from some
// provider object or function, properly synchronizing this
// access if in different goroutines.
//
// Calling the Future.Get on a valid Future blocks the goroutine
// until the provider makes the shared state ready or error if
// overtime which will return false.
//
// Calling the Future.Set sets the val member if the state is
// idle.
type Future[T any] struct {
	// The val is the result which be retrieved.
	val T
	// The state indicates that whether the val has been retrieved.
	state int32
	// The ctx is used for timeout processing.
	ctx context.Context
	// The cancel is used for timeout processing.
	cancel context.CancelFunc
}

func MakeFuture[T any]() *Future[T] {
	return &Future[T]{
		state:  FUTURE_IDLE,
		ctx:    context.Background(),
		cancel: func() {},
	}
}

func MakeFutureWithTimeout[T any](t time.Duration) *Future[T] {
	f := MakeFuture[T]()
	f.ctx, f.cancel = context.WithTimeout(context.Background(), t)
	return f
}

func (f *Future[T]) Get() (T, bool) {
	for {
		select {
		case <-f.ctx.Done():
			if atomic.CompareAndSwapInt32(&f.state, FUTURE_IDLE, FUTURE_ERROR) {
				return f.val, false
			}
			state := atomic.LoadInt32(&f.state)
			if state == FUTURE_ERROR {
				return f.val, false
			}
			if state == FUTURE_READY {
				return f.val, true
			}
		default:
			state := atomic.LoadInt32(&f.state)
			if state == FUTURE_READY {
				return f.val, true
			}
			if state == FUTURE_ERROR {
				return f.val, false
			}
		}
		runtime.Gosched()
	}
}

func (f *Future[T]) Set(val T) bool {
	if atomic.CompareAndSwapInt32(&f.state, FUTURE_IDLE, FUTURE_WAIT) {
		f.val = val
		atomic.StoreInt32(&f.state, FUTURE_READY)
		return true
	}
	return false
}
