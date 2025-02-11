package orders

import (
	"gorm.io/gorm"

	"usdtpay/infr/mysql/dao"
)

func GetOrdersByAmount(tx *gorm.DB, appId string, address string, amount int64) ([]*dao.TradeOrders, error) {
	var res []*dao.TradeOrders
	err := tx.Where(
		"address = ? and app_id = ? and amount = ? and status = ?",
		address,
		appId,
		amount,
		dao.StatusWait,
	).Order("money asc").Limit(100).Select("id,amount").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func OrderByTradeId(tx *gorm.DB, appid string, tradeId string) (*dao.TradeOrders, bool, error) {
	var order dao.TradeOrders
	err := tx.Where("app_id = ? and trade_id = ?", appid, tradeId).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}

func OrderByOrderId(tx *gorm.DB, appid string, orderId string) (*dao.TradeOrders, bool, error) {
	var order dao.TradeOrders
	err := tx.Where("app_id = ? and order_id = ?", appid, orderId).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}

func OrderById(tx *gorm.DB, id int64) (*dao.TradeOrders, bool, error) {
	var order dao.TradeOrders
	err := tx.Where("id = ?", id).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &order, true, nil
}
