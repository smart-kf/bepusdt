package transaction

import (
	"gorm.io/gorm"

	"usdtpay/infr/mysql/dao"
)

func GetLastAddressTransactionId(tx *gorm.DB, address string) *dao.AddressTransaction {
	var tr dao.AddressTransaction
	err := tx.Where("address = ? and finger_print != ''", address).Order("create_time desc").First(&tr).Error
	if err != nil {
		return nil
	}
	return &tr
}
