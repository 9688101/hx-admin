package server

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/9688101/hx-admin/core/logger"
	"github.com/9688101/hx-admin/initialize"
	"github.com/9688101/hx-admin/model"
	"github.com/9688101/hx-admin/utils"
)

const (
	RoleGuestUser  = 0
	RoleCommonUser = 1
	RoleAdminUser  = 10
	RoleRootUser   = 100
)

const (
	UserStatusEnabled  = 1 // don't use 0, 0 is the default value!
	UserStatusDisabled = 2 // also don't use 0
	UserStatusDeleted  = 3
)

func GetMaxUserId() int {
	u := model.NewUser()
	initialize.DB.Last(u)
	return u.Id
}

func GetAllUsers(startIdx int, num int, order string) (users []model.User, err error) {
	query := initialize.DB.Limit(num).Offset(startIdx).Omit("password").Where("status != ?", UserStatusDeleted)

	switch order {
	// case "quota":
	// 	query = query.Order("quota desc")
	// case "used_quota":
	// 	query = query.Order("used_quota desc")
	case "request_count":
		query = query.Order("request_count desc")
	default:
		query = query.Order("id desc")
	}

	err = query.Find(&users).Error
	return users, err
}

func SearchUsers(keyword string) (users []model.User, err error) {
	if !initialize.UsingPostgreSQL {
		err = initialize.DB.Omit("password").Where("id = ? or username LIKE ? or email LIKE ? or display_name LIKE ?", keyword, keyword+"%", keyword+"%", keyword+"%").Find(&users).Error
	} else {
		err = initialize.DB.Omit("password").Where("username LIKE ? or email LIKE ? or display_name LIKE ?", keyword+"%", keyword+"%", keyword+"%").Find(&users).Error
	}
	return users, err
}

func GetUserById(id int, selectAll bool) (*model.User, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	u := model.NewUserById(id)
	var err error = nil
	if selectAll {
		err = initialize.DB.First(u, "id = ?", id).Error
	} else {
		err = initialize.DB.Omit("password", "access_token").First(u, "id = ?", id).Error
	}
	return u, err
}

func GetUserIdByAffCode(affCode string) (int, error) {
	if affCode == "" {
		return 0, errors.New("affCode 为空！")
	}
	u := model.NewUser()
	err := initialize.DB.Select("id").First(u, "aff_code = ?", affCode).Error
	return u.Id, err
}

func DeleteUserById(id int) (err error) {
	if id == 0 {
		return errors.New("id 为空！")
	}
	u := model.NewUserById(id)
	return DeleteUser(u)
}

func UpdateUser(u *model.User, updatePassword bool) error {
	var err error
	if updatePassword {
		u.Password, err = utils.Password2Hash(u.Password)
		if err != nil {
			return err
		}
	}
	if u.Status == UserStatusDisabled {
		utils.BanUser(u.Id)
	} else if u.Status == UserStatusEnabled {
		utils.UnbanUser(u.Id)
	}
	err = initialize.DB.Model(u).Updates(u).Error
	return err
}

func DeleteUser(user *model.User) error {
	if user.Id == 0 {
		return errors.New("id 为空！")
	}
	utils.BanUser(user.Id)
	user.Username = fmt.Sprintf("deleted_%s", utils.GetUUID())
	user.Status = UserStatusDeleted
	err := initialize.DB.Model(user).Updates(user).Error
	return err
}
func InsertUser(ctx context.Context, user *model.User, inviterId int) error {
	var err error
	if user.Password != "" {
		user.Password, err = utils.Password2Hash(user.Password)
		if err != nil {
			return err
		}
	}
	// user.Quota = config.QuotaForNewUser
	user.AccessToken = utils.GetUUID()
	user.AffCode = utils.GetRandomString(4)
	result := initialize.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	// if config.QuotaForNewUser > 0 {
	// 	RecordLog(ctx, user.Id, LogTypeSystem, fmt.Sprintf("新用户注册赠送 %s", common.LogQuota(config.QuotaForNewUser)))
	// }
	// if inviterId != 0 {
	// 	if config.QuotaForInvitee > 0 {
	// 		_ = IncreaseUserQuota(user.Id, config.QuotaForInvitee)
	// 		RecordLog(ctx, user.Id, LogTypeSystem, fmt.Sprintf("使用邀请码赠送 %s", common.LogQuota(config.QuotaForInvitee)))
	// 	}
	// 	if config.QuotaForInviter > 0 {
	// 		_ = IncreaseUserQuota(inviterId, config.QuotaForInviter)
	// 		RecordLog(ctx, inviterId, LogTypeSystem, fmt.Sprintf("邀请用户赠送 %s", common.LogQuota(config.QuotaForInviter)))
	// 	}
	// }
	// create default token
	cleanToken := model.Token{
		UserId:       user.Id,
		Name:         "default",
		Key:          utils.GenerateKey(),
		CreatedTime:  utils.GetTimestamp(),
		AccessedTime: utils.GetTimestamp(),
		ExpiredTime:  -1,
		// RemainQuota:    -1,
		// UnlimitedQuota: true,
	}
	result.Error = InsertToken(&cleanToken)
	if result.Error != nil {
		// do not block
		logger.SysError(fmt.Sprintf("create default token for user %d failed: %s", user.Id, result.Error.Error()))
	}
	return nil
}

// ValidateAndFill check password & user status
func ValidateAndFill(user *model.User) (err error) {
	// When querying with struct, GORM will only query with non-zero fields,
	// that means if your field’s value is 0, '', false or other zero values,
	// it won’t be used to build query conditions
	password := user.Password
	if user.Username == "" || password == "" {
		return errors.New("用户名或密码为空")
	}
	err = initialize.DB.Where("username = ?", user.Username).First(user).Error
	if err != nil {
		// we must make sure check username firstly
		// consider this case: a malicious user set his username as other's email
		err := initialize.DB.Where("email = ?", user.Username).First(user).Error
		if err != nil {
			return errors.New("用户名或密码错误，或用户已被封禁")
		}
	}
	okay := utils.ValidatePasswordAndHash(password, user.Password)
	if !okay || user.Status != UserStatusEnabled {
		return errors.New("用户名或密码错误，或用户已被封禁")
	}
	return nil
}

func FillUserById(u *model.User) error {
	if u.Id == 0 {
		return errors.New("id 为空！")
	}
	initialize.DB.Where(model.NewUserById(u.Id)).First(u)
	return nil
}

func FillUserByEmail(u *model.User) error {
	if u.Email == "" {
		return errors.New("email 为空！")
	}
	initialize.DB.Where(model.NewUserByEmail(u.Email)).First(u)
	return nil
}

func FillUserByGitHubId(u *model.User) error {
	if u.GitHubId == "" {
		return errors.New("GitHub id 为空！")
	}
	initialize.DB.Where(model.NewUserByGitHubId(u.GitHubId)).First(u)
	return nil
}

func FillUserByLarkId(u *model.User) error {
	if u.LarkId == "" {
		return errors.New("lark id 为空！")
	}
	initialize.DB.Where(model.NewUserByLarkId(u.LarkId)).First(u)
	return nil
}

func FillUserByOidcId(u *model.User) error {
	if u.OidcId == "" {
		return errors.New("oidc id 为空！")
	}
	initialize.DB.Where(model.NewUserByOidcId(u.OidcId))
	return nil
}

func FillUserByWeChatId(u *model.User) error {
	if u.WeChatId == "" {
		return errors.New("WeChat id 为空！")
	}
	initialize.DB.Where(model.NewUserByWeChatId(u.WeChatId)).First(u)
	return nil
}

func FillUserByUsername(u *model.User) error {
	if u.Username == "" {
		return errors.New("username 为空！")
	}
	initialize.DB.Where(model.NewUserByUsername(u.Username)).First(u)
	return nil
}

func IsEmailAlreadyTaken(email string) bool {
	return initialize.DB.Where("email = ?", email).Find(model.NewUser()).RowsAffected == 1
}

func IsWeChatIdAlreadyTaken(wechatId string) bool {
	return initialize.DB.Where("wechat_id = ?", wechatId).Find(model.NewUser()).RowsAffected == 1
}

func IsGitHubIdAlreadyTaken(githubId string) bool {
	return initialize.DB.Where("github_id = ?", githubId).Find(model.NewUser()).RowsAffected == 1
}

func IsLarkIdAlreadyTaken(githubId string) bool {
	return initialize.DB.Where("lark_id = ?", githubId).Find(model.NewUser()).RowsAffected == 1
}

func IsOidcIdAlreadyTaken(oidcId string) bool {
	return initialize.DB.Where("oidc_id = ?", oidcId).Find(model.NewUser()).RowsAffected == 1
}

func IsUsernameAlreadyTaken(username string) bool {
	return initialize.DB.Where("username = ?", username).Find(model.NewUser()).RowsAffected == 1
}

func ResetUserPasswordByEmail(email string, password string) error {
	if email == "" || password == "" {
		return errors.New("邮箱地址或密码为空！")
	}
	hashedPassword, err := utils.Password2Hash(password)
	if err != nil {
		return err
	}
	err = initialize.DB.Model(model.NewUser()).Where("email = ?", email).Update("password", hashedPassword).Error
	return err
}

func IsAdmin(userId int) bool {
	if userId == 0 {
		return false
	}
	var user model.User
	err := initialize.DB.Where("id = ?", userId).Select("role").Find(&user).Error
	if err != nil {
		logger.SysError("no such user " + err.Error())
		return false
	}
	return user.Role >= RoleAdminUser
}

func IsUserEnabled(userId int) (bool, error) {
	if userId == 0 {
		return false, errors.New("user id is empty")
	}
	var user model.User
	err := initialize.DB.Where("id = ?", userId).Select("status").Find(&user).Error
	if err != nil {
		return false, err
	}
	return user.Status == UserStatusEnabled, nil
}

func ValidateAccessToken(token string) (user *model.User) {
	if token == "" {
		return nil
	}
	token = strings.Replace(token, "Bearer ", "", 1)
	user = model.NewUser()
	if initialize.DB.Where("access_token = ?", token).First(user).RowsAffected == 1 {
		return user
	}
	return nil
}

// func GetUserQuota(id int) (quota int64, err error) {
// 	err = initialize.DB.Model(model.NewUser()).Where("id = ?", id).Select("quota").Find(&quota).Error
// 	return quota, err
// }

// func GetUserUsedQuota(id int) (quota int64, err error) {
// 	err = initialize.DB.Model(model.NewUser()).Where("id = ?", id).Select("used_quota").Find(&quota).Error
// 	return quota, err
// }

func GetUserEmail(id int) (email string, err error) {
	err = initialize.DB.Model(model.NewUser()).Where("id = ?", id).Select("email").Find(&email).Error
	return email, err
}

func GetUserGroup(id int) (group string, err error) {
	groupCol := "`group`"
	if initialize.UsingPostgreSQL {
		groupCol = `"group"`
	}

	err = initialize.DB.Model(model.NewUser()).Where("id = ?", id).Select(groupCol).Find(&group).Error
	return group, err
}

// func IncreaseUserQuota(id int, quota int64) (err error) {
// 	if quota < 0 {
// 		return errors.New("quota 不能为负数！")
// 	}
// 	if config.BatchUpdateEnabled {
// 		addNewRecord(BatchUpdateTypeUserQuota, id, quota)
// 		return nil
// 	}
// 	return increaseUserQuota(id, quota)
// }

// func increaseUserQuota(id int, quota int64) (err error) {
// 	err = initialize.DB.Model(model.NewUser()).Where("id = ?", id).Update("quota", gorm.Expr("quota + ?", quota)).Error
// 	return err
// }

// func DecreaseUserQuota(id int, quota int64) (err error) {
// 	if quota < 0 {
// 		return errors.New("quota 不能为负数！")
// 	}
// 	if config.BatchUpdateEnabled {
// 		addNewRecord(BatchUpdateTypeUserQuota, id, -quota)
// 		return nil
// 	}
// 	return decreaseUserQuota(id, quota)
// }

// func decreaseUserQuota(id int, quota int64) (err error) {
// 	err = initialize.DB.Model(model.NewUser()).Where("id = ?", id).Update("quota", gorm.Expr("quota - ?", quota)).Error
// 	return err
// }

func GetRootUserEmail() (email string) {
	initialize.DB.Model(model.NewUser()).Where("role = ?", RoleRootUser).Select("email").Find(&email)
	return email
}

// func UpdateUserUsedQuotaAndRequestCount(id int, quota int64) {
// 	if config.BatchUpdateEnabled {
// 		addNewRecord(BatchUpdateTypeUsedQuota, id, quota)
// 		addNewRecord(BatchUpdateTypeRequestCount, id, 1)
// 		return
// 	}
// 	updateUserUsedQuotaAndRequestCount(id, quota, 1)
// }

// func updateUserUsedQuotaAndRequestCount(id int, quota int64, count int) {
// 	err := initialize.DB.Model(model.NewUser()).Where("id = ?", id).Updates(
// 		map[string]interface{}{
// 			"used_quota":    gorm.Expr("used_quota + ?", quota),
// 			"request_count": gorm.Expr("request_count + ?", count),
// 		},
// 	).Error
// 	if err != nil {
// 		logger.SysError("failed to update user used quota and request count: " + err.Error())
// 	}
// }

// func updateUserUsedQuota(id int, quota int64) {
// 	err := initialize.DB.Model(model.NewUser()).Where("id = ?", id).Updates(
// 		map[string]interface{}{
// 			"used_quota": gorm.Expr("used_quota + ?", quota),
// 		},
// 	).Error
// 	if err != nil {
// 		logger.SysError("failed to update user used quota: " + err.Error())
// 	}
// }

// func updateUserRequestCount(id int, count int) {
// 	err := initialize.DB.Model(model.NewUser()).Where("id = ?", id).Update("request_count", gorm.Expr("request_count + ?", count)).Error
// 	if err != nil {
// 		logger.SysError("failed to update user request count: " + err.Error())
// 	}
// }

func GetUsernameById(id int) (username string) {
	initialize.DB.Model(model.NewUser()).Where("id = ?", id).Select("username").Find(&username)
	return username
}
