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

package log

import (
	"log"
	"os"
)

var (
	indent  = NewIndent(2)
	Verbose = New(new(verbose), "", 0, indent)
	Std     = New(os.Stdout, "", 0, indent)
	Err     = New(os.Stderr, "", log.Llongfile, indent)
	Ignore  = New(devNull, "", 0, indent)
)

func Print(v interface{}) {
	Std.Print(v)
}

func Println(v interface{}) {
	Std.Println(v)
}

func Printf(format string, v ...interface{}) {
	Std.Printf(format, v...)
}

func IndentIn() {
	indent.In()
}

func IndentOut() {
	indent.Out()
}
