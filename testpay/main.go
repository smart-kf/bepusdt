package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/v03413/bepusdt/app/help/sign"
)

func main() {
	g := gin.Default()
	g.Any(
		"/notify", func(ctx *gin.Context) {
			ctx.String(200, "success")
		},
	)

	g.Any(
		"/", func(ctx *gin.Context) {

		},
	)

	g.Any(
		"/test", func(ctx *gin.Context) {
			req := CreateOrderRequest{
				OrderId:     fmt.Sprintf("NO%d", time.Now().Unix()),
				Name:        "testing",
				Amount:      float64(rand.Intn(10)),
				NotifyUrl:   "http://localhost:8083/",
				RedirectUrl: "http://localhost:8083/",
				TradeType:   "usdt-usdt",
			}
			var mp = make(map[string]interface{})
			data, _ := json.Marshal(req)
			json.Unmarshal(data, &mp)
			sign.SetSign(mp, "123")
			data, _ = json.Marshal(mp)
			httpReq, _ := http.NewRequest(
				http.MethodPost, "http://localhost:8082/api/v1/order/create",
				bytes.NewReader(data),
			)
			httpReq.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(httpReq)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			w := io.MultiWriter(os.Stdout, ctx.Writer)
			io.Copy(w, resp.Body)
		},
	)

	g.Run(":8083")
}

type CreateOrderRequest struct {
	OrderId     string  `json:"order_id" binding:"required"`
	Name        string  `json:"name" binding:"name"`
	Amount      float64 `json:"amount" binding:"required"`
	NotifyUrl   string  `json:"notify_url" binding:"required"`
	RedirectUrl string  `json:"redirect_url" binding:"required"`
	TradeType   string  `json:"tradeType" binding:"require"` // usdt-usdt  || cny-usdt || cny-trx || trx-trx
	Sign        string  `json:"sign" binding:"required"`     // 签名.
}
