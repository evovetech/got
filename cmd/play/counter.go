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
	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
	"github.com/spf13/cobra"
)

var counterCmd = &cobra.Command{
	Use:   "counter",
	Short: "counter",
	RunE: func(cmd *cobra.Command, a []string) error {
		counter := util.NewCounter()
		size, num := 5, 10
		res := make([]*Result, size)
		for i := 0; i < size; i++ {
			r := newRes(i+1, counter)
			res[i] = r
			r.run(num)
		}
		await(res)
		log.Printf("Final Count: %d", counter.Get())
		return nil
	},
}
