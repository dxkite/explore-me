package actions

import (
	"encoding/json"
	"net/http"
	"os"
	"path"

	"dxkite.cn/explorer/src/core/config"
	"github.com/gin-gonic/gin"
)

func Exts(c *gin.Context) {
	cfg := config.GetConfig()

	f := path.Join(cfg.DataRoot, cfg.ScanConfig.ExtListFile)
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

	vv := createMapItemArray(v)
	c.JSON(http.StatusOK, vv)
}
