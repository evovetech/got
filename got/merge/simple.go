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

package merge

import (
	"github.com/evovetech/got/git"
)

type simple Merger

func (s *simple) Run() (RunStep, error) {
	if err := s.HeadRef.Checkout(); err != nil {
		return nil, err
	}

	m := s.MergeRef.Merge()
	m.IgnoreAllSpace()
	if err := m.Run(); err != nil {
		if err = git.AbortMerge(); err != nil {
			return nil, err
		}
		return (*multi)(s).Run()
	}

	return nil, nil
}
