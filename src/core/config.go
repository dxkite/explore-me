package core

import (
	"os"

	goget "dxkite.cn/explorer/src/middleware/go-get"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

// 基础配置
type ScanConfig struct {
	// .* 忽略文件名
	IgnoreNameExpr string `yaml:"ignore_name_expr" default:"^\\."`
	// 忽略文件扩展名
	IgnoreExt []string `yaml:"ignore_ext" default:"[\"c\",\"h\",\"php\",\"h\"]"`
	// 忽略文件
	IgnoreName []string `yaml:"ignore_name" default:"[\".git\"]"`
	// [(.*)] 标签表达式
	TagExpr string `yaml:"tag_expr" default:"\\[(.+?)\\]"`

	ReadmeFile  string `yaml:"readme_file" default:"README.md"`
	TagListFile string `yaml:"tag_list_file" default:"tags.json"`
	ExtListFile string `yaml:"ext_list_file" default:"exts.json"`
	IndexFile   string `yaml:"index_file" default:"index.json"`
	MetaFile    string `yaml:"meta_file" default:"meta.json"`
}

type Config struct {
	Listen string `yaml:"listen" default:":8080"`
	// 网站目录
	WebRoot string `yaml:"web_root" default:"./web"`
	// 资源目录
	SrcRoot string `yaml:"src_root" default:"./src"`
	// 数据目录
	DataRoot string `yaml:"data_root" default:"./data"`
	// 解析配置
	ScanConfig ScanConfig `yaml:"scan_config"`
	// 自动刷新时间 60s
	AsyncLoad int `yaml:"async_time" default:"60"`
	// go-get
	GoGetConfig goget.PackageConfig `yaml:"go_get_config"`
}

func LoadConfig(filename string) (*Config, error) {
	cfg := &Config{}
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	c, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(c, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		c := &Config{}
		defaults.Set(c)
		return c
	}
	return cfg
}

func InitConfig(filename string) error {
	c, err := LoadConfig(filename)
	if err != nil {
		return err
	}
	cfg = c
	return nil
}
