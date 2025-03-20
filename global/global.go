package global

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/9688101/hx-admin/utils/env" // 环境变量处理工具
	"github.com/google/uuid"                // UUID生成工具
)

// 系统基础配置
var SystemName = "One API"                  // 系统名称
var ServerAddress = "http://localhost:3000" // 服务器基础地址
var Footer = ""                             // 页面底部信息
var Logo = ""                               // 系统Logo地址
var TopUpLink = ""                          // 充值链接
var ChatLink = ""                           // 聊天页面链接
var QuotaPerUnit = 500 * 1000.0             // 每单位配额对应的token数量（对应$0.002/1K tokens）
var DisplayInCurrencyEnabled = true         // 是否显示货币单位
var DisplayTokenStatEnabled = true          // 是否显示token统计

// 主题配置
var Theme = env.String("THEME", "default") // 当前系统主题
var ValidThemes = map[string]bool{         // 有效主题列表
	"default": true, "berry": true, "air": true,
}

// 注意：包含"Secret"/"Token"的配置项不会通过GetOptions接口返回

// 会话安全配置
var SessionSecret = uuid.New().String() // 会话加密密钥（自动生成）

// 动态配置存储（需配合读写锁）
var OptionMap map[string]string   // 系统配置项键值对存储
var OptionMapRWMutex sync.RWMutex // 配置项的读写锁

// 分页与显示设置
var ItemsPerPage = 10    // 每页显示条目数
var MaxRecentItems = 100 // 最大最近记录显示数

// 认证功能开关
var PasswordLoginEnabled = true      // 密码登录开关
var PasswordRegisterEnabled = true   // 密码注册开关
var EmailVerificationEnabled = false // 邮箱验证开关
var GitHubOAuthEnabled = false       // GitHub OAuth开关
var OidcEnabled = false              // OIDC认证开关
var WeChatAuthEnabled = false        // 微信认证开关
var TurnstileCheckEnabled = false    // Turnstile验证开关
var RegisterEnabled = true           // 注册功能总开关

// 调试相关配置
var DebugEnabled = strings.ToLower(os.Getenv("DEBUG")) == "true"                      // 调试模式开关
var DebugSQLEnabled = strings.ToLower(os.Getenv("DEBUG_SQL")) == "true"               // SQL调试开关
var MemoryCacheEnabled = strings.ToLower(os.Getenv("MEMORY_CACHE_ENABLED")) == "true" // 内存缓存开关

// 日志配置
var LogConsumeEnabled = true // 日志记录开关

// 配额管理配置
var QuotaForNewUser int64 = 0              // 新用户初始配额
var QuotaForInviter int64 = 0              // 邀请人奖励配额
var QuotaForInvitee int64 = 0              // 被邀请人奖励配额
var ChannelDisableThreshold = 5.0          // 渠道自动禁用阈值（错误率）
var AutomaticDisableChannelEnabled = false // 自动禁用渠道开关
var AutomaticEnableChannelEnabled = false  // 自动启用渠道开关
var QuotaRemindThreshold int64 = 1000      // 配额提醒阈值
var PreConsumedQuota int64 = 500           // 预消耗配额
var ApproximateTokenEnabled = false        // 启用近似token计算
var RetryTimes = 0                         // 请求重试次数

// 根用户配置
var RootUserEmail = "" // 根用户邮箱

// 节点类型配置
var IsMasterNode = os.Getenv("NODE_TYPE") != "slave" // 是否为主节点

// 请求间隔配置
var requestInterval, _ = strconv.Atoi(os.Getenv("POLLING_INTERVAL"))
var RequestInterval = time.Duration(requestInterval) * time.Second // 轮询间隔时间

// 数据同步配置
var SyncFrequency = env.Int("SYNC_FREQUENCY", 10*60) // 数据同步频率（秒）

// 批量更新配置
var BatchUpdateEnabled = false                                // 批量更新开关
var BatchUpdateInterval = env.Int("BATCH_UPDATE_INTERVAL", 5) // 批量更新间隔（秒）

// 请求超时配置
var RelayTimeout = env.Int("RELAY_TIMEOUT", 0) // 转发超时时间（秒）

// Gemini安全配置
var GeminiSafetySetting = env.String("GEMINI_SAFETY_SETTING", "BLOCK_NONE") // 安全设置级别

// 速率限制配置（单位：秒）
var (
	GlobalApiRateLimitNum            = env.Int("GLOBAL_API_RATE_LIMIT", 480) // API全局请求限制数
	GlobalApiRateLimitDuration int64 = 3 * 60                                // API限制时间窗口（180秒）

	GlobalWebRateLimitNum            = env.Int("GLOBAL_WEB_RATE_LIMIT", 240) // Web请求限制数
	GlobalWebRateLimitDuration int64 = 3 * 60                                // Web限制时间窗口

	UploadRateLimitNum            = 10 // 上传操作限制数
	UploadRateLimitDuration int64 = 60 // 上传限制时间窗口

	DownloadRateLimitNum            = 10 // 下载操作限制数
	DownloadRateLimitDuration int64 = 60 // 下载限制时间窗口

	CriticalRateLimitNum            = 20      // 关键操作限制数
	CriticalRateLimitDuration int64 = 20 * 60 // 关键操作限制时间窗口
)

var RateLimitKeyExpirationDuration = 20 * time.Minute // 限速key过期时间

// 监控指标配置
var EnableMetric = env.Bool("ENABLE_METRIC", false)                                // 监控开关
var MetricQueueSize = env.Int("METRIC_QUEUE_SIZE", 10)                             // 指标队列大小
var MetricSuccessRateThreshold = env.Float64("METRIC_SUCCESS_RATE_THRESHOLD", 0.8) // 成功率阈值
var MetricSuccessChanSize = env.Int("METRIC_SUCCESS_CHAN_SIZE", 1024)              // 成功指标通道大小
var MetricFailChanSize = env.Int("METRIC_FAIL_CHAN_SIZE", 128)                     // 失败指标通道大小

// 初始化令牌配置
var InitialRootToken = os.Getenv("INITIAL_ROOT_TOKEN")              // 根用户初始令牌
var InitialRootAccessToken = os.Getenv("INITIAL_ROOT_ACCESS_TOKEN") // 根用户初始访问令牌

// Gemini版本配置
var GeminiVersion = env.String("GEMINI_VERSION", "v1") // Gemini接口版本

// 日志文件配置
var OnlyOneLogFile = env.Bool("ONLY_ONE_LOG_FILE", false) // 单日志文件模式

// 代理配置
var RelayProxy = env.String("RELAY_PROXY", "")                              // 请求转发代理
var UserContentRequestProxy = env.String("USER_CONTENT_REQUEST_PROXY", "")  // 用户内容请求代理
var UserContentRequestTimeout = env.Int("USER_CONTENT_REQUEST_TIMEOUT", 30) // 用户内容请求超时（秒）

// 特殊功能配置
var EnforceIncludeUsage = env.Bool("ENFORCE_INCLUDE_USAGE", false)                                          // 强制包含用量信息
var TestPrompt = env.String("TEST_PROMPT", "Output only your specific model name with no additional text.") // 模型测试提示词
