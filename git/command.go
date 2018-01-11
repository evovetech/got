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

package git

import (
	"os/exec"
)

type Cmd struct {
	Name    string
	Options []string
	Args    []string
}

func Command(name string, args ...string) *Cmd {
	return &Cmd{Name: name, Args: args}
}

func (c *Cmd) AddOption(name string, args ...string) {
	options := append([]string{name}, args...)
	c.Options = append(c.Options, options...)
}

func (c *Cmd) AddArg(arg string) {
	c.Args = append(c.Args, arg)
}

func (c *Cmd) AddArgs(args ...string) {
	c.Args = append(c.Args, args...)
}

func (c *Cmd) AllArgs() []string {
	args := []string{c.Name}
	if len(c.Options) > 0 {
		args = append(args, c.Options...)
	}
	return append(args, c.Args...)
}

func (c *Cmd) Build() *exec.Cmd {
	return exec.Command("git", c.AllArgs()...)
}

func (c *Cmd) Run() error {
	return Run(c.Build())
}

func (c *Cmd) OutputBytes() ([]byte, error) {
	return OutputBytes(c.Build())
}

func (c *Cmd) Output() (string, error) {
	return Output(c.Build())
}

func (c *Cmd) OutputString() string {
	return OutputString(c.Build())
}

func (c *Cmd) OutputLines() []string {
	return OutputLines(c.Build())
}

func (c *Cmd) ForEachLine(f func(string) error) error {
	return ForEachLine(c.Build(), f)
}
