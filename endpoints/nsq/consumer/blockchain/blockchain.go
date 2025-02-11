package blockchain

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

type BlockChainConsumer struct {
	nsqConsumer *nsq.Consumer
}

func NewBlockChainConsumer() *BlockChainConsumer {
	consumer := &BlockChainConsumer{}
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

	nsqConsumer, err := nsq.NewConsumer(c.NSQ.BlockChainTopic, c.NSQ.BlockChainGroup, cfg)
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

func (b *BlockChainConsumer) HandleMessage(message *nsq.Message) error {
	var transfer event.BlockChainEvent
	err := json.Unmarshal(message.Body, &transfer)
	if err != nil {
		return err
	}
	svc := service.NewConfirmOrderService(&transfer)
	err = svc.Confirm() // 直接消费，错误报错
	if err != nil {
		xlogger.Error(context.Background(), "BlockChainConsumer error", xlogger.Err(err))
	}
	return nil
}

func (c *BlockChainConsumer) Stop() {
	c.nsqConsumer.Stop()
}
