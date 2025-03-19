package router

// import (
// 	"github.com/gin-contrib/gzip"       // 引入 gzip 中间件，用于响应压缩，提高数据传输效率
// 	"github.com/gin-gonic/gin"          // 引入 Gin 框架，用于处理 HTTP 请求
// 	// "github.com/songquanpeng/one-api/controller" // 控制器模块（当前被注释掉，可能需要手动解除注释）
// 	"github.com/songquanpeng/one-api/middleware" // 中间件模块，用于权限控制、限流等功能
// )

// // SetDashboardRouter 配置仪表盘相关的 API 路由
// func SetDashboardRouter(router *gin.Engine) {
// 	// 创建根路由组（"/"），用于仪表盘相关的 API
// 	apiRouter := router.Group("/")

// 	// 启用 CORS（跨域资源共享）中间件，允许跨域请求
// 	apiRouter.Use(middleware.CORS())

// 	// 启用 gzip 压缩，提高数据传输效率
// 	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))

// 	// 启用全局 API 速率限制，防止滥用
// 	apiRouter.Use(middleware.GlobalAPIRateLimit())

// 	// 启用 Token 认证，确保访问 API 时需要提供有效的 Token
// 	apiRouter.Use(middleware.TokenAuth())

// 	{
// 		// 获取订阅信息（旧版本接口）
// 		apiRouter.GET("/dashboard/billing/subscription", controller.GetSubscription)

// 		// 获取订阅信息（v1 版本接口）
// 		apiRouter.GET("/v1/dashboard/billing/subscription", controller.GetSubscription)

// 		// 获取账单使用情况（旧版本接口）
// 		apiRouter.GET("/dashboard/billing/usage", controller.GetUsage)

// 		// 获取账单使用情况（v1 版本接口）
// 		apiRouter
