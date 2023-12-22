package firstrun

import (
	"os"
	"path"

	"errors"

	"dxkite.cn/explore-me/src/core/config"
	"dxkite.cn/explore-me/src/core/scan"
	"gopkg.in/yaml.v3"
)

var themeConfig = `name: "Explore Me"
logo: "/dxkite.png"
copyrightName: "dxkite"
websiteRecord: ""
websiteRecordLink: ""
websitePoliceRecord: ""
websitePoliceLink: ""
textViewExt: ["txt", "cpp", "c", "js", "yaml"]
videoViewExt: ["3gpp","3gp","ts","mp4","mpeg","mpg","mov","webm","flv","m4v","mng","asx","asf","wmv","avi"]
markdownRawExt: ["jpg","jpeg","gif","png","svg","webp","bmp","ico",""]
`

func Init(root string) error {
	configPath := path.Join(root, ".explore-me", "config.yaml")
	themeConfigFile := path.Join(root, ".explore-me", "theme-config.yaml")

	if isExist(configPath) {
		return nil
	}

	initConfig := &config.Config{}
	initConfig.Listen = ":80"
	initConfig.AsyncLoad = 60
	initConfig.DataRoot = path.Join(root, ".explore-me/data")
	initConfig.WebRoot = path.Join(root, ".explore-me/web")
	initConfig.SrcRoot = root
	initConfig.ThemeConfig = themeConfigFile
	initConfig.DirConfig = scan.DirConfig{
		ConfigName: ".dir-config.yaml",
		MetaName:   ".meta.yaml",
		TagExpr:    "\\[(.+?)\\]",
		IgnoreName: []string{
			"^\\..+$",
		},
	}

	out, err := yaml.Marshal(initConfig)
	if err != nil {
		return errors.Join(errors.New("config config error"), err)
	}

	if err := os.MkdirAll(path.Join(root, ".explore-me"), os.ModePerm); err != nil {
		return errors.Join(errors.New("make dir error"), err)
	}

	if err := os.WriteFile(configPath, out, os.ModePerm); err != nil {
		return errors.Join(errors.New("create config error"), err)
	}

	if err := os.WriteFile(themeConfigFile, []byte(themeConfig), os.ModePerm); err != nil {
		return errors.Join(errors.New("create theme config error"), err)
	}

	return nil
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}
