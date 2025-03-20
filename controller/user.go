package controller

type topUpRequest struct {
	Key string `json:"key"`
}

// func TopUp(c *gin.Context) {
// 	// ctx := c.Request.Context()
// 	req := topUpRequest{}
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	// id := c.GetInt("id")
// 	// quota, err := model.Redeem(ctx, req.Key, id)
// 	// if err != nil {
// 	// 	c.JSON(http.StatusOK, gin.H{
// 	// 		"success": false,
// 	// 		"message": err.Error(),
// 	// 	})
// 	// 	return
// 	// }
// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "",
// 		// "data":    quota,
// 	})
// 	return
// }

type adminTopUpRequest struct {
	UserId int    `json:"user_id"`
	Quota  int    `json:"quota"`
	Remark string `json:"remark"`
}

// func AdminTopUp(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	req := adminTopUpRequest{}
// 	err := c.ShouldBindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	err = model.IncreaseUserQuota(req.UserId, int64(req.Quota))
// 	if err != nil {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": false,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	if req.Remark == "" {
// 		req.Remark = fmt.Sprintf("通过 API 充值 %s", common.LogQuota(int64(req.Quota)))
// 	}
// 	model.RecordTopupLog(ctx, req.UserId, req.Remark, req.Quota)
// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "",
// 	})
// 	return
// }
