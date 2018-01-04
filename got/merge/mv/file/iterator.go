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

package file

import (
	"github.com/emirpasic/gods/containers"
	"github.com/emirpasic/gods/trees/avltree"
)

type Iterator interface {
	containers.ReverseIteratorWithKey
	Path() Path
	Entry() Entry
}

func ItBegin(it Iterator)     { it.Begin() }
func ItEnd(it Iterator)       { it.End() }
func ItNext(it Iterator) bool { return it.Next() }
func ItPrev(it Iterator) bool { return it.Prev() }

type iterator struct {
	containers.ReverseIteratorWithKey
}

func newIterator(tree *avltree.Tree) Iterator {
	return &iterator{tree.Iterator()}
}

func (it *iterator) Path() Path {
	if p, ok := it.Key().(Path); ok {
		return p
	}
	return nil
}

func (it *iterator) Entry() Entry {
	if e, ok := it.Value().(Entry); ok {
		return e
	}
	return nil
}

type emptyIterator uint8

const noEntries = emptyIterator(0)

func (emptyIterator) Prev() bool         { return false }
func (emptyIterator) End()               {}
func (emptyIterator) Last() bool         { return false }
func (emptyIterator) Next() bool         { return false }
func (emptyIterator) Value() interface{} { return nil }
func (emptyIterator) Key() interface{}   { return nil }
func (emptyIterator) Begin()             {}
func (emptyIterator) First() bool        { return false }
func (emptyIterator) Path() Path         { return nil }
func (emptyIterator) Entry() Entry       { return nil }
