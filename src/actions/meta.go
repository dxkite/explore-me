package actions

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"dxkite.cn/explorer/src/core"
	"github.com/gin-gonic/gin"
)

func Meta(cfg *core.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		p := c.Param("path")
		pathname := path.Join(cfg.SrcRoot, p)

		log.Println(pathname)

		fi, err := os.Stat(pathname)
		if err != nil {
			if os.IsNotExist(err) {
				c.Status(http.StatusNotFound)
				return
			}
		}

		absRoot, _ := filepath.Abs(cfg.SrcRoot)
		absPathname, _ := filepath.Abs(pathname)
		if !strings.HasPrefix(absPathname, absRoot) {
			c.Status(http.StatusNotFound)
			return
		}

		m := createMeta(cfg, pathname, fi)

		if m.IsDir {
			if ch, err := getDir(cfg, pathname); err != nil {
				log.Println(err.Error())
			} else {
				m.Children = ch
			}
		}
		c.JSON(http.StatusOK, m)
	}
}

type MetaData struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Tags     []string    `json:"tags,omitempty"`
	Ext      string      `json:"ext,omitempty"`
	IsDir    bool        `json:"is_dir"`
	Readme   string      `json:"readme"`
	Children []*MetaData `json:"children,omitempty"`
}

func createMeta(cfg *core.Config, pathname string, fi fs.FileInfo) *MetaData {
	m := &MetaData{}
	m.Name = fi.Name()
	m.Path = core.NormalizePath(cfg.SrcRoot, pathname)
	m.Tags, _ = core.ParseTag(cfg, m.Name)
	m.Ext = core.GetExt(m.Name)
	m.IsDir = fi.IsDir()

	if m.IsDir {
		rm := path.Join(pathname, cfg.ScanConfig.ReadmeFile)
		if isExist(rm) {
			m.Readme = core.NormalizePath(cfg.SrcRoot, rm)
		}
	} else {
		rm := pathname + "." + cfg.ScanConfig.ReadmeFile
		if isExist(rm) {
			m.Readme = core.NormalizePath(cfg.SrcRoot, rm)
		}
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

func getDir(cfg *core.Config, dirname string) ([]*MetaData, error) {
	rd, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	md := []*MetaData{}
	for _, de := range rd {
		pathname := path.Join(dirname, de.Name())
		fi, _ := de.Info()
		mdi := createMeta(cfg, pathname, fi)
		md = append(md, mdi)
	}
	return md, nil
}
