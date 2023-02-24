package bitMex

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

var (
	bitmex *ClientBitMex
)

func init() {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
		IsTest:   false,
	}
	bitmex = NewClientBitMex(conf)
}

func TestGetSymbols(t *testing.T) {
	res := bitmex.GetSymbols()
	fmt.Println(res)
}
