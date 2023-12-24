package firstrun

import (
	"os"
	"path"

	"errors"
)

var themeConfig = `name: "Explore Me"
logo: "/logo.png"
copyrightName: "dxkite"
websiteRecord: ""
websiteRecordLink: ""
websitePoliceRecord: ""
websitePoliceLink: ""
textViewExt: ["txt", "cpp", "c", "js", "yaml"]
videoViewExt: ["3gpp","3gp","ts","mp4","mpeg","mpg","mov","webm","flv","m4v","mng","asx","asf","wmv","avi"]
markdownRawExt: ["jpg","jpeg","gif","png","svg","webp","bmp","ico",""]
`

var defaultConfig = `listen: :80
web_root: .explore-me/web
src_root: ./
web_index: /index.html
data_root: .explore-me/data
async_time: 60
dir_config:
    config_name: .dir-config.yaml
    meta_name: .meta.yaml
    ignore_name:
        - ^\..+$
theme_config: .explore-me/theme-config.yaml
`

func Init() error {
	root := "./"

	configPath := path.Join(root, ".explore-me", "config.yaml")
	themeConfigFile := path.Join(root, ".explore-me", "theme-config.yaml")

	if isExist(configPath) {
		return nil
	}

	exploreMe := path.Join(root, ".explore-me")
	if err := os.MkdirAll(exploreMe, os.ModePerm); err != nil {
		return errors.Join(errors.New("make dir error"), err)
	}

	_ = Hide(exploreMe)

	if err := os.MkdirAll(path.Join(exploreMe, "web"), os.ModePerm); err != nil {
		return errors.Join(errors.New("make dir error"), err)
	}

	if err := os.WriteFile(configPath, []byte(defaultConfig), os.ModePerm); err != nil {
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
