package controller

import (
	"time"

	"github.com/gin-gonic/gin"

	"usdtpay/config"
	"usdtpay/infr/mysql/orders"
)

type PayRequest struct {
	TradeId string `json:"trade_id" form:"trade_id" binding:"required"`
	AppId   string `json:"app_id" form:"app_id" binding:"required"`
}

func Pay(ctx *gin.Context) {
	var req PayRequest
	if err := ctx.BindQuery(&req); err != nil {
		ctx.String(200, "参数错误")
		return
	}
	order, ok, err := orders.OrderByTradeId(config.Setting.MysqlClient, req.AppId, req.TradeId)
	if err != nil {
		ctx.String(200, "参数错误")
		return
	}
	if !ok {
		ctx.String(200, "参数错误")
		return
	}
	money := float64(order.Money) / 1e6
	expire := order.ExpiredAt.Unix() - time.Now().Unix()
	ctx.HTML(
		200, "usdt.trc20.html", gin.H{
			"order":  order,
			"money":  money,
			"expire": expire,
		},
	)
}
