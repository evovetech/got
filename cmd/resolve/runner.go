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

package resolve

import (
	"fmt"
	"regexp"

	"github.com/evovetech/got/git"
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/log"
)

var reDD = regexp.MustCompile("^(DD)")
var reDeletedOurs = regexp.MustCompile("^(D|UA)")
var reDeletedTheirs = regexp.MustCompile("^(.D|AU)")

func Run(st merge.Strategy) error {
	//git diff --name-only --diff-filter=UXB
	var errors []error
	diff := git.Command("diff", "--name-only", "--diff-filter=UXB")
	for _, file := range diff.OutputLines() {
		if err := resolveFile(file, st); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		var errString string
		for _, err := range errors {
			errString += fmt.Sprintln(err.Error())
		}
		return fmt.Errorf("%s", errString)
	}
	return nil
}

func resolveFile(file string, st merge.Strategy) error {
	var err error
	status, err := git.Status(file).Output()
	if err != nil {
		return err
	}
	log.Print(status)
	switch {
	case reDD.MatchString(status):
		err = git.Add(file).Run()
	case st == merge.OURS:
		err = resolveOurs(file, status)
	case st == merge.THEIRS:
		err = resolveTheirs(file, status)
	default:
		err = fmt.Errorf("unknown strategy: ")
	}
	return err
}

func resolveOurs(file string, status string) error {
	switch {
	case reDeletedOurs.MatchString(status):
		return git.ResolveRm(file).Run()
	default:
		return git.ResolveCheckout(file, merge.OURS).Run()
	}
}

func resolveTheirs(file string, status string) error {
	switch {
	case reDeletedTheirs.MatchString(status):
		return git.ResolveRm(file).Run()
	default:
		return git.ResolveCheckout(file, merge.THEIRS).Run()
	}
}
