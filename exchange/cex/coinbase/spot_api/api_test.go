package spot_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

var (
	apiClient *ApiClient
	ProxyUrl  = "http://127.0.0.1:9999"
)

func init() {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		//ProxyUrl: ProxyUrl,
		//AccessKey: config.ExchangeConfig.ExchangeList["coinbase"].ApiKeyConfig.AccessKey,
		//SecretKey: config.ExchangeConfig.ExchangeList["coinbase"].ApiKeyConfig.SecretKey,
	}
	apiClient = NewApiClient(conf)
}

func TestGetBook(t *testing.T) {
	res, err := apiClient.GetBook("BTC-USD", 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
}

func TestGetProducts(t *testing.T) {
	var symbols []string
	res, err := apiClient.GetProducts()
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range res.ProductInfos {
		symbols = append(symbols, item.DisplayName)
	}

	fmt.Println(symbols)
}
