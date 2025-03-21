package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/9688101/hx-admin/core/i18n"
	"github.com/9688101/hx-admin/global"
	"github.com/9688101/hx-admin/server"
	"github.com/9688101/hx-admin/utils"
	"github.com/9688101/hx-admin/utils/message"
	"github.com/gin-gonic/gin"
)

func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"version":            global.Version,
			"start_time":         global.StartTime,
			"email_verification": global.EmailVerificationEnabled,
			"github_oauth":       global.GitHubOAuthEnabled,
			"github_client_id":   global.GitHubClientId,
			"lark_client_id":     global.LarkClientId,
			"system_name":        global.SystemName,
			"logo":               global.Logo,
			"footer_html":        global.Footer,
			"wechat_qrcode":      global.WeChatAccountQRCodeImageURL,
			"wechat_login":       global.WeChatAuthEnabled,
			"server_address":     global.ServerAddress,
			"turnstile_check":    global.TurnstileCheckEnabled,
			"turnstile_site_key": global.TurnstileSiteKey,
			// "top_up_link":                 global.TopUpLink,
			"chat_link": global.ChatLink,
			// "quota_per_unit":              global.QuotaPerUnit,
			"display_in_currency":         global.DisplayInCurrencyEnabled,
			"oidc":                        global.OidcEnabled,
			"oidc_client_id":              global.OidcClientId,
			"oidc_well_known":             global.OidcWellKnown,
			"oidc_authorization_endpoint": global.OidcAuthorizationEndpoint,
			"oidc_token_endpoint":         global.OidcTokenEndpoint,
			"oidc_userinfo_endpoint":      global.OidcUserinfoEndpoint,
		},
	})
	return
}

func GetNotice(c *gin.Context) {
	global.OptionMapRWMutex.RLock()
	defer global.OptionMapRWMutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    global.OptionMap["Notice"],
	})
	return
}

func GetAbout(c *gin.Context) {
	global.OptionMapRWMutex.RLock()
	defer global.OptionMapRWMutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    global.OptionMap["About"],
	})
	return
}

func GetHomePageContent(c *gin.Context) {
	global.OptionMapRWMutex.RLock()
	defer global.OptionMapRWMutex.RUnlock()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    global.OptionMap["HomePageContent"],
	})
	return
}

func SendEmailVerification(c *gin.Context) {
	email := c.Query("email")
	if err := utils.Validate.Var(email, "required,email"); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": i18n.Translate(c, "invalid_parameter"),
		})
		return
	}
	if global.EmailDomainRestrictionEnabled {
		allowed := false
		for _, domain := range global.EmailDomainWhitelist {
			if strings.HasSuffix(email, "@"+domain) {
				allowed = true
				break
			}
		}
		if !allowed {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "管理员启用了邮箱域名白名单，您的邮箱地址的域名不在白名单中",
			})
			return
		}
	}
	if server.IsEmailAlreadyTaken(email) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "邮箱地址已被占用",
		})
		return
	}
	code := utils.GenerateVerificationCode(6)
	utils.RegisterVerificationCodeWithKey(email, code, utils.EmailVerificationPurpose)
	subject := fmt.Sprintf("%s 邮箱验证邮件", global.SystemName)
	content := message.EmailTemplate(
		subject,
		fmt.Sprintf(`
			<p>您好！</p>
			<p>您正在进行 %s 邮箱验证。</p>
			<p>您的验证码为：</p>
			<p style="font-size: 24px; font-weight: bold; color: #333; background-color: #f8f8f8; padding: 10px; text-align: center; border-radius: 4px;">%s</p>
			<p style="color: #666;">验证码 %d 分钟内有效，如果不是本人操作，请忽略。</p>
		`, global.SystemName, code, utils.VerificationValidMinutes),
	)
	err := message.SendEmail(subject, email, content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

func SendPasswordResetEmail(c *gin.Context) {
	email := c.Query("email")
	if err := utils.Validate.Var(email, "required,email"); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": i18n.Translate(c, "invalid_parameter"),
		})
		return
	}
	if !server.IsEmailAlreadyTaken(email) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "该邮箱地址未注册",
		})
		return
	}
	code := utils.GenerateVerificationCode(0)
	utils.RegisterVerificationCodeWithKey(email, code, utils.PasswordResetPurpose)
	link := fmt.Sprintf("%s/user/reset?email=%s&token=%s", global.ServerAddress, email, code)
	subject := fmt.Sprintf("%s 密码重置", global.SystemName)
	content := message.EmailTemplate(
		subject,
		fmt.Sprintf(`
			<p>您好！</p>
			<p>您正在进行 %s 密码重置。</p>
			<p>请点击下面的按钮进行密码重置：</p>
			<p style="text-align: center; margin: 30px 0;">
				<a href="%s" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; display: inline-block;">重置密码</a>
			</p>
			<p style="color: #666;">如果按钮无法点击，请复制以下链接到浏览器中打开：</p>
			<p style="background-color: #f8f8f8; padding: 10px; border-radius: 4px; word-break: break-all;">%s</p>
			<p style="color: #666;">重置链接 %d 分钟内有效，如果不是本人操作，请忽略。</p>
		`, global.SystemName, link, link, utils.VerificationValidMinutes),
	)
	err := message.SendEmail(subject, email, content)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": fmt.Sprintf("%s%s", i18n.Translate(c, "send_email_failed"), err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

type PasswordResetRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func ResetPassword(c *gin.Context) {
	var req PasswordResetRequest
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if req.Email == "" || req.Token == "" {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": i18n.Translate(c, "invalid_parameter"),
		})
		return
	}
	if !utils.VerifyCodeWithKey(req.Email, req.Token, utils.PasswordResetPurpose) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "重置链接非法或已过期",
		})
		return
	}
	password := utils.GenerateVerificationCode(12)
	err = server.ResetUserPasswordByEmail(req.Email, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	utils.DeleteKey(req.Email, utils.PasswordResetPurpose)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    password,
	})
	return
}
