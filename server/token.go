package server

import (
	"errors"

	"github.com/9688101/hx-admin/initialize"
	"github.com/9688101/hx-admin/model"
)

const (
	TokenStatusEnabled   = 1 // don't use 0, 0 is the default value!
	TokenStatusDisabled  = 2 // also don't use 0
	TokenStatusExpired   = 3
	TokenStatusExhausted = 4
)

func GetTokenByIds(id int, userId int) (*model.Token, error) {
	if id == 0 || userId == 0 {
		return nil, errors.New("id 或 userId 为空！")
	}
	t := model.NewTokenByUserId(id, userId)
	var err error = nil
	err = initialize.DB.First(t, "id = ? and user_id = ?", id, userId).Error
	return t, err
}

func GetTokenById(id int) (*model.Token, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	t := model.NewTokenById(id)
	var err error = nil
	err = initialize.DB.First(&t, "id = ?", id).Error
	return t, err
}

func InsertToken(t *model.Token) error {
	var err error
	err = initialize.DB.Create(t).Error
	return err
}

// Update Make sure your token's fields is completed, because this will update non-zero values
func UpdateToken(t *model.Token) error {
	var err error
	err = initialize.DB.Model(t).Select("name", "status", "expired_time", "remain_quota", "unlimited_quota", "models", "subnet").Updates(t).Error
	return err
}

func SelectUpdateToken(t *model.Token) error {
	// This can update zero values
	return initialize.DB.Model(t).Select("accessed_time", "status").Updates(t).Error
}

func DeleteToken(t *model.Token) error {
	var err error
	err = initialize.DB.Delete(t).Error
	return err
}

// func (t *Token) GetModels() string {
// 	if t == nil {
// 		return ""
// 	}
// 	if t.Models == nil {
// 		return ""
// 	}
// 	return *t.Models
// }

func DeleteTokenById(id int, userId int) (err error) {
	// Why we need userId here? In case user want to delete other's token.
	if id == 0 || userId == 0 {
		return errors.New("id 或 userId 为空！")
	}
	token := model.Token{Id: id, UserId: userId}
	err = initialize.DB.Where(token).First(&token).Error
	if err != nil {
		return err
	}
	return DeleteToken(&token)
}
