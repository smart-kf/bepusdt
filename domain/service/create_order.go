package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"usdtpay/config"
	"usdtpay/domain/dto"
	"usdtpay/domain/utils"
	"usdtpay/infr/mysql/address"
	"usdtpay/infr/mysql/dao"
	"usdtpay/infr/mysql/orders"
)

type CreateOrderService struct {
	tx      *gorm.DB            // 事务.
	order   *dao.TradeOrders    // 订单
	address *dao.Address        // 地址.
	data    *dto.CreateOrderDTO // amount
	money   int64
}

func NewCreateOrderService(data dto.CreateOrderDTO) (*CreateOrderService, error) {
	service := CreateOrderService{}
	service.data = &data

	if err := service.init(); err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *CreateOrderService) init() error {
	// 1. 初始化事务.
	s.tx = config.Setting.MysqlClient.Begin()

	// 2. 查找一个可用的地址.
	if err := s.initAvailableAddress(); err != nil {
		return err
	}

	// 3. 获取浮点位
	if err := s.getAmount(); err != nil {
		return err
	}

	return nil
}

func (s *CreateOrderService) initAvailableAddress() error {
	addr, err := address.GetNextAddress(s.tx, s.data.AppId)
	if err != nil {
		return err
	}
	s.address = addr
	return nil
}

func (s *CreateOrderService) getAmount() error {
	orders, err := orders.GetOrdersByAmount(s.tx, s.data.AppId, s.address.Address, s.data.Amount)
	if err != nil {
		return nil
	}
	if len(orders) == 0 {
		s.money = int64(s.data.Amount * 1e6)
		return nil
	}
	if len(orders) == 100 {
		return errors.New("no available address")
	}
	money := int64(s.data.Amount * 1e6)
	floor := int64(0.01 * 1000000)
	for _, o := range orders {
		if o.Money != money {
			s.money = money
			break
		} else {
			money += floor
		}
	}
	return nil
}

func (s *CreateOrderService) CreateOrder() (*dao.TradeOrders, error) {
	order, err := s.buildOrder()
	if err != nil {
		s.tx.Rollback()
		return nil, err
	}
	s.tx.Commit()
	return order, err
}

func (s *CreateOrderService) buildOrder() (*dao.TradeOrders, error) {
	var order = &dao.TradeOrders{
		OrderId:     s.data.OrderId,
		AppId:       s.data.AppId,
		Amount:      s.data.Amount,
		Money:       s.money,
		Address:     s.address.Address,
		FromAddress: s.data.FromAddress,
		Status:      dao.StatusWait,
		Name:        s.data.GoodName,
		ExpiredAt:   time.Now().Add(s.data.ExpireIn),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	order.TradeId = utils.Sha1Object(order)
	if err := s.tx.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}
