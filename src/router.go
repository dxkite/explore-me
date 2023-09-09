package src

import (
	"net/http"

	"dxkite.cn/explorer/src/actions"
	"dxkite.cn/explorer/src/core/config"
	"dxkite.cn/explorer/src/core/storage"
	goget "dxkite.cn/explorer/src/middleware/go-get"
	"dxkite.cn/explorer/static"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {

	r := gin.Default()
	//获取文件元信息
	r.GET("/api/explore/meta/*path", actions.Meta)

	//获取标签信息
	r.GET("/api/explore/tags", actions.Tags)

	//获取扩展信息
	r.GET("/api/explore/exts", actions.Exts)

	//搜索文件
	r.GET("/api/explore/search", actions.Search)

	// 获取原始文件内容
	r.StaticFS("/api/explore/raw", http.Dir(cfg.SrcRoot))

	mtx := http.NewServeMux()

	// API
	mtx.Handle("/api/", r.Handler())

	// web根目录
	webStatic := storage.NewPrefix("/dist", http.FS(static.Web))

	// 配置了web根目录
	if cfg.WebRoot != "" {
		webRoot := http.FileSystem(http.Dir(cfg.WebRoot))
		webStatic = storage.NewMultiFileSystem(webStatic, webRoot)
	}

	// 单页应用
	if cfg.SingleIndex != "" {
		webStatic = storage.NewSingleIndex(webStatic, cfg.SingleIndex)
	}

	mtx.Handle("/", goget.Middleware(func() *goget.PackageConfig {
		return &config.GetConfig().GoGetConfig
	}, http.FileServer(webStatic)))

	return http.ListenAndServe(cfg.Listen, mtx)
}
