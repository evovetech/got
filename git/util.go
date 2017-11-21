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

func SymbolicRef(ref string) string {
	cmd := Command("symbolic-ref", "--short", ref)
	return cmd.OutputString()
}

func Status(file string) *Cmd {
	return Command("status", "-s", "--", file)
}

func Add(file string) *Cmd {
	return Command("add", "--", file)
}

func Checkout(args ...string) *Cmd {
	return Command("checkout", args...)
}

func ResolveRm(file string) *CmdGroup {
	return Group(
		exec.Command("rm", file),
		Add(file).Build(),
	)
}

func ResolveCheckout(file string, s merge.Strategy) *CmdGroup {
	return Group(
		Checkout(s.Option(), "--", file).Build(),
		Add(file).Build(),
	)
}
