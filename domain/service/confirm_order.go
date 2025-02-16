package service

import (
	"context"
	"encoding/json"
	"time"

	xlogger "github.com/clearcodecn/log"
	"gorm.io/gorm"

	"usdtpay/config"
	"usdtpay/domain/event"
	"usdtpay/infr/mysql/dao"
)

type ConfirmOrderService struct {
	tx          *gorm.DB
	transaction *dao.AddressTransaction
	event       *event.BlockChainEvent
	order       *dao.TradeOrders
}

func NewConfirmOrderService(event *event.BlockChainEvent) *ConfirmOrderService {
	return &ConfirmOrderService{
		event: event,
	}
}

func (c *ConfirmOrderService) init() error {
	c.tx = config.Setting.MysqlClient.Begin()
	return c.initOrder()
}

func (c *ConfirmOrderService) initOrder() error {
	var (
		order dao.TradeOrders
		tx    *dao.AddressTransaction
	)

	err := c.tx.Where("transaction_id = ?", c.event.TransactionId).First(&tx).Error
	if err != nil {
		return err
	}
	c.transaction = tx

	err = c.tx.Where(
		"address = ? and from_address = ? and status = ? and created_at < ? and money = ?",
		tx.Address,
		tx.FromAddress,
		dao.StatusWait,
		time.Unix(tx.BlockTime/1000, 0),
		tx.Value,
	).Order("created_at asc").First(&order).Error

	if err != nil {
		return err
	}
	c.order = &order
	return nil
}

func (c *ConfirmOrderService) Confirm() error {
	if err := c.init(); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		xlogger.Error(context.Background(), "Confirm order failed", xlogger.Err(err), xlogger.Any("data", c.event))
		c.tx.Rollback()
		return err
	}
	// 2. 更新订单状态。
	if err := c.updateOrder(); err != nil {
		c.tx.Rollback()
		return err
	}

	// :: 这里先提交事务，因为先推送队列，会导致队列消费查询订单，事务还没来得及提交的话，
	// 订单数据还是旧的.
	c.tx.Commit()
	// 3. 通知调用方, 推送调用队列里面.
	e := event.OrderNotify{
		Id:     c.order.Id,
		Status: c.order.Status,
	}
	data, _ := json.Marshal(e)
	if err := config.Setting.NsqProducer.Publish(config.Setting.NSQ.NotifyTopic, data); err != nil {
		xlogger.Error(context.Background(), "notifyPublish failed", xlogger.Err(err))
		return err
	}
	return nil
}

func (c *ConfirmOrderService) updateOrder() error {
	confirm := time.Unix(c.transaction.BlockTime/1000, 0)
	c.order.TradeHash = c.transaction.TransactionId
	c.order.ConfirmedAt = &confirm
	c.order.Status = dao.StatusSuccess

	err := c.tx.Where("id = ?", c.order.Id).Select(
		[]string{
			"trade_hash",
			"confirmed_at",
			"status",
		},
	).Updates(c.order).Error

	return err
}
