package global

// SMTP邮件服务配置
var SMTPServer = ""  // SMTP服务器地址
var SMTPPort = 587   // SMTP端口号
var SMTPAccount = "" // SMTP账号
var SMTPFrom = ""    // 发件人邮箱
var SMTPToken = ""   // SMTP授权令牌

// 第三方服务配置
var GitHubClientId = ""     // GitHub OAuth客户端ID
var GitHubClientSecret = "" // GitHub OAuth客户端密钥

// Lark飞书配置
var LarkClientId = ""     // 飞书客户端ID
var LarkClientSecret = "" // 飞书客户端密钥

// OIDC配置
var OidcClientId = ""              // OIDC客户端ID
var OidcClientSecret = ""          // OIDC客户端密钥
var OidcWellKnown = ""             // OIDC发现文档地址
var OidcAuthorizationEndpoint = "" // 授权端点
var OidcTokenEndpoint = ""         // Token端点
var OidcUserinfoEndpoint = ""      // 用户信息端点

// 微信服务配置
var WeChatServerAddress = ""         // 微信服务地址
var WeChatServerToken = ""           // 微信服务令牌
var WeChatAccountQRCodeImageURL = "" // 微信公众号二维码地址

// 消息推送配置
var MessagePusherAddress = "" // 消息推送服务地址
var MessagePusherToken = ""   // 消息推送服务令牌

// Turnstile验证配置
var TurnstileSiteKey = ""   // 站点验证密钥
var TurnstileSecretKey = "" // 服务端验证密钥

// 邮箱域名限制配置
var EmailDomainRestrictionEnabled = false // 邮箱域名限制开关
var EmailDomainWhitelist = []string{      // 允许的邮箱域名列表
	"gmail.com", "163.com", "126.com", "qq.com",
	"outlook.com", "hotmail.com", "icloud.com",
	"yahoo.com", "foxmail.com",
}
