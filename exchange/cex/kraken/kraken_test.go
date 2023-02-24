package kraken

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

var (
	kraken *ClientKraken
)

func init() {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
		IsTest:   false,
	}
	kraken = NewClientKraken(conf)
}

func TestGetSymbols(t *testing.T) {
	res := kraken.GetSymbols()
	fmt.Println(res)
}
