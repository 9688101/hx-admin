package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/9688101/hx-admin/global"
	"github.com/9688101/hx-admin/initialize"
	"github.com/9688101/hx-admin/utils"
)

var timeFormat = "2006-01-02T15:04:05.000Z"

var inMemoryRateLimiter utils.InMemoryRateLimiter

func redisRateLimiter(c *gin.Context, maxRequestNum int, duration int64, mark string) {
	ctx := context.Background()
	rdb := initialize.RDB
	key := "rateLimit:" + mark + c.ClientIP()
	listLength, err := rdb.LLen(ctx, key).Result()
	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		c.Abort()
		return
	}
	if listLength < int64(maxRequestNum) {
		rdb.LPush(ctx, key, time.Now().Format(timeFormat))
		rdb.Expire(ctx, key, global.RateLimitKeyExpirationDuration)
	} else {
		oldTimeStr, _ := rdb.LIndex(ctx, key, -1).Result()
		oldTime, err := time.Parse(timeFormat, oldTimeStr)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		nowTimeStr := time.Now().Format(timeFormat)
		nowTime, err := time.Parse(timeFormat, nowTimeStr)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		// time.Since will return negative number!
		// See: https://stackoverflow.com/questions/50970900/why-is-time-since-returning-negative-durations-on-windows
		if int64(nowTime.Sub(oldTime).Seconds()) < duration {
			rdb.Expire(ctx, key, global.RateLimitKeyExpirationDuration)
			c.Status(http.StatusTooManyRequests)
			c.Abort()
			return
		} else {
			rdb.LPush(ctx, key, time.Now().Format(timeFormat))
			rdb.LTrim(ctx, key, 0, int64(maxRequestNum-1))
			rdb.Expire(ctx, key, global.RateLimitKeyExpirationDuration)
		}
	}
}

func memoryRateLimiter(c *gin.Context, maxRequestNum int, duration int64, mark string) {
	key := mark + c.ClientIP()
	if !inMemoryRateLimiter.Request(key, maxRequestNum, duration) {
		c.Status(http.StatusTooManyRequests)
		c.Abort()
		return
	}
}

func rateLimitFactory(maxRequestNum int, duration int64, mark string) func(c *gin.Context) {
	if maxRequestNum == 0 || global.DebugEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	if initialize.RedisEnabled {
		return func(c *gin.Context) {
			redisRateLimiter(c, maxRequestNum, duration, mark)
		}
	} else {
		// It's safe to call multi times.
		inMemoryRateLimiter.Init(global.RateLimitKeyExpirationDuration)
		return func(c *gin.Context) {
			memoryRateLimiter(c, maxRequestNum, duration, mark)
		}
	}
}

func GlobalWebRateLimit() func(c *gin.Context) {
	return rateLimitFactory(global.GlobalWebRateLimitNum, global.GlobalWebRateLimitDuration, "GW")
}

func GlobalAPIRateLimit() func(c *gin.Context) {
	return rateLimitFactory(global.GlobalApiRateLimitNum, global.GlobalApiRateLimitDuration, "GA")
}

func CriticalRateLimit() func(c *gin.Context) {
	return rateLimitFactory(global.CriticalRateLimitNum, global.CriticalRateLimitDuration, "CT")
}

func DownloadRateLimit() func(c *gin.Context) {
	return rateLimitFactory(global.DownloadRateLimitNum, global.DownloadRateLimitDuration, "DW")
}

func UploadRateLimit() func(c *gin.Context) {
	return rateLimitFactory(global.UploadRateLimitNum, global.UploadRateLimitDuration, "UP")
}
