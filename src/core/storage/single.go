package storage

import (
	"net/http"
	"os"

	"dxkite.cn/log"
)

type single struct {
	src   http.FileSystem
	index string
}

func NewSingleIndex(src http.FileSystem, index string) http.FileSystem {
	return &single{src: src, index: index}
}

func (s *single) Open(name string) (http.File, error) {
	f, err := s.src.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("use single", s.index)
			return s.src.Open(s.index)
		}
		return nil, err
	}
	return f, nil
}
