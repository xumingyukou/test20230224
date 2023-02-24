package c_api

import (
	"clients/exchange/cex/base"
	"fmt"
	"net/url"
	"testing"
)

func TestGetFutureContractInfo(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	res, err := a.GetFutureContractInfo(&url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetFutureContractIndex(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	options := url.Values{}
	options.Add("symbol", "BTC")
	res, err := a.GetFutureContractIndex(&options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetFutureDepth(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	res, err := a.GetFutureDepth("BTC_CW", "step0")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostFutureContractFee(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	post_params := ReqPostFutureContractFee{}
	res, err := a.PostFutureContractFee(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostFutureBalanceValuation(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	post_params := ReqPostFutureBalanceValuation{}
	res, err := a.PostFutureBalanceValuation(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostFutureOrder(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	post_params := ReqPostFutureOrder{
		Volume:         1,
		Direction:      "buy",
		Offset:         "close",
		LeverRate:      1,
		OrderPriceType: "limit",
	}
	res, err := a.PostFutureOrder(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostFutureCancelOrder(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	post_params := ReqPostFutureCancelOrder{
		Symbol: "btc",
	}
	res, err := a.PostFutureCancelOrder(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostFutureContractOrderInfo(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	post_params := ReqPostFutureContractOrderInfo{
		Symbol: "btc",
	}
	res, err := a.PostFutureContractOrderInfo(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostFutureContractOpenorders(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewCApiClient(conf)
	post_params := ReqPostFutureContractOpenorders{
		Symbol: "btc",
	}
	res, err := a.PostFutureContractOpenorders(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
