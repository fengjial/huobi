package config

import (
	"fmt"
	"github.com/go-gcfg/gcfg"
	"sync"
)

var cfg *ProjectConf
var once sync.Once
var prefix *string

// API KEY
const (
	ACCESS_KEY string = ""
	SECRET_KEY string = ""
)

// API请求地址, 不要带最后的/
const (
	MARKET_URL string = "https://api.huobi.pro"
	TRADE_URL  string = "https://api.huobi.pro"
)

type ProjectConf struct {
	Log struct {
		Path  string
		Level string
		Save  int
	}
}

func Read() *ProjectConf {
	Load()
	return cfg
}

func Load() {
	once.Do(func() {
		cfg = &ProjectConf{}
		err := gcfg.ReadFileInto(cfg, *prefix+"/config.gcfg")
		fmt.Println(cfg)
		if err != nil {
			fmt.Println(err)
		}
	})
}

func SetPrefix(input *string) {
	prefix = input
}
