package core

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type FileInfo struct {
	Name string   `json:"n,omitempty"`
	Path string   `json:"p,omitempty"`
	Tags []string `json:"t,omitempty"`
	Ext  string   `json:"e,omitempty"`
}

type MetaData struct {
	LastUpdate time.Time `json:"last_update"`
	CreateTime time.Time `json:"create_time"`
}

type SearchParams struct {
	Name string
	Tag  string
	Ext  string
}

type ExtValue struct {
	Count  int  `json:"count"`
	Ignore bool `json:"ignore"`
}

type IndexCreator struct {
	Config        *ScanConfig
	ignoreNameMap map[string]bool
	extMap        map[string]ExtValue
	tagMap        map[string]int
}

func InitIndex(cfg *Config) error {
	ic := NewIndexCreator(&cfg.ScanConfig)
	return ic.Create(cfg.SrcRoot, cfg.DataRoot)
}

func NewIndexCreator(cfg *ScanConfig) *IndexCreator {
	ic := &IndexCreator{}
	ic.Config = cfg

	ic.ignoreNameMap = map[string]bool{}
	for _, v := range cfg.IgnoreName {
		ic.ignoreNameMap[v] = true
	}

	ic.extMap = map[string]ExtValue{}
	for _, v := range cfg.IgnoreExt {
		ic.extMap[v] = ExtValue{Ignore: true}
	}

	ic.tagMap = map[string]int{}
	return ic
}

// 扫描目录
func (ic *IndexCreator) Create(root, dataRoot string) error {
	meta := ic.getMeta(dataRoot)

	// 修改时间没有变化
	if fi, err := os.Stat(root); err == nil {
		if fi.ModTime() == meta.LastUpdate {
			return nil
		}

		meta.CreateTime = time.Now()
		meta.LastUpdate = fi.ModTime()
	}

	if err := ic.createIndexFile(root, dataRoot); err != nil {
		return err
	}

	if err := ic.createExtListFile(dataRoot); err != nil {
		return err
	}

	if err := ic.createTagListFile(dataRoot); err != nil {
		return err
	}

	if err := writeJsonFile(path.Join(dataRoot, ic.Config.MetaFile), meta); err != nil {
		return err
	}
	return nil
}

func (ic *IndexCreator) createIndexFile(root, dataRoot string) error {
	reg, err := regexp.Compile(ic.Config.TagExpr)
	if err != nil {
		log.Panicln("compile reg expr error", ic.Config.TagExpr, err)
		return err
	}

	index := path.Join(dataRoot, ic.Config.IndexFile)
	if err := os.MkdirAll(dataRoot, os.ModePerm); err != nil {
		log.Panicln("mkdir all", dataRoot, err)
		return err
	}

	idx, err := os.OpenFile(index, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return err
	}

	defer idx.Close()

	absRootPath, _ := filepath.Abs(root)

	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		name := info.Name()
		ext := GetExt(name)

		if ic.ignoreNameMap[name] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if v, ok := ic.extMap[ext]; ok {
			v.Count++
			ic.extMap[ext] = v
		} else {
			ic.extMap[ext] = ExtValue{Count: 1}
		}

		if v, ok := ic.extMap[ext]; ok && v.Ignore {
			return nil
		}

		tags, err := parseTag(name, reg)
		if err != nil {
			return err
		}

		for _, v := range tags {
			ic.tagMap[v]++
		}

		filePath := strings.TrimPrefix(path, absRootPath)
		filePath = normalizePath(filePath)
		v := FileInfo{
			Name: name,
			Path: filePath,
			Tags: tags,
			Ext:  ext,
		}

		if b, err := json.Marshal(v); err != nil {
			return err
		} else {
			if _, err := idx.Write(b); err != nil {
				return err
			}
			if _, err := idx.Write([]byte{'\n'}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (ic *IndexCreator) getMeta(dataRoot string) *MetaData {
	meta := &MetaData{}
	filename := path.Join(dataRoot, ic.Config.MetaFile)
	b, err := os.ReadFile(filename)
	if err != nil {
		return meta
	}

	if err := json.Unmarshal(b, meta); err != nil {
		return meta
	}
	return meta
}

func (ic *IndexCreator) createExtListFile(dataRoot string) error {
	filename := path.Join(dataRoot, ic.Config.ExtListFile)
	if err := writeJsonFile(filename, ic.extMap); err != nil {
		return err
	}
	return nil
}

func (ic *IndexCreator) createTagListFile(dataRoot string) error {
	filename := path.Join(dataRoot, ic.Config.TagListFile)
	if err := writeJsonFile(filename, ic.tagMap); err != nil {
		return err
	}
	return nil
}

func writeJsonFile(filename string, v interface{}) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()
	if data, err := json.Marshal(v); err != nil {
		return err
	} else {
		if _, err := f.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func ParseTag(cfg *Config, name string) ([]string, error) {
	reg, err := regexp.Compile(cfg.ScanConfig.TagExpr)
	if err != nil {
		return []string{}, err
	}

	return parseTag(name, reg)
}

func parseTag(name string, reg *regexp.Regexp) ([]string, error) {
	matches := reg.FindAllStringSubmatch(name, -1)
	tags := []string{}
	for _, m := range matches {
		tags = append(tags, m[1])
	}
	return tags, nil
}

func normalizePath(filename string) string {
	return strings.ReplaceAll(filename, "\\", "/")
}

func GetExt(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) <= 1 {
		return ""
	}
	return strings.ToLower(ext[1:])
}

func NormalizePath(root, filename string) string {
	absRoot, _ := filepath.Abs(root)
	absName, _ := filepath.Abs(filename)
	newName := strings.TrimPrefix(absName, absRoot)
	return normalizePath(newName)
}
