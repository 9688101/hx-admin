package model

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/9688101/hx-admin/common/logger"
	"github.com/9688101/hx-admin/global"
	// billingratio "github.com/songquanpeng/one-api/relay/billing/ratio"
)

type Option struct {
	Key   string `json:"key" gorm:"primaryKey"`
	Value string `json:"value"`
}

func AllOption() ([]*Option, error) {
	var options []*Option
	var err error
	err = DB.Find(&options).Error
	return options, err
}

func InitOptionMap() {
	global.OptionMapRWMutex.Lock()
	global.OptionMap = make(map[string]string)
	global.OptionMap["PasswordLoginEnabled"] = strconv.FormatBool(global.PasswordLoginEnabled)
	global.OptionMap["PasswordRegisterEnabled"] = strconv.FormatBool(global.PasswordRegisterEnabled)
	global.OptionMap["EmailVerificationEnabled"] = strconv.FormatBool(global.EmailVerificationEnabled)
	global.OptionMap["GitHubOAuthEnabled"] = strconv.FormatBool(global.GitHubOAuthEnabled)
	global.OptionMap["OidcEnabled"] = strconv.FormatBool(global.OidcEnabled)
	global.OptionMap["WeChatAuthEnabled"] = strconv.FormatBool(global.WeChatAuthEnabled)
	global.OptionMap["TurnstileCheckEnabled"] = strconv.FormatBool(global.TurnstileCheckEnabled)
	global.OptionMap["RegisterEnabled"] = strconv.FormatBool(global.RegisterEnabled)
	global.OptionMap["AutomaticDisableChannelEnabled"] = strconv.FormatBool(global.AutomaticDisableChannelEnabled)
	global.OptionMap["AutomaticEnableChannelEnabled"] = strconv.FormatBool(global.AutomaticEnableChannelEnabled)
	global.OptionMap["ApproximateTokenEnabled"] = strconv.FormatBool(global.ApproximateTokenEnabled)
	global.OptionMap["LogConsumeEnabled"] = strconv.FormatBool(global.LogConsumeEnabled)
	global.OptionMap["DisplayInCurrencyEnabled"] = strconv.FormatBool(global.DisplayInCurrencyEnabled)
	global.OptionMap["DisplayTokenStatEnabled"] = strconv.FormatBool(global.DisplayTokenStatEnabled)
	global.OptionMap["ChannelDisableThreshold"] = strconv.FormatFloat(global.ChannelDisableThreshold, 'f', -1, 64)
	global.OptionMap["EmailDomainRestrictionEnabled"] = strconv.FormatBool(global.EmailDomainRestrictionEnabled)
	global.OptionMap["EmailDomainWhitelist"] = strings.Join(global.EmailDomainWhitelist, ",")
	global.OptionMap["SMTPServer"] = ""
	global.OptionMap["SMTPFrom"] = ""
	global.OptionMap["SMTPPort"] = strconv.Itoa(global.SMTPPort)
	global.OptionMap["SMTPAccount"] = ""
	global.OptionMap["SMTPToken"] = ""
	global.OptionMap["Notice"] = ""
	global.OptionMap["About"] = ""
	global.OptionMap["HomePageContent"] = ""
	global.OptionMap["Footer"] = global.Footer
	global.OptionMap["SystemName"] = global.SystemName
	global.OptionMap["Logo"] = global.Logo
	global.OptionMap["ServerAddress"] = ""
	global.OptionMap["GitHubClientId"] = ""
	global.OptionMap["GitHubClientSecret"] = ""
	global.OptionMap["WeChatServerAddress"] = ""
	global.OptionMap["WeChatServerToken"] = ""
	global.OptionMap["WeChatAccountQRCodeImageURL"] = ""
	global.OptionMap["MessagePusherAddress"] = ""
	global.OptionMap["MessagePusherToken"] = ""
	global.OptionMap["TurnstileSiteKey"] = ""
	global.OptionMap["TurnstileSecretKey"] = ""
	global.OptionMap["QuotaForNewUser"] = strconv.FormatInt(global.QuotaForNewUser, 10)
	global.OptionMap["QuotaForInviter"] = strconv.FormatInt(global.QuotaForInviter, 10)
	global.OptionMap["QuotaForInvitee"] = strconv.FormatInt(global.QuotaForInvitee, 10)
	global.OptionMap["QuotaRemindThreshold"] = strconv.FormatInt(global.QuotaRemindThreshold, 10)
	global.OptionMap["PreConsumedQuota"] = strconv.FormatInt(global.PreConsumedQuota, 10)
	// global.OptionMap["ModelRatio"] = billingratio.ModelRatio2JSONString()
	// global.OptionMap["GroupRatio"] = billingratio.GroupRatio2JSONString()
	global.OptionMap["CompletionRatio"] = CompletionRatio2JSONString()
	global.OptionMap["TopUpLink"] = global.TopUpLink
	global.OptionMap["ChatLink"] = global.ChatLink
	global.OptionMap["QuotaPerUnit"] = strconv.FormatFloat(global.QuotaPerUnit, 'f', -1, 64)
	global.OptionMap["RetryTimes"] = strconv.Itoa(global.RetryTimes)
	global.OptionMap["Theme"] = global.Theme
	global.OptionMapRWMutex.Unlock()
	loadOptionsFromDatabase()
}

func loadOptionsFromDatabase() {
	options, _ := AllOption()
	for _, option := range options {
		// if option.Key == "ModelRatio" {
		// 	// option.Value = billingratio.AddNewMissingRatio(option.Value)
		// }
		err := updateOptionMap(option.Key, option.Value)
		if err != nil {
			logger.SysError("failed to update option map: " + err.Error())
		}
	}
}

func SyncOptions(frequency int) {
	for {
		time.Sleep(time.Duration(frequency) * time.Second)
		logger.SysLog("syncing options from database")
		loadOptionsFromDatabase()
	}
}

func UpdateOption(key string, value string) error {
	// Save to database first
	option := Option{
		Key: key,
	}
	// https://gorm.io/docs/update.html#Save-All-Fields
	DB.FirstOrCreate(&option, Option{Key: key})
	option.Value = value
	// Save is a combination function.
	// If save value does not contain primary key, it will execute Create,
	// otherwise it will execute Update (with all fields).
	DB.Save(&option)
	// Update OptionMap
	return updateOptionMap(key, value)
}

func updateOptionMap(key string, value string) (err error) {
	global.OptionMapRWMutex.Lock()
	defer global.OptionMapRWMutex.Unlock()
	global.OptionMap[key] = value
	if strings.HasSuffix(key, "Enabled") {
		boolValue := value == "true"
		switch key {
		case "PasswordRegisterEnabled":
			global.PasswordRegisterEnabled = boolValue
		case "PasswordLoginEnabled":
			global.PasswordLoginEnabled = boolValue
		case "EmailVerificationEnabled":
			global.EmailVerificationEnabled = boolValue
		case "GitHubOAuthEnabled":
			global.GitHubOAuthEnabled = boolValue
		case "OidcEnabled":
			global.OidcEnabled = boolValue
		case "WeChatAuthEnabled":
			global.WeChatAuthEnabled = boolValue
		case "TurnstileCheckEnabled":
			global.TurnstileCheckEnabled = boolValue
		case "RegisterEnabled":
			global.RegisterEnabled = boolValue
		case "EmailDomainRestrictionEnabled":
			global.EmailDomainRestrictionEnabled = boolValue
		case "AutomaticDisableChannelEnabled":
			global.AutomaticDisableChannelEnabled = boolValue
		case "AutomaticEnableChannelEnabled":
			global.AutomaticEnableChannelEnabled = boolValue
		case "ApproximateTokenEnabled":
			global.ApproximateTokenEnabled = boolValue
		case "LogConsumeEnabled":
			global.LogConsumeEnabled = boolValue
		case "DisplayInCurrencyEnabled":
			global.DisplayInCurrencyEnabled = boolValue
		case "DisplayTokenStatEnabled":
			global.DisplayTokenStatEnabled = boolValue
		}
	}
	switch key {
	case "EmailDomainWhitelist":
		global.EmailDomainWhitelist = strings.Split(value, ",")
	case "SMTPServer":
		global.SMTPServer = value
	case "SMTPPort":
		intValue, _ := strconv.Atoi(value)
		global.SMTPPort = intValue
	case "SMTPAccount":
		global.SMTPAccount = value
	case "SMTPFrom":
		global.SMTPFrom = value
	case "SMTPToken":
		global.SMTPToken = value
	case "ServerAddress":
		global.ServerAddress = value
	case "GitHubClientId":
		global.GitHubClientId = value
	case "GitHubClientSecret":
		global.GitHubClientSecret = value
	case "LarkClientId":
		global.LarkClientId = value
	case "LarkClientSecret":
		global.LarkClientSecret = value
	case "OidcClientId":
		global.OidcClientId = value
	case "OidcClientSecret":
		global.OidcClientSecret = value
	case "OidcWellKnown":
		global.OidcWellKnown = value
	case "OidcAuthorizationEndpoint":
		global.OidcAuthorizationEndpoint = value
	case "OidcTokenEndpoint":
		global.OidcTokenEndpoint = value
	case "OidcUserinfoEndpoint":
		global.OidcUserinfoEndpoint = value
	case "Footer":
		global.Footer = value
	case "SystemName":
		global.SystemName = value
	case "Logo":
		global.Logo = value
	case "WeChatServerAddress":
		global.WeChatServerAddress = value
	case "WeChatServerToken":
		global.WeChatServerToken = value
	case "WeChatAccountQRCodeImageURL":
		global.WeChatAccountQRCodeImageURL = value
	case "MessagePusherAddress":
		global.MessagePusherAddress = value
	case "MessagePusherToken":
		global.MessagePusherToken = value
	case "TurnstileSiteKey":
		global.TurnstileSiteKey = value
	case "TurnstileSecretKey":
		global.TurnstileSecretKey = value
	case "QuotaForNewUser":
		global.QuotaForNewUser, _ = strconv.ParseInt(value, 10, 64)
	case "QuotaForInviter":
		global.QuotaForInviter, _ = strconv.ParseInt(value, 10, 64)
	case "QuotaForInvitee":
		global.QuotaForInvitee, _ = strconv.ParseInt(value, 10, 64)
	case "QuotaRemindThreshold":
		global.QuotaRemindThreshold, _ = strconv.ParseInt(value, 10, 64)
	case "PreConsumedQuota":
		global.PreConsumedQuota, _ = strconv.ParseInt(value, 10, 64)
	case "RetryTimes":
		global.RetryTimes, _ = strconv.Atoi(value)
	// case "ModelRatio":
	// 	err = billingratio.UpdateModelRatioByJSONString(value)
	// case "GroupRatio":
	// 	err = billingratio.UpdateGroupRatioByJSONString(value)
	case "CompletionRatio":
		err = UpdateCompletionRatioByJSONString(value)
	case "TopUpLink":
		global.TopUpLink = value
	case "ChatLink":
		global.ChatLink = value
	case "ChannelDisableThreshold":
		global.ChannelDisableThreshold, _ = strconv.ParseFloat(value, 64)
	case "QuotaPerUnit":
		global.QuotaPerUnit, _ = strconv.ParseFloat(value, 64)
	case "Theme":
		global.Theme = value
	}
	return err
}

var CompletionRatio = map[string]float64{
	// aws llama3
	"llama3-8b-8192(33)":  0.0006 / 0.0003,
	"llama3-70b-8192(33)": 0.0035 / 0.00265,
	// whisper
	"whisper-1": 0, // only count input tokens
	// deepseek
	"deepseek-chat":     0.28 / 0.14,
	"deepseek-reasoner": 2.19 / 0.55,
}

func CompletionRatio2JSONString() string {
	jsonBytes, err := json.Marshal(CompletionRatio)
	if err != nil {
		logger.SysError("error marshalling completion ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateCompletionRatioByJSONString(jsonStr string) error {
	CompletionRatio = make(map[string]float64)
	return json.Unmarshal([]byte(jsonStr), &CompletionRatio)
}
