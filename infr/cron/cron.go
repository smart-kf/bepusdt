package cron

import (
	"github.com/jasonlvhit/gocron"

	"usdtpay/domain/service"
)

func StartCronTask(stopChan chan struct{}) {
	// apiKey 限制一天10w次请求  86400 ( 10w / 地址数量 )
	gocron.Every(5).Second().Do(TransactionMonitor)

	ch := gocron.Start()

	<-stopChan
	close(ch)
}

func TransactionMonitor() {
	s := service.NewTransactionService()
	s.RunOnce()
}
