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
	"strings"

	"github.com/evovetech/got/util"
)

type RefName struct {
	Symbolic string
	Full     string
	Short    string
	Remote   string
}

func (r RefName) String() string {
	return util.String(r)
}

func (r RefName) Checkout() error {
	return Checkout(r.Short)
}

func ParseRefName(fullName string) RefName {
	var symbolic, shortName, remote string
	if success, name := util.Omit(fullName, "refs/heads/"); success {
		symbolic = name
		shortName = name
	} else if success, name := util.Omit(fullName, "refs/remotes/"); success {
		index := strings.Index(name, "/")
		symbolic = name
		shortName = name[index+1:]
		remote = name[0:index]
	}
	return RefName{
		Symbolic: symbolic,
		Full:     fullName,
		Short:    shortName,
		Remote:   remote,
	}
}
