package scan

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"

	"dxkite.cn/log"

	"github.com/dlclark/regexp2"

	"dxkite.cn/explorer/src/core/storage"
	"gopkg.in/yaml.v3"
)

type DirConfig struct {
	// 配置名
	ConfigName string `yaml:"config_name" default:".dir-config.yaml"`
	// 配置名
	MetaName string `yaml:"meta_name" default:".meta.yaml"`
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

	infos, err := ReadDir(ctx, fs, name)
	if err != nil {
		log.Println("readDirNameFromFs", err)
		return err
	}

	for _, item := range infos {
		filename := path.Join(name, item.Name())
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
	// log.Println("ignore", cfg.IgnoreName, name, "test")
	for _, expr := range cfg.IgnoreName {
		exp, err := loadExpr(expr)
		if err != nil {
			log.Println(expr, "reg error", err)
		} else {
			if match, _ := exp.MatchString(name); match {
				return true
			}
		}
	}
	// log.Println("ignore", cfg.IgnoreName, name, "test failed")
	return false
}

func readDir(ctx context.Context, fs storage.FileSystem, name string) ([]fs.FileInfo, error) {
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

func ReadDir(ctx context.Context, src storage.FileSystem, name string) ([]fs.FileInfo, error) {
	infos, err := readDir(ctx, src, name)
	if err != nil {
		return nil, err
	}
	if c, err := createContextFromDir(ctx, src, name); err == nil {
		ctx = c
	}
	cfg := getConfigFromContext(ctx)
	log.Println("ReadDir", cfg)
	sortNames(cfg, infos)
	newInfos := []fs.FileInfo{}
	for _, item := range infos {
		if isIgnoreName(cfg, item.Name()) {
			continue
		}
		newInfos = append(newInfos, item)
	}
	return newInfos, nil
}

func createContextFromDir(ctx context.Context, fs storage.FileSystem, name string) (context.Context, error) {
	cfg := getConfigFromContext(ctx)
	log.Println("read default config", name, cfg)

	if cfg == nil {
		return ctx, nil
	}
	cfgPath := path.Join(name, cfg.ConfigName)
	cfg, _ = LoadConfig(ctx, fs, cfg, cfgPath)

	log.Println("dir config", name, cfg)
	ctx = context.WithValue(ctx, DirConfigKey, cfg)
	return ctx, nil
}

func LoadConfig(ctx context.Context, fs storage.FileSystem, defCfg *DirConfig, filename string) (*DirConfig, error) {
	r, err := fs.OpenFile(ctx, filename, os.O_RDONLY, 0)
	if err != nil {
		return defCfg, err
	}
	defer r.Close()
	b, err := io.ReadAll(r)
	if err != nil {
		return defCfg, err
	}
	cfg := *defCfg
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return defCfg, err
	}
	return &cfg, nil
}

func LoadConfigForDir(ctx context.Context, fs storage.FileSystem, defCfg *DirConfig, dirname, cfgName string) *DirConfig {
	dirname = "/" + dirname
	for {
		dirname = path.Clean(dirname)
		cfgPath := path.Join(dirname, cfgName)
		log.Println("LoadConfigForDir", dirname, cfgPath)
		if cfg, err := LoadConfig(ctx, fs, defCfg, cfgPath); err == nil {
			return cfg
		} else {
			log.Println("LoadConfigForDirErr", dirname, cfgPath, err)
		}
		dirname = path.Dir(dirname)
		if dirname == "/" {
			break
		}
	}
	log.Println("LoadConfigForDirErr", dirname, "useDefaultConfig", defCfg)
	return defCfg
}

func getConfigFromContext(ctx context.Context) *DirConfig {
	if v, ok := ctx.Value(DirConfigKey).(*DirConfig); ok {
		return v
	}
	if v, ok := ctx.Value(DirConfigKey).(DirConfig); ok {
		return &v
	}
	return &DirConfig{
		ConfigName: ".dir-config.yaml",
		MetaName:   ".meta.yaml",
	}
}

func init() {
	exprCache = map[string]*regexp2.Regexp{}
}

var exprCache map[string]*regexp2.Regexp

func loadExpr(expr string) (*regexp2.Regexp, error) {
	if v, ok := exprCache[expr]; ok {
		return v, nil
	}

	if v, err := regexp2.Compile(expr, 0); err != nil {
		return nil, err
	} else {
		exprCache[expr] = v
		return v, nil
	}
}
