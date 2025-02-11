package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"usdtpay/example_client/sdk"
)

type GoodInfo struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// this is a example client to run usdt pay.
func main() {
	g := gin.Default()
	g.LoadHTMLGlob("example_client/views/*.gohtml")
	g.GET(
		"/", func(ctx *gin.Context) {
			var goods = []GoodInfo{
				{
					Id:    1,
					Name:  "测试商品",
					Price: 10,
				},
			}
			ctx.HTML(
				200, "index.gohtml", gin.H{
					"goods": goods,
				},
			)
		},
	)
	g.GET(
		"/pay", func(ctx *gin.Context) {
			price, _ := strconv.ParseFloat(ctx.Query("price"), 64)
			fromAddress := ctx.Query("fromAddress")
			cli := sdk.NewUsdtPaymentClient("http://localhost:8082", "1234", 3*time.Second)
			rsp, err := cli.CreateOrder(
				context.Background(), &sdk.CreateOrderRequest{
					AppId:       "app1",
					OrderId:     "123456" + time.Now().Format(time.RFC3339),
					Name:        "购买商品",
					Amount:      price,
					FromAddress: fromAddress,
					Expire:      600,
				},
			)
			if err != nil {
				ctx.String(200, "创建订单失败:"+err.Error())
				return
			}

			fmt.Println("新订单: " + rsp.TradeId)
			ctx.Redirect(302, rsp.PayUrl)
		},
	)
	g.Any(
		"/api/notify", func(ctx *gin.Context) {
			bd, _ := ioutil.ReadAll(ctx.Request.Body)
			fmt.Println("收到回调", string(bd))
			ctx.String(200, "ok")
		},
	)

	g.Any(
		"/return", func(ctx *gin.Context) {
			ctx.String(200, "支付成功")
		},
	)

	g.Run(":9092")
}
