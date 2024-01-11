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

// ParseHandler defines the interface for parsing robots.txt.
type ParseHandler interface {
	HandleStart()
	HandleEnd()
	HandleUserAgent(lineNum int, value string)
	HandleAllow(lineNum int, value string)
	HandleDisallow(lineNum int, value string)
	HandleSitemap(lineNum int, value string)
	HandleUnknownAction(lineNum int, key, value string)
}

// DefaultParseHandler is the default implementation of ParseHandler.
type DefaultParseHandler struct {
	currentRuleGroup *rulesGroup
	ruleGroup        []*rulesGroup
	seenSeparator    bool
}

func (h *DefaultParseHandler) HandleStart() {
	h.currentRuleGroup = new(rulesGroup)
	h.ruleGroup = nil
	h.seenSeparator = false
}

func (h *DefaultParseHandler) HandleEnd() {
	h.ruleGroup = append(h.ruleGroup, h.currentRuleGroup)
	h.currentRuleGroup = nil
	h.seenSeparator = false
}

func (h *DefaultParseHandler) HandleUserAgent(lineNum int, value string) {
	if h.seenSeparator {
		h.ruleGroup = append(h.ruleGroup, h.currentRuleGroup)
		h.currentRuleGroup = new(rulesGroup)
		h.seenSeparator = false
	}
	h.currentRuleGroup.AddUserAgent(value)
}

func (h *DefaultParseHandler) HandleAllow(lineNum int, value string) {
	if !h.seenAnyAgent() {
		return
	}
	h.seenSeparator = true
	h.currentRuleGroup.AddRule(keyTypeAllow, value, lineNum)
}

func (h *DefaultParseHandler) HandleDisallow(lineNum int, value string) {
	if !h.seenAnyAgent() {
		return
	}
	h.seenSeparator = true
	h.currentRuleGroup.AddRule(keyTypeDisallow, value, lineNum)
}

func (h *DefaultParseHandler) HandleSitemap(lineNum int, value string) {
	if !h.seenAnyAgent() {
		return
	}
	h.seenSeparator = true
	h.currentRuleGroup.AddRule(keyTypeSiteMap, value, lineNum)
}

func (h *DefaultParseHandler) HandleUnknownAction(lineNum int, key, value string) {
	// Ignore unknown actions.
}

func (h *DefaultParseHandler) seenAnyAgent() bool {
	return h.currentRuleGroup.GlobalAgent || len(h.currentRuleGroup.SpecificAgent) > 0
}
