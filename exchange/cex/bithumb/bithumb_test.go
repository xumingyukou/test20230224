package bithumb

import (
	"clients/exchange/cex/base"
	"errors"
	"fmt"
	"testing"
)

var (
	ProxyUrl         = "http://127.0.0.1:1087"
	TimeOffset int64 = 30
	conf       base.APIConf
)

func init() {
	// config.LoadExchangeConfig("./conf/exchange.toml")
	conf = base.APIConf{
		ReadTimeout: TimeOffset,
		ProxyUrl:    ProxyUrl,
		EndPoint:    "",
		AccessKey:   "",
		Passphrase:  "",
		SecretKey:   "",
		IsTest:      true,
	}
}

func TestLoadConfig(t *testing.T) {
	fmt.Println(conf)
}

func TestGetSymbols(t *testing.T) {
	client := NewClientBithumb(conf)
	symbols := client.GetSymbols()
	if len(symbols) == 0 {
		t.Fatal(errors.New("no symbols"))
	}
	fmt.Println(symbols)
	fmt.Println(len(symbols))
}
