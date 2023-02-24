package bitfinex

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetSymbols(t *testing.T) {
	conf := base.APIConf{
		// ProxyUrl: "http://127.0.0.1:7890"
	}
	a := NewClientBitfinex(conf)
	symbols := a.GetSymbols()
	fmt.Println(len(symbols))
	b, _ := json.Marshal(symbols)
	fmt.Println(string(b))
}
