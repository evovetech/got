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

package object

import (
	"fmt"
	"github.com/gyuho/goraph"
)

type Edge interface {
	goraph.Edge
	Parent() Node
	Child() Node

	AddTo(*Graph) error
}

func NewEdge(parent, child Node) Edge {
	return &edge{goraph.NewEdge(parent, child, 0)}
}

type edge struct {
	goraph.Edge
}

func (e *edge) Parent() Node {
	return e.Source().(Node)
}

func (e *edge) Child() Node {
	return e.Target().(Node)
}

func (e *edge) AddTo(g *Graph) error {
	return g.AddEdge(e.Parent().ID(), e.Child().ID(), 0)
}

func (e *edge) String() string {
	return fmt.Sprintf("%s -→ %s\n", e.Source(), e.Target())
}
