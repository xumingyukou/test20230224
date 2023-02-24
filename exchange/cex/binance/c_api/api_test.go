package c_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/spot_api"
	"encoding/json"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"net/url"
	"testing"
)

var (
	proxyUrl = "http://127.0.0.1:7890"
)

func TestExchangeInfo(t *testing.T) {
	/*
		"BTCUSD_PERP", "BTCUSD_220930", "BTCUSD_221230", "ETHUSD_PERP", "ETHUSD_220930", "ETHUSD_221230", "LINKUSD_PERP", "BNBUSD_PERP", "TRXUSD_PERP", "DOTUSD_PERP", "ADAUSD_PERP", "EOSUSD_PERP", "LTCUSD_PERP", "BCHUSD_PERP", "XRPUSD_PERP", "ETCUSD_PERP", "FILUSD_PERP", "EGLDUSD_PERP", "DOGEUSD_PERP", "UNIUSD_PERP", "THETAUSD_PERP", "XLMUSD_PERP", "SOLUSD_PERP", "FTMUSD_PERP", "SANDUSD_PERP", "MANAUSD_PERP", "AVAXUSD_PERP", "GALAUSD_PERP", "MATICUSD_PERP", "NEARUSD_PERP", "ATOMUSD_PERP", "AAVEUSD_PERP", "AXSUSD_PERP", "ROSEUSD_PERP", "XTZUSD_PERP", "ICXUSD_PERP", "ALGOUSD_PERP", "RUNEUSD_PERP", "ADAUSD_220930", "LINKUSD_220930", "BCHUSD_220930", "DOTUSD_220930", "XRPUSD_220930", "LTCUSD_220930", "BNBUSD_220930", "APEUSD_PERP", "VETUSD_PERP", "ZILUSD_PERP", "KNCUSD_PERP", "XMRUSD_PERP", "GMTUSD_PERP", "ADAUSD_221230", "LINKUSD_221230", "BCHUSD_221230", "DOTUSD_221230", "XRPUSD_221230", "LTCUSD_221230", "BNBUSD_221230", "OPUSD_PERP", "ENSUSD_PERP", --- PASS: TestExchangeInfo (0.56s)
	*/
	api := NewCApiClient(base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
	})
	res, err := api.ExchangeInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Timezone, res.ServerTime)
	for _, rateLimit := range res.RateLimits {
		fmt.Printf("%#v\n", rateLimit)
	}
	swapSym := make([]string, 0)
	futureSym := make([]string, 0)
	for _, symbols := range res.Symbols {
		sym, market, symType := GetContractType(symbols.Symbol)
		fmt.Println(symbols.Symbol, sym, market, symType)
		if market == common.Market_SWAP_COIN {
			swapSym = append(swapSym, fmt.Sprintf("%s_%d_%d", sym, market.Number(), symType.Number()))
		} else if market == common.Market_FUTURE_COIN {
			futureSym = append(futureSym, fmt.Sprintf("%s_%d_%d", sym, market.Number(), symType.Number()))
		}
	}
	a, _ := json.Marshal(swapSym)
	fmt.Println(len(swapSym))
	fmt.Println(string(a))

	b, _ := json.Marshal(futureSym)
	fmt.Println(len(futureSym))
	fmt.Println(string(b))
}

func TestGetDepth(t *testing.T) {
	api := NewCApiClient(base.APIConf{
		ProxyUrl: proxyUrl,
	})
	res, err := api.GetDepth("BTC/USD_PERP", 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetTrades(t *testing.T) {
	api := NewCApiClient(base.APIConf{ProxyUrl: proxyUrl})
	res, err := api.GetTrades("BTC/USD_PERP", 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPremiumIndex(t *testing.T) {
	api := NewCApiClient(base.APIConf{ProxyUrl: proxyUrl})
	res, err := api.PremiumIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	op := &url.Values{}
	op.Add("symbol", "BTCUSDT")
	res, err = api.PremiumIndex(op)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestCommissionRate(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewCApiClient(base.APIConf{
		ProxyUrl:  proxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	})
	res, err := api.CommissionRate("BTCUSD_220930")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestUOrder(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewCApiClient(base.APIConf{
		ProxyUrl:  proxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	})
	options := &url.Values{}
	options.Add("quantity", "1")
	options.Add("timeInForce", "GTC")
	options.Add("price", "1900")
	//options.Add("positionSide", "LONG")
	//options.Add("reduceOnly", "false")

	res, err := api.Order("ETHUSD_PERP", spot_api.SIDE_TYPE_BUY, spot_api.ORDER_TYPE_LIMIT, options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetOrder(t *testing.T) {

	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewCApiClient(base.APIConf{
		ProxyUrl:  proxyUrl,
		EndPoint:  spot_api.UBASE_TEST_BASE_URL,
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	})
	options := &url.Values{}
	options.Add("orderId", "491141")
	res, err := api.GetOrder("BTCUSDT_220930", options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetPositionSideDual(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewCApiClient(base.APIConf{
		ProxyUrl:  proxyUrl,
		EndPoint:  spot_api.CBASE_API_BASE_URL,
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	})
	options := &url.Values{}
	options.Add("orderId", "491141")
	res, err := api.PositionSideDual("false")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestCBaseBalance(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewCApiClient(base.APIConf{
		ProxyUrl:  proxyUrl,
		EndPoint:  spot_api.CBASE_API_BASE_URL,
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	})
	options := &url.Values{}
	options.Add("orderId", "491141")
	res, err := api.Balance()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range *res {
		fmt.Printf("%#v\n", i)
	}
}
