package nsq

import (
	"log"
	"os"
	"time"

	"github.com/nsqio/go-nsq"

	"usdtpay/config"
)

var (
	nsqConsumers []*nsq.Consumer
)

func InitNSQ() {
	nsqConsumers = append(nsqConsumers, messageConsumer())
}

func StartConsume(stopChan chan struct{}) {
	<-stopChan
	for _, c := range nsqConsumers {
		c.Stop()
		<-c.StopChan
	}
}

func messageConsumer() *nsq.Consumer {
	hostname, _ := os.Hostname()
	conf := config.GetConfig().NSQ

	cfg := nsq.NewConfig()
	cfg.DialTimeout = time.Duration(conf.Timeout) * time.Second
	cfg.ReadTimeout = time.Duration(conf.Timeout) * time.Second
	cfg.WriteTimeout = time.Duration(conf.Timeout) * time.Second
	cfg.ClientID = hostname
	cfg.Hostname = hostname + "-message-consumer"
	cfg.UserAgent = "go-" + hostname + "-message-consumer"

	consumer, err := nsq.NewConsumer(conf.MessageTopic, conf.MessageTopicGroup, cfg)
	if err != nil {
		panic(err)
	}
	// TODO:: consumer.
	// consumer.AddHandler(&consumer2.MessageConsumer{})
	err = consumer.ConnectToNSQDs(conf.Addrs)
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}
