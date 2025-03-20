package core

// 导入需要的包
import (
	"embed"
	"fmt"
	"os"
	"strconv"

	"github.com/9688101/hx-admin/core/i18n"
	"github.com/9688101/hx-admin/core/logger"
	"github.com/9688101/hx-admin/global"
	"github.com/9688101/hx-admin/initialize"
	"github.com/9688101/hx-admin/middleware"
	"github.com/9688101/hx-admin/router"
	"github.com/9688101/hx-admin/server"
	"github.com/9688101/hx-admin/source"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// 嵌入前端构建产物到二进制文件中
var buildFS embed.FS

// 主程序入口
func RunServer() {
	source.Init()                                        // 初始化通用配置
	logger.SetupLogger()                                 // 初始化日志系统
	logger.SysLogf("One API %s started", global.Version) // 记录启动日志

	// 设置Gin运行模式
	if os.Getenv("GIN_MODE") != gin.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	if global.DebugEnabled {
		logger.SysLog("running in debug mode") // 调试模式日志
	}

	// 初始化数据库
	initialize.InitDB()    // 主数据库初始化
	initialize.InitLogDB() // 日志数据库初始化

	// 创建根账户（如果需要）
	var err error
	err = initialize.CreateRootAccountIfNeed()
	if err != nil {
		logger.FatalLog("database init error: " + err.Error())
	}
	defer func() { // 确保程序退出时关闭数据库连接
		err := initialize.CloseDB()
		if err != nil {
			logger.FatalLog("failed to close database: " + err.Error())
		}
	}()

	// 初始化Redis客户端
	err = initialize.InitRedisClient()
	if err != nil {
		logger.FatalLog("failed to initialize Redis: " + err.Error())
	}

	// 初始化系统配置选项
	server.InitOptionMap()
	logger.SysLog(fmt.Sprintf("using theme %s", global.Theme)) // 记录主题信息

	// 配置缓存设置
	if initialize.RedisEnabled {
		global.MemoryCacheEnabled = true // Redis启用时强制开启内存缓存
	}
	if global.MemoryCacheEnabled {
		logger.SysLog("memory cache enabled")
		logger.SysLog(fmt.Sprintf("sync frequency: %d seconds", global.SyncFrequency))
	}

	// 初始化API客户端
	source.Init()

	// 初始化国际化支持
	if err := i18n.Init(); err != nil {
		logger.FatalLog("failed to initialize i18n: " + err.Error())
	}

	// 创建Gin引擎实例
	server := gin.New()
	server.Use(gin.Recovery())         // 添加崩溃恢复中间件
	server.Use(middleware.RequestId()) // 添加请求ID中间件
	server.Use(middleware.Language())  // 添加语言中间件
	middleware.SetUpLogger(server)     // 设置日志中间件

	// 配置会话存储
	store := cookie.NewStore([]byte(global.SessionSecret))
	server.Use(sessions.Sessions("session", store)) // 添加会话中间件

	router.SetRouter(server, buildFS) // 设置路由并传入前端构建文件

	// 获取并设置服务端口
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*source.Port) // 使用默认端口
	}
	logger.SysLogf("server started on http://localhost:%s", port)

	// 启动HTTP服务器
	err = server.Run(":" + port)
	if err != nil {
		logger.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
