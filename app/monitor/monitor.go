package monitor

import (
	"log"
	"time"

	"github.com/shopspring/decimal"

	"github.com/v03413/bepusdt/app/config"
)

type task struct {
	Tick     time.Duration
	Callback func(tick time.Duration)
}

var scheduleList []task

var _paymentMinAmount, _paymentMaxAmount decimal.Decimal

func InitPayment() {
	_paymentMinAmount = config.GetPaymentMinAmount()
	_paymentMaxAmount = config.GetPaymentMaxAmount()
}

func Start() {
	for _, t := range scheduleList {
		schedu := t
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()
			schedu.Callback(schedu.Tick)
		}()
	}
}

func RegisterSchedule(tick time.Duration, callback func(tick time.Duration)) {
	scheduleList = append(
		scheduleList, task{
			Tick:     tick,
			Callback: callback,
		},
	)
}

func inPaymentAmountRange(payAmount decimal.Decimal) bool {
	if payAmount.GreaterThan(_paymentMaxAmount) {
		return false
	}

	if payAmount.LessThan(_paymentMinAmount) {
		return false
	}

	return true
}
