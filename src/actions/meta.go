package actions

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"dxkite.cn/explorer/src/core/config"
	"dxkite.cn/explorer/src/core/scan"
	"dxkite.cn/explorer/src/core/storage"
	"github.com/gin-gonic/gin"
)

func Meta(c *gin.Context) {
	cfg := config.GetConfig()
	fs := storage.Local(cfg.SrcRoot)

	pathname := c.Param("path")

	log.Println("Meta", pathname)

	fi, err := fs.Stat(c, pathname)
	if err != nil {
		if os.IsNotExist(err) {
			c.Status(http.StatusNotFound)
			return
		}
	}

	m := createMeta(cfg, c, fs, pathname, fi)

	if m.IsDir {
		ch, rm, _ := getDir(cfg, c, fs, pathname)

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
