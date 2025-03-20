package server

import (
	"fmt"

	// "sync"
	"time"

	"github.com/9688101/hx-admin/core/logger"
	"github.com/9688101/hx-admin/global"
	"github.com/9688101/hx-admin/initialize"
)

var (
	TokenCacheSeconds         = global.SyncFrequency
	UserId2GroupCacheSeconds  = global.SyncFrequency
	UserId2QuotaCacheSeconds  = global.SyncFrequency
	UserId2StatusCacheSeconds = global.SyncFrequency
	GroupModelsCacheSeconds   = global.SyncFrequency
)

func CacheGetUserGroup(id int) (group string, err error) {
	if !initialize.RedisEnabled {
		return GetUserGroup(id)
	}
	group, err = initialize.RedisGet(fmt.Sprintf("user_group:%d", id))
	if err != nil {
		group, err = GetUserGroup(id)
		if err != nil {
			return "", err
		}
		err = initialize.RedisSet(fmt.Sprintf("user_group:%d", id), group, time.Duration(UserId2GroupCacheSeconds)*time.Second)
		if err != nil {
			logger.SysError("Redis set user group error: " + err.Error())
		}
	}
	return group, err
}

func CacheIsUserEnabled(userId int) (bool, error) {
	if !initialize.RedisEnabled {
		return IsUserEnabled(userId)
	}
	enabled, err := initialize.RedisGet(fmt.Sprintf("user_enabled:%d", userId))
	if err == nil {
		return enabled == "1", nil
	}

	userEnabled, err := IsUserEnabled(userId)
	if err != nil {
		return false, err
	}
	enabled = "0"
	if userEnabled {
		enabled = "1"
	}
	err = initialize.RedisSet(fmt.Sprintf("user_enabled:%d", userId), enabled, time.Duration(UserId2StatusCacheSeconds)*time.Second)
	if err != nil {
		logger.SysError("Redis set user enabled error: " + err.Error())
	}
	return userEnabled, err
}
