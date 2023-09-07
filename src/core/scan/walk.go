package scan

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path"
	"regexp"
	"sort"

	"dxkite.cn/explorer/src/core/storage"
	"gopkg.in/yaml.v3"
)

type DirConfig struct {
	// 配置名
	ConfigName string `yaml:"config_name" default:".config.yml"`
	// 配置名
	MetaName string `yaml:"config_name" default:".meta.yml"`
	// 忽略文件名
	IgnoreName []string `yaml:"ignore_name"`
	// 置顶
	Pin []string `yaml:"pin"`
	// [(.*)] 标签表达式
	TagExpr string `yaml:"tag_expr" default:"\\[(.+?)\\]"`
}

type dirConfig string

const DirConfigKey dirConfig = "DirConfig"

var SkipDir error = fs.SkipDir
var SkipAll error = fs.SkipAll

type WalkCallback func(ctx context.Context, fs storage.FileSystem, path string, info fs.FileInfo, err error) error

func Walk(ctx context.Context, fs storage.FileSystem, root string, fn WalkCallback) error {
	info, err := fs.Stat(ctx, root)
	if err != nil {
		err = fn(ctx, fs, root, nil, err)
	} else {
		err = walk(ctx, fs, root, info, fn)
	}
	if err == SkipDir || err == SkipAll {
		return nil
	}
	return err
}

func walk(ctx context.Context, fs storage.FileSystem, name string, info fs.FileInfo, fn WalkCallback) error {
	if !info.IsDir() {
		return fn(ctx, fs, name, info, nil)
	}

	ctx, err := createContextFromDir(ctx, fs, name)
	if err != nil {
		log.Println("createContextFromDir", err)
		return err
	}

	infos, err := readDirNameFromFs(ctx, fs, name)
	if err != nil {
		log.Println("readDirNameFromFs", err)
		return err
	}

	cfg := getConfigFromContext(ctx)

	sortNames(cfg, infos)

	for _, item := range infos {
		filename := path.Join(name, item.Name())
		if isIgnoreName(cfg, item.Name()) {
			continue
		}

		err := walk(ctx, fs, filename, item, fn)
		if err != nil {
			if !item.IsDir() || err != SkipDir {
				return err
			}
		}
	}
	return nil
}

func sortNames(cfg *DirConfig, infos []fs.FileInfo) {
	pinIdx := map[string]int{}

	pin := false
	if cfg != nil && len(cfg.Pin) > 0 {
		for i, pin := range cfg.Pin {
			pinIdx[pin] = i
		}
		pin = true
	}

	if !pin {
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Name() < infos[j].Name()
		})
	}

	sort.Slice(infos, func(i, j int) bool {
		ni, nj := infos[i].Name(), infos[j].Name()

		ii, oki := pinIdx[ni]
		ij, okj := pinIdx[nj]

		if oki && okj {
			return ii < ij
		}

		if oki {
			return false
		}

		if okj {
			return true
		}

		return ni < nj
	})
}

func isIgnoreName(cfg *DirConfig, name string) bool {

	if cfg == nil {
		return false
	}

	if len(cfg.IgnoreName) == 0 {
		return false
	}

	for _, expr := range cfg.IgnoreName {
		exp, _ := loadExpr(expr)
		if exp != nil && exp.Match([]byte(name)) {
			return true
		}
	}

	return false
}

func readDirNameFromFs(ctx context.Context, fs storage.FileSystem, name string) ([]fs.FileInfo, error) {
	f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func createContextFromDir(ctx context.Context, fs storage.FileSystem, name string) (context.Context, error) {
	cfg := getConfigFromContext(ctx)
	if cfg == nil {
		return ctx, nil
	}

	cfgPath := path.Join(name, cfg.ConfigName)
	d, err := os.ReadFile(cfgPath)
	if err != nil {
		return ctx, err
	}

	if err := yaml.Unmarshal(d, cfg); err != nil {
		return ctx, err
	}

	return ctx, nil
}

func getConfigFromContext(ctx context.Context) *DirConfig {
	if v, ok := ctx.Value(DirConfigKey).(*DirConfig); ok {
		return v
	}
	return nil
}

var exprCache map[string]*regexp.Regexp

func loadExpr(expr string) (*regexp.Regexp, error) {
	if v, ok := exprCache[expr]; ok {
		return v, nil
	}

	if v, err := regexp.Compile(expr); err != nil {
		return nil, err
	} else {
		exprCache[expr] = v
		return v, nil
	}
}
