package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	xlogger "github.com/clearcodecn/log"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"

	"usdtpay/config"
	"usdtpay/domain/event"
	"usdtpay/infr/mysql/dao"
	"usdtpay/infr/mysql/orders"
)

type NotifyService struct {
	event  *event.OrderNotify
	tx     *gorm.DB
	order  *dao.TradeOrders
	status int // 通知的状态.
	app    config.App
}

func NewNotifyService(e *event.OrderNotify) *NotifyService {
	return &NotifyService{
		event: e,
	}
}

func (s *NotifyService) init() error {
	s.tx = config.Setting.MysqlClient

	if err := s.initOrder(); err != nil {
		return err
	}

	app := config.Setting.FindApp(s.order.AppId)
	if app.AppId == "" {
		return errors.New("app not found")
	}
	s.app = app
	return nil
}

func (s *NotifyService) initOrder() error {
	order, ok, err := orders.OrderById(s.tx, s.event.Id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("order not found")
	}
	s.order = order
	return nil
}

func (s *NotifyService) Notify() error {
	if err := s.init(); err != nil {
		xlogger.Error(context.Background(), "Notify error", xlogger.Err(err), xlogger.Any("data", s.event))
		return err
	}
	var (
		retry     bool
		retryTime time.Duration
	)
	// 置为失败
	if s.order.NotifyNum >= s.app.NotifyNumber {
		s.order.NotifyState = dao.NotifyFailed
	} else {
		err := s.notify()
		if err != nil {
			s.order.NotifyState = dao.NotifyFailed
			s.order.NotifyNum++
			retry = true
			retryTime = time.Duration(s.order.NotifyNum) * time.Second
			if retryTime >= 60*time.Second {
				retryTime = 60 * time.Second
			}
		} else {
			s.order.NotifyState = dao.NotifySuccess
		}
	}
	err := s.updateOrder()
	if err != nil {
		return err
	}
	if retry {
		data, _ := json.Marshal(s.event)
		if err := config.Setting.NsqProducer.DeferredPublish(
			config.Setting.NSQ.NotifyTopic, retryTime,
			data,
		); err != nil {
			xlogger.Error(context.Background(), "retry Publish failed", xlogger.Err(err), xlogger.Any("event", s.event))
			return nil
		}
	}
	xlogger.Info(
		context.Background(), "异步通知成功", xlogger.Any(
			"data", map[string]interface{}{
				"app_id":   s.app.AppId,
				"order_id": s.order.OrderId,
				"trade_id": s.order.TradeId,
				"hash":     s.order.TradeHash,
			},
		),
	)
	return nil
}

func (s *NotifyService) updateOrder() error {
	err := s.tx.Where("id = ?", s.order.Id).Select(
		[]string{
			"notify_state",
			"notify_num",
		},
	).Updates(s.order).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *NotifyService) notify() error {
	httpClient := &http.Client{
		Timeout: time.Duration(config.Setting.HttpClient.Timeout) * time.Second,
	}
	var confirmAt int64 = 0
	if s.order.ConfirmedAt != nil {
		confirmAt = s.order.ConfirmedAt.Unix()
	}
	r := resty.NewWithClient(httpClient)
	rsp, err := r.R().SetBody(
		map[string]interface{}{
			"trade_id":  s.order.TradeId,
			"order_id":  s.order.OrderId,
			"status":    s.order.Status,
			"timestamp": confirmAt,
		},
	).
		SetHeader("Authorization", config.Setting.Token).
		Post(s.app.NotifyUrl)

	if err != nil {
		return err
	}
	if rsp.StatusCode() != 200 {
		return errors.New("notify failed, status not 200")
	}
	return nil
}
