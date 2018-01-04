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

package cmd

import (
	"github.com/evovetech/got/got/resolve"
	"github.com/spf13/cobra"
)

var r = struct {
	args resolve.Args
	cmd  *cobra.Command
}{
	cmd: &cobra.Command{
		Use:   "resolve",
		Short: "Resolve conflicts",
	},
}

func resolveCmd() *cobra.Command {
	return r.cmd
}

func init() {
	args, cmd := r.args, r.cmd
	cmd.RunE = func(cmd *cobra.Command, a []string) error {
		if err := args.Parse(a); err != nil {
			return err
		}
		return resolve.Run(args.Strategy)
	}
}
