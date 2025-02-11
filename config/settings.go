package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"usdtpay/infr/mysql/dao"
)

var (
	Setting *Config
)

func init() {
	initConfig()
}

func getConfigPath() string {
	wd, _ := os.Getwd()
	idx := strings.Index(wd, "usdtpay")
	path := string(wd[:idx]) + "usdtpay/config.yaml"
	return path
}

func initConfig() {
	configFilepath := os.Getenv("CONFIG_FILE")
	if configFilepath == "" {
		configFilepath = getConfigPath()
	}

	Setting = Load(configFilepath)
	Setting.initMysql()
	Setting.initNsqProducer()
}

func (c *Config) initMysql() {
	db, err := gorm.Open(mysql.Open(c.DB.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	c.MysqlClient = db
	syncTable(db)
	go initAddress(db)
}

func (c *Config) initNsqProducer() {
	timeout := time.Duration(c.NSQ.Timeout) * time.Second
	hostname, _ := os.Hostname()
	cfg := nsq.NewConfig()
	cfg.DialTimeout = timeout
	cfg.ReadTimeout = timeout
	cfg.WriteTimeout = timeout
	cfg.ClientID = hostname
	cfg.Hostname = hostname + "-usdt-payment"
	cfg.UserAgent = "go-" + hostname + "-usdt-payment"
	p, err := nsq.NewProducer(c.NSQ.Addrs[0], cfg)
	if err != nil {
		panic(err)
	}
	c.NsqProducer = p
}

func syncTable(db *gorm.DB) {
	if err := db.AutoMigrate(
		&dao.Address{},
		&dao.TradeOrders{},
		&dao.AddressTransaction{},
	); err != nil {
		log.Fatal(err)
	}
}

func initAddress(db *gorm.DB) {
	tx := db.Begin()
	address := Setting.AddressList
	for _, a := range address {
		var err error
		var addr dao.Address
		err = tx.Where("app_id = ? and address = ?", a.AppId, a.Address).First(&addr).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				addr.Address = a.Address
				addr.Enable = a.Enable
				addr.AppId = a.AppId
				err = tx.Model(addr).Save(&addr).Error
				continue
			}
			tx.Rollback()
			panic(err)
		}
		if addr.Enable != a.Enable {
			addr.Enable = a.Enable
			err = tx.Where("id = ?", addr.Id).Save(&addr).Error
			if err != nil {
				tx.Rollback()
				panic(err)
			}
			continue
		}
	}
	tx.Commit()
}
