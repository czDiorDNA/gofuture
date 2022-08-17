// Copyright (c) 2022 DiorDNA. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package gofuture

import "time"

func MakeFutureTask[T any, F func() T](task F) *Future[T] {
	f := MakeFuture[T]()
	go func ()  {
		f.Set(task())
	}()
	return f
}

func MakeFutureTaskWithTimeout[T any, F func() T](t time.Duration, task F) *Future[T] {
	f := MakeFutureWithTimeout[T](t)
	go func ()  {
		f.Set(task())
	}()
	return f
}
