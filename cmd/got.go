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

package cmd

import (
	"github.com/evovetech/got/cmd/merge"
	"github.com/evovetech/got/cmd/play"
	"github.com/evovetech/got/cmd/resolve"
	"github.com/evovetech/got/cmd/version"
	"github.com/evovetech/got/options"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "got",
}

func init() {
	Cmd.AddCommand(
		merge.Cmd,
		resolve.Cmd,
		version.Cmd,
		play.Cmd,
	)
	options.AddTo(Cmd)
}
