package bitMex_api

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

var (
	c *ClientBitMex
)

func init() {
	conf := base.APIConf{
		ProxyUrl: ProxyUrl,
		IsTest:   false,
	}
	c = NewClientBitMex(conf)
}

func Test_Symbol_Info(t *testing.T) {
	params := url.Values{}
	params.Add("typ", "IFXXXP")
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

func TestDepth(t *testing.T) {
	params := url.Values{}
	params.Add("symbol", "ETHUSDT")
	res, err := c.DepthInfo(&params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	a, _ := json.Marshal(res)
	fmt.Println(string(a))
}
