package controller

import (
	"errors"
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

func CheckStatus(ctx *gin.Context) {
	id := ctx.Param("tradeId")
	appId := ctx.Query("appid")
	app := config.Setting.FindApp(appId)
	if app.AppId == "" {
		sendError(
			ctx, errors.New("500 internal server error"),
		)
		return
	}
	order, ok, err := orders.OrderByTradeId(config.Setting.MysqlClient, appId, id)
	if !ok {
		sendError(
			ctx, errors.New("500 internal server error"),
		)
		return
	}
	if err != nil {
		sendError(
			ctx, errors.New("500 internal server error"),
		)
		return
	}
	sendSuccess(
		ctx, gin.H{
			"status":     order.Status,
			"return_url": app.ReturnUrl,
		},
	)
}
