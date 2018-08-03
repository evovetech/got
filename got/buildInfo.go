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

package got

import (
	"sync"
	"fmt"
)

var build buildInfo

func Init(
	version string,
	commit string,
	date string,
) BuildInfo {
	return build.init(version, commit, date)
}

func GetBuildInfo() BuildInfo {
	return build
}

type BuildInfo interface {
	Version() string
	Commit() string
	Date() string
	String() string
}

type buildInfo struct {
	version string
	commit  string
	date    string

	once sync.Once
}

func (b *buildInfo) init(
	version string,
	commit string,
	date string,
) BuildInfo {
	b.once.Do(func() {
		b.version = version
		b.commit = commit
		b.date = date
	})
	return b
}

func (b buildInfo) Version() string {
	return b.version
}

func (b buildInfo) Commit() string {
	return b.commit
}

func (b buildInfo) Date() string {
	return b.date
}

func (b buildInfo) String() string {
	return fmt.Sprintf("BuildInfo{\n"+
		"  version: %v\n"+
		"  commit:  %v\n"+
		"  date:    %v\n"+
		"}",
		b.version,
		b.commit,
		b.date,
	)
}
