package config

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type TtSysConfig struct {
	MaxProcess int `toml:"max_process"` // 最大占用cpu核数
}

type TtConfig struct {
	LogLevel    string      `toml:"log_level"`
	Sys         TtSysConfig `toml:"sys"`
	LogLevelInt logrus.Level
}

type ExchangeWeightInfo struct {
	Spot  map[string]int64 `toml:"spot"`
	UBase map[string]int64 `toml:"u_base"`
	CBase map[string]int64 `toml:"c_base"`
}

type ExchangeInfo struct {
	ApiKeyPath     string             `toml:"api_key_path"`
	SpotConfigPath string             `toml:"spot_config_path"`
	Weight         ExchangeWeightInfo `toml:"weight"`
	SpotConfig     *ExchangeSpotConfig
	ApiKeyConfig   *ApiConfig
}

type PairConfigInfo struct {
	Address string `toml:"address"`
	Reverse bool   `toml:"reverse"`
}

type FactoryConfigInfo struct {
	Address     string `toml:"address"`
	Width       int64  `toml:"width"`
	TickSpacing int64  `toml:"tickspacing"`
	Pairs       map[string]*PairConfigInfo
}

type TtExchangeConfig struct {
	ExchangeList map[string]*ExchangeInfo `toml:"exchange"`
}

type TtDexDepthConfig struct {
	FactoryList map[string]map[string]*FactoryConfigInfo `toml:"factory"` // chain:factory:symbol
}

type TtTokenConfig struct {
	Address string `toml:"address"`
	Decimal int64  `toml:"decimal"`
}

type TtChainItem struct {
	Endpoint         string                    `toml:"endpoint"`
	WsEndpoint       string                    `toml:"ws_endpoint"`
	MultiCallAddress string                    `json:"multicalladdress"`
	LoopTime         int64                     `json:"looptime"`   //second
	Identifier       string                    `json:"identifier"` //pending, confirmed, failed
	Router           map[string]string         `toml:"router"`
	Factory          map[string]string         `toml:"factory"`
	Tokens           map[string]*TtTokenConfig `toml:"tokens"`
	WrappedAddress   string                    `toml:"wrapped_address"`
	KeyPath          string                    `toml:"key_path"`
	ApiKey           *ApiConfig                `toml:"apikey"`
}

type TtChainConfig struct {
	ChainList map[string]*TtChainItem `toml:"chain"`
}

func LoadConfig(file string) {
	var c TtConfig
	err := LoadConfigFromFile(file, &c)
	if err != nil {
		panic(err)
	}

	switch c.LogLevel {
	case "debug":
		c.LogLevelInt = logrus.DebugLevel
	case "info":
		c.LogLevelInt = logrus.InfoLevel
	case "warning":
		c.LogLevelInt = logrus.WarnLevel
	case "error":
		c.LogLevelInt = logrus.ErrorLevel
	case "fatal":
		c.LogLevelInt = logrus.FatalLevel
	}

	mtx.Lock()
	GlobalConfig = &c
	mtx.Unlock()
}

func LoadChainConfig(file string) {
	var c TtChainConfig
	err := LoadConfigFromFile(file, &c)
	if err != nil {
		panic(err)
	}

	for _, conf := range c.ChainList {
		var (
			spotConfData []byte
		)
		if conf.KeyPath != "" {
			spotConfData, err = os.ReadFile(conf.KeyPath)
			if err != nil {
				fmt.Println("get config path err:", conf.KeyPath)
			} else {
				conf.ApiKey = new(ApiConfig)
				err = json.Unmarshal(spotConfData, conf.ApiKey)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
			}
		}
	}
	mtx.Lock()
	ChainConfig = &c
	mtx.Unlock()
}

func LoadExchangeConfig(file string) {
	var c TtExchangeConfig
	err := LoadConfigFromFile(file, &c)
	if err != nil {
		panic(err)
	}
	for _, exchangeConf := range c.ExchangeList {
		var (
			spotConfData, ApiConfData []byte
		)

		if exchangeConf.SpotConfigPath != "" {
			spotConfData, err = os.ReadFile(exchangeConf.SpotConfigPath)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(spotConfData, &exchangeConf.SpotConfig)
			if err != nil {
				panic(err)
			}
		}
		if exchangeConf.ApiKeyPath != "" {
			ApiConfData, err = os.ReadFile(exchangeConf.ApiKeyPath)
			if err != nil {
				dir, _ := os.Getwd()
				fmt.Println(dir)
				panic(err)
			}
			err = json.Unmarshal(ApiConfData, &exchangeConf.ApiKeyConfig)
			if err != nil {
				panic(err)
			}
		}
	}
	mtx.Lock()
	ExchangeConfig = &c
	mtx.Unlock()
}

func LoadDexDepthConfig(file string) {
	var c TtDexDepthConfig
	err := LoadConfigFromFile(file, &c)
	if err != nil {
		panic(err)
	}
	mtx.Lock()
	DexDepthConfig = &c
	mtx.Unlock()
}

var mtx sync.Mutex
var GlobalConfig *TtConfig
var ChainConfig *TtChainConfig
var ExchangeConfig *TtExchangeConfig
var DexDepthConfig *TtDexDepthConfig

/*
func init() {
	LoadConfig("./conf/config.toml")
	LoadChainConfig("./conf/chain.toml")
	LoadExchangeConfig("./conf/exchange.toml")
	LoadDexDepthConfig("./conf/exchange.toml")
}
*/
