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

type Node interface {
	goraph.Node

	WithParent(Node) Edge
	WithChild(Node) Edge
	AddTo(*Graph) bool
	Populate(*Graph) error

	// private
	info() Info
}

type node struct {
	Info
	populated bool
}

func NewNode(commit types.Sha) (n Node, ok bool) {
	if info := NewInfo(commit); info != nil {
		n, ok = &node{Info: *info}, true
	}
	//log.Printf("NewNode(%s) -> (ok=%v, n=%s)", commit, ok, n)
	return
}

func (n *node) ID() goraph.ID {
	return n.sha
}

func (n *node) WithParent(parent Node) Edge {
	return NewEdge(parent, n)
}

func (n *node) WithChild(child Node) Edge {
	return NewEdge(n, child)
}

func (n *node) AddTo(g *Graph) (ok bool) {
	ok = g.AddNode(n)
	//log.Printf("%s.AddTo(graph) -> %v", n, ok)
	return
}

func (n *node) Populate(g *Graph) error {
	if n.populated {
		return nil
	}

	n.populated = true
	var errors []error
	for _, p := range n.parents {
		if pn, _ := g.GetOrAdd(p); pn != nil {
			if err := pn.WithChild(n).AddTo(g); err != nil {
				n.populated = false
				errors = append(errors, err)
			}
		}
	}
	return util.CompositeError(errors)
}

func (n *node) info() Info {
	return n.Info
}

func (n node) String() string {
	return n.sha.String()
}
