package router

import (
	"github.com/9688101/hx-admin/controller"
	"github.com/9688101/hx-admin/controller/auth"
	"github.com/9688101/hx-admin/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// SetApiRouter 配置 API 路由
func SetApiRouter(router *gin.Engine) {
	// 创建 /api 组路由
	apiRouter := router.Group("/api")

	// 启用 gzip 压缩中间件，提高数据传输效率
	apiRouter.Use(gzip.Gzip(gzip.DefaultCompression))

	// 启用全局 API 速率限制，防止接口滥用
	apiRouter.Use(middleware.GlobalAPIRateLimit())

	{
		// 获取 API 状态
		apiRouter.GET("/status", controller.GetStatus)

		// 获取公告信息
		apiRouter.GET("/notice", controller.GetNotice)

		// 获取关于我们页面的信息
		apiRouter.GET("/about", controller.GetAbout)

		// 获取首页内容
		apiRouter.GET("/home_page_content", controller.GetHomePageContent)

		// 发送邮箱验证码，受严格限流并启用 Turnstile 机制防止滥用
		apiRouter.GET("/verification", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendEmailVerification)

		// 发送重置密码的邮件，同样受严格限流与 Turnstile 保护
		apiRouter.GET("/reset_password", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendPasswordResetEmail)

		// 处理用户密码重置请求
		apiRouter.POST("/user/reset", middleware.CriticalRateLimit(), controller.ResetPassword)

		// GitHub OAuth 登录
		apiRouter.GET("/oauth/github", middleware.CriticalRateLimit(), auth.GitHubOAuth)

		// OIDC（OpenID Connect）认证
		apiRouter.GET("/oauth/oidc", middleware.CriticalRateLimit(), auth.OidcAuth)

		// Lark OAuth 登录
		apiRouter.GET("/oauth/lark", middleware.CriticalRateLimit(), auth.LarkOAuth)

		// 生成 OAuth 状态码
		apiRouter.GET("/oauth/state", middleware.CriticalRateLimit(), auth.GenerateOAuthCode)

		// 微信 OAuth 登录
		apiRouter.GET("/oauth/wechat", middleware.CriticalRateLimit(), auth.WeChatAuth)

		// 绑定微信账号，需用户身份验证
		apiRouter.GET("/oauth/wechat/bind", middleware.CriticalRateLimit(), middleware.UserAuth(), auth.WeChatBind)

		// 绑定邮箱账号，需用户身份验证
		apiRouter.GET("/oauth/email/bind", middleware.CriticalRateLimit(), middleware.UserAuth(), controller.EmailBind)

		// 充值功能，仅限管理员，当前被注释掉
		// apiRouter.POST("/topup", middleware.AdminAuth(), controller.AdminTopUp)

		// 用户相关路由组
		userRoute := apiRouter.Group("/user")
		{
			// 用户注册，受限流和 Turnstile 保护
			userRoute.POST("/register", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.Register)

			// 用户登录，受限流保护
			userRoute.POST("/login", middleware.CriticalRateLimit(), controller.Login)

			// 用户登出
			userRoute.GET("/logout", controller.Logout)

			// 用户自管理相关路由，需登录认证
			selfRoute := userRoute.Group("/")
			selfRoute.Use(middleware.UserAuth())
			{
				// 获取用户仪表盘信息
				selfRoute.GET("/dashboard", controller.GetUserDashboard)

				// 获取当前用户信息
				selfRoute.GET("/self", controller.GetSelf)

				// 更新用户信息
				selfRoute.PUT("/self", controller.UpdateSelf)

				// 删除用户账户
				selfRoute.DELETE("/self", controller.DeleteSelf)

				// 生成访问令牌
				selfRoute.GET("/token", controller.GenerateAccessToken)

				// 获取推广码
				selfRoute.GET("/aff", controller.GetAffCode)

				// 充值功能（普通用户），当前被注释掉
				// selfRoute.POST("/topup", controller.TopUp)

				// 获取用户可用的模型列表，当前被注释掉
				// selfRoute.GET("/available_models", controller.GetUserAvailableModels)
			}

			// 管理员管理用户相关路由
			adminRoute := userRoute.Group("/")
			adminRoute.Use(middleware.AdminAuth())
			{
				// 获取所有用户信息
				adminRoute.GET("/", controller.GetAllUsers)

				// 搜索用户
				adminRoute.GET("/search", controller.SearchUsers)

				// 根据 ID 获取用户信息
				adminRoute.GET("/:id", controller.GetUser)

				// 创建新用户
				adminRoute.POST("/", controller.CreateUser)

				// 管理用户
				adminRoute.POST("/manage", controller.ManageUser)

				// 更新用户信息
				adminRoute.PUT("/", controller.UpdateUser)

				// 删除用户
				adminRoute.DELETE("/:id", controller.DeleteUser)
			}
		}

		// 系统配置相关路由，仅超级管理员可访问
		optionRoute := apiRouter.Group("/option")
		optionRoute.Use(middleware.RootAuth())
		{
			// 获取系统配置
			optionRoute.GET("/", controller.GetOptions)

			// 更新系统配置
			optionRoute.PUT("/", controller.UpdateOption)
		}

		// 支付渠道管理路由，当前被注释掉
		// channelRoute := apiRouter.Group("/channel")
		// channelRoute.Use(middleware.AdminAuth())
		// {
		// 	// 获取所有支付渠道
		// 	channelRoute.GET("/", controller.GetAllChannels)
		// 	// 搜索支付渠道
		// 	channelRoute.GET("/search", controller.SearchChannels)
		// 	// 获取所有支付模型
		// 	channelRoute.GET("/models", controller.ListAllModels)
		// 	// 根据 ID 获取支付渠道
		// 	channelRoute.GET("/:id", controller.GetChannel)
		// }

		// API 访问令牌管理，当前被注释掉
		// tokenRoute := apiRouter.Group("/token")
		// tokenRoute.Use(middleware.UserAuth())
		// {
		// 	// 获取所有访问令牌
		// 	tokenRoute.GET("/", controller.GetAllTokens)
		// 	// 搜索访问令牌
		// 	tokenRoute.GET("/search", controller.SearchTokens)
		// }

		// 兑换码管理，当前被注释掉
		// redemptionRoute := apiRouter.Group("/redemption")
		// redemptionRoute.Use(middleware.AdminAuth())
		// {
		// 	// 获取所有兑换码
		// 	redemptionRoute.GET("/", controller.GetAllRedemptions)
		// }

		// 日志管理，当前被注释掉
		// logRoute := apiRouter.Group("/log")
		// logRoute.GET("/", middleware.AdminAuth(), controller.GetAllLogs)

		// 用户分组管理，仅管理员可访问
		groupRoute := apiRouter.Group("/group")
		groupRoute.Use(middleware.AdminAuth())
		{
			// 获取用户组信息，当前被注释掉
			// groupRoute.GET("/", controller.GetGroups)
		}
	}
}
