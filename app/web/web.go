package web

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/v03413/bepusdt/app/config"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/static"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)

	listen := config.GetListen()

	r := loadStatic(gin.New())
	r.Use(gin.LoggerWithWriter(log.GetWriter()), gin.Recovery())
	r.Use(
		func(ctx *gin.Context) {
			// 解析请求地址
			_host := "http://" + ctx.Request.Host
			if ctx.Request.TLS != nil {
				_host = "https://" + ctx.Request.Host
			}
			_host = config.GetAppUri(_host)

			ctx.Set("HTTP_HOST", _host)
		},
	)
	r.GET(
		"/", func(c *gin.Context) {
			c.HTML(
				200,
				"index.html",
				gin.H{"title": "一款更易用的USDT收款网关", "url": "https://github.com/v03413/bepusdt"},
			)
		},
	)

	payRoute := r.Group("/pay")
	{
		// 收银台
		payRoute.GET("/checkout-counter/:trade_id", checkoutCounter)
		// 状态检测
		payRoute.GET("/check-status/:trade_id", checkStatus)
	}

	// 创建订单
	orderRoute := r.Group("/api/v1/order")
	{
		orderRoute.POST("/create", createTransaction)
	}

	// 易支付兼容
	// r.POST("/submit.php", epaySubmit)

	log.Info("WEB尝试启动 Listen: ", listen)
	go func() {
		err := r.Run(listen)
		if err != nil {
			log.Error("Web启动失败", err)
		}
	}()
}

// 加载静态资源
func loadStatic(engine *gin.Engine) *gin.Engine {
	staticPath := config.GetStaticPath()
	if staticPath != "" {
		engine.Static("/img", config.GetStaticPath()+"/img")
		engine.Static("/css", config.GetStaticPath()+"/css")
		engine.Static("/js", config.GetStaticPath()+"/js")
		engine.LoadHTMLGlob(config.GetStaticPath() + "/views/*")

		return engine
	}

	engine.StaticFS("/img", http.FS(subFs(static.Img, "img")))
	engine.StaticFS("/css", http.FS(subFs(static.Css, "css")))
	engine.StaticFS("/js", http.FS(subFs(static.Js, "js")))
	engine.SetHTMLTemplate(template.Must(template.New("").ParseFS(static.Views, "views/*.html")))

	return engine
}

func subFs(src fs.FS, dir string) fs.FS {
	subFS, _ := fs.Sub(src, dir)

	return subFS
}
