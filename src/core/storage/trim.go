package storage

import (
	"net/http"
	"path"
)

type trimPrefixFS struct {
	src    http.FileSystem
	prefix string
}

func (s *trimPrefixFS) Open(name string) (http.File, error) {
	name = path.Join(s.prefix, name)
	return s.src.Open(name)
}

func NewPrefix(prefix string, src http.FileSystem) http.FileSystem {
	return &trimPrefixFS{src: src, prefix: prefix}
}
