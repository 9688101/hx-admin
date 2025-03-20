package router

import (
	"embed" // 用于嵌入静态资源文件
	"fmt"   // 格式化字符串
	"strings"

	"github.com/9688101/hx-admin/controller"
	"github.com/9688101/hx-admin/global" // 引入配置模块
	"github.com/9688101/hx-admin/utils"
	"github.com/gin-contrib/gzip"   // Gin 的 gzip 中间件，用于压缩 HTTP 响应
	"github.com/gin-contrib/static" // Gin 的静态文件服务中间件
	"github.com/gin-gonic/gin"      // Gin 框架，用于处理 HTTP 请求

	"net/http" // 提供 HTTP 相关功能

	"github.com/9688101/hx-admin/middleware" // 引入中间件模块
	// "strings" // 用于字符串操作（当前被注释掉，可能用于路径判断）
)

// SetWebRouter 配置 Web 前端的路由
func SetWebRouter(router *gin.Engine, buildFS embed.FS) {
	// 读取前端 `index.html` 文件内容，作为前端入口文件
	indexPageData, _ := buildFS.ReadFile(fmt.Sprintf("web/build/%s/index.html", global.Theme))

	// 启用 gzip 压缩，提高数据传输效率
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// 启用 Web 全局速率限制，防止恶意请求
	router.Use(middleware.GlobalWebRateLimit())

	// 启用缓存中间件，提高访问性能
	router.Use(middleware.Cache())

	// 配置静态文件服务，将前端 `web/build/{主题名}` 目录中的文件嵌入到 Gin 路由中
	router.Use(static.Serve("/", utils.EmbedFolder(buildFS, fmt.Sprintf("web/build/%s", "default"))))

	// 配置默认的 404 处理（未匹配到的路由）
	router.NoRoute(func(c *gin.Context) {
		// 如果请求路径以 `/v1` 或 `/api` 开头，则可能是 API 请求，返回自定义的 404 处理（目前被注释）
		if strings.HasPrefix(c.Request.RequestURI, "/v1") || strings.HasPrefix(c.Request.RequestURI, "/api") {
			controller.RelayNotFound(c)
			return
		}

		// 取消缓存，确保前端文件最新
		c.Header("Cache-Control", "no-cache")

		// 返回 `index.html`，支持前端路由（前端框架如 React/Vue 通常使用前端路由）
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPageData)
	})
}
