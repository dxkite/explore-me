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

var defaultConfig = `# 监听地址
listen: :80
# 网站根目录
# 可以将静态文件放到这个目录下面能直接访问
web_root: .explore-me/web
# 文件目录
src_root: ./
# 数据目录
data_root: .explore-me/data
# 配置扫描间隔
async_time: 60
# 目录配置
dir_config:
# 显示的时候忽略的文件
  ignore_name:
    - ^\..+$
    - ^explore-me.+$
    - .*\.meta\.yaml$
    - \.dir-config\.yaml$
theme_config: .explore-me/theme-config.yaml
# 扫描配置
scan_config:
  # 忽略指定扩展名
  ignore_ext:
    - php
    - c
    - h
    - js
    - css
    - vue
    - html
    - asm
    - ttf
    - woff
    - sql
  ignore_name:
    - .git
    - .DS_Store
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
