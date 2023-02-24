package huobi

import (
	"clients/config"
	"clients/exchange/cex/base"
	"fmt"
	"testing"

	"github.com/valyala/fasthttp"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

var (
	ProxyUrl   = "http://127.0.0.1:1080"
	huobi      *ClientHuobi
	timeOffset int64 = 10
	conf       base.APIConf
)

func init() {
	config.LoadExchangeConfig("../../../conf/exchange.toml")
	conf = base.APIConf{
		ReadTimeout: timeOffset,
		ProxyUrl:    ProxyUrl,
		IsTest:      true,
		AccessKey:   config.ExchangeConfig.ExchangeList["huobi"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["huobi"].ApiKeyConfig.SecretKey,
	}
	huobi = NewClientHuobi(conf)
}

func TestGetSymbols(t *testing.T) {
	symbols := huobi.GetSymbols()
	fmt.Println(symbols)
}

func TestGetMarginIsolatedBalance(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewClientHuobi(conf)
	balance, err := a.GetMarginIsolatedBalance()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(balance)
}

func TestGetTradeFee(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewClientHuobi(conf)
	res, err := a.GetTradeFee("btcusdt", "ethbtc")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.TradeFeeList)
}

func TestGetPrecision(t *testing.T) {
	conf := base.APIConf{}
	a := NewClientHuobi(conf)
	res, err := a.GetPrecision()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetDepth(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewClientHuobi(conf)
	symbol := &client.SymbolInfo{
		Symbol: "USTC/USDT",
	}
	dep, err := a.GetDepth(symbol, 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dep)
}

func TestExchangeEnable(t *testing.T) {
	conf := base.APIConf{}
	a := NewClientHuobi(conf)
	res1 := a.IsExchangeEnable()
	fmt.Println(res1)
}

func TestFundingFee(t *testing.T) {
	//urlC := "https://api.hbdm.com/swap-api/v1/swap_batch_funding_rate"
	urlU := "https://api.hbdm.com/linear-swap-api/v1/swap_batch_funding_rate"
	_, resp, err := fasthttp.Get(nil, urlU)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resp))
}

func TestGetTransferFee(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewClientHuobi(conf)
	res, err := a.GetTransferFee(1, "btc", "eth")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetBalance(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewClientHuobi(conf)
	res, err := a.GetBalance()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetFutureSymbols(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewClientHuobi(conf)
	res := a.GetFutureSymbols(common.Market_SWAP_COIN)
	fmt.Println(res)
}

func TestGetFutureDepth(t *testing.T) {
	symbolInfo := &client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_FUTURE_THIS_WEEK,
	}
	res, err := huobi.GetFutureDepth(symbolInfo, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestMoveAsset(t *testing.T) {
	req := &order.OrderMove{
		Asset:  "usdt",
		Amount: 10,
		Source: common.Market_FUTURE,
		Target: common.Market_SPOT,
	}
	res, err := huobi.MoveAsset(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
