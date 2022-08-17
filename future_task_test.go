// Copyright (c) 2022 DiorDNA. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

package gofuture_test

import (
	"testing"
	"time"

	"github.com/czDiorDNA/gofuture"
)

func TestFutureTask(t *testing.T) {
	const RETURN_VAL = 1
	f := gofuture.MakeFutureTask(func () int {
		return RETURN_VAL
	})
	if val, ok := f.Get(); !ok || val != RETURN_VAL {
		t.Fatalf("unexpected result: ok: %t, val: %d", ok, val)
	}
}

func TestFutureTaskWithDuration(t *testing.T) {
	const RETURN_VAL = 1
	f := gofuture.MakeFutureTaskWithTimeout(10 * time.Millisecond, func () int {
		return RETURN_VAL
	})
	if val, ok := f.Get(); !ok || val != RETURN_VAL {
		t.Fatalf("unexpected result: ok: %t, val: %d", ok, val)
	}
	time.Sleep(1 * time.Second)
	if val, ok := f.Get(); !ok || val != RETURN_VAL {
		t.Fatalf("unexpected result: ok: %t, val: %d", ok, val)
	}
}

func TestFutureTaskWithTimeout(t *testing.T) {
	const RETURN_VAL = 1
	f := gofuture.MakeFutureTaskWithTimeout(10 * time.Millisecond, func () int {
		time.Sleep(1 * time.Second)
		return RETURN_VAL
	})
	if val, ok := f.Get(); ok {
		t.Fatalf("unexpected result: ok: %t, val: %d", ok, val)
	}
}
