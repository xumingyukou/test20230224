package bitbank

import (
	"clients/exchange/cex/base"
	"errors"
	"fmt"
	"testing"
)

var (
	symbols               = []string{"BTC/JPY", "ETH/BTC", "BAT/BTC", "DOGE/JPY"}
	symbolsExchange       = []string{"btc_jpy", "eth_btc", "bat_btc", "doge_jpy"}
	ProxyUrl              = "http://127.0.0.1:1087"
	TimeOffset      int64 = 30
	conf            base.APIConf
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

func TestCanonical2Exchange(t *testing.T) {
	for i := 0; i < len(symbols); i++ {
		if Canonical2Exchange(symbols[i]) != symbolsExchange[i] {
			t.Fatal(errors.New(""))
		}
	}
}

func TestExchange2Canonical(t *testing.T) {
	for i := 0; i < len(symbols); i++ {
		if Exchange2Canonical(symbolsExchange[i]) != symbols[i] {
			t.Fatal(errors.New(""))
		}
	}
}

func TestGetSymbols(t *testing.T) {
	client := NewClientBitbank(conf)
	symbols := client.GetSymbols()
	if len(symbols) == 0 {
		t.Fatal(errors.New("no symbols"))
	}
	fmt.Println(symbols)
	fmt.Println(len(symbols))
}
