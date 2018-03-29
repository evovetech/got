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

package merge

import (
	"fmt"
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/util"
	"github.com/spf13/cobra"
)

type Args struct {
	Strategy      git.MergeStrategy
	FollowRenames bool
	Branch        string
}

func (args *Args) Init(cmd *cobra.Command) {
	args.Strategy.AddTo(cmd.Flags())
	cmd.Flags().BoolVar(&args.FollowRenames, "followRenames", false, "follow renames (default false)")
}

func (args *Args) Parse(a []string) error {
	if len(a) != 1 {
		return fmt.Errorf("wront number of args. expecting BRANCH, instead got %s", a)
	}
	args.Branch = a[0]
	return nil
}

func (args Args) String() string {
	return util.String(args)
}
