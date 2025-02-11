package dto

import "time"

type CreateOrderDTO struct {
	AppId       string        // 业务id
	OrderId     string        // 订单id
	Amount      float64       // 金额, 实际金额，没有乘以 1000000
	FromAddress string        // 支付地址
	GoodName    string        // 商品名称
	ExpireIn    time.Duration // 超时时间
}
