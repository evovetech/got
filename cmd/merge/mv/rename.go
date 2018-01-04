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
	"github.com/evovetech/got/cmd/merge/mv/file"
	"github.com/evovetech/got/git"
)

type Rename struct {
	From file.Path
	To   file.Path
}

func renameRun(mv file.MovePath) error {
	rename := Rename{mv.FromPath(), mv.ToPath()}
	return rename.run()
}

func (p *Rename) run() error {
	if err := p.To.Dir().Mkdirs(); err != nil {
		return err
	}
	return git.Command("mv", p.From.OsPath(), p.To.OsPath()).Run()
}
