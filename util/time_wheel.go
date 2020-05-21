// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"errors"
	"time"
)

// Task means handle unit in time wheel
type Task struct {
	delay    time.Duration
	key      interface{}
	round    int // optimize time wheel to handle delay  bigger than bucketsNum * tick
	callback func()
}

// TimeWheel means time wheel
type TimeWheel struct {
	tick   time.Duration
	ticker *time.Ticker

	bucketsNum    int
	buckets       []map[interface{}]*Task // key: added item, value: *Task
	bucketIndexes map[interface{}]int     // key: added item, value: bucket position

	currentIndex int

	addC    chan *Task
	removeC chan interface{}
	stopC   chan struct{}
}

// NewTimeWheel create new time wheel
func NewTimeWheel(tick time.Duration, bucketsNum int) (*TimeWheel, error) {
	if bucketsNum <= 0 {
		return nil, errors.New("bucket number must be greater than 0")
	}
	if int(tick.Seconds()) < 1 {
		return nil, errors.New("tick cannot be less than 1s")
	}

	tw := &TimeWheel{
		tick:          tick,
		bucketsNum:    bucketsNum,
		bucketIndexes: make(map[interface{}]int, 1024),
		buckets:       make([]map[interface{}]*Task, bucketsNum),
		currentIndex:  0,
		addC:          make(chan *Task, 1024),
		removeC:       make(chan interface{}, 1024),
		stopC:         make(chan struct{}),
	}

	for i := 0; i < bucketsNum; i++ {
		tw.buckets[i] = make(map[interface{}]*Task, 16)
	}

	return tw, nil
}

// Start start the time wheel
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.tick)
	go tw.start()
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.handleTick()
		case task := <-tw.addC:
			tw.add(task)
		case key := <-tw.removeC:
			tw.remove(key)
		case <-tw.stopC:
			tw.ticker.Stop()
			return
		}
	}
}

// Stop stop the time wheel
func (tw *TimeWheel) Stop() {
	tw.stopC <- struct{}{}
}

func (tw *TimeWheel) handleTick() {
	bucket := tw.buckets[tw.currentIndex]
	for k := range bucket {
		if bucket[k].round > 0 {
			bucket[k].round--
			continue
		}
		go bucket[k].callback()
		delete(bucket, k)
		delete(tw.bucketIndexes, k)
	}
	if tw.currentIndex == tw.bucketsNum-1 {
		tw.currentIndex = 0
		return
	}
	tw.currentIndex++
}

// Add add an item into time wheel
func (tw *TimeWheel) Add(delay time.Duration, key interface{}, callback func()) error {
	if delay <= 0 || key == nil {
		return errors.New("invalid params")
	}
	tw.addC <- &Task{delay: delay, key: key, callback: callback}
	return nil
}

func (tw *TimeWheel) add(task *Task) {
	round := tw.calculateRound(task.delay)
	index := tw.calculateIndex(task.delay)
	task.round = round
	if originIndex, ok := tw.bucketIndexes[task.key]; ok {
		delete(tw.buckets[originIndex], task.key)
	}
	tw.bucketIndexes[task.key] = index
	tw.buckets[index][task.key] = task
}

func (tw *TimeWheel) calculateRound(delay time.Duration) (round int) {
	delaySeconds := int(delay.Seconds())
	tickSeconds := int(tw.tick.Seconds())
	round = delaySeconds / tickSeconds / tw.bucketsNum
	return
}

func (tw *TimeWheel) calculateIndex(delay time.Duration) (index int) {
	delaySeconds := int(delay.Seconds())
	tickSeconds := int(tw.tick.Seconds())
	index = (tw.currentIndex + delaySeconds/tickSeconds) % tw.bucketsNum
	return
}

// Remove remove an item from time wheel
func (tw *TimeWheel) Remove(key interface{}) error {
	if key == nil {
		return errors.New("invalid params")
	}
	tw.removeC <- key
	return nil
}

// don't need to call callback
func (tw *TimeWheel) remove(key interface{}) {
	if index, ok := tw.bucketIndexes[key]; ok {
		delete(tw.bucketIndexes, key)
		delete(tw.buckets[index], key)
	}
	return
}
