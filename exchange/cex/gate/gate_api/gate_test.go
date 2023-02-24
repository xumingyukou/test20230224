package gate_api

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

var (
	c        *ClientGate
	ProxyUrl = "http://127.0.0.1:7890"
)

func init() {
	conf := base.APIConf{
		ProxyUrl: ProxyUrl,
		IsTest:   false,
	}
	c = NewClientGate(conf)
}

func Test_Market_Books(t *testing.T) {
	params := url.Values{}
	params.Add("limit", "10")
	params.Add("with_id", "false")
	res, err := c.Market_Books_Info("BTC_USDT", &params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	a, _ := json.Marshal(res)
	fmt.Println(string(a))
}

func Test_GetSymbols(t *testing.T) {
	params := url.Values{}
	res, err := c.SymbolsInfo(&params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	a, _ := json.Marshal(res)
	fmt.Println(string(a))
}
