package actions

import (
	"context"
	"log"
	"net/http"
	"os"
	"path"

	"dxkite.cn/explorer/src/core"
	"dxkite.cn/explorer/src/core/config"
	"dxkite.cn/explorer/src/core/storage"
	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Path string `form:"path"`
	Name string `form:"name"`
	Ext  string `form:"ext"`
	Tag  string `form:"tag"`

	Offset int64 `form:"offset"`
	Limit  int64 `form:"limit"`
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
		Path: req.Path,
	}

	log.Println("search", param)

	rst, err := core.SearchFile(idx, param, req.Offset, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Error{Code: InternalError, Message: err.Error()})
		return
	}

	mdl := createMetaList(cfg, rst)
	c.JSON(http.StatusOK, mdl)
}

func createMetaList(cfg *config.Config, fia []*core.SearchFileInfo) []*MetaData {
	md := []*MetaData{}
	src := storage.Local(cfg.SrcRoot)

	for _, f := range fia {
		filename := path.Join(cfg.SrcRoot, f.Path)

		fi, err := os.Stat(filename)
		if err != nil {
			log.Println("createMetaList:Stat", filename, err)
			continue
		}

		mdi := createMeta(cfg, context.TODO(), src, filename, fi)
		mdi.Id = f.Id
		md = append(md, mdi)
	}

	return md
}
