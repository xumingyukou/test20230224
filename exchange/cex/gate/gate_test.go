package gate

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"testing"
)

var (
	gate *ClientGate
)

func init() {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
		IsTest:   false,
	}
	gate = NewClientGate(conf)
}

func TestGetSymbols(t *testing.T) {
	res := gate.GetSymbols()
	fmt.Println(len(res))
	b, _ := json.Marshal(res)
	fmt.Println(string(b))
}
