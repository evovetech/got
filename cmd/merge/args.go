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
	"github.com/evovetech/got/git/merge"
	"github.com/evovetech/got/util"
	"github.com/spf13/cobra"
)

type Args struct {
	Strategy merge.Strategy
	Branch   string
}

func (args *Args) Init(cmd *cobra.Command) {
	args.Strategy.AddTo(cmd.Flags())
}

func (args *Args) Parse(a []string) error {
	args.Branch = a[0]
	return nil
}

func (args Args) String() string {
	return util.String(args)
}
