/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/XiaoMi/Gaea/util/sync2"
)

var lastID, count sync2.AtomicInt64

type TestResource struct {
	num    int64
	closed bool
}

func (tr *TestResource) Close() {
	if !tr.closed {
		count.Add(-1)
		tr.closed = true
	}
}

func PoolFactory() (Resource, error) {
	count.Add(1)
	return &TestResource{lastID.Add(1), false}, nil
}

func FailFactory() (Resource, error) {
	return nil, errors.New("Failed")
}

func SlowFailFactory() (Resource, error) {
	time.Sleep(10 * time.Millisecond)
	return nil, errors.New("Failed")
}

func TestOpen(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(PoolFactory, 6, 6, time.Second)
	p.SetDynamic(false)
	p.ScaleCapacity(5)
	var resources [10]Resource

	// Test Get
	for i := 0; i < 5; i++ {
		r, err := p.Get(ctx)
		resources[i] = r
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if p.Available() != int64(5-i-1) {
			t.Errorf("expecting %d, received %d", 5-i-1, p.Available())
		}
		if p.WaitCount() != 0 {
			t.Errorf("expecting 0, received %d", p.WaitCount())
		}
		if p.WaitTime() != 0 {
			t.Errorf("expecting 0, received %d", p.WaitTime())
		}

		// Since all connection will be connected at start and `PoolFactory` has been called 6 times
		if lastID.Get() != 6 {
			t.Errorf("Expecting %d, received %d", 6, lastID.Get())
		}

		// All connection will be connected at start, so the count is current capacity = 6 -1 = 5
		if count.Get() != 5 {
			t.Errorf("Expecting %d, received %d", 5, count.Get())
		}
	}

	// Test that Get waits
	ch := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			r, err := p.Get(ctx)
			if err != nil {
				t.Errorf("Get failed: %v", err)
			}
			resources[i] = r
		}
		for i := 0; i < 5; i++ {
			p.Put(resources[i])
		}
		ch <- true
	}()
	for i := 0; i < 5; i++ {
		// Sleep to ensure the goroutine waits
		time.Sleep(10 * time.Millisecond)
		p.Put(resources[i])
	}
	<-ch
	if p.WaitCount() != 5 {
		t.Errorf("Expecting 5, received %d", p.WaitCount())
	}
	if p.WaitTime() == 0 {
		t.Errorf("Expecting non-zero")
	}
	if lastID.Get() != 6 {
		t.Errorf("Expecting 5, received %d", lastID.Get())
	}

	// Test Close resource
	r, err := p.Get(ctx)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	r.Close()
	p.Put(nil)
	if count.Get() != 4 {
		t.Errorf("Expecting 4, received %d", count.Get())
	}
	for i := 0; i < 5; i++ {
		r, err := p.Get(ctx)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		resources[i] = r
	}
	for i := 0; i < 5; i++ {
		p.Put(resources[i])
	}
	if count.Get() != 5 {
		t.Errorf("Expecting 5, received %d", count.Get())
	}
	// the last values is 6 and since close resource at line 129 then call p.Get() 5 times, and one of five need connect
	if lastID.Get() != 7 {
		t.Errorf("Expecting 7, received %d", lastID.Get())
	}

	// ScaleCapacity
	p.ScaleCapacity(3)
	if count.Get() != 3 {
		t.Errorf("Expecting 3, received %d", count.Get())
	}
	if lastID.Get() != 7 {
		t.Errorf("Expecting 6, received %d", lastID.Get())
	}
	if p.Capacity() != 3 {
		t.Errorf("Expecting 3, received %d", p.Capacity())
	}
	if p.Available() != 3 {
		t.Errorf("Expecting 3, received %d", p.Available())
	}
	p.ScaleCapacity(6)
	if p.Capacity() != 6 {
		t.Errorf("Expecting 6, received %d", p.Capacity())
	}
	if p.Available() != 6 {
		t.Errorf("Expecting 6, received %d", p.Available())
	}
	for i := 0; i < 6; i++ {
		r, err := p.Get(ctx)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		resources[i] = r
	}
	for i := 0; i < 6; i++ {
		p.Put(resources[i])
	}
	if count.Get() != 6 {
		t.Errorf("Expecting 5, received %d", count.Get())
	}
	if lastID.Get() != 10 {
		t.Errorf("Expecting 9, received %d", lastID.Get())
	}

	// Close
	p.Close()
	if p.Capacity() != 0 {
		t.Errorf("Expecting 0, received %d", p.Capacity())
	}
	if p.Available() != 0 {
		t.Errorf("Expecting 0, received %d", p.Available())
	}
	if count.Get() != 0 {
		t.Errorf("Expecting 0, received %d", count.Get())
	}
}

func TestOpenDynamic(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(PoolFactory, 6, 10, time.Second)
	p.ScaleCapacity(5)
	p.SetDynamic(true)
	var resources [10]Resource

	// Test Get
	for i := 0; i < 7; i++ {
		r, err := p.Get(ctx)
		resources[i] = r
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if i < 5 {
			if p.Available() != int64(5-i-1) {
				t.Errorf("expecting %d, received %d", 5-i-1, p.Available())
			}
		} else {
			if p.Available() != 0 {
				t.Errorf("expecting %d, received %d", 0, p.Available())
			}
		}

		if p.WaitCount() != 0 {
			t.Errorf("expecting 0, received %d", p.WaitCount())
		}
		if p.WaitTime() != 0 {
			t.Errorf("expecting 0, received %d", p.WaitTime())
		}

		if i < 5 {
			if lastID.Get() != 6 {
				t.Errorf("Expecting %d, received %d", i+1, lastID.Get())
			}
			if count.Get() != 5 {
				t.Errorf("Expecting %d, received %d", i+1, count.Get())
			}
		} else {
			if lastID.Get() != int64(i+2) {
				t.Errorf("Expecting %d, received %d", i+2, lastID.Get())
			}
			if count.Get() != int64(i+1) {
				t.Errorf("Expecting %d, received %d", i+1, count.Get())
			}
		}
	}

	// Test that Get waits
	ch := make(chan bool)
	go func() {
		for i := 0; i < 7; i++ {
			r, err := p.Get(ctx)
			if err != nil {
				t.Errorf("Get failed: %v", err)
			}
			resources[i] = r
		}
		for i := 0; i < 7; i++ {
			p.Put(resources[i])
		}
		ch <- true
	}()
	for i := 0; i < 7; i++ {
		// Sleep to ensure the goroutine waits
		time.Sleep(10 * time.Millisecond)
		p.Put(resources[i])
	}
	<-ch
	if p.WaitCount() != 4 {
		t.Errorf("Expecting 4, received %d", p.WaitCount())
	}
	if p.WaitTime() == 0 {
		t.Errorf("Expecting non-zero")
	}
	if lastID.Get() != 11 {
		t.Errorf("Expecting 11, received %d", lastID.Get())
	}

	// Test Close resource
	r, err := p.Get(ctx)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	r.Close()
	p.Put(nil)
	if count.Get() != 9 {
		t.Errorf("Expecting 9, received %d", count.Get())
	}
	for i := 0; i < 5; i++ {
		r, err := p.Get(ctx)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		resources[i] = r
	}
	for i := 0; i < 5; i++ {
		p.Put(resources[i])
	}
	if count.Get() != 9 {
		t.Errorf("Expecting 9, received %d", count.Get())
	}
	if lastID.Get() != 11 {
		t.Errorf("Expecting 10, received %d", lastID.Get())
	}
}

func TestShrinking(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(PoolFactory, 5, 5, time.Second)
	p.SetDynamic(false)
	var resources [10]Resource
	// Leave one empty slot in the pool
	for i := 0; i < 4; i++ {
		r, err := p.Get(ctx)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		resources[i] = r
	}
	done := make(chan bool)
	go func() {
		p.ScaleCapacity(3)
		done <- true
	}()
	expected := `{"Capacity": 3, "Available": 0, "Active": 4, "InUse": 4, "MaxCapacity": 5, "WaitCount": 0, "WaitTime": 0, "IdleTimeout": 1000000000, "IdleClosed": 0}`
	for i := 0; i < 10; i++ {
		time.Sleep(10 * time.Millisecond)
		stats := p.StatsJSON()
		if stats != expected {
			if i == 9 {
				t.Errorf(`expecting '%s', received '%s'`, expected, stats)
			}
		}
	}
	// There are already 2 resources available in the pool.
	// So, returning one should be enough for ScaleCapacity to complete.
	p.Put(resources[3])
	<-done
	// Return the rest of the resources
	for i := 0; i < 3; i++ {
		p.Put(resources[i])
	}
	stats := p.StatsJSON()
	expected = `{"Capacity": 3, "Available": 3, "Active": 3, "InUse": 0, "MaxCapacity": 5, "WaitCount": 0, "WaitTime": 0, "IdleTimeout": 1000000000, "IdleClosed": 0}`
	if stats != expected {
		t.Errorf(`expecting '%s', received '%s'`, expected, stats)
	}
	if count.Get() != 3 {
		t.Errorf("Expecting 3, received %d", count.Get())
	}

	// Ensure no deadlock if ScaleCapacity is called after we start
	// waiting for a resource
	var err error
	for i := 0; i < 3; i++ {
		resources[i], err = p.Get(ctx)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
	// This will wait because pool is empty
	go func() {
		r, err := p.Get(ctx)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		p.Put(r)
		done <- true
	}()

	// This will also wait
	go func() {
		p.ScaleCapacity(2)
		done <- true
	}()
	time.Sleep(10 * time.Millisecond)

	// This should not hang
	for i := 0; i < 3; i++ {
		p.Put(resources[i])
	}
	<-done
	<-done
	if p.Capacity() != 2 {
		t.Errorf("Expecting 2, received %d", p.Capacity())
	}
	if p.Available() != 2 {
		t.Errorf("Expecting 2, received %d", p.Available())
	}
	if p.WaitCount() != 1 {
		t.Errorf("Expecting 1, received %d", p.WaitCount())
	}
	if count.Get() != 2 {
		t.Errorf("Expecting 2, received %d", count.Get())
	}

	// Test race condition of ScaleCapacity with itself
	p.ScaleCapacity(3)
	for i := 0; i < 3; i++ {
		resources[i], err = p.Get(ctx)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	}
	// This will wait because pool is empty
	go func() {
		r, err := p.Get(ctx)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		p.Put(r)
		done <- true
	}()
	time.Sleep(10 * time.Millisecond)

	// This will wait till we Put
	go p.ScaleCapacity(2)
	time.Sleep(10 * time.Millisecond)
	go p.ScaleCapacity(4)
	time.Sleep(10 * time.Millisecond)

	// This should not hang
	for i := 0; i < 3; i++ {
		p.Put(resources[i])
	}
	<-done

	err = p.ScaleCapacity(-1)
	if err == nil {
		t.Errorf("Expecting error")
	}
	err = p.ScaleCapacity(255555)
	if err == nil {
		t.Errorf("Expecting error")
	}

	if p.Capacity() != 4 {
		t.Errorf("Expecting 4, received %d", p.Capacity())
	}
	if p.Available() != 4 {
		t.Errorf("Expecting 4, received %d", p.Available())
	}
}

func TestClosing(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(PoolFactory, 5, 5, time.Second)
	p.SetDynamic(false)
	var resources [10]Resource
	for i := 0; i < 5; i++ {
		r, err := p.Get(ctx)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		resources[i] = r
	}
	ch := make(chan bool)
	go func() {
		p.Close()
		ch <- true
	}()

	// Wait for goroutine to call Close
	time.Sleep(10 * time.Millisecond)
	stats := p.StatsJSON()
	expected := `{"Capacity": 0, "Available": 0, "Active": 5, "InUse": 5, "MaxCapacity": 5, "WaitCount": 0, "WaitTime": 0, "IdleTimeout": 1000000000, "IdleClosed": 0}`
	if stats != expected {
		t.Errorf(`expecting '%s', received '%s'`, expected, stats)
	}

	// Put is allowed when closing
	for i := 0; i < 5; i++ {
		p.Put(resources[i])
	}

	// Wait for Close to return
	<-ch

	// ScaleCapacity must be ignored after Close
	err := p.ScaleCapacity(1)
	if err == nil {
		t.Errorf("expecting error")
	}

	stats = p.StatsJSON()
	expected = `{"Capacity": 0, "Available": 0, "Active": 0, "InUse": 0, "MaxCapacity": 5, "WaitCount": 0, "WaitTime": 0, "IdleTimeout": 1000000000, "IdleClosed": 0}`
	if stats != expected {
		t.Errorf(`expecting '%s', received '%s'`, expected, stats)
	}
	if lastID.Get() != 5 {
		t.Errorf("Expecting 5, received %d", count.Get())
	}
	if count.Get() != 0 {
		t.Errorf("Expecting 0, received %d", count.Get())
	}
}

func TestIdleTimeout(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(PoolFactory, 1, 1, 10*time.Millisecond)
	p.SetDynamic(false)
	defer p.Close()

	r, err := p.Get(ctx)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if count.Get() != 1 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if p.IdleClosed() != 0 {
		t.Errorf("Expecting 0, received %d", p.IdleClosed())
	}
	p.Put(r)
	if lastID.Get() != 1 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if count.Get() != 1 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if p.IdleClosed() != 0 {
		t.Errorf("Expecting 0, received %d", p.IdleClosed())
	}
	time.Sleep(20 * time.Millisecond)

	if count.Get() != 0 {
		t.Errorf("Expecting 0, received %d", count.Get())
	}
	if p.IdleClosed() != 1 {
		t.Errorf("Expecting 1, received %d", p.IdleClosed())
	}
	r, err = p.Get(ctx)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if lastID.Get() != 2 {
		t.Errorf("Expecting 2, received %d", count.Get())
	}
	if count.Get() != 1 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if p.IdleClosed() != 1 {
		t.Errorf("Expecting 1, received %d", p.IdleClosed())
	}

	// sleep to let the idle closer run while all resources are in use
	// then make sure things are still as we expect
	time.Sleep(20 * time.Millisecond)
	if lastID.Get() != 2 {
		t.Errorf("Expecting 2, received %d", count.Get())
	}
	if count.Get() != 1 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if p.IdleClosed() != 1 {
		t.Errorf("Expecting 1, received %d", p.IdleClosed())
	}
	p.Put(r)
	r, err = p.Get(ctx)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if lastID.Get() != 2 {
		t.Errorf("Expecting 2, received %d", count.Get())
	}
	if count.Get() != 1 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if p.IdleClosed() != 1 {
		t.Errorf("Expecting 1, received %d", p.IdleClosed())
	}

	// the idle close thread wakes up every 1/100 of the idle time, so ensure
	// the timeout change applies to newly added resources
	p.SetIdleTimeout(1000 * time.Millisecond)
	p.Put(r)

	time.Sleep(20 * time.Millisecond)
	if lastID.Get() != 2 {
		t.Errorf("Expecting 2, received %d", count.Get())
	}
	if count.Get() != 1 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if p.IdleClosed() != 1 {
		t.Errorf("Expecting 1, received %d", p.IdleClosed())
	}

	p.SetIdleTimeout(10 * time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	if lastID.Get() != 2 {
		t.Errorf("Expecting 2, received %d", count.Get())
	}
	if count.Get() != 0 {
		t.Errorf("Expecting 1, received %d", count.Get())
	}
	if p.IdleClosed() != 2 {
		t.Errorf("Expecting 2, received %d", p.IdleClosed())
	}
}

func TestCreateFail(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(FailFactory, 5, 5, time.Second)
	p.SetDynamic(false)
	defer p.Close()
	if _, err := p.Get(ctx); err.Error() != "Failed" {
		t.Errorf("Expecting Failed, received %v", err)
	}
	stats := p.StatsJSON()
	expected := `{"Capacity": 5, "Available": 5, "Active": 0, "InUse": 0, "MaxCapacity": 5, "WaitCount": 0, "WaitTime": 0, "IdleTimeout": 1000000000, "IdleClosed": 0}`
	if stats != expected {
		t.Errorf(`expecting '%s', received '%s'`, expected, stats)
	}
}

func TestSlowCreateFail(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(SlowFailFactory, 2, 2, time.Second)
	p.SetDynamic(false)
	defer p.Close()
	ch := make(chan bool)
	// The third Get should not wait indefinitely
	for i := 0; i < 3; i++ {
		go func() {
			p.Get(ctx)
			ch <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-ch
	}
	if p.Available() != 2 {
		t.Errorf("Expecting 2, received %d", p.Available())
	}
}

func TestTimeout(t *testing.T) {
	ctx := context.Background()
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(PoolFactory, 1, 1, time.Second)
	p.SetDynamic(false)
	defer p.Close()
	r, err := p.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	newctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	_, err = p.Get(newctx)
	cancel()
	want := "resource pool timed out"
	if err == nil || err.Error() != want {
		t.Errorf("got %v, want %s", err, want)
	}
	p.Put(r)
}

func TestExpired(t *testing.T) {
	lastID.Set(0)
	count.Set(0)
	p, _ := NewResourcePool(PoolFactory, 1, 1, time.Second)
	p.SetDynamic(false)
	defer p.Close()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	r, err := p.Get(ctx)
	if err == nil {
		p.Put(r)
	}
	cancel()
	want := "resource pool timed out"
	if err == nil || err.Error() != want {
		t.Errorf("got %v, want %s", err, want)
	}
}
