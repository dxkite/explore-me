package actions

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"dxkite.cn/log"

	"dxkite.cn/explorer/src/core/config"
	"dxkite.cn/explorer/src/core/scan"
	"dxkite.cn/explorer/src/core/storage"
	"github.com/gin-gonic/gin"
)

func Meta(c *gin.Context) {
	cfg := config.GetConfig()
	fs := storage.Local(cfg.SrcRoot)

	pathname := c.Param("path")

	log.Println("LoadMeta", pathname, "24")

	fi, err := fs.Stat(c, pathname)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusNotFound)
		return
	}

	m := createMeta(cfg, context.TODO(), fs, pathname, fi)
	if m.IsDir {
		ch, rm, err := getDir(cfg, context.TODO(), fs, pathname)
		if err != nil {
			log.Println("getDir error", pathname, err)
		}
		if ch != nil {
			m.Children = ch
		}
		if rm != nil {
			m.Readme = path.Join(pathname, rm.Name())
		}
	}
	c.JSON(http.StatusOK, m)
}

type MetaData struct {
	Id       int64       `json:"id,omitempty"`
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Tags     []string    `json:"tags,omitempty"`
	Ext      string      `json:"ext,omitempty"`
	IsDir    bool        `json:"is_dir"`
	Readme   string      `json:"readme"`
	ModTime  string      `json:"mod_time"`
	RawUrl   string      `json:"raw_url,omitempty"`
	Children []*MetaData `json:"children,omitempty"`
}

func createMeta(cfg *config.Config, ctx context.Context, fs storage.FileSystem, pathname string, fi fs.FileInfo) *MetaData {
	meta := scan.GetFileMeta(ctx, fs, pathname, fi)
	m := &MetaData{}
	m.Name = meta.Name
	m.Path = pathname
	m.Tags = meta.Tags
	m.Ext = scan.GetExt(fi.Name())
	m.IsDir = fi.IsDir()
	m.ModTime = meta.ModTime
	if !m.IsDir {
		m.RawUrl = path.Join(config.RawUrlRoot, m.Path)
	}
	return m
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

func getDir(cfg *config.Config, ctx context.Context, src storage.FileSystem, dirname string) ([]*MetaData, fs.FileInfo, error) {
	dirCfg := scan.LoadConfigForDir(ctx, src, &cfg.DirConfig, dirname, cfg.DirConfig.ConfigName)
	log.Println("dirConfig", dirname, dirCfg)

	ctx = context.WithValue(ctx, scan.DirConfigKey, dirCfg)

	dirInfo, err := scan.ReadDir(ctx, src, dirname)
	if err != nil {
		return nil, nil, err
	}

	md := []*MetaData{}

	var readme fs.FileInfo
	for _, di := range dirInfo {
		pathname := path.Join(dirname, di.Name())
		mdi := createMeta(cfg, ctx, src, pathname, di)
		md = append(md, mdi)
		if strings.ToLower(di.Name()) == "readme.md" {
			readme = di
		}
	}
	return md, readme, nil
}
