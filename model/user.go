package model

import (
	"encoding/json"
)

// User if you add sensitive fields, don't forget to clean them in setupLogin function.
// Otherwise, the sensitive information will be saved on local storage in plain text!
type User struct {
	Id               int    `json:"id"`
	Username         string `json:"username" gorm:"unique;index" validate:"max=12"`
	Password         string `json:"password" gorm:"not null;" validate:"min=8,max=20"`
	DisplayName      string `json:"display_name" gorm:"index" validate:"max=20"`
	Role             int    `json:"role" gorm:"type:int;default:1"`   // admin, util
	Status           int    `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Email            string `json:"email" gorm:"index" validate:"max=50"`
	GitHubId         string `json:"github_id" gorm:"column:github_id;index"`
	WeChatId         string `json:"wechat_id" gorm:"column:wechat_id;index"`
	LarkId           string `json:"lark_id" gorm:"column:lark_id;index"`
	OidcId           string `json:"oidc_id" gorm:"column:oidc_id;index"`
	VerificationCode string `json:"verification_code" gorm:"-:all"`                                    // this field is only for Email verification, don't save it to database!
	AccessToken      string `json:"access_token" gorm:"type:char(32);column:access_token;uniqueIndex"` // this token is for system management
	// Quota            int64  `json:"quota" gorm:"bigint;default:0"`
	UsedQuota int64 `json:"used_quota" gorm:"bigint;default:0;column:used_quota"` // used quota
	// RequestCount int    `json:"request_count" gorm:"type:int;default:0;"`             // request number
	// Group        string `json:"group" gorm:"type:varchar(32);default:'default'"`
	AffCode   string `json:"aff_code" gorm:"type:varchar(32);column:aff_code;uniqueIndex"`
	InviterId int    `json:"inviter_id" gorm:"type:int;column:inviter_id;index"`
}

func NewUser() *User {
	return &User{}
}

func NewUserById(id int) *User {
	return &User{Id: id}
}
func NewUserByUsername(key string) *User {
	return &User{Username: key}

}
func NewUserByEmail(key string) *User {
	return &User{Email: key}

}
func NewUserByGitHubId(key string) *User {
	return &User{GitHubId: key}
}
func NewUserByWeChatId(key string) *User {
	return &User{WeChatId: key}

}
func NewUserByLarkId(key string) *User {
	return &User{LarkId: key}
}
func NewUserByOidcId(key string) *User {
	return &User{Email: key}

}

// ToJSON 将 User 转换为 JSON 字符串
func (u *User) ToJSON() (string, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从 JSON 解析为 User
func FromJSON(data string) (*User, error) {
	var user User
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
