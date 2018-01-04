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

package play

import (
	"github.com/evovetech/got/got/merge"
	"github.com/evovetech/got/log"
	"github.com/spf13/cobra"
)

type mergeArgs struct {
	merge.Args
}

var args mergeArgs
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge",
	RunE: func(cmd *cobra.Command, a []string) error {
		if err := args.Parse(a); err != nil {
			return err
		}
		return args.run()
	},
}

func MergeCmd() *cobra.Command {
	return mergeCmd
}

func init() {
	args.Init(mergeCmd)
}

func (m *mergeArgs) run() error {
	log.Printf("args: %s", m.Args)
	return nil
}
