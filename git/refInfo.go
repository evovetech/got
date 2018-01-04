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

package git

import (
	"encoding/json"
	"fmt"

	"github.com/evovetech/got/util"
)

type RefInfo struct {
	Type        string
	Author      string
	AuthorEmail string
	FullName    string
	ShortName   string
	RefName     string
	Push        string
	Upstream    string
}

var (
	refTemplate = RefInfo{
		Type:        "%(objecttype)",
		Author:      "%(authorname)",
		AuthorEmail: "%(authoremail)",
		FullName:    "%(objectname)",
		ShortName:   "%(objectname:short)",
		RefName:     "%(refname)",
		Push:        "%(push)",
		Upstream:    "%(upstream)",
	}
	refFormat = refTemplate.Json()
)

func (r RefInfo) Json() string {
	return util.Json(r)
}

func (r RefInfo) String() string {
	return util.String(r)
}

func GetRefInfo(ref string) []RefInfo {
	cmd := Command("for-each-ref")
	cmd.AddOption("--points-at", ref)
	return outputRefs(cmd)
}

func GetBranchInfo(branch string) []RefInfo {
	cmd := Command("branch", "--list", "--all")
	cmd.AddOption("--points-at", branch)
	return outputRefs(cmd)
}

func outputRefs(cmd *Cmd) (refs []RefInfo) {
	cmd.AddOption("--format", refFormat)
	for _, str := range cmd.OutputLines() {
		var ref RefInfo
		if err := json.Unmarshal([]byte(str), &ref); err == nil {
			refs = append(refs, ref)
		} else {
			fmt.Print(err)
		}
	}
	return
}
