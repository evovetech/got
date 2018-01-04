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

package mv

import (
	"github.com/evovetech/got/git"
	"github.com/evovetech/got/got/merge/mv/file"
	"github.com/evovetech/got/log"
)

type FileMoves struct {
	Renames []file.MovePath
	//errs    []*Group
}

func GetFileMoves() (*FileMoves, bool) {
	//errs, renames := NewMap().Run()
	renames := NewMap().Run()
	if len(renames) == 0 {
		//if len(errs) != 0 {
		//	log.Verbose.Printf("errors: %s", util.String(errs))
		//}
		return nil, true
	}
	return &FileMoves{
		Renames: renames,
		//errs:    errs,
	}, false
}

func (m *FileMoves) Run() error {
	// abort merge, move files
	git.AbortMerge()
	for _, mv := range m.Renames {
		if err := renameRun(mv); err != nil {
			log.Err.Printf("error: %s", err.Error())
		}
	}
	//if len(m.errs) > 0 {
	//	log.Verbose.Printf("errors: %s", util.String(m.errs))
	//}
	return git.Command("commit", "-m", "moving files to prepare for merge").Run()
}
