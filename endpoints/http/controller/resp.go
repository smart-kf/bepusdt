package controller

import "github.com/gin-gonic/gin"

func sendSuccess(ctx *gin.Context, data any) {
	ctx.JSON(
		200, gin.H{
			"data": data,
			"code": 0,
			"msg":  "",
		},
	)
}

func sendError(ctx *gin.Context, err error) {
	ctx.JSON(
		200, gin.H{
			"data": nil,
			"code": -1,
			"msg":  err.Error(),
		},
	)
}
