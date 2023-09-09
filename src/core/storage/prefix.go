package storage

import (
	"net/http"
	"path"
)

type prefixFS struct {
	src    http.FileSystem
	prefix string
}

func (s *prefixFS) Open(name string) (http.File, error) {
	name = path.Join(s.prefix, name)
	return s.src.Open(name)
}

func NewPrefix(prefix string, src http.FileSystem) http.FileSystem {
	return &prefixFS{src: src, prefix: prefix}
}
