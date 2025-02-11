package controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"usdtpay/domain/dto"
	"usdtpay/domain/service"
)

type CreateOrderRequest struct {
	AppId       string  `json:"app_id" binding:"required"`
	OrderId     string  `json:"order_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	FromAddress string  `json:"from_address" binding:"required"`
	Expire      int     `json:"expire" binding:"required"` // 过期秒数
}

func (r *CreateOrderRequest) Validate() error {
	return nil
}

func CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.Bind(&req); err != nil {
		sendError(ctx, fmt.Errorf("参数错误:"+err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		sendError(ctx, fmt.Errorf("参数错误:"+err.Error()))
		return
	}

	svc, err := service.NewCreateOrderService(
		dto.CreateOrderDTO{
			AppId:       req.AppId,
			OrderId:     req.OrderId,
			Amount:      req.Amount,
			FromAddress: req.FromAddress,
			GoodName:    req.Name,
			ExpireIn:    time.Duration(req.Expire) * time.Second,
		},
	)
	if err != nil {
		sendError(ctx, fmt.Errorf("系统错误:"+err.Error()))
		return
	}
	order, err := svc.CreateOrder()
	if err != nil {
		sendError(ctx, fmt.Errorf("系统错误,创建订单失败:"+err.Error()))
		return
	}
	// 返回响应数据
	sendSuccess(
		ctx, gin.H{
			"trade_id": order.TradeId,
		},
	)
}
