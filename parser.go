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
	"bufio"
	"io"
	"strings"
	"unicode"
)

// Parse parses robots.txt and returns a Matcher.
func Parse(robotsBody io.Reader, options ...ParseOption) Matcher {
	config := new(ParseConfig)
	config.applyDefault()
	for _, option := range options {
		option.Apply(config)
	}
	handler := config.Handler
	scanner := bufio.NewScanner(robotsBody)
	handler.HandleStart()
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if lineNum == 1 {
			strings.TrimPrefix(line, "\ufeff")
		}
		if line == "" {
			continue
		}
		key, value, ok := getKeyAndValueFrom(line)
		if !ok {
			continue
		}
		switch parseKey(key) {
		case keyTypeUserAgent:
			handler.HandleUserAgent(lineNum, value)
		case keyTypeAllow:
			handler.HandleAllow(lineNum, maybeEscapePattern(value))
		case keyTypeDisallow:
			handler.HandleDisallow(lineNum, maybeEscapePattern(value))
		case keyTypeSiteMap:
			handler.HandleSitemap(lineNum, value)
		default:
			handler.HandleUnknownAction(lineNum, key, maybeEscapePattern(value))
		}
	}
	handler.HandleEnd()
	m := new(defaultMatcher)
	m.groups = handler.(*DefaultParseHandler).ruleGroup
	m.matchStrategy = config.MatchStrategy
	return m
}

type ParseOption interface {
	Apply(*ParseConfig)
}

type ParseConfig struct {
	Handler       ParseHandler
	MatchStrategy MatchStrategy
}

func (c *ParseConfig) applyDefault() {
	if c.Handler == nil {
		c.Handler = new(DefaultParseHandler)
	}
	if c.MatchStrategy == nil {
		c.MatchStrategy = new(DefaultMatchStrategy)
	}
}

var kHexDigits = []byte("0123456789ABCDEF")

// maybeEscapePattern 规范化允许/拒绝的路径
func maybeEscapePattern(src string) string {
	numToEscape := 0
	needCapitalize := false

	for _, char := range src {
		if char == '%' && isHexDigit(src[1]) && isHexDigit(src[2]) {
			if isLower(src[1]) || isLower(src[2]) {
				needCapitalize = true
			}
			src = src[3:]
		} else if char&0x80 != 0 {
			numToEscape++
		}
	}

	if numToEscape == 0 && !needCapitalize {
		return src
	}

	dst := make([]byte, len(src)+numToEscape*2)
	j := 0
	for i := 0; i < len(src); i++ {
		if src[i] == '%' && isHexDigit(src[i+1]) && isHexDigit(src[i+2]) {
			dst[j] = src[i]
			dst[j+1] = toUpper(src[i+1])
			dst[j+2] = toUpper(src[i+2])
			i += 2
		} else if src[i]&0x80 != 0 {
			dst[j] = '%'
			dst[j+1] = kHexDigits[(src[i]>>4)&0xf]
			dst[j+2] = kHexDigits[src[i]&0xf]
		} else {
			dst[j] = src[i]
		}
		j++
	}
	return string(dst)
}

func isHexDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ('A' <= ch && ch <= 'F') || ('a' <= ch && ch <= 'f')
}

func isLower(ch byte) bool {
	return 'a' <= ch && ch <= 'z'
}

func toUpper(ch byte) byte {
	if 'a' <= ch && ch <= 'z' {
		return ch - 'a' + 'A'
	}
	return ch
}

func getKeyAndValueFrom(line string) (key, value string, ok bool) {
	idx := strings.IndexFunc(line, func(char rune) bool {
		return char == ':' || unicode.IsSpace(char)
	})
	if idx == -1 {
		return "", "", false
	}

	key = strings.TrimSpace(line[:idx])
	value = strings.TrimSpace(line[idx+1:])
	return key, value, true
}

var disallowedTypos = []string{
	"dissallow", "dissalow", "disalow", "diasllow", "disallaw",
}

const allowFrequentTypos = true

type keyType int

const (
	keyTypeUserAgent keyType = iota
	keyTypeSiteMap
	keyTypeAllow
	keyTypeDisallow
	keyTypeUnknown = 128
)

func parseKey(key string) keyType {
	lowerKey := strings.ToLower(key)
	switch lowerKey {
	case "user-agent":
		return keyTypeUserAgent
	case "allow":
		return keyTypeAllow
	case "disallow":
		return keyTypeDisallow
	case "sitemap":
		return keyTypeSiteMap
	case "site-map":
		return keyTypeSiteMap
	}
	if !allowFrequentTypos {
		return keyTypeUnknown
	}
	if lowerKey == "useragent" || lowerKey == "user agent" {
		return keyTypeUserAgent
	}
	for _, typo := range disallowedTypos {
		if lowerKey == typo {
			return keyTypeDisallow
		}
	}
	return keyTypeUnknown
}
