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

package merge

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/log"
)

type PreRunStep interface {
	PreRun() error
}
type RunStep interface {
	Run() (RunStep, error)
}
type beginStep struct {
	PreRunStep
	RunStep

	args Args
}

func NewMerger(args Args) (merger *Merger, err error) {
	log.Printf("args: %s\n", args)
	var headRef, mergeRef git.Ref
	if mergeRef, err = git.ParseRef(args.Branch); err != nil {
		return
	}
	if headRef, err = git.ParseRef("HEAD"); err != nil {
		return
	}

	merger = &Merger{headRef, mergeRef, args}

	log.Verbose.Printf("merger: %s\n", merger)
	return
}

func Run(args Args) (err error) {
	next := beginMerge(args)
	for next != nil {
		if pre, ok := next.(PreRunStep); ok {
			if err := pre.PreRun(); err != nil {
				return err
			}
		}
		if next, err = next.Run(); err != nil {
			return err
		}
	}
	return nil
}

func beginMerge(args Args) RunStep {
	return &beginStep{args: args}
}

func (s *beginStep) PreRun() error {
	return git.CheckStatus()
}

func (s *beginStep) Run() (RunStep, error) {
	return NewMerger(s.args)
}
