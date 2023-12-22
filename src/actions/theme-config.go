package actions

import (
	"net/http"
	"os"

	"dxkite.cn/explore-me/src/core/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func ThemeConfig(c *gin.Context) {
	cfg := config.GetConfig()
	f := cfg.ThemeConfig
	data, err := os.ReadFile(f)

	v := map[string]interface{}{}

	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, v)
			return
		}
		c.JSON(http.StatusInternalServerError, Error{InternalError, err.Error()})
		return
	}

	if err := yaml.Unmarshal(data, &v); err != nil {
		c.JSON(http.StatusInternalServerError, Error{InternalError, err.Error()})
		return
	}

	c.JSON(http.StatusOK, v)
}
