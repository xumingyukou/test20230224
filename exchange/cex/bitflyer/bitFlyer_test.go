package bitFlyer

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

var (
	bitFlyer *ClientBitFlyer
)

func init() {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
		IsTest:   false,
	}
	bitFlyer = NewClientBitFlyer(conf)
}

func TestGetSymbols(t *testing.T) {
	res := bitFlyer.GetSymbols()
	fmt.Println(res)
}
