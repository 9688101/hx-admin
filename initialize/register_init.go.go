package initialize

// import (
// 	"github.com/9688101/hx-admin/common"
// 	"github.com/9688101/hx-admin/common/helper"
// 	"github.com/9688101/hx-admin/common/logger"
// 	"github.com/9688101/hx-admin/common/random"
// 	"github.com/9688101/hx-admin/global"
// )

// func CreateRootAccountIfNeed() error {
// 	var user User
// 	//if user.Status != util.UserStatusEnabled {
// 	if err := DB.First(&user).Error; err != nil {
// 		logger.SysLog("no user exists, creating a root user for you: username is root, password is 123456")
// 		hashedPassword, err := common.Password2Hash("123456")
// 		if err != nil {
// 			return err
// 		}
// 		accessToken := random.GetUUID()
// 		if global.InitialRootAccessToken != "" {
// 			accessToken = global.InitialRootAccessToken
// 		}
// 		rootUser := User{
// 			Username:    "root",
// 			Password:    hashedPassword,
// 			Role:        RoleRootUser,
// 			Status:      UserStatusEnabled,
// 			DisplayName: "Root User",
// 			AccessToken: accessToken,
// 			// Quota:       500000000000000,
// 		}
// 		DB.Create(&rootUser)
// 		if global.InitialRootToken != "" {
// 			logger.SysLog("creating initial root token as requested")
// 			token := Token{
// 				Id:     1,
// 				UserId: rootUser.Id,
// 				Key:    global.InitialRootToken,
// 				// Status:       TokenStatusEnabled,
// 				Name:         "Initial Root Token",
// 				CreatedTime:  helper.GetTimestamp(),
// 				AccessedTime: helper.GetTimestamp(),
// 				ExpiredTime:  -1,
// 				// RemainQuota:    500000000000000,
// 				// UnlimitedQuota: true,
// 			}
// 			DB.Create(&token)
// 		}
// 	}
// 	return nil
// }
