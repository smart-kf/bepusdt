package http

import (
	"context"
	"html/template"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"usdtpay/config"
	"usdtpay/endpoints/http/controller"
	"usdtpay/static"
)

func StartHttpServer(stopChan chan struct{}) {
	g := gin.New()
	g.Use(gin.Recovery())
	var logConfig xlogger.GinLogConfigure
	logConfig.LogIP(ClientIP)
	logConfig.SkipPrefix("/static", "/favico.ico")
	if config.Setting.Debug {
		logConfig.EnableRequestBody()
	}
	g.Use(xlogger.GinLog(logConfig))

	g.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
				AllowHeaders:     []string{"authorization,x-requested-with,content-type"},
				AllowCredentials: false,
				MaxAge:           24 * time.Hour,
			},
		),
	)

	if config.Setting.Debug {
		g.Static("/img", config.Setting.Web.StaticDir+"/img")
		g.Static("/css", config.Setting.Web.StaticDir+"/css")
		g.Static("/js", config.Setting.Web.StaticDir+"/js")
		g.LoadHTMLGlob(config.Setting.Web.StaticDir + "/views/*.html")
	} else {
		g.StaticFS("/img", http.FS(subFs(static.Img, "img")))
		g.StaticFS("/css", http.FS(subFs(static.Css, "css")))
		g.StaticFS("/js", http.FS(subFs(static.Js, "js")))
		g.SetHTMLTemplate(template.Must(template.New("").ParseFS(static.Views, "views/*.html")))
	}
	orderRoute := g.Group("/api/v1/order")
	orderRoute.Use(
		func(ctx *gin.Context) {
			authorization := ctx.Request.Header.Get("Authorization")
			if authorization != config.Setting.Token {
				ctx.AbortWithStatus(403)
				return
			}
		},
	)
	{
		orderRoute.POST("/create", controller.CreateOrder)
		orderRoute.POST("/query", controller.Query)
		orderRoute.POST("/mail", controller.SendMail)
	}

	g.GET("/pay", controller.Pay)
	g.GET("/pay/check-status/:tradeId", controller.CheckStatus)
	server := &http.Server{Addr: config.Setting.Web.String(), Handler: g}

	go func() {
		server.ListenAndServe()
	}()

	<-stopChan
	server.Shutdown(context.Background())
}

func ClientIP(ctx *gin.Context) string {
	if ip := ctx.Request.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := ctx.Request.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr))
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

func subFs(src fs.FS, dir string) fs.FS {
	subFS, _ := fs.Sub(src, dir)

	return subFS
}
