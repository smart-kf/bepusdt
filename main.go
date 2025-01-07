package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/v03413/bepusdt/app"
	"github.com/v03413/bepusdt/app/config"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/monitor"
	"github.com/v03413/bepusdt/app/rate"
	"github.com/v03413/bepusdt/app/telegram"
	"github.com/v03413/bepusdt/app/web"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "config.yaml", "配置文件地址")
}

func Init() error {
	if err := model.Init(); err != nil {
		return fmt.Errorf("数据库初始化失败：" + err.Error())
	}

	if bot := config.GetTgBot(); bot.Enable {
		telegram.InitBot(bot.Token)
		monitor.RegisterSchedule(0, monitor.BotStart)
	}

	rate.InitRates()
	monitor.InitPayment()

	return nil
}

func main() {
	flag.Parse()
	config.Load(configFile)

	if err := Init(); err != nil {
		panic(err)
	}

	monitor.Start()

	web.Start()

	fmt.Println("Bepusdt 启动成功，当前版本：" + app.Version)

	{
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		runtime.GC()
	}
}
