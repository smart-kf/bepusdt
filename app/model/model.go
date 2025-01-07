package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/v03413/bepusdt/app/config"
)

var (
	DB   *gorm.DB
	_err error
)

func Init() error {
	db, err := gorm.Open(mysql.Open(config.GetConfig().DB.Dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	if _err = AutoMigrate(); _err != nil {
		return _err
	}

	addStartWalletAddress()

	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(&WalletAddress{}, &TradeOrders{}, &NotifyRecord{}, &Config{})
}
