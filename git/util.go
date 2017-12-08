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
	"regexp"

	"fmt"
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/util"
)

var reUU = regexp.MustCompile("^\\?\\?\\s+(.*)")

func SymbolicRef(ref string) string {
	cmd := Command("symbolic-ref", "--short", ref)
	return cmd.OutputString()
}

func StatusCmd(file string) *Cmd {
	return Command("status", "-s", "--", file)
}

func AddCmd(file string, options ...string) *Cmd {
	options = append(options, "--", file)
	return Command("add", options...)
}

func Add(file string, options ...string) error {
	return AddCmd(file, options...).Run()
}

func CheckoutCmd(args ...string) *Cmd {
	return Command("checkout", args...)
}

func Checkout(args ...string) error {
	err := FuncGroup(
		CheckStatus,
		CheckoutCmd(args...).Run,
		Command("reset", "--soft", "HEAD").Run,
		CheckStatus,
	).Run()
	return err
}

func ResolveRmCmd(file string) Runner {
	return Group(
		exec.Command("rm", file),
		AddCmd(file, "-A"),
	)
}

func ResolveCheckoutCmd(file string, s merge.Strategy) Runner {
	return Group(
		CheckoutCmd(s.Option(), "--", file),
		AddCmd(file, "-A"),
	)
}

func AbortMerge() error {
	return FuncGroup(
		Merge().Abort,
		RemoveUntracked,
	).Run()
}

func RemoveUntracked() error {
	var errors []error
	diff := Command("status", "-s", "--untracked-files=all")
	for _, status := range diff.OutputLines() {
		switch {
		case reUU.MatchString(status):
			match := reUU.FindStringSubmatch(status)
			cmd := exec.Command("rm", match[1])
			if err := Run(cmd); err != nil {
				errors = append(errors, err)
			}
		}
	}
	return util.CompositeError(errors)
}

func CheckStatus() error {
	// check git status/diff on HEAD and bail if there are changes
	status, err := Command("status", "-s", "--untracked-files=all").Output()
	if err != nil {
		return err
	} else if status != "" {
		return fmt.Errorf("please stash or commit changes before merging")
	}
	return nil
}
