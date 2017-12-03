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

type funcGroup []func() error
type cmdGroup []Runner

func FuncGroup(funcs ...func() error) Runner {
	return funcGroup(funcs)
}

func Group(cmds ...Runner) Runner {
	return cmdGroup(cmds)
}

func (g cmdGroup) Run() error {
	for _, cmd := range g {
		if e := Run(cmd); e != nil {
			return e
		}
	}
	return nil
}

func (g funcGroup) Run() error {
	for _, f := range g {
		if e := Run(ToRunner(f)); e != nil {
			return e
		}
	}
	return nil
}
