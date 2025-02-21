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

// Package util ResourcePool provides functionality to manage and reuse resources
// like connections.
package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
	"golang.org/x/time/rate"
	"sync"
	"time"

	log "github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/util/sync2"
	"github.com/XiaoMi/Gaea/util/timer"
	"github.com/google/uuid"
)

var (
	// ErrClosed is returned if ResourcePool is used when it's closed.
	ErrClosed = errors.New("resource pool is closed")

	// ErrTimeout is returned if a resource get times out.
	ErrTimeout = errors.New("resource pool timed out")
)

// Factory is a function that can be used to create a resource.
type Factory func() (Resource, error)

// Resource defines the interface that every resource must provide.
// Thread synchronization between Close() and IsClosed()
// is the responsibility of the caller.
type Resource interface {
	Close()
}

// ResourcePool allows you to use a pool of resources.
// ResourcePool允许你使用各种资源池，需要根据提供的factory创建特定的资源，比如连接
type ResourcePool struct {
	resources   chan resourceWrapper
	factory     Factory
	capacity    sync2.AtomicInt64
	idleTimeout sync2.AtomicDuration
	idleTimer   *timer.Timer
	capTimer    *timer.Timer

	// stats
	available    sync2.AtomicInt64
	active       sync2.AtomicInt64
	inUse        sync2.AtomicInt64
	waitCount    sync2.AtomicInt64
	waitTime     sync2.AtomicDuration
	idleClosed   sync2.AtomicInt64
	baseCapacity sync2.AtomicInt64
	maxCapacity  sync2.AtomicInt64
	lock         *sync.Mutex
	scaleOutTime int64
	scaleInTodo  chan int8
	Dynamic      bool
}

type resourceWrapper struct {
	resource Resource
	timeUsed time.Time
}

// NewResourcePool creates a new ResourcePool pool.
// capacity is the number of possible resources in the pool:
// there can be up to 'capacity' of these at a given time.
// maxCap specifies the extent to which the pool can be resized
// in the future through the SetCapacity function.
// You cannot resize the pool beyond maxCap.
// If a resource is unused beyond idleTimeout, it's discarded.
// An idleTimeout of 0 means that there is no timeout.
// 创建一个资源池子，capacity是池子中可用资源数量
// maxCap代表最大资源数量
// 超过设定空闲时间的连接会被丢弃
// 资源池会根据传入的factory进行具体资源的初始化，比如建立与mysql的连接
func NewResourcePool(factory Factory, capacity, maxCap int, idleTimeout time.Duration) (*ResourcePool, error) {
	if capacity <= 0 || maxCap <= 0 || capacity > maxCap {
		return nil, fmt.Errorf("invalid/out of range capacity")
	}
	rp := &ResourcePool{
		resources:    make(chan resourceWrapper, maxCap),
		factory:      factory,
		available:    sync2.NewAtomicInt64(int64(capacity)),
		capacity:     sync2.NewAtomicInt64(int64(capacity)),
		idleTimeout:  sync2.NewAtomicDuration(idleTimeout),
		baseCapacity: sync2.NewAtomicInt64(int64(capacity)),
		maxCapacity:  sync2.NewAtomicInt64(int64(maxCap)),
		lock:         &sync.Mutex{},
		scaleInTodo:  make(chan int8, 1),
		Dynamic:      true, // 动态扩展连接池
	}

	for i := 0; i < capacity; i++ {
		rp.resources <- resourceWrapper{}
	}

	if idleTimeout != 0 {
		rp.idleTimer = timer.NewTimer(idleTimeout / 10)
		rp.idleTimer.Start(rp.closeIdleResources)
	}
	rp.capTimer = timer.NewTimer(5 * time.Second)
	rp.capTimer.Start(rp.scaleInResources)
	return rp, nil
}

// Close empties the pool calling Close on all its resources.
// You can call Close while there are outstanding resources.
// It waits for all resources to be returned (Put).
// After a Close, Get is not allowed.
func (rp *ResourcePool) Close() {
	if rp.idleTimer != nil {
		rp.idleTimer.Stop()
	}
	if rp.capTimer != nil {
		rp.capTimer.Stop()
	}
	_ = rp.ScaleCapacity(0)
}

func (rp *ResourcePool) SetDynamic(value bool) {
	rp.Dynamic = value
}

// IsClosed returns true if the resource pool is closed.
func (rp *ResourcePool) IsClosed() (closed bool) {
	return rp.capacity.Get() == 0
}

// closeIdleResources scans the pool for idle resources
// 定期回收超过IdleTimeout的资源
func (rp *ResourcePool) closeIdleResources() {
	available := int(rp.Available())
	idleTimeout := rp.IdleTimeout()

	for i := 0; i < available; i++ {
		var wrapper resourceWrapper
		select {
		case wrapper, _ = <-rp.resources:
		default:
			// stop early if we don't get anything new from the pool
			return
		}

		if wrapper.resource != nil && idleTimeout > 0 && wrapper.timeUsed.Add(idleTimeout).Sub(time.Now()) < 0 {
			wrapper.resource.Close()
			wrapper.resource = nil
			rp.idleClosed.Add(1)
			rp.active.Add(-1)
		}

		rp.resources <- wrapper
	}
}

// Get will return the next available resource. If capacity
// has not been reached, it will create a new one using the factory. Otherwise,
// it will wait till the next resource becomes available or a timeout.
// A timeout of 0 is an indefinite wait.
// Get会返回下一个可用的资源
// 如果容量没有达到上线，它会根据factory创建一个新的资源，否则会一直等待直到资源可用或超时
func (rp *ResourcePool) Get(ctx context.Context) (resource Resource, err error) {
	return rp.get(ctx, true)
}

func (rp *ResourcePool) get(ctx context.Context, wait bool) (resource Resource, err error) {
	// If ctx has already expired, avoid racing with rp's resource channel.
	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	default:
	}

	// Fetch
	var wrapper resourceWrapper
	var ok bool
	select {
	case wrapper, ok = <-rp.resources:
	default:
		var newWrapper resourceWrapper
		if rp.Dynamic {
			newWrapper, ok = rp.scaleOutResources()
		}
		if ok {
			wrapper = newWrapper
		} else {
			if !wait {
				return nil, nil
			}
			startTime := time.Now()
			select {
			case wrapper, ok = <-rp.resources:
			case <-ctx.Done():
				rp.recordWait(startTime)
				return nil, ErrTimeout
			}
			endTime := time.Now()
			if startTime.UnixNano()/100000 != endTime.UnixNano()/100000 {
				rp.recordWait(startTime)
			}
		}
	}
	if !ok {
		return nil, ErrClosed
	}

	if wrapper.resource == nil {
		wrapper.resource, err = rp.createResourceWithRetry(ctx)
		if err != nil {
			rp.resources <- resourceWrapper{}
			return nil, err
		}
		rp.active.Add(1)
	}
	rp.available.Add(-1)
	rp.inUse.Add(1)
	return wrapper.resource, err
}

func (rp *ResourcePool) createResourceWithRetry(ctx context.Context) (Resource, error) {
	var resource Resource
	var err error

	// Generate a UUID for this connection attempt
	connectionUUID, uuidErr := uuid.NewRandom()
	if uuidErr != nil {
		// Handle UUID generation error
		log.Warn("failed to generate UUID for connection attempt: %v", uuidErr)
		connectionUUID = uuid.UUID{} // Zero UUID
	}

	// Initialize retry count
	retryCount := 0
	var startTime = time.Now()

	for i := 0; i < 3; i++ {
		// Increment retry count
		retryCount = i + 1

		// Create a channel to receive the factory result
		resultChan := make(chan struct {
			resource Resource
			err      error
		}, 1)

		// Start a goroutine to call rp.factory()
		go func() {
			res, err := rp.factory()
			resultChan <- struct {
				resource Resource
				err      error
			}{resource: res, err: err}
		}()

		// Wait for either the factory to finish or the context to be canceled
		select {
		case <-ctx.Done():
			rp.recordWait(startTime)
			log.Warn("connection attempt canceled. UUID: %s, Retry: %d, Error: %v", connectionUUID.String(), retryCount, ctx.Err())
			return nil, fmt.Errorf("new connection timeout: %s", ctx.Err())
		case result := <-resultChan:
			resource = result.resource
			err = result.err
		}

		if err == nil {
			// Successfully created the resource
			break
		}
		// If the first attempt fails, start the wait timer
		// Log the error with UUID and retry count
		log.Warn("connection attempt failed. UUID: %s, Retry: %d, Error: %v", connectionUUID.String(), retryCount, err)
	}

	// If we needed to wait (first attempt failed), record the wait time once
	if retryCount > 1 {
		rp.recordWait(startTime)
	}

	if err != nil {
		log.Warn("[ERROR] all connection attempts failed. UUID: %s, Total Retries: %d, Last Error: %v", connectionUUID.String(), retryCount, err)
		return nil, err
	}

	return resource, nil
}

// Put will return a resource to the pool. For every successful Get,
// a corresponding Put is required. If you no longer need a resource,
// you will need to call Put(nil) instead of returning the closed resource.
// The will eventually cause a new resource to be created in its place.
func (rp *ResourcePool) Put(resource Resource) {
	var wrapper resourceWrapper
	if resource != nil {
		wrapper = resourceWrapper{resource, time.Now()}
	} else {
		rp.active.Add(-1)
	}
	select {
	case rp.resources <- wrapper:
	default:
		panic(errors.New("attempt to Put into a full ResourcePool"))
	}
	rp.inUse.Add(-1)
	rp.available.Add(1)
}

func (rp *ResourcePool) SetCapacity(capacity int) error {
	oldcap := rp.baseCapacity.Get()
	rp.baseCapacity.CompareAndSwap(oldcap, int64(capacity))
	if int(oldcap) < capacity {
		rp.ScaleCapacity(capacity)
	}
	return nil
}

// SetCapacity changes the capacity of the pool.
// You can use it to shrink or expand, but not beyond
// the max capacity. If the change requires the pool
// to be shrunk, SetCapacity waits till the necessary
// number of resources are returned to the pool.
// A SetCapacity of 0 is equivalent to closing the ResourcePool.
func (rp *ResourcePool) ScaleCapacity(capacity int) error {
	if capacity < 0 || capacity > int(rp.maxCapacity.Get()) {
		return fmt.Errorf("capacity %d is out of range", capacity)
	}

	// Atomically swap new capacity with old, but only
	// if old capacity is non-zero.
	var oldcap int
	for {
		oldcap = int(rp.capacity.Get())
		if oldcap == 0 {
			return ErrClosed
		}
		if oldcap == capacity {
			return nil
		}
		if rp.capacity.CompareAndSwap(int64(oldcap), int64(capacity)) {
			break
		}
	}

	if capacity < oldcap {
		for i := 0; i < oldcap-capacity; i++ {
			wrapper := <-rp.resources
			if wrapper.resource != nil {
				wrapper.resource.Close()
				rp.active.Add(-1)
			}
			rp.available.Add(-1)
		}
	} else {
		for i := 0; i < capacity-oldcap; i++ {
			rp.resources <- resourceWrapper{}
			rp.available.Add(1)
		}
	}
	if capacity == 0 {
		close(rp.resources)
	}
	return nil
}

// 扩容
func (rp *ResourcePool) scaleOutResources() (resourceWrapper, bool) {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	if rp.capacity.Get() < rp.maxCapacity.Get() {
		wrapper, ok := rp.AddCapacityResource()
		rp.scaleOutTime = time.Now().Unix()
		return wrapper, ok
	}
	return resourceWrapper{}, false
}

// 扩容并获取连接, 外层加锁了，所以这边不加锁
func (rp *ResourcePool) AddCapacityResource() (resourceWrapper, bool) {
	capacity := int(rp.capacity.Get())
	if capacity < 0 || capacity >= int(rp.maxCapacity.Get()) {
		return resourceWrapper{}, false
	}
	rp.capacity.Add(1)
	rp.available.Add(1)
	return resourceWrapper{}, true
}

// 缩容
func (rp *ResourcePool) scaleInResources() {
	rp.lock.Lock()
	defer rp.lock.Unlock()
	if rp.capacity.Get() > rp.baseCapacity.Get() && time.Now().Unix()-rp.scaleOutTime > 60 {
		select {
		case rp.scaleInTodo <- 0:
			go func() {
				rp.ScaleCapacity(int(rp.capacity.Get()) - 1)
				<-rp.scaleInTodo
			}()
		default:
			return
		}
	}
}

func (rp *ResourcePool) recordWait(start time.Time) {
	rp.waitCount.Add(1)
	rp.waitTime.Add(time.Now().Sub(start))
}

// SetIdleTimeout sets the idle timeout. It can only be used if there was an
// idle timeout set when the pool was created.
func (rp *ResourcePool) SetIdleTimeout(idleTimeout time.Duration) {
	if rp.idleTimer == nil {
		panic("SetIdleTimeout called when timer not initialized")
	}

	rp.idleTimeout.Set(idleTimeout)
	rp.idleTimer.SetInterval(idleTimeout / 10)
}

// StatsJSON returns the stats in JSON format.
func (rp *ResourcePool) StatsJSON() string {
	return fmt.Sprintf(`{"Capacity": %v, "Available": %v, "Active": %v, "InUse": %v, "MaxCapacity": %v, "WaitCount": %v, "WaitTime": %v, "IdleTimeout": %v, "IdleClosed": %v}`,
		rp.Capacity(),
		rp.Available(),
		rp.Active(),
		rp.InUse(),
		rp.MaxCap(),
		rp.WaitCount(),
		rp.WaitTime().Nanoseconds(),
		rp.IdleTimeout().Nanoseconds(),
		rp.IdleClosed(),
	)
}

// Capacity returns the capacity.
func (rp *ResourcePool) Capacity() int64 {
	return rp.capacity.Get()
}

// Available returns the number of currently unused and available resources.
func (rp *ResourcePool) Available() int64 {
	return rp.available.Get()
}

// Active returns the number of active (i.e. non-nil) resources either in the
// pool or claimed for use
func (rp *ResourcePool) Active() int64 {
	return rp.active.Get()
}

// InUse returns the number of claimed resources from the pool
func (rp *ResourcePool) InUse() int64 {
	return rp.inUse.Get()
}

// MaxCap returns the max capacity.
func (rp *ResourcePool) MaxCap() int64 {
	return int64(cap(rp.resources))
}

// WaitCount returns the total number of waits.
func (rp *ResourcePool) WaitCount() int64 {
	return rp.waitCount.Get()
}

// WaitTime returns the total wait time.
func (rp *ResourcePool) WaitTime() time.Duration {
	return rp.waitTime.Get()
}

// IdleTimeout returns the idle timeout.
func (rp *ResourcePool) IdleTimeout() time.Duration {
	return rp.idleTimeout.Get()
}

// IdleClosed returns the count of resources closed due to idle timeout.
func (rp *ResourcePool) IdleClosed() int64 {
	return rp.idleClosed.Get()
}

type ResourceLimitWrapper struct {
	sync.Mutex
	factory Factory
	limiter *rate.Limiter
	addr    string
	running sync2.AtomicInt64
}

// NewResourceLimitWrapper returns a new ResourceLimitWrapper.
func NewResourceLimitWrapper(addr string, resourceLimit int, factory Factory) *ResourceLimitWrapper {
	return &ResourceLimitWrapper{
		addr:    addr,
		factory: factory,
		limiter: rate.NewLimiter(rate.Limit(resourceLimit), resourceLimit),
		running: sync2.NewAtomicInt64(0),
	}
}

/*
每秒允许创建有限个 resource
符合预期的场景，当某个时刻实例突然不可用了，连接被阻塞了，能够快速识别出来，而不是一直等待
因此，当同时有 limit 个数的 resouce 在创建中被阻塞了，则任务数据库实例存在问题，返回 mysql.ConnectionError 设置实例状态为 down
假设 limit 为 100，每秒种允许新建 100 个连接
如果只有 golang.org/x/time/rate 包，则会有个问题，并发量高的情况下，比如 1000，如果前面 100 个连接创建成功了，则后面连接无法继续创建了
即
如果 每秒超过了 建立 100 个连接，但是正在建连的线程数 < 100 仍然允许继续建连
如果 正在建立连接的线程数为 100，仍然下一秒后，可以以每秒 100 的数量继续建连
当然如果 正在建立连接的线程数为 100， 当前秒已经建立了 100个连接，该秒内不再被允许建连，返回 mysql.ConnectionError 设置实例状态为 down
*/
func (r *ResourceLimitWrapper) Factory() (Resource, error) {
	ok1 := r.limiter.Allow()
	ok2 := r.runningAllow()
	if ok2 {
		defer r.runningFinish()
	}
	if ok1 || ok2 {
		return r.factory()
	}
	return nil, &mysql.ConnectionError{
		Addr: r.addr,
		Msg:  "resource limit exceeded",
	}
}

func (r *ResourceLimitWrapper) runningAllow() bool {
	r.Lock()
	defer r.Unlock()
	if r.running.Get() < int64(r.limiter.Limit()) {
		r.running.Add(1)
		return true
	}
	return false
}

func (r *ResourceLimitWrapper) runningFinish() {
	r.Lock()
	defer r.Unlock()
	r.running.Add(-1)
	if r.running.Get() < 0 {
		r.running.Set(0)
	}
}

func (r *ResourceLimitWrapper) GetRunningNum() int {
	r.Lock()
	defer r.Unlock()
	return int(r.running.Get())
}
