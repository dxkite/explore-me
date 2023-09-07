package core

import (
	"io"
	"os"
	"strings"

	"dxkite.cn/explorer/src/core/scan"
	"dxkite.cn/explorer/src/core/stream"
)

type SearchParams struct {
	Name string
	Tag  string
	Ext  string
	Path string
}

type SearchFileInfo struct {
	Id int64 `json:"id"`
	*scan.Index
}

func SearchFile(filename string, match SearchParams, offset, limit int64) ([]*SearchFileInfo, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	s := stream.NewJsonStream(f)
	v := createSearchParam(match)

	rst := []*SearchFileInfo{}

	var take int64

	for {
		offset, info, err := s.ScanNext(&scan.Index{}, v)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		fi := info.(*scan.Index)

		if !isMatchSearch(fi, match) {
			continue
		}

		rst = append(rst, &SearchFileInfo{Id: offset, Index: fi})
		take++
		if limit == -1 {
			continue
		}

		if take >= limit {
			break
		}

	}
	return rst, nil
}

func createSearchParam(match SearchParams) [][]string {
	param := [][]string{}
	cond := []string{}
	if match.Name != "" {
		cond = append(cond, match.Name)
	}

	if match.Ext != "" {
		cond = append(cond, match.Ext)
	}

	if match.Tag != "" {
		cond = append(cond, match.Tag)
	}

	if match.Path != "" {
		cond = append(cond, match.Path)
	}

	param = append(param, cond)
	return param
}

// 强匹配
func isMatchSearch(fi *scan.Index, match SearchParams) bool {
	if match.Path != "" {
		if strings.Index(fi.Path, match.Path) == -1 {
			return false
		}
	}

	if match.Name != "" {
		if strings.Index(fi.Name, match.Name) == -1 {
			return false
		}
	}

	if match.Ext != "" {
		if fi.Ext != match.Ext {
			return false
		}
	}

	if match.Tag != "" {
		mm := false
		for _, t := range fi.Tags {
			if t == match.Tag {
				mm = true
				break
			}
		}
		if !mm {
			return false
		}
	}

	return true
}
