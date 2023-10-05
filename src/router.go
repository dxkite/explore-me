package src

import (
	"net/http"
	"path"

	"dxkite.cn/explorer/src/actions"
	"dxkite.cn/explorer/src/core/client"
	"dxkite.cn/explorer/src/core/config"
	"dxkite.cn/explorer/src/core/storage"
	"dxkite.cn/explorer/src/middleware/clientid"
	goget "dxkite.cn/explorer/src/middleware/go-get"
	"dxkite.cn/explorer/static"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

// theme -> web_root -> inner static
func createFs(cfg *config.Config) http.FileSystem {

	webFs := []http.FileSystem{}

	// 配置了web根目录
	if cfg.Theme != "" {
		themeRoot := path.Join(cfg.ThemeRoot, cfg.Theme)
		themeFs := http.FileSystem(http.Dir(themeRoot))
		webFs = append(webFs, themeFs)
	}

	// 配置了web根目录
	if cfg.WebRoot != "" {
		webRoot := http.FileSystem(http.Dir(cfg.WebRoot))
		webFs = append(webFs, webRoot)
	}

	// web根目录
	webStatic := storage.NewPrefix("/dist", http.FS(static.Web))
	webFs = append(webFs, webStatic)

	return storage.NewMultiFileSystem(webFs...)
}

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
	r.StaticFS(config.RawUrlRoot, http.Dir(cfg.SrcRoot))

	mtx := http.NewServeMux()

	clientPool := client.NewClientPool()
	clientPool.GetClientId = func(c *websocket.Conn) string {
		return c.Request().Header.Get(cfg.ClientIdKey)
	}

	// API
	mtx.Handle("/api/", r.Handler())
	// WebSocket
	mtx.Handle("/api/websocket/client", websocket.Handler(clientPool.HandleClient))

	// 目录读取
	webStatic := createFs(cfg)
	// 单页应用
	webStatic = storage.NewSingleIndex(webStatic, cfg.SingleIndex)
	mtx.Handle("/", goget.Middleware(func() *goget.PackageConfig {
		return &config.GetConfig().GoGetConfig
	}, http.FileServer(webStatic)))

	return http.ListenAndServe(cfg.Listen, clientid.Middleware(mtx, cfg.ClientIdKey))
}
