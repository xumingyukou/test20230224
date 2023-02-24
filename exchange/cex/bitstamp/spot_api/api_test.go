package spot_api

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

var (
	apiClient *ApiClient
	proxyUrl  = "http://127.0.0.1:9999"
)

func init() {
	//config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		//ProxyUrl: proxyUrl,
		//AccessKey: config.ExchangeConfig.ExchangeList["bitstamp"].ApiKeyConfig.AccessKey,
		//SecretKey: config.ExchangeConfig.ExchangeList["bitstamp"].ApiKeyConfig.SecretKey,
	}
	apiClient = NewApiClient(conf)
}

func TestWSTokens(t *testing.T) {
	fmt.Println(apiClient.GetWebsocketsToken())
}

func TestTradingPairs(t *testing.T) {
	fmt.Println(apiClient.GetTradingPairsInfo())
}

func TestOrderbook(t *testing.T) {
	res, _ := apiClient.GetOrderbook("ETH/USD")
	fmt.Printf("%#v\n", res)
}
