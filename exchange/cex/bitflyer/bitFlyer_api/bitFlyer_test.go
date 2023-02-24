package bitFlyer_api

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

var (
	c *ClientBitFlyer
)

func init() {
	conf := base.APIConf{
		ProxyUrl: ProxyUrl,
		IsTest:   false,
	}
	c = NewClientBitFlyer(conf)
}

func Test_Market_Books(t *testing.T) {
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
