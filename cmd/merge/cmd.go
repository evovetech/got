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
	"github.com/spf13/cobra"
)

var strategy merge.Strategy
var Cmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge [BRANCH] into HEAD",
	Args:  cobra.ExactArgs(1),
	RunE:  RunE,
}

func init() {
	strategy.AddTo(Cmd.Flags())
}

func RunE(cmd *cobra.Command, args []string) error {
	m := NewRunner(cmd, strategy, args[0])
	return m.RunE()
}
