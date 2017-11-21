/*
 * Copyright 2017 evove.tech
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package git

type RevParse string

func (c RevParse) OutputString(option string) string {
	cmd := Command("rev-parse", option, string(c))
	return cmd.OutputString()
}

func (c RevParse) Short() string {
	return c.OutputString("--short")
}

func (c RevParse) Symbolic() string {
	return c.OutputString("--symbolic")
}

func (c RevParse) SymbolicFullName() string {
	return c.OutputString("--symbolic-full-name")
}
