# robotstxt
Robots.txt parser and matcher library in go according to https://github.com/google/robotstxt.

[![GoDoc](https://godoc.org/go.qingyu31.com/robotstxt?status.svg)](https://godoc.org/go.qingyu31.com/robotstxt)
[![Go Report Card](https://goreportcard.com/badge/go.qingyu31.com/robotstxt)](https://goreportcard.com/report/go.qingyu31.com/robotstxt)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)]()
## About the library
The Robots Exclusion Protocol (REP) is a standard that enables website owners to control which URLs may be accessed by automated clients (i.e. crawlers) through a simple text file with a specific syntax. It's one of the basic building blocks of the internet as we know it and what allows search engines to operate.

Because the REP was only a de-facto standard for the past 25 years, different implementers implement parsing of robots.txt slightly differently, leading to confusion. Google aims to fix that by releasing the [parser](https://github.com/google/robotstxt) that Google uses.

This library is a native go implement of the Google parser with most compatible.

## Installation

```shell
go install go.qingyu31.com/robotstxt
```

## Usage
```go
package main

import (
	"bytes"
	"fmt"
	"go.qingyu31.com/robotstxt"
)

func main() {
	robotsTxt := "User-agent: *\nDisallow: /search"
	matcher := robotstxt.Parse(bytes.NewBufferString(robotsTxt))
	fmt.Println(matcher.OneAgentAllowedByRobots("/search", "Googlebot"))
	fmt.Println(matcher.OneAgentAllowedByRobots("/search", "Baiduspider"))
}
```