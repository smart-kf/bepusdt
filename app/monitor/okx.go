package monitor

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cast"
	"github.com/tidwall/gjson"

	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/rate"
)

func init() {
	RegisterSchedule(0, OkxUsdtRateStart)
	RegisterSchedule(0, OkxTrxUsdtRateStart)
}

// OkxUsdtRateStart Okx USDT_CNY 汇率监控
func OkxUsdtRateStart(dur time.Duration) {
	if dur < time.Minute {
		dur = time.Minute
	}
	tk := time.NewTicker(dur)
	defer tk.Stop()

	for range tk.C {
		rawRate, err := getOkxUsdtCnySellPrice()
		if err != nil {
			log.Error("Okx USDT_CNY 汇率获取失败", err)
		} else {
			rate.SetRate(rate.Cny2Usdt, rawRate)
		}
	}
}

// OkxTrxUsdtRateStart  Okx TRX_USDT 汇率监控
func OkxTrxUsdtRateStart(dur time.Duration) {
	if dur < time.Minute {
		dur = time.Minute
	}
	tk := time.NewTicker(dur)
	defer tk.Stop()
	for range tk.C {
		rawRate, err := getOkxTrxUsdtRate()
		if err != nil {
			log.Error("Okx TRX_USDT 汇率获取失败", err)
		} else {
			rate.SetRate(rate.Cny2Trx, rawRate)
		}
	}
}

// GetOkxTrxUsdtRate 获取 Okx TRX 汇率 https://www.okx.com/zh-hans/trade-spot/trx-usdt
func getOkxTrxUsdtRate() (float64, error) {
	url := "https://www.okx.com/priapi/v5/market/candles?instId=TRX-USDT&before=1727143156000&bar=4H&limit=1&t=" + cast.ToString(time.Now().UnixNano())
	client := http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("referer", "https://www.okx.com/zh-hans/trade-spot/trx-usdt")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
	)
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("okx resp error:" + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, errors.New("okx resp status code:" + strconv.Itoa(resp.StatusCode))
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("okx resp read body error:" + err.Error())
	}

	result := gjson.ParseBytes(all)
	if result.Get("data").Exists() {
		data := result.Get("data").Array()
		if len(data) > 0 {
			return data[0].Get("1").Float(), nil
		}
	}

	return 0, errors.New("okx resp json data not found")
}

// getOkxUsdtCnySellPrice  Okx  C2C快捷交易 USDT出售 实时汇率
func getOkxUsdtCnySellPrice() (float64, error) {
	t := strconv.Itoa(int(time.Now().Unix()))
	okxApi := "https://www.okx.com/v4/c2c/express/price?crypto=USDT&fiat=CNY&side=sell&t=" + t
	client := http.Client{Timeout: time.Second * 5}
	req, _ := http.NewRequest("GET", okxApi, nil)
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	)
	resp, err := client.Do(req)
	if err != nil {
		return 0, errors.New("okx resp error:" + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, errors.New("okx resp status code:" + strconv.Itoa(resp.StatusCode))
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, errors.New("okx resp read error:" + err.Error())
	}

	result := gjson.ParseBytes(all)
	if result.Get("error_code").Int() != 0 {
		return 0, errors.New("json parse error:" + result.Get("error_message").String())
	}

	if result.Get("data.price").Exists() {
		_ret := result.Get("data.price").Float()
		if _ret <= 0 {
			return 0, errors.New("okx resp json data.price <= 0")
		}

		return cast.ToFloat64(_ret), nil
	}

	return 0, errors.New("okx resp json data.price not found")
}
