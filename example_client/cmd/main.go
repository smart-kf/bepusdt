package main

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
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
					Price: 1,
				},
			}
			ctx.HTML(
				200, "index.gohtml", gin.H{
					"goods": goods,
				},
			)
		},
	)
	g.Any(
		"/api/notify", func(ctx *gin.Context) {
			bd, _ := ioutil.ReadAll(ctx.Request.Body)
			fmt.Println("收到回调", string(bd))
			ctx.String(200, "ok")
		},
	)

	g.Run(":9092")
}
