package actions

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"dxkite.cn/explorer/src/core"
	"dxkite.cn/explorer/src/core/config"
	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Name string `form:"name"`
	Ext  string `form:"ext"`
	Tag  string `form:"tag"`

	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

func Search(c *gin.Context) {
	cfg := config.GetConfig()

	req := &SearchRequest{}

	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, Error{Code: ParamError, Message: err.Error()})
		return
	}

	idx := path.Join(cfg.DataRoot, cfg.ScanConfig.IndexFile)

	param := core.SearchParams{
		Name: req.Name,
		Ext:  req.Ext,
		Tag:  req.Tag,
	}
	fmt.Println(param)
	rst, err := core.SearchFile(idx, param, req.Offset, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Code: InternalError, Message: err.Error()})
		return
	}

	mdl := createMetaList(cfg, rst)
	c.JSON(http.StatusOK, mdl)
}

func createMetaList(cfg *config.Config, fia []*core.FileInfo) []*MetaData {
	md := []*MetaData{}

	for _, f := range fia {
		filename := path.Join(cfg.SrcRoot, f.Path)

		fi, err := os.Stat(filename)
		if err != nil {
			log.Println("error", err)
			continue
		}

		mdi := createMeta(cfg, filename, fi)
		md = append(md, mdi)
	}

	return md
}
