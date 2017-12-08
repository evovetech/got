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
	"github.com/evovetech/got/util"
)

var reDD = regexp.MustCompile("^(DD)")
var reDeletedOurs = regexp.MustCompile("^(D|UA)")
var reDeletedTheirs = regexp.MustCompile("^(.D|AU)")

func Run(st merge.Strategy) error {
	//git.Command("update-index", "--really-refresh", "--again", "--verbose").Run()
	var errors []error
	diff := git.Command("diff", "--name-only", "--diff-filter=UXB")
	for _, file := range diff.OutputLines() {
		if err := resolveFile(file, st); err != nil {
			errors = append(errors, err)
		}
	}
	err := util.CompositeError(errors)
	if err == nil {
		err = git.RemoveUntracked()
	}
	return err
}

func resolveFile(file string, st merge.Strategy) error {
	var err error
	status, err := git.StatusCmd(file).Output()
	if err != nil {
		return err
	}
	switch {
	case reDD.MatchString(status):
		err = git.Add(file, "-u")
	case st == merge.OURS:
		err = resolveOurs(file, status)
	case st == merge.THEIRS:
		err = resolveTheirs(file, status)
	default:
		err = fmt.Errorf("unknown strategy: ")
	}
	status2 := git.StatusCmd(file).OutputString()
	if status2 == "" {
		status2 = "  "
	}
	log.Printf("'%s' -> '%s' %s", status[:2], status2[:2], file)
	return err
}

func resolveOurs(file string, status string) error {
	switch {
	case reDeletedOurs.MatchString(status):
		return git.ResolveRmCmd(file).Run()
	default:
		return git.ResolveCheckoutCmd(file, merge.OURS).Run()
	}
}

func resolveTheirs(file string, status string) error {
	switch {
	case reDeletedTheirs.MatchString(status):
		return git.ResolveRmCmd(file).Run()
	default:
		return git.ResolveCheckoutCmd(file, merge.THEIRS).Run()
	}
}