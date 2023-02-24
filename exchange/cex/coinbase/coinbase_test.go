package coinbase

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

var (
	ProxyUrl = "http://127.0.0.1:7890"
)

func TestGetSymbols(t *testing.T) {
	var timeOffset int64 = 30
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		//ProxyUrl:    ProxyUrl,
	}
	coinbase := NewClientCoinbase(conf)
	res := coinbase.GetSymbols()
	fmt.Println(res)
}
