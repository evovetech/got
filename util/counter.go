/*
 * Copyright 2018 evove.tech
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import "sync/atomic"

type Counter interface {
	Get() uint32
	IncrementAndGet() uint32
	GetAndIncrement() uint32
	DecrementAndGet() uint32
	TryDecrementAndGet() (uint32, bool)
}

func NewCounter() Counter {
	return new(counter)
}

func NewCounterWithInitialValue(initial uint32) Counter {
	c := counter(initial)
	return &c
}

type counter uint32

func (c *counter) ptr() *uint32 {
	return (*uint32)(c)
}

func (c *counter) Get() uint32 {
	return atomic.LoadUint32(c.ptr())
}

func (c *counter) IncrementAndGet() uint32 {
	return c.GetAndIncrement() + 1
}

func (c *counter) GetAndIncrement() uint32 {
	if cur := c.Get(); c.CompareAndSwap(cur, cur+1) {
		return cur
	}
	return c.GetAndIncrement()
}

func (c *counter) DecrementAndGet() uint32 {
	if v, ok := c.TryDecrementAndGet(); ok {
		return v
	}
	return c.DecrementAndGet()
}

func (c *counter) TryDecrementAndGet() (uint32, bool) {
	if cur := c.Get(); cur > 0 {
		next := cur + ^uint32(0)
		return next, c.CompareAndSwap(cur, next)
	}
	return 0, false
}

func (c *counter) CompareAndSwap(cur uint32, next uint32) bool {
	return atomic.CompareAndSwapUint32(c.ptr(), cur, next)
}
