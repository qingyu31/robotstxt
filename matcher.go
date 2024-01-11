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

type Matcher interface {
	AllowedByRobots(userAgent []string, path string) bool
	OneAgentAllowedByRobots(userAgent string, path string) bool
}

type defaultMatcher struct {
	groups        []*rulesGroup
	matchStrategy MatchStrategy
}

func (m defaultMatcher) AllowedByRobots(userAgent []string, path string) bool {
	allowed := newMatchHierarchy()
	disallowed := newMatchHierarchy()
	for _, group := range m.groups {
		isSpecificGroup := false
		for _, agent := range group.SpecificAgent {
			for _, ua := range userAgent {
				if ua == agent {
					isSpecificGroup = true
					break
				}
			}
		}
		if !isSpecificGroup && !group.GlobalAgent {
			continue
		}
		for _, r := range group.Rules {
			if r.key == keyTypeDisallow {
				priority := m.matchStrategy.MatchDisallow(path, r.value)
				if isSpecificGroup {
					disallowed.UpdateSpecific(&match{priority, r.line})
				} else {
					disallowed.UpdateGlobal(&match{priority, r.line})
				}
			} else if r.key == keyTypeAllow {
				priority := m.matchStrategy.MatchAllow(path, r.value)
				if isSpecificGroup {
					allowed.UpdateSpecific(&match{priority, r.line})
				} else {
					allowed.UpdateGlobal(&match{priority, r.line})
				}
			}
		}
	}
	if allowed.specific.priority > 0 || disallowed.specific.priority > 0 {
		return allowed.specific.priority >= disallowed.specific.priority
	}
	if allowed.global.priority > 0 || disallowed.global.priority > 0 {
		return allowed.global.priority >= disallowed.global.priority
	}
	return true
}

func (m defaultMatcher) OneAgentAllowedByRobots(userAgent string, path string) bool {
	return m.AllowedByRobots([]string{userAgent}, path)
}

const kNoMatchPriority = -1

type match struct {
	priority int
	line     int
}

func (m *match) Clear() {
	m.Set(kNoMatchPriority, 0)
}

func (m *match) Set(priority, line int) {
	m.priority = priority
	m.line = line
}

func highPriorityMatch(a, b *match) *match {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority > b.priority {
		return a
	}
	return b
}

type matchHierarchy struct {
	global   *match
	specific *match
}

func newMatchHierarchy() *matchHierarchy {
	return &matchHierarchy{
		global:   &match{priority: kNoMatchPriority},
		specific: &match{priority: kNoMatchPriority},
	}
}

func (h *matchHierarchy) UpdateGlobal(m *match) {
	h.global = highPriorityMatch(h.global, m)
}

func (h *matchHierarchy) UpdateSpecific(m *match) {
	h.specific = highPriorityMatch(h.specific, m)
}

func (h *matchHierarchy) Clear() {
	h.global.Clear()
	h.specific.Clear()
}
