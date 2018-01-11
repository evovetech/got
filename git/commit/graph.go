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

package commit

import (
	"github.com/evovetech/got/types"
	"github.com/evovetech/got/util"
	"gopkg.in/gyuho/goraph.v2"
)

type Graph struct {
	goraph.Graph
}

func NewGraph() *Graph {
	return &Graph{goraph.NewGraph()}
}

func (g *Graph) Get(commit types.Sha) Node {
	if n := g.GetNode(commit); n != nil {
		return n.(Node)
	}
	return nil
}

func (g *Graph) GetOrAdd(commit types.Sha) (n Node, created bool) {
	if n = g.Get(commit); n == nil {
		n, created = g.Add(commit)
	}
	return
}

func (g *Graph) Add(commit types.Sha) (Node, bool) {
	if n, ok := NewNode(commit); ok {
		return n, n.AddTo(g)
	}
	return nil, false
}

func (g *Graph) Populate(commit string, num int) (error, bool) {
	if num <= 0 {
		return nil, false
	}

	var errors []error
	if n, _ := g.GetOrAdd(types.Sha(commit)); n != nil {
		if err := n.Populate(g); err != nil {
			errors = append(errors, err)
		} else {
			for l := n.info().Parents(); l != nil; {
				p := l.value
				if err, _ := g.Populate(p.String(), num-1); err != nil {
					errors = append(errors, err)
				}
			}
		}
	}
	return util.CompositeError(errors), true
}
