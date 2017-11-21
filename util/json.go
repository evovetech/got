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

package util

import "encoding/json"

func Marshal(v interface{}, indent string) (str string, err error) {
	if b, err := marshalFunc(indent)(v); err == nil {
		str = string(b)
	}
	return
}

func MarshalQuiet(v interface{}, indent string) string {
	str, _ := Marshal(v, indent)
	return str
}

func Json(v interface{}) string {
	return MarshalQuiet(v, "")
}

func String(v interface{}) string {
	return MarshalQuiet(v, "  ")
}

func marshalFunc(indent string) func(v interface{}) ([]byte, error) {
	if indent == "" {
		return json.Marshal
	}
	return func(v interface{}) ([]byte, error) {
		return json.MarshalIndent(v, "", indent)
	}
}
