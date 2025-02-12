package dao

import "time"

const (
	StatusWait = iota + 1
	StatusSuccess
	StatusFailed
)

const (
	NotifySuccess = iota + 1
	NotifyFailed  = 2
)

type TradeOrders struct {
	Id          int64      `gorm:"primary_key;AUTO_INCREMENT;comment:id"`
	OrderId     string     `gorm:"type:varchar(255);not null;unique;color:blue;comment:客户订单ID"`
	AppId       string     `gorm:"column:app_id;type:varchar(255)"`
	TradeId     string     `gorm:"type:varchar(255);not null;unique;color:blue;comment:本地订单ID"`
	TradeHash   string     `gorm:"type:varchar(64);default:'';comment:交易哈希"`
	Amount      int64      `gorm:"default:0;comment:USDT交易数额,实际金额"`
	Money       int64      `gorm:"default:0;comment:订单交易金额,乘以1e6"`
	Address     string     `gorm:"type:varchar(34);not null;comment:收款地址"`
	FromAddress string     `gorm:"type:varchar(34);not null;default:'';comment:支付地址"`
	Status      int        `gorm:"type:tinyint(1);not null;default:0;comment:交易状态 1：等待支付 2：支付成功 3：订单过期"`
	Name        string     `gorm:"type:varchar(64);not null;default:'';comment:商品名称"`
	NotifyNum   int        `gorm:"type:int(11);not null;default:0;comment:回调次数"`
	NotifyState int        `gorm:"type:tinyint(1);not null;default:0;comment:回调状态 1：成功 0：失败"`
	ExpiredAt   time.Time  `gorm:"type:timestamp;not null;comment:订单失效时间"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;type:timestamp;not null;comment:创建时间"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime;type:timestamp;not null;comment:更新时间"`
	ConfirmedAt *time.Time `gorm:"type:timestamp;null;comment:交易确认时间"`
}

func (TradeOrders) TableName() string {
	return "trade_orders"
}
