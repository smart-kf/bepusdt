package notify

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/nsqio/go-nsq"

	"usdtpay/config"
	"usdtpay/domain/event"
	"usdtpay/domain/service"
)

type NotifyConsumer struct {
	nsqConsumer *nsq.Consumer
}

func NewNotifyConsumer() *NotifyConsumer {
	consumer := &NotifyConsumer{}
	c := config.Setting

	timeout := time.Duration(c.NSQ.Timeout) * time.Second
	hostname, _ := os.Hostname()
	cfg := nsq.NewConfig()
	cfg.DialTimeout = timeout
	cfg.ReadTimeout = timeout
	cfg.WriteTimeout = timeout
	cfg.ClientID = hostname
	cfg.Hostname = hostname + "-blockchain-consumer"
	cfg.UserAgent = "go-" + hostname + "-blockchain-consumer"

	nsqConsumer, err := nsq.NewConsumer(c.NSQ.NotifyTopic, c.NSQ.NotifyGroup, cfg)
	if err != nil {
		panic(err)
	}
	nsqConsumer.AddHandler(consumer)
	err = nsqConsumer.ConnectToNSQDs(c.NSQ.Addrs)
	if err != nil {
		log.Fatal(err)
	}
	consumer.nsqConsumer = nsqConsumer
	return consumer
}

func (b *NotifyConsumer) HandleMessage(message *nsq.Message) error {
	var transfer event.OrderNotify
	err := json.Unmarshal(message.Body, &transfer)
	if err != nil {
		return err
	}
	svc := service.NewNotifyService(&transfer)
	err = svc.Notify()
	if err != nil {
		xlogger.Error(context.Background(), "NotifyConsumer error", xlogger.Err(err))
	}
	return nil
}

func (c *NotifyConsumer) Stop() {
	c.nsqConsumer.Stop()
}
