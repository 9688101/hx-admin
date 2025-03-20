package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/9688101/hx-admin/utils/helper"
)

func RequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := helper.GenRequestID()
		c.Set(helper.RequestIdKey, id)
		ctx := helper.SetRequestID(c.Request.Context(), id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(helper.RequestIdKey, id)
		c.Next()
	}
}
