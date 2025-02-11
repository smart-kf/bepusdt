package config

import (
	"fmt"

	"github.com/nsqio/go-nsq"
	"gorm.io/gorm"
)

type Config struct {
	Debug       bool       `json:"debug"`
	Web         Web        `json:"web"`
	Log         Log        `json:"log"`
	DB          Db         `json:"db"`
	AddressList []Address  `json:"address_list"`
	NSQ         NSQ        `json:"nsq"`
	HttpClient  HttpClient `json:"httpClient"`
	Token       string     `json:"token"`
	Tron        Tron       `json:"tron"`

	MysqlClient *gorm.DB      `json:"-"`
	NsqProducer *nsq.Producer `json:"-"`
	Apps        []App         `json:"apps"`
}

func (c Config) AddressMap() map[string]struct{} {
	var res = make(map[string]struct{})
	for _, v := range c.AddressList {
		if v.Enable {
			res[v.Address] = struct{}{}
		}
	}
	return res
}

type App struct {
	AppId        string `json:"app_id"`
	Token        string `json:"token"`
	ReturnUrl    string `json:"return_url"`
	NotifyUrl    string `json:"notify_url"`
	NotifyNumber int    `json:"notify_number"`
}

type Address struct {
	AppId   string `json:"app_id"`
	Address string `json:"address"`
	Enable  bool   `json:"enable"`
}

func (c *Config) FindApp(id string) App {
	for _, a := range c.Apps {
		if a.AppId == id {
			return a
		}
	}
	return App{}
}

type Web struct {
	Addr      string `json:"addr" default:"127.0.0.1"`
	Port      int    `json:"port" default:"8081"`
	StaticDir string `json:"staticDir" default:"static"`
	AppHost   string `json:"app_host" default:"http://localhost:8082"` // 访问域名
}

func (w Web) String() string {
	return fmt.Sprintf("%s:%d", w.Addr, w.Port)
}

type Db struct {
	Dsn    string `json:"dsn"`    // 连接
	Driver string `json:"driver"` // 默认 sqlite3
}

type Log struct {
	Level  string `json:"level" default:"info"`
	Format string `json:"format" default:"json"`
	File   string `json:"file"`
}

type NSQ struct {
	Addrs           []string `json:"addrs"`
	Timeout         int      `json:"timeout" default:"60"`
	BlockChainTopic string   `json:"block_chain_topic" default:"block-chain"`
	BlockChainGroup string   `json:"block_chain_group" default:"block-chain-group"`
	NotifyTopic     string   `json:"notify_topic" default:"notify-topic"`
	NotifyGroup     string   `json:"notify_group" default:"notify-group"`
}

type HttpClient struct {
	SocketServerClient string `json:"socketServerAddress"`
	Timeout            int    `json:"timeout"`
	Proxy              string `json:"proxy"`
}

type Tron struct {
	ApiHost             string `json:"apiHost"`
	ApiKey              string `json:"apiKey"`
	UsdtContractAddress string `json:"usdtContractAddress"`
	Proxy               string `json:"proxy"`
	Timeout             int    `json:"timeout"`
	CronSecond          int    `json:"cron_second"` // 定时任务执行秒数间隔
}
