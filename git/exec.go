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
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Run(cmd *exec.Cmd) (err error) {
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	fmt.Printf("$ %s\n", strings.Join(cmd.Args, " "))
	if err = cmd.Run(); err != nil {
		fmt.Print(errOut.String())
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

func OutputLines(cmd *exec.Cmd) []string {
	b, err := OutputBytes(cmd)
	if err != nil {
		return []string{}
	}

	s := bufio.NewScanner(bytes.NewReader(b))
	s.Split(bufio.ScanLines)
	var output []string
	for s.Scan() {
		output = append(output, s.Text())
	}
	return output
}
