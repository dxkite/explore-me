package config

import (
	"os"

	"dxkite.cn/explorer/src/core/scan"
	goget "dxkite.cn/explorer/src/middleware/go-get"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

const RawUrlRoot = "/api/explore/raw"

type Config struct {
	Listen string `yaml:"listen" default:":8080"`
	// 网站目录
	WebRoot string `yaml:"web_root" default:"./web"`
	// 单页应用
	SingleIndex string `yaml:"web_index" default:"/index.html"`
	// 资源目录
	SrcRoot string `yaml:"src_root" default:"./src"`
	// 数据目录
	DataRoot string `yaml:"data_root" default:"./data"`
	// 自动刷新时间 60s
	AsyncLoad int `yaml:"async_time" default:"60"`
	// 目录配置
	DirConfig scan.DirConfig `yaml:"dir_config"`
	// 主题目录
	ThemeRoot string `yaml:"theme_root" default:"./theme"`
	// 主题
	Theme string `yaml:"theme" default:"./default"`
	// go-get
	GoGetConfig goget.PackageConfig `yaml:"go_get_config"`
	// client-id
	ClientIdKey string `yaml:"client_id_key" default:"client-id"`
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
