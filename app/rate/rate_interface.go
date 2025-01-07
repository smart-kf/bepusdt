package rate

import (
	"github.com/spf13/cast"

	"github.com/v03413/bepusdt/app/config"
	"github.com/v03413/bepusdt/app/log"
)

var (
	globalRates = map[string]RateInterface{}
)

const (
	Usdt      = "usdt"
	Cny       = "cny"
	Trx       = "trx"
	Usdt2Usdt = "usdt-usdt"
	Cny2Usdt  = "cny-usdt"
	Trx2Trx   = "trx-trx"
	Cny2Trx   = "cny-trx"
)

func InitRates() {
	globalRates[Cny2Usdt] = &Cny2UsdtRate{
		configRate: cast.ToFloat64(config.GetUsdtRate()),
	}
	globalRates[Cny2Trx] = &Cny2TrxRate{
		configRate: cast.ToFloat64(config.GetUsdtRate()),
	}
	globalRates[Usdt2Usdt] = &UsdtUsdtRate{}
	globalRates[Trx2Trx] = &TrxTrxRate{}
}

func ConvertRate(tradeType string, amount float64) float64 {
	impl, ok := globalRates[tradeType]
	if !ok {
		panic("rate impl not found: " + tradeType)
	}
	return impl.ConvertRate(amount)
}

func SetRate(tradeType string, rate float64) {
	impl, ok := globalRates[tradeType]
	if !ok {
		return
	}
	log.Info("设置okex汇率: %s, %v", tradeType, rate)
	impl.SetOkexRate(rate)
}

type RateInterface interface {
	ConvertRate(from float64) float64
	SetOkexRate(rate float64)
}

type Cny2UsdtRate struct {
	configRate float64
	okexRate   float64
}

func (u *Cny2UsdtRate) ConvertRate(from float64) float64 {
	if u.configRate > 0 {
		return from / u.configRate
	}

	if u.okexRate == 0 {
		panic("please start Okex Client to get usdt rate")
	}

	return from / u.okexRate
}

func (u *Cny2UsdtRate) SetOkexRate(rate float64) {
	u.okexRate = rate
}

type UsdtUsdtRate struct{}

func (u UsdtUsdtRate) ConvertRate(from float64) float64 {
	return from
}

func (u UsdtUsdtRate) SetOkexRate(rate float64) {}

type Cny2TrxRate struct {
	configRate float64
	okexRate   float64
}

func (r *Cny2TrxRate) ConvertRate(from float64) float64 {
	if r.configRate > 0 {
		return from / r.configRate
	}
	if r.okexRate == 0 {
		panic("please start Okex Client to get trx rate")
	}
	return from / r.okexRate
}

func (u *Cny2TrxRate) SetOkexRate(rate float64) {
	u.okexRate = rate
}

type TrxTrxRate struct{}

func (r *TrxTrxRate) ConvertRate(from float64) float64 {
	return from
}

func (u TrxTrxRate) SetOkexRate(rate float64) {}
