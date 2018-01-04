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

package file

import "github.com/evovetech/got/util"

func nextMatch(from Path, to Path) (f int, t int, ok bool) {
	if len(from) > len(to) {
		f, t, ok = from.nextMatch(to)
	} else {
		t, f, ok = to.nextMatch(from)
	}
	return
}

func (p Path) nextMatch(o Path) (int, int, bool) {
	oLen := len(o)
	if len(p) == 0 || oLen == 0 {
		return -1, -1, false
	}
	for pi := range p {
		oMax := util.MinInt(pi+1, oLen)
		oi := oMax - 1
		for i := 0; i < pi; i++ {
			if p[i] == o[oi] {
				return i, oi, true
			}
		}
		for oi := 0; oi < oMax; oi++ {
			if p[pi] == o[oi] {
				return pi, oi, true
			}
		}
	}
	return -1, -1, false
}
