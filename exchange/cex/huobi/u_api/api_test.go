package u_api

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

func TestPostSwapSwitchAccountType(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewUApiClient(conf)
	post_params := ReqPostSwapSwitchAccountType{
		AccountType: 1,
	}
	res, err := a.PostSwapSwitchAccountType(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetSwapContractInfo(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewUApiClient(conf)
	res, err := a.GetSwapContractInfo("swap")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
