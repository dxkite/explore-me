package scan

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"dxkite.cn/explorer/src/core/storage"
	"gopkg.in/yaml.v3"
)

const (
	ExtIndex  = "exts.json"
	TagIndex  = "tags.json"
	MetaIndex = "index.json"
	LockFile  = "scan.lock"
)

var DefaultTagExpr *regexp.Regexp

func init() {
	DefaultTagExpr, _ = regexp.Compile("\\[(.+?)\\]")
}

type Index struct {
	Name string   `json:"n,omitempty"`
	Path string   `json:"p,omitempty"`
	Tags []string `json:"t,omitempty"`
	Ext  string   `json:"e,omitempty"`
}

type FileMeta struct {
	Name    string   `yaml:"name"`
	ModTime string   `yaml:"mod_time"`
	Tags    []string `yaml:"tags"`
}

type Scanner struct {
	fs     storage.FileSystem
	tagMap map[string]int
	extMap map[string]int
	idx    io.WriteCloser
	output string
}

func NewScanner(output string) *Scanner {
	return &Scanner{
		tagMap: map[string]int{},
		extMap: map[string]int{},
		output: output,
	}
}

func (s *Scanner) Scan(ctx context.Context, fs storage.FileSystem) error {
	if err := os.MkdirAll(s.output, os.ModePerm); err != nil {
		return err
	}
	lockFile := path.Join(s.output, LockFile)
	upTime, update := s.isUpdate(ctx, fs, lockFile)
	if update == false {
		return nil
	}

	idxFile := path.Join(s.output, MetaIndex)
	f, err := os.OpenFile(idxFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	s.idx = f
	if err := Walk(ctx, fs, ".", s.scanIndex); err != nil {
		return err
	}

	extFile := path.Join(s.output, ExtIndex)
	if err := writeJsonFile(extFile, s.extMap); err != nil {
		return err
	}

	tagFile := path.Join(s.output, TagIndex)
	if err := writeJsonFile(tagFile, s.tagMap); err != nil {
		return err
	}

	writeFile(lockFile, []byte(upTime.Format(time.DateTime)))
	return nil
}

func (s *Scanner) isUpdate(ctx context.Context, fs storage.FileSystem, lockFile string) (time.Time, bool) {
	t := getLockTime(lockFile)
	info, _ := fs.Stat(ctx, ".")
	if t == nil || info == nil {
		return time.Now(), true
	}
	if info.ModTime().After(*t) {
		return info.ModTime(), true
	}
	return info.ModTime(), false
}

func getLockTime(p string) *time.Time {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil
	}
	t, err := time.Parse(time.DateTime, string(b))
	if err != nil {
		return nil
	}
	return &t
}

func (s *Scanner) scanIndex(ctx context.Context, fs storage.FileSystem, name string, info fs.FileInfo, err error) error {
	if err != nil {
		return nil
	}

	if info.IsDir() {
		return nil
	}

	meta := GetFileMeta(ctx, fs, name, info)

	ext := GetExt(info.Name())
	s.extMap[ext]++

	for _, v := range meta.Tags {
		s.tagMap[v]++
	}

	v := Index{
		Name: meta.Name,
		Path: name,
		Tags: meta.Tags,
		Ext:  ext,
	}

	if b, err := json.Marshal(v); err != nil {
		return err
	} else {
		if _, err := s.idx.Write(b); err != nil {
			return err
		}
		if _, err := s.idx.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	return nil
}

func GetFileMeta(ctx context.Context, fs storage.FileSystem, name string, info fs.FileInfo) *FileMeta {
	meta := &FileMeta{}
	meta.ModTime = info.ModTime().Format(time.DateTime)
	meta.Name = info.Name()
	meta.Tags = []string{}
	tagExpr := DefaultTagExpr
	var metaCfg string

	cfg := getConfigFromContext(ctx)
	if cfg != nil {
		if cfg.TagExpr != "" {
			if expr, err := loadExpr(cfg.TagExpr); err == nil {
				tagExpr = expr
			}
		}
		metaCfg = cfg.MetaName
	}

	if tags, err := parseTag(info.Name(), tagExpr); err == nil {
		meta.Tags = tags
	}

	if len(metaCfg) > 0 {
		fn := name + metaCfg
		meta = loadMetaFrom(ctx, fs, fn, meta)
	}

	return meta
}

func parseTag(name string, reg *regexp.Regexp) ([]string, error) {
	matches := reg.FindAllStringSubmatch(name, -1)
	// log.Println(matches)
	tags := []string{}
	for _, m := range matches {
		if len(m) >= 1 {
			tags = append(tags, m[1])
		}
	}
	return tags, nil
}

func isExist(ctx context.Context, fs storage.FileSystem, filename string) bool {
	_, err := fs.Stat(ctx, filename)
	if err != nil {
		return false
	}
	return true
}

func loadMetaFrom(ctx context.Context, fs storage.FileSystem, filename string, defVal *FileMeta) *FileMeta {
	log.Println("read", filename)

	if !isExist(ctx, fs, filename) {
		// log.Println("read", filename, "isExist")
		return defVal
	}

	r, err := fs.OpenFile(ctx, filename, os.O_RDONLY, 0)
	if err != nil {
		// log.Println("read", err)
		return defVal
	}

	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		// log.Println("read", err)
		return defVal
	}
	if err := yaml.Unmarshal(b, defVal); err != nil {
		return defVal
	}
	return defVal
}

func GetExt(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) <= 1 {
		return ""
	}
	return strings.ToLower(ext[1:])
}

func writeJsonFile(filename string, v interface{}) error {
	if data, err := json.Marshal(v); err != nil {
		return err
	} else {
		return writeFile(filename, data)
	}
}

func writeFile(filename string, d []byte) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()
	if _, err := f.Write(d); err != nil {
		return err
	}
	return nil
}
