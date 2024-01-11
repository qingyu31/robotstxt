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

// MatchStrategy defines the interface for matching Allow/Disallow patterns.
type MatchStrategy interface {
	MatchAllow(path, pattern string) int
	MatchDisallow(path, pattern string) int
}

// DefaultMatchStrategy is the default match strategy for Allow/Disallow patterns
// which is the official way of Google crawler to match robots.txt.
type DefaultMatchStrategy struct{}

func (s *DefaultMatchStrategy) MatchAllow(path, pattern string) int {
	if matches(path, pattern) {
		return len(pattern)
	}
	return -1
}

func (s *DefaultMatchStrategy) MatchDisallow(path, pattern string) int {
	if matches(path, pattern) {
		return len(pattern)
	}
	return -1
}

func matches(path, pattern string) bool {
	pathLen := len(path)
	pos := make([]int, pathLen+1)
	numpos := 1

	pos[0] = 0

	for _, pat := range pattern {
		if pat == '$' && pattern[len(pattern)-1] == '$' {
			return pos[numpos-1] == pathLen
		}
		if pat == '*' {
			numpos = pathLen - pos[0] + 1
			for i := 1; i < numpos; i++ {
				pos[i] = pos[i-1] + 1
			}
		} else {
			newnumpos := 0
			for i := 0; i < numpos; i++ {
				if pos[i] < pathLen && path[pos[i]] == byte(pat) {
					pos[newnumpos] = pos[i] + 1
					newnumpos++
				}
			}
			numpos = newnumpos
			if numpos == 0 {
				return false
			}
		}
	}

	return true
}
