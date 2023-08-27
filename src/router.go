package src

import (
	"net/http"

	"dxkite.cn/explorer/src/actions"
	"dxkite.cn/explorer/src/core"
	"github.com/gin-gonic/gin"
)

func Run(cfg *core.Config) error {
	r := gin.Default()

	// 获取原始文件内容
	r.StaticFS("/api/explore/raw", http.Dir(cfg.SrcRoot))

	//获取文件元信息
	r.GET("/api/explore/meta/*path")

	//获取标签信息
	r.GET("/api/explore/tags", actions.Tags(cfg))

	//获取扩展信息
	r.GET("/api/explore/exts", actions.Exts(cfg))

	//搜索文件
	r.GET("/api/explore/search", actions.Search(cfg))

	return r.Run(cfg.Listen)
}
