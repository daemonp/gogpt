package exporter

import (
	"bufio"
	"bytes"
	"regexp"
)

type ContentFilter struct {
	excludePattern *regexp.Regexp
}

func NewContentFilter(pattern string) (*ContentFilter, error) {
	if pattern == "" {
		return &ContentFilter{}, nil
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &ContentFilter{
		excludePattern: regex,
	}, nil
}

func (cf *ContentFilter) Filter(content []byte) []byte {
	if cf.excludePattern == nil {
		return content
	}

	var filteredLines [][]byte
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Bytes()
		if !cf.excludePattern.Match(line) {
			filteredLines = append(filteredLines, line)
		}
	}

	return bytes.Join(filteredLines, []byte("\n"))
}
