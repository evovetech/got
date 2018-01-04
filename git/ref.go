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
	"fmt"

	"github.com/evovetech/got/log"
	"github.com/evovetech/got/util"
)

type Ref struct {
	Name   RefName
	Commit CommitSha
	Info   RefInfo
}

func (r Ref) IsEmpty() bool {
	return r.Commit.Full == ""
}

func (r Ref) SymbolicName() string {
	return r.Name.Symbolic
}

func (r Ref) ShortName() string {
	return r.Name.Short
}

func (r Ref) ShortSha() string {
	return r.Commit.Short
}

func (r Ref) TreeRef() string {
	return r.ShortSha() + "^{tree}"
}

func (r Ref) String() string {
	return util.String(r)
}

func (r Ref) Checkout() error {
	return r.Name.Checkout()
}

func (r Ref) Reset() Reset {
	return Reset(r.ShortSha())
}

func (r Ref) Merge() *MergeCmd {
	cmd := Merge()
	cmd.MergeRef = r.ShortSha()
	return cmd
}

func (r Ref) Delete() error {
	return Command("branch", "-D", r.ShortName()).Run()
}

func (r *Ref) Update() (err error) {
	var up Ref
	if up, err = ParseRef(r.Name.Full); err == nil {
		*r = up
	}
	return
}

func ParseRef(ref string) (r Ref, err error) {
	parse := RevParse(ref)
	refName := parse.SymbolicFullName()
	for _, refInfo := range GetRefInfo(ref) {
		if refInfo.RefName == refName {
			r = Ref{
				Name:   ParseRefName(refName),
				Commit: CommitSha{refInfo.FullName, refInfo.ShortName},
				Info:   refInfo,
			}
			log.Verbose.Printf("Ref %s\n", r)
			return
		}
	}
	err = fmt.Errorf("can't find ref info for %s", ref)
	return
}
