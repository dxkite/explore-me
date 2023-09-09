package storage

import (
	"net/http"
	"os"
)

type multiFileSystem struct {
	src       []http.FileSystem
	defaultFs int
}

func NewMultiFileSystem(src ...http.FileSystem) http.FileSystem {
	return &multiFileSystem{src: src, defaultFs: len(src) - 1}
}

func (m *multiFileSystem) Open(name string) (http.File, error) {
	for _, src := range m.src {
		f, err := src.Open(name)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		return f, nil
	}

	return m.src[m.defaultFs].Open(name)
}
