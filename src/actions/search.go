package actions

import (
	"fmt"
	"net/http"
	"path"

	"dxkite.cn/explorer/src/core"
	"github.com/gin-gonic/gin"
)

type SearchRequest struct {
	Name string `form:"name"`
	Ext  string `form:"ext"`
	Tag  string `form:"tag"`

	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

func Search(cfg *core.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
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

		c.JSON(http.StatusOK, rst)
	}
}
