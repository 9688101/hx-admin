package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/9688101/hx-admin/utils"
)

func RequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := utils.GenRequestID()
		c.Set(utils.RequestIdKey, id)
		ctx := utils.SetRequestID(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(utils.RequestIdKey, id)
		c.Next()
	}
}
