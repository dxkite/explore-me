package firstrun

import (
	"os"
	"path"

	"errors"

	"dxkite.cn/explore-me/src/core/config"
	"dxkite.cn/explore-me/src/core/scan"
	"gopkg.in/yaml.v3"
)

func Init(root string) error {
	initConfig := &config.Config{}
	initConfig.Listen = ":80"
	initConfig.AsyncLoad = 60
	initConfig.DataRoot = path.Join(root, ".explore-me/data")
	initConfig.WebRoot = path.Join(root, ".explore-me/web")
	initConfig.SrcRoot = root
	initConfig.DirConfig = scan.DirConfig{
		IgnoreName: []string{
			".explore-me",
			".git",
		},
	}

	out, err := yaml.Marshal(initConfig)
	if err != nil {
		return errors.Join(errors.New("config config error"), err)
	}

	if err := os.MkdirAll(path.Join(root, ".explore-me"), os.ModePerm); err != nil {
		return errors.Join(errors.New("make dir error"), err)
	}

	if err := os.WriteFile(path.Join(root, ".explore-me", "config.yaml"), out, os.ModePerm); err != nil {
		return errors.Join(errors.New("create config error"), err)
	}

	return nil
}
