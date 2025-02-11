package dao

import "time"

type AddressTransaction struct {
	Id            int64     `json:"id" gorm:"primaryKey"`
	Address       string    `json:"address"` // 地址.
	TransactionId string    `json:"transaction_id" gorm:"unique"`
	BlockTime     int64     `json:"block_time"`
	FromAddress   string    `json:"from_address"`
	Type          string    `json:"type"`    // 必须=Transfer
	Token         string    `json:"token"`   // address=合约地址，symbol=USDT
	Value         int64     `json:"value"`   // 值
	Decimal       int       `json:"decimal"` // 小数点位数
	FingerPrint   string    `json:"finger_print" gorm:"column:finger_print"`
	CreateTime    time.Time `json:"create_time"`
}

func (a *AddressTransaction) TableName() string {
	return "address_transaction"
}
