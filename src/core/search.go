package core

import (
	"io"
	"os"
	"strings"

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
	*FileInfo
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
		offset, info, err := s.ScanNext(&FileInfo{}, v)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		fi := info.(*FileInfo)

		if !isMatchSearch(fi, match) {
			continue
		}

		rst = append(rst, &SearchFileInfo{Id: offset, FileInfo: fi})
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
	if match.Name != "" {
		param = append(param, []string{match.Name})
	}

	if match.Ext != "" {
		param = append(param, []string{match.Ext})
	}

	if match.Tag != "" {
		param = append(param, []string{match.Tag})
	}

	if match.Path != "" {
		param = append(param, []string{match.Path})
	}

	return param
}

// 强匹配
func isMatchSearch(fi *FileInfo, match SearchParams) bool {
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
			if t == match.Ext {
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
