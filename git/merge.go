/*
 * Copyright 2017 evove.tech
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package git

import (
	"os/exec"

	"github.com/evovetech/got/git/merge"
)

type MergeCmd struct {
	cmd *Cmd

	Strategy        merge.Strategy
	StrategyOptions []string
	MergeRef        string
}

func Merge() *MergeCmd {
	return &MergeCmd{
		cmd: &Cmd{Name: "merge"},
	}
}

func (m *MergeCmd) Abort() error {
	m.cmd.AddOption("--abort")
	return m.cmd.Run()
}

func (m *MergeCmd) Commit() {
	m.cmd.AddOption("--commit")
}

func (m *MergeCmd) NoCommit() {
	m.cmd.AddOption("--no-commit")
}

func (m *MergeCmd) AddStrategyOption(option string) {
	m.StrategyOptions = append(m.StrategyOptions, option)
}

func (m *MergeCmd) IgnoreAllSpace() {
	m.AddStrategyOption("ignore-all-space")
}

func (m *MergeCmd) FFOnly() {
	m.cmd.AddOption("--ff-only")
}

func (m *MergeCmd) Build() *exec.Cmd {
	m.cmd.AddOption("-s", "recursive")
	if st := m.Strategy; st != merge.NONE {
		m.cmd.AddOption("-X", st.String())
	}
	for _, o := range m.StrategyOptions {
		m.cmd.AddOption("-X", o)
	}
	m.cmd.AddOption("-X", "diff-algorithm=histogram")
	m.cmd.AddOption("-X", "find-renames=75%")
	m.cmd.AddOption("-X", "ignore-all-space")
	m.cmd.AddArg(m.MergeRef)
	return m.cmd.Build()
}

func (m *MergeCmd) Run() error {
	return Run(m.Build())
}
