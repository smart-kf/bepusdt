package controller

import (
	"errors"

	"github.com/gin-gonic/gin"

	"usdtpay/config"
	"usdtpay/infr/mysql/orders"
)

type QueryOrderRequest struct {
	AppId   string `json:"app_id"`
	OrderId string `json:"order_id"`
}

func Query(ctx *gin.Context) {
	var req QueryOrderRequest
	if err := ctx.ShouldBind(&req); err != nil {
		sendError(ctx, err)
		return
	}
	order, ok, err := orders.OrderByTradeId(config.Setting.MysqlClient, req.AppId, req.OrderId)
	if err != nil {
		sendError(ctx, err)
		return
	}
	if !ok {
		sendError(ctx, errors.New("order not exists"))
		return
	}
	sendSuccess(ctx, order)
}
