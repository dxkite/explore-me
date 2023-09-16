package actions

import (
	"context"
	"net/http"
	"path"

	"dxkite.cn/log"

	"dxkite.cn/explorer/src/core"
	"dxkite.cn/explorer/src/core/config"
	"dxkite.cn/explorer/src/core/scan"
	"dxkite.cn/explorer/src/core/storage"
	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Path string `form:"path"`
	Name string `form:"name"`
	Ext  string `form:"ext"`
	Tag  string `form:"tag"`

	Recent bool  `form:"recent"`
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

	idx := path.Join(cfg.DataRoot, scan.MetaIndex)

	if req.Recent {
		idx = path.Join(cfg.DataRoot, scan.RecentIndex)
	}

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
		filename := f.Path
		fi, err := src.Stat(context.TODO(), filename)

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
