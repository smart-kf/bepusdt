package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type UsdtPaymentClient struct {
	host    string
	token   string
	timeout time.Duration
}

func NewUsdtPaymentClient(host, token string, timeout time.Duration) *UsdtPaymentClient {
	c := &UsdtPaymentClient{
		host:    host,
		token:   token,
		timeout: timeout,
	}
	return c
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func NewParamsError(msg string) error {
	return &Error{
		Code:    400,
		Message: msg,
	}
}

func NewCodeError(code int, msg string) error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

type CreateOrderRequest struct {
	AppId       string  `json:"app_id"`
	OrderId     string  `json:"order_id"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	FromAddress string  `json:"from_address"`
	Expire      int     `json:"expire"` // 过期秒数
}

func (r *CreateOrderRequest) Validate() error {
	if r.AppId == "" {
		return NewParamsError("appid 不能为空")
	}
	if r.OrderId == "" {
		return NewParamsError("orderid 不能为空")
	}
	if r.Name == "" {
		return NewParamsError("name 不能为空")
	}
	if r.FromAddress == "" {
		return NewParamsError("fromAddress不能为空")
	}
	if r.Expire == 0 {
		return NewParamsError("过期时间不能为空")
	}
	return nil
}

type CreateOrderResponse struct {
	Data CreateOrderData `json:"data"`
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
}

type CreateOrderData struct {
	TradeId string `json:"trade_id"`
	PayUrl  string `json:"pay_url"`
}

func (c *UsdtPaymentClient) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderData, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	r := resty.New().R().SetBody(req).SetHeader("Authorization", c.token)
	rsp, err := r.Post(fmt.Sprintf("%s/api/v1/order/create", c.host))
	if err != nil {
		return nil, NewCodeError(500, err.Error())
	}
	var cor CreateOrderResponse
	err = json.Unmarshal(rsp.Body(), &cor)
	if err != nil {
		return nil, err
	}
	if cor.Code != 0 {
		return nil, NewCodeError(cor.Code, cor.Msg)
	}
	return &cor.Data, nil
}
