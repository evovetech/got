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
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/evovetech/got/log"
	"github.com/evovetech/got/options"
	"github.com/evovetech/got/util"
)

func Run(cmd Runner) error {
	if c, ok := cmd.(*Cmd); ok {
		return RunCmd(c)
	} else if c, ok := cmd.(*exec.Cmd); ok {
		return RunExecCmd(c)
	}
	return cmd.Run()
}

func RunCmd(cmd *Cmd) (err error) {
	return RunExecCmd(cmd.Build())
}

func RunExecCmd(cmd *exec.Cmd) (err error) {
	if options.Verbose {
		log.Printf("$ %s\n", strings.Join(cmd.Args, " "))
		log.IndentIn()
		defer log.IndentOut()
	}

	var errOut bytes.Buffer
	cmd.Stderr = util.CompositeWriter(&errOut, log.Verbose)
	cmd.Stdout = util.CompositeWriter(cmd.Stdout, log.Verbose)
	if err = cmd.Run(); err != nil {
		if options.Verbose {
			return
		}
		for s := bufio.NewScanner(&errOut); s.Scan(); {
			log.Err.Write(s.Bytes())
		}
	}
	return
}

func OutputBytes(cmd *exec.Cmd) (b []byte, err error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = Run(cmd); err == nil {
		b = out.Bytes()
	}
	return
}

func Output(cmd *exec.Cmd) (str string, err error) {
	var b []byte
	if b, err = OutputBytes(cmd); err == nil {
		str = strings.TrimSpace(string(b))
	}
	return
}

func OutputString(cmd *exec.Cmd) string {
	str, _ := Output(cmd)
	return str
}

func OutputLines(cmd *exec.Cmd) (output []string) {
	ForEachLine(cmd, func(line string) error {
		output = append(output, line)
		return nil
	})
	return
}

func ForEachLine(cmd *exec.Cmd, f func(string) error) error {
	b, err := OutputBytes(cmd)
	if err != nil {
		return err
	}

	var errors []error
	s := bufio.NewScanner(bytes.NewReader(b))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if err := f(s.Text()); err != nil {
			errors = append(errors, err)
		}
	}
	return util.CompositeError(errors)
}
