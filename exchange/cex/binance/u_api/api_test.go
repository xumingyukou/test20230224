package u_api

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

func TestExchangeInfo(t *testing.T) {
	api := NewUApiClient(base.APIConf{
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
		if market == common.Market_SWAP {
			swapSym = append(swapSym, fmt.Sprintf("%s_%d_%d", sym, market.Number(), symType.Number()))
		} else if market == common.Market_FUTURE {
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
	api := NewUApiClient(base.APIConf{})
	res, err := api.GetDepth("BTCUSDT_221230", 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetTrades(t *testing.T) {
	api := NewUApiClient(base.APIConf{})
	res, err := api.GetTrades("BTC/USDT", 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPremiumIndex(t *testing.T) {
	api := NewUApiClient(base.APIConf{})
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
	api := NewUApiClient(base.APIConf{
		ProxyUrl:  "http://127.0.0.1:7890",
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	})
	res, err := api.CommissionRate("BTCUSDT_220930")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestUOrder(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewUApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	})
	options := &url.Values{}
	options.Add("quantity", "5")
	options.Add("timeInForce", "GTX")
	options.Add("price", "19299")
	options.Add("positionSide", "LONG")
	res, err := api.Order("BTCUSDT", spot_api.SIDE_TYPE_SELL, spot_api.ORDER_TYPE_LIMIT, options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetOrder(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewUApiClient(base.APIConf{
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
