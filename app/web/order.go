package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/v03413/bepusdt/app/config"
	"github.com/v03413/bepusdt/app/epay"
	"github.com/v03413/bepusdt/app/help"
	"github.com/v03413/bepusdt/app/help/sign"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/rate"
)

type CreateOrderRequest struct {
	OrderId     string  `json:"order_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	NotifyUrl   string  `json:"notify_url" binding:"required"`
	RedirectUrl string  `json:"redirect_url" binding:"required"`
	TradeType   string  `json:"tradeType" binding:"required"` // usdt-usdt  || cny-usdt || cny-trx || trx-trx
	Sign        string  `json:"sign" binding:"required"`      // 签名.
}

func (r *CreateOrderRequest) Validate() error {
	var mp = make(map[string]interface{})
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &mp)
	if err != nil {
		return err
	}
	if sign.ValidateSign(mp, config.GetAuthToken()) {
		return nil
	}

	// tradeType.
	switch r.TradeType {
	case rate.Usdt2Usdt, rate.Cny2Trx, rate.Trx2Trx:
	default:
		return errors.New("tradeType 错误")
	}

	return fmt.Errorf("签名错误")
}

// createTransaction 创建订单
func createTransaction(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(200, respFailJson(fmt.Errorf("参数错误:"+err.Error())))
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(200, respFailJson(fmt.Errorf("参数错误:"+err.Error())))
		return
	}

	// 解析请求地址
	host := "http://" + ctx.Request.Host
	if ctx.Request.TLS != nil {
		host = "https://" + ctx.Request.Host
	}

	order, err := buildOrder(
		req.Amount,
		model.OrderApiTypeEpusdt,
		req.OrderId,
		req.TradeType,
		req.RedirectUrl,
		req.NotifyUrl,
		req.Name,
	)
	if err != nil {
		ctx.JSON(200, respFailJson(fmt.Errorf("订单创建失败：%w", err)))
		return
	}

	// 返回响应数据
	ctx.JSON(
		200, respSuccJson(
			gin.H{
				"trade_id":        order.TradeId,
				"order_id":        req.OrderId,
				"amount":          req.Amount,
				"actual_amount":   order.Amount,
				"token":           order.Address,
				"expiration_time": int64(order.ExpiredAt.Sub(time.Now()).Seconds()),
				"payment_url":     fmt.Sprintf("%s/pay/checkout-counter/%s", config.GetAppUri(host), order.TradeId),
			},
		),
	)
	log.Info(fmt.Sprintf("订单创建成功，商户订单号：%s", req.OrderId))
}

func buildOrder(money float64, apiType, orderId, tradeType, redirectUrl, notifyUrl, name string) (
	model.TradeOrders,
	error,
) {
	var order model.TradeOrders

	// 获取钱包地址
	wallet := model.GetAvailableAddress()
	if len(wallet) == 0 {
		log.Error("订单创建失败：还没有配置收款地址")
		return order, fmt.Errorf("还没有配置收款地址")
	}

	// 计算交易金额
	address, amount := model.CalcTradeAmount(wallet, money, tradeType)
	tradeId, err := help.GenerateTradeId()
	if err != nil {
		return order, err
	}

	// 创建交易订单
	expiredAt := time.Now().Add(config.GetExpireTime())
	tradeOrder := model.TradeOrders{
		OrderId:     orderId,
		TradeId:     tradeId,
		TradeHash:   tradeId, // 这里默认填充一个本地交易ID，等支付成功后再更新为实际交易哈希
		TradeType:   tradeType,
		Amount:      amount,
		Money:       money,
		Address:     address.Address,
		Status:      model.OrderStatusWaiting,
		Name:        name,
		ApiType:     apiType,
		ReturnUrl:   redirectUrl,
		NotifyUrl:   notifyUrl,
		NotifyNum:   0,
		NotifyState: model.OrderNotifyStateFail,
		ExpiredAt:   expiredAt,
	}
	res := model.DB.Create(&tradeOrder)
	if res.Error != nil {
		log.Error("订单创建失败：", res.Error.Error())

		return order, res.Error
	}

	return tradeOrder, nil
}

func checkoutCounter(ctx *gin.Context) {
	tradeId := ctx.Param("trade_id")
	order, ok := model.GetTradeOrder(tradeId)
	if !ok {
		ctx.String(200, "订单不存在")

		return
	}

	uri, err := url.ParseRequestURI(order.ReturnUrl)
	if err != nil {
		ctx.String(200, "同步地址错误")
		log.Error("同步地址解析错误", err.Error())

		return
	}

	ctx.HTML(
		200, order.TradeType+".html", gin.H{
			"http_host":  uri.Host,
			"trade_id":   tradeId,
			"amount":     order.Amount,
			"address":    order.Address,
			"expire":     int64(order.ExpiredAt.Sub(time.Now()).Seconds()),
			"return_url": order.ReturnUrl,
			"usdt_rate":  order.TradeRate,
		},
	)
}

func checkStatus(ctx *gin.Context) {
	tradeId := ctx.Param("trade_id")
	order, ok := model.GetTradeOrder(tradeId)
	if !ok {
		ctx.JSON(200, respFailJson(fmt.Errorf("订单不存在")))

		return
	}

	var returnUrl string
	if order.Status == model.OrderStatusSuccess {
		returnUrl = order.ReturnUrl
		if order.ApiType == model.OrderApiTypeEpay {
			// 易支付兼容
			returnUrl = fmt.Sprintf("%s?%s", returnUrl, epay.BuildNotifyParams(order))
		}
	}

	// 返回响应数据
	ctx.JSON(200, gin.H{"trade_id": tradeId, "status": order.Status, "return_url": returnUrl})
}
