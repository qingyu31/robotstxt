//Copyright 2024 qingyu31
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package robotstxt

import (
	"strings"
	"unicode"
)

type rulesGroup struct {
	GlobalAgent   bool
	SpecificAgent []string
	Rules         []rule
}

func (r *rulesGroup) AddUserAgent(agent string) {
	if len(agent) >= 1 && agent[0] == '*' && (len(agent) == 1 || unicode.IsSpace(rune(agent[1]))) {
		r.GlobalAgent = true
		return
	}
	idx := strings.IndexFunc(agent, func(char rune) bool {
		return !unicode.IsLetter(char) || char == '-' || char == '_' || char == '.'
	})
	if idx > 0 {
		agent = agent[:idx]
	}
	r.SpecificAgent = append(r.SpecificAgent, agent)
}

func (r *rulesGroup) AddRule(key keyType, value string, line int) {
	r.Rules = append(r.Rules, rule{key, value, line})
}

type rule struct {
	key   keyType
	value string
	line  int
}
