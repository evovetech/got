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

import "os/exec"

type CmdGroup struct {
	cmds []*exec.Cmd
}

func Group(cmds ...*exec.Cmd) *CmdGroup {
	return &CmdGroup{cmds}
}

func (g *CmdGroup) Run() error {
	for _, cmd := range g.cmds {
		if e := Run(cmd); e != nil {
			return e
		}
	}
	return nil
}
