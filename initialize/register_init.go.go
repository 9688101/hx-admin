package initialize

import (
	"github.com/9688101/hx-admin/core/logger"
	"github.com/9688101/hx-admin/global"
	"github.com/9688101/hx-admin/model"
	"github.com/9688101/hx-admin/utils"
)

func CreateRootAccountIfNeed() error {
	u := model.NewUser()
	//if user.Status != util.UserStatusEnabled {
	if err := DB.First(u).Error; err != nil {
		logger.SysLog("no user exists, creating a root user for you: username is root, password is 123456")
		hashedPassword, err := utils.Password2Hash("123456")
		if err != nil {
			return err
		}
		accessToken := utils.GetUUID()
		if global.InitialRootAccessToken != "" {
			accessToken = global.InitialRootAccessToken
		}
		u = &model.User{
			Username:    "root",
			Password:    hashedPassword,
			Role:        100,
			Status:      1,
			DisplayName: "Root User",
			AccessToken: accessToken,
		}
		DB.Create(u)
		if global.InitialRootToken != "" {
			logger.SysLog("creating initial root token as requested")
			t := model.Token{
				Id:           1,
				UserId:       u.Id,
				Key:          global.InitialRootToken,
				Status:       1,
				Name:         "Initial Root Token",
				CreatedTime:  utils.GetTimestamp(),
				AccessedTime: utils.GetTimestamp(),
				ExpiredTime:  -1,
			}
			DB.Create(&t)
		}
	}
	return nil
}
