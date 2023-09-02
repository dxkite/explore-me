package actions

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"sort"

	"dxkite.cn/explorer/src/core"
	"github.com/gin-gonic/gin"
)

type TagItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

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

		vv := createTagList(v)
		c.JSON(http.StatusOK, vv)
	}
}

func createTagList(v map[string]int) []TagItem {
	t := []TagItem{}
	for v, c := range v {
		t = append(t, TagItem{Name: v, Count: c})
	}
	sort.Slice(t, func(i, j int) bool {
		return t[i].Name > t[j].Name
	})
	return t
}
