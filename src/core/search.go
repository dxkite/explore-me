package core

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

func SearchFile(filename string, match SearchParams, offset, limit int) ([]*FileInfo, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	rst := []*FileInfo{}
	s := bufio.NewScanner(f)

	skip := 0
	take := 0
	for s.Scan() {
		line := s.Text()

		if !isFuzzyMatchSearch(line, match) {
			continue
		}
		fi, ok := isMatchSearch(line, match)

		if !ok {
			continue
		}

		if skip < offset {
			skip++
			continue
		}

		if take >= limit {
			break
		}

		take++

		rst = append(rst, fi)
	}
	return rst, nil
}

// 模糊匹配
func isFuzzyMatchSearch(text string, match SearchParams) bool {
	if match.Name != "" && strings.Index(text, match.Name) >= 0 {
		return true
	}
	if match.Ext != "" && strings.Index(text, match.Ext) >= 0 {
		return true
	}
	if match.Tag != "" && strings.Index(text, match.Tag) >= 0 {
		return true
	}
	return false
}

// 强匹配
func isMatchSearch(text string, match SearchParams) (*FileInfo, bool) {
	fi := &FileInfo{}

	if err := json.Unmarshal([]byte(text), fi); err != nil {
		return nil, false
	}

	if match.Name != "" && strings.Index(fi.Name, match.Name) >= 0 {
		return fi, true
	}

	if match.Ext != "" && fi.Ext == match.Ext {
		return fi, true
	}

	for _, t := range fi.Tags {
		if match.Tag != "" && t == match.Ext {
			return fi, true
		}
	}
	return nil, false
}
