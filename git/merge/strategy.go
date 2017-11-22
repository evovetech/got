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

package merge

import (
	"fmt"

	"encoding/json"
	"github.com/spf13/pflag"
	"strings"
)

type Strategy int

const (
	NONE Strategy = iota
	THEIRS
	OURS
)

func GetStrategy(str string) Strategy {
	switch strings.ToLower(str) {
	case "theirs":
		return THEIRS
	case "ours":
		return OURS
	}
	return NONE
}

func (s *Strategy) String() string {
	switch *s {
	case THEIRS:
		return "theirs"
	case OURS:
		return "ours"
	}
	return ""
}

func (s *Strategy) Option() string {
	st := s.String()
	if st == "" {
		return st
	}
	return fmt.Sprintf("--%s", st)
}

func (s *Strategy) Set(val string) error {
	if st := GetStrategy(val); st != NONE {
		*s = st
		return nil
	}
	return fmt.Errorf("error parsing strategy: '%s'\n", val)
}

func (s *Strategy) Type() string {
	return "string"
}

func (s *Strategy) AddTo(f *pflag.FlagSet) *pflag.Flag {
	// set default
	*s = THEIRS
	return f.VarPF(s, "strategy", "s", "strategy")
}

// Json
func (s Strategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
