package address

import (
	"errors"
	"math/rand"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"usdtpay/infr/mysql/caches"
	"usdtpay/infr/mysql/dao"
)

var (
	addressCache = caches.NewCache[*dao.Address]()
)

func GetNextAddress(tx *gorm.DB, appId string) (*dao.Address, error) {
	var add []*dao.Address
	err := tx.Where("enable = ? and app_id = ?", true, appId).Find(&add).Error
	if err != nil {
		return nil, err
	}
	if len(add) == 0 {
		return nil, errors.New("address is not found")
	}
	return add[rand.Intn(len(add))], nil
}

func GetAddressList(tx *gorm.DB) ([]*dao.Address, error) {
	address, ok := addressCache.Get("address_list")
	if !ok {
		var add []*dao.Address
		err := tx.Where("enable = ? ", true).Find(&add).Error
		if err != nil {
			return nil, err
		}
		addressCache.Set("address_list", add)
		return add, nil
	}
	return address, nil
}

func GetAddressMap(tx *gorm.DB) (map[string]struct{}, error) {
	list, err := GetAddressList(tx)
	if err != nil {
		return nil, err
	}
	return lo.SliceToMap(
		list, func(item *dao.Address) (string, struct{}) {
			return item.Address, struct{}{}
		},
	), nil
}
