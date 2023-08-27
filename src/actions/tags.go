package actions

import (
	"encoding/json"
	"net/http"
	"os"
	"path"

	"dxkite.cn/explorer/src/core"
	"github.com/gin-gonic/gin"
)

func Tags(cfg *core.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		f := path.Join(cfg.DataRoot, cfg.ScanConfig.TagListFile)
		data, err := os.ReadFile(f)

		v := map[string]int{}

		if err != nil {
			if os.IsNotExist(err) {
				c.JSON(http.StatusOK, v)
				return
			}
			c.JSON(http.StatusInternalServerError, Error{InternalError, err.Error()})
			return
		}

		if err := json.Unmarshal(data, &v); err != nil {
			c.JSON(http.StatusInternalServerError, Error{InternalError, err.Error()})
			return
		}

		c.JSON(http.StatusOK, v)
	}
}
