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

type TypeCount map[Type]int

func (c TypeCount) add(t Type, i int) {
	c[t] = c[t] + i
}

func (c TypeCount) addAll(o TypeCount) {
	for t, n := range o {
		c.add(t, n)
	}
}

func (c TypeCount) All(t Type) (n int) {
	for k, v := range c {
		if k.HasFlag(t) {
			n += v
		}
	}
	return
}

func (c TypeCount) AllAdd() int {
	return c.All(Add)
}

func (c TypeCount) AllDel() int {
	return c.All(Del)
}
