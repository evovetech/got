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
	"github.com/evovetech/got/got/play"
	"github.com/spf13/cobra"
)

var pl = struct {
	cmd     *cobra.Command
	merge   *cobra.Command
	counter *cobra.Command
}{
	cmd: &cobra.Command{
		Use:   "play",
		Short: "Play",
	},
	merge:   play.MergeCmd(),
	counter: play.CounterCmd(),
}

func playCmd() *cobra.Command {
	return pl.cmd
}

func init() {
	pl.cmd.AddCommand(pl.merge, pl.counter)
}
