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
	"github.com/evovetech/got/git/object"
	"github.com/evovetech/got/got/merge"
	"github.com/evovetech/got/log"
	"github.com/spf13/cobra"
)

type merger struct {
	*merge.Merger
	args merge.Args
}

var m merger
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "merge",
	RunE: func(cmd *cobra.Command, a []string) error {
		if err := m.parse(a); err != nil {
			return err
		}
		return m.run()
	},
}

func MergeCmd() *cobra.Command {
	return mergeCmd
}

func init() {
	m.args.Init(mergeCmd)
}

func (m *merger) parse(a []string) (err error) {
	if err = m.args.Parse(a); err != nil {
		return
	}
	m.Merger, err = merge.NewMerger(m.args)
	return
}

func (m *merger) run() error {
	log.Printf("merge: %s", m)
	if fork, ok := object.FindForkCommit(m.HeadRef, m.MergeRef); ok {
		log.Printf("fork: %s", fork.Id())
		log.Printf("tree: %s", fork.Tree())
	}

	//g := commit.NewGraph()
	//if err, pop := g.Populate(m.HeadRef.Commit.Full, 5); err != nil {
	//	log.Printf("head: pop=%v, err=%s", pop, err.Error())
	//}
	//if err, pop := g.Populate(m.MergeRef.Commit.Full, 5); err != nil {
	//	log.Printf("merge: pop=%v, err=%s", pop, err.Error())
	//}
	//kruskal, _ := goraph.Kruskal(g)
	//map2 := make(map[string]struct{}, len(kruskal))
	//for k, v := range kruskal {
	//	map2[k.String()] = v
	//}
	//if ts, ok := goraph.TopologicalSort(g); ok {
	//	log.Printf("graph: %s", util.String(ts))
	//}
	//if ref := m.HeadRef.GetCommits(5); ref != nil {
	//	log.Printf("head: %s", util.String(ref))
	//}
	//if ref := m.MergeRef.GetCommits(5); ref != nil {
	//	log.Printf("merge: %s", util.String(ref))
	//}
	return nil
}
