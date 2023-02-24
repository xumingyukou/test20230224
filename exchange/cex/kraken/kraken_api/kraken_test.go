package kraken_api

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
)

var (
	c *ClientKraken
)

func init() {
	conf := base.APIConf{
		ProxyUrl: ProxyUrl,
		IsTest:   false,
	}
	c = NewClientKraken(conf)
}

func Test_Asset_Info(t *testing.T) {
	params := url.Values{}
	//params.Add()
	res, err := c.Asset_Pairs_Info(&params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	a, _ := json.Marshal(res)
	fmt.Println(string(a))
}
