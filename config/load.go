package config

import (
	"fmt"
	"sync"

	"github.com/make-money-fast/xconfig"
	"gopkg.in/yaml.v2"
)

var (
	o sync.Once
)

func Load(file string) *Config {
	var conf Config
	err := xconfig.ParseFromFile(file, &conf)
	if err != nil {
		panic(err)
	}
	o.Do(
		func() {
			data, _ := yaml.Marshal(conf)
			fmt.Println("================ config start ===========")
			fmt.Println(string(data))
			fmt.Println("================ config end ===========")
		},
	)
	return &conf
}
