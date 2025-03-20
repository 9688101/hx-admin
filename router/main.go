package router

import (
	"embed"    // 用于嵌入静态资源文件
	"fmt"      // 用于字符串格式化
	"net/http" // 提供 HTTP 相关功能
	"os"       // 处理系统环境变量
	"strings"  // 处理字符串操作

	"github.com/9688101/hx-admin/common/logger" // 引入日志模块
	"github.com/9688101/hx-admin/global"        // 引入配置模块
	"github.com/gin-gonic/gin"                  // 引入 Gin 框架，用于处理 HTTP 请求
)

// SetRouter 配置整个 Web 服务器的路由
func SetRouter(router *gin.Engine, buildFS embed.FS) {
	// 设置 API 相关的路由
	SetApiRouter(router)

	// 设置仪表盘相关的路由（当前被注释掉，可能需要手动解除）
	// SetDashboardRouter(router)

	// 设置代理相关的路由（当前被注释掉，可能用于请求转发）
	// SetRelayRouter(router)

	// 获取前端的基本 URL（通常用于前后端分离架构）
	frontendBaseUrl := os.Getenv("FRONTEND_BASE_URL")

	// 如果当前节点是主节点（Master Node），则忽略前端 URL 并记录日志
	if global.IsMasterNode && frontendBaseUrl != "" {
		frontendBaseUrl = ""
		logger.SysLog("FRONTEND_BASE_URL is ignored on master node")
	}

	// 如果 `FRONTEND_BASE_URL` 为空，则启用内置的 Web 前端路由
	if frontendBaseUrl == "" {
		SetWebRouter(router, buildFS)
	} else {
		// 处理 `FRONTEND_BASE_URL`，确保它没有尾部斜杠
		frontendBaseUrl = strings.TrimSuffix(frontendBaseUrl, "/")

		// 设置默认路由（NoRoute），将所有未匹配的请求重定向到前端地址
		router.NoRoute(func(c *gin.Context) {
			// 301 重定向到指定的前端地址
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s", frontendBaseUrl, c.Request.RequestURI))
		})
	}
}
