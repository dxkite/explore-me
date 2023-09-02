package src

import (
	"net/http"

	"dxkite.cn/explorer/src/actions"
	"dxkite.cn/explorer/src/core"
	goget "dxkite.cn/explorer/src/middleware/go-get"
	"github.com/gin-gonic/gin"
)

func Run(cfg *core.Config) error {

	r := gin.Default()
	//获取文件元信息
	r.GET("/api/explore/meta/*path", actions.Meta(cfg))

	//获取标签信息
	r.GET("/api/explore/tags", actions.Tags(cfg))

	//获取扩展信息
	r.GET("/api/explore/exts", actions.Exts(cfg))

	//搜索文件
	r.GET("/api/explore/search", actions.Search(cfg))

	// 获取原始文件内容
	r.StaticFS("/api/explore/raw", http.Dir(cfg.SrcRoot))

	mtx := http.NewServeMux()

	// API
	mtx.Handle("/api/", r.Handler())

	// web根目录
	mtx.Handle("/", goget.Middleware(func() *goget.PackageConfig {
		return &core.GetConfig().GoGetConfig
	}, http.FileServer(http.Dir(cfg.WebRoot))))

	return http.ListenAndServe(cfg.Listen, mtx)
}
