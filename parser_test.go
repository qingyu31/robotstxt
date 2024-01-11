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
	"testing"
)

var parseTests = []struct {
	txt      string
	agent    string
	path     string
	expected bool
}{
	{"", "", "", true},
	{"", "FooBot", "", true},
	{"user-agent: FooBot\ndisallow: /\n", "", "", true},
	{"user-agent: FooBot\ndisallow: /\n", "FooBot", "/x/y", false},
	{"foo: FooBot\nbar: /\n", "FooBot", "/x/y", true},
	{"user-agent FooBot\ndisallow /\n", "FooBot", "/x/y", false},
	{multiGroupTxt, "FooBot", "/x/b", true},
	{multiGroupTxt, "FooBot", "/z/d", true},
	{multiGroupTxt, "FooBot", "/y/c", false},
	{multiGroupTxt, "BarBot", "/y/c", true},
	{multiGroupTxt, "BarBot", "/w/a", true},
	{multiGroupTxt, "BarBot", "/z/d", false},
	{multiGroupTxt, "BazBot", "/z/d", true},
	{multiGroupTxt, "FooBot", "/foo/bar", false},
	{multiGroupTxt, "BarBot", "/foo/bar", false},
	{multiGroupTxt, "BazBot", "/foo/bar", false},
}

func TestParseAndMatch(t *testing.T) {
	for _, tt := range parseTests {
		matcher := Parse(strings.NewReader(tt.txt))
		if matcher.OneAgentAllowedByRobots(tt.agent, tt.path) != tt.expected {
			t.Fatalf("Parse(%q) .Agent(%s) .Allow(%s) = %v, want %v", tt.txt, tt.agent, tt.path, !tt.expected, tt.expected)
		}
	}
}

const multiGroupTxt = `allow: /foo/bar/

user-agent: FooBot
disallow: /
allow: /x/
user-agent: BarBot
disallow: /
allow: /y/


allow: /w/
user-agent: BazBot

user-agent: FooBot
allow: /z/
disallow: /
`

func TestParseMultiGroup(t *testing.T) {
	reader := strings.NewReader(multiGroupTxt)
	groups := Parse(reader).(*defaultMatcher).groups
	if len(groups) != 3 {
		t.Fatalf("Parse() = %d groups, want 3", len(groups))
	}
	expected := []*rulesGroup{
		{
			SpecificAgent: []string{"FooBot"},
			Rules: []rule{
				{key: keyTypeDisallow, value: "/", line: 4},
				{key: keyTypeAllow, value: "/x/", line: 5},
			},
		},
		{
			SpecificAgent: []string{"BarBot"},
			Rules: []rule{
				{key: keyTypeDisallow, value: "/", line: 7},
				{key: keyTypeAllow, value: "/y/", line: 8},
				{key: keyTypeAllow, value: "/w/", line: 11},
			},
		},
		{
			SpecificAgent: []string{"BazBot", "FooBot"},
			Rules: []rule{
				{key: keyTypeAllow, value: "/z/", line: 15},
				{key: keyTypeDisallow, value: "/", line: 16},
			},
		},
	}
	for i, group := range groups {
		if len(group.SpecificAgent) != len(expected[i].SpecificAgent) {
			t.Fatalf("Parse() = %d agents, want %d", len(group.SpecificAgent), len(expected[i].SpecificAgent))
		}
		for j, agent := range group.SpecificAgent {
			if agent != expected[i].SpecificAgent[j] {
				t.Fatalf("Parse() = %v, want %v", agent, expected[i].SpecificAgent[j])
			}
		}
		if len(group.Rules) != len(expected[i].Rules) {
			t.Fatalf("Parse() = %d rules, want %d", len(group.Rules), len(expected[i].Rules))
		}
		for j, r := range group.Rules {
			if r.key != expected[i].Rules[j].key || r.value != expected[i].Rules[j].value || r.line != expected[i].Rules[j].line {
				t.Fatalf("Parse() = %v, want %v", r, expected[i].Rules[j])
			}
		}
	}
}
