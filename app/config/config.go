package config

import (
	"errors"
	"log"
	"math"
	"strings"
	"time"

	"github.com/make-money-fast/xconfig"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"

	"github.com/v03413/bepusdt/app/help"
)

var (
	_config *config
)

type config struct {
	TronGrpcNode  string      `json:"tronGrpcNode" default:"18.141.79.38:50051"` //  TRON_GRPC_NODE
	UsdtAtom      string      `json:"usdtAtom" default:"0.01"`                   // 精度
	TrxAtom       string      `json:"trxAtom" default:"0.01"`                    // trx 精度
	MinPay        float64     `json:"minPay" default:"0.01"`                     // 最低支付金额
	MaxPay        float64     `json:"maxPay" default:"99999"`                    // 最多支付
	ExpireTime    int         `json:"expireTime" default:"expireTime"`           // 默认支付超时时间
	UsdtRate      string      `json:"usdtRate"`                                  // usdt Rate, 配置了则是默认比例，否则从交易所获取
	TrxRate       string      `json:"trxRate"`                                   // trx rate 配置了则是默认比例，否则从交易所获取
	AuthToken     string      `json:"authToken"`                                 // 交互的 token
	ListenAddress string      `json:"listenAddress" default:"127.0.0.1:8082"`    // 监听地址
	TradeConfirm  bool        `json:"tradeConfirm" default:"false"`              // 是否需要交易确认，保险
	AppUrl        string      `json:"appUrl"`                                    // 收银台地址
	StaticPath    string      `json:"staticPath" default:"./static"`             // 静态目录
	WalletAddress []string    `json:"walletAddress"`                             // 钱包地址
	TelegramBot   TelegramBot `json:"telegramBot"`
	DB            DB          `json:"db"`
}

type DB struct {
	Dsn string `json:"dsn"`
}

func GetConfig() *config {
	return _config
}

type TelegramBot struct {
	Enable  bool   `json:"enable"`  // 是否启用
	AdminId string `json:"adminId"` // 管理员id
	Token   string `json:"token"`   // 机器人token
	GroupId string `json:"groupId"` // 群id
}

func GetTronGrpcNode() string {
	return _config.TronGrpcNode
}

func GetUsdtAtomicity() (decimal.Decimal, int) {
	atom, exp, err := parseAtomicity(_config.UsdtAtom)
	if err != nil {
		log.Fatal(err)
	}
	return atom, exp
}

func GetTrxAtomicity() (decimal.Decimal, int) {
	atom, exp, err := parseAtomicity(_config.TrxAtom)
	if err == nil {
		log.Fatal(err)
	}
	return atom, exp
}

func GetPaymentMinAmount() decimal.Decimal {
	_min := decimal.NewFromFloat(_config.MinPay)
	return _min
}

func GetPaymentMaxAmount() decimal.Decimal {
	max := decimal.NewFromFloat(_config.MaxPay)
	return max
}

func GetExpireTime() time.Duration {
	return time.Duration(_config.ExpireTime) * time.Second
}

func GetUsdtRate() string {
	return _config.UsdtRate
}

func GetTrxRate() string {
	return _config.TrxRate
}

func GetAuthToken() string {
	return _config.AuthToken
}

func GetListen() string {
	return _config.ListenAddress
}

func GetTradeConfirmed() bool {
	return _config.TradeConfirm
}

func GetAppUri(host string) string {
	return _config.AppUrl
}

func GetTGBotToken() string {
	if data := help.GetEnv("TG_BOT_TOKEN"); data != "" {
		return strings.TrimSpace(data)
	}

	return ""
}

func GetTgBot() TelegramBot {
	return _config.TelegramBot
}

func GetOutputLog() string {
	return "/tmp/bepusdt.log"
}

func GetStaticPath() string {
	return _config.StaticPath
}

func GetInitWalletAddress() []string {
	return _config.WalletAddress
}

func parseAtomicity(data string) (decimal.Decimal, int, error) {
	atom, err := decimal.NewFromString(data)
	if err != nil {
		return decimal.Zero, 0, err
	}

	// 如果大于0，且小数点后位数大于0
	if atom.GreaterThan(decimal.Zero) && atom.Exponent() < 0 {
		return atom, cast.ToInt(math.Abs(float64(atom.Exponent()))), nil
	}

	return decimal.Zero, 0, errors.New("原子精度参数不合法")
}

func Load(filename string) {
	var c config
	err := xconfig.ParseFromFile(filename, &c)
	if err != nil {
		panic(err)
	}
	_config = &c
}
