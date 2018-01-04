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

package util

import "io"

type compositeWriter struct {
	writers []io.Writer
}

func (w *compositeWriter) Write(b []byte) (n int, err error) {
	main := w.writers[0]
	if n, err = main.Write(b); err == nil {
		for _, w := range w.writers[1:] {
			w.Write(b)
		}
	}
	return
}

func CompositeWriter(writers ...io.Writer) io.Writer {
	var ret []io.Writer
	for _, w := range writers {
		if w != nil {
			ret = append(ret, w)
		}
	}
	switch len(ret) {
	case 0:
		return nil
	case 1:
		return ret[0]
	default:
		return &compositeWriter{ret}
	}
}
