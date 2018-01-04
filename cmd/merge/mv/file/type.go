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

type Type uint32

func (t Type) String() string {
	if t == Mv {
		return "MV"
	}
	var str string
	if t.HasFlag(Rn) {
		str += "R"
	}
	switch {
	case t.HasFlag(Add):
		str += "A"
	case t.HasFlag(Del):
		str += "D"
	default:
		str += "?"
	}
	return str
}

func (t Type) HasFlag(flag Type) bool {
	return t&flag != 0
}

const (
	Add Type = 1 << iota
	Del
	Rn
	Mv = Add | Del | Rn
)
