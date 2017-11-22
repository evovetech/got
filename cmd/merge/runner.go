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

package merge

import (
	"fmt"

	"github.com/evovetech/got/git"
)

type Runner int8

var (
	runner Runner
)

func (r *Runner) Run(args Args) error {
	fmt.Printf("args: %s\n", args)

	var err error
	if err = CheckStatus(); err != nil {
		return err
	}

	var headRef, mergeRef git.Ref
	if mergeRef, err = git.ParseRef(args.Branch); err != nil {
		return err
	}
	if headRef, err = git.ParseRef("HEAD"); err != nil {
		return err
	}
	merger := &Merger{headRef, mergeRef, args.Strategy}

	fmt.Printf("merger: %s\n", merger)

	return merger.RunE()
}

func CheckStatus() error {
	// check git status/diff on HEAD and bail if there are changes
	status, err := git.Command("status", "-s").Output()
	if err != nil {
		return err
	} else if status != "" {
		return fmt.Errorf("please stash or commit changes before merging")
	}
	return nil
}
