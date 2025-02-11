package main

import (
	"os"
	"os/signal"
	"sync"

	xlogger "github.com/clearcodecn/log"

	"usdtpay/config"
	"usdtpay/endpoints/http"
	"usdtpay/endpoints/nsq/consumer/blockchain"
	"usdtpay/endpoints/nsq/consumer/notify"
	"usdtpay/infr/cron"
)

var (
	wg sync.WaitGroup
)

func main() {
	stopChan := make(chan struct{})
	initLogger(config.Setting)

	wg.Add(1)
	go func() {
		defer wg.Done()

		http.StartHttpServer(stopChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c := blockchain.NewBlockChainConsumer()
		<-stopChan
		c.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c := notify.NewNotifyConsumer()
		<-stopChan
		c.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cron.StartCronTask(stopChan)
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)
	<-sig
	close(stopChan)
	wg.Wait()
}

func initLogger(conf *config.Config) {
	// xlogger.AddHook(func(ctx context.Context) xlogger.Field {
	//	reqid, ok := ctx.Value("reqid").(string)
	//	if !ok {
	//		return xlogger.Field{}
	//	}
	//	return xlogger.Any("reqid", reqid)
	// })
	logger, err := xlogger.NewLog(
		xlogger.Config{
			Level:  conf.Log.Level,
			Format: conf.Log.Format,
			File:   conf.Log.File,
		},
	)

	if err != nil {
		panic(err)
	}

	xlogger.SetGlobal(logger)
}
