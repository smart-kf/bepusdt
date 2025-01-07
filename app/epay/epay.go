package epay

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/config"
	"github.com/v03413/bepusdt/app/help"
	"github.com/v03413/bepusdt/app/model"
)

const Pid = "1000" // 固定商户号

func Sign(params map[string]string, key string) string {
	// 提取 keys 并排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	// 构建签名字符串
	signStr := ""
	for _, k := range keys {
		if k != "sign" && k != "sign_type" && params[k] != "" {
			signStr += fmt.Sprintf("%s=%s&", k, params[k])
		}
	}
	signStr = signStr[:len(signStr)-1] // 移除最后一个 '&'
	signStr += key                     // 添加密钥

	// 计算 MD5
	hash := md5.New()
	hash.Write([]byte(signStr))
	md5sum := hex.EncodeToString(hash.Sum(nil))

	return md5sum
}

func BuildNotifyParams(order model.TradeOrders) string {
	query := fmt.Sprintf("money=%s&name=%s&out_trade_no=%s&pid=%s&trade_no=%s&trade_status=TRADE_SUCCESS&type=%s",
		cast.ToString(order.Money), order.Name, order.OrderId, Pid, order.TradeId, order.TradeType)

	sign := help.Md5String(query + config.GetAuthToken())

	return fmt.Sprintf("%s&sign=%s", query, sign)
}
