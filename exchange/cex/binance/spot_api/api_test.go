package spot_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/logger"
	"fmt"
	"net/url"
	"strconv"
	"testing"
	"time"
)

var (
	apiClient *ApiClient
	proxyUrl  = "http://127.0.0.1:7890"
)

func init() {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl:  proxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	}
	apiClient = NewApiClient(conf, config.ExchangeConfig.ExchangeList["binance"].Weight)
}

func TestPing(t *testing.T) {
	fmt.Println(apiClient.Ping())
}

func TestExchangeInfo(t *testing.T) {
	res, err := apiClient.ExchangeInfo()
	if err != nil {
		t.Fatal(err)
	}
	logger.SaveToFile("exchange.json", res)
}

func TestGetAPIKEY(t *testing.T) {
	binanceConf := config.ExchangeConfig.ExchangeList["binance"]
	fmt.Println(binanceConf.SpotConfig)
}

func TestCapitalConfigGetAll(t *testing.T) {
	configAll, err := apiClient.CapitalConfigGetAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range *configAll {
		for _, network := range item.NetworkList {
			var (
				fee float64
			)
			fee, err = strconv.ParseFloat(network.WithdrawFee, 64)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println(network.Coin, network.Network, fee)
		}
	}
}

func TestAllNetwork(t *testing.T) {
	configAll, err := apiClient.CapitalConfigGetAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range *configAll {
		logger.SaveToFile("network.json", item.NetworkList)
	}
}

func TestTrade(t *testing.T) {
	var (
		params = url.Values{}
	)
	params.Add("quantity", "20")
	res, err := apiClient.Order("BUSDUSDT", SIDE_TYPE_SELL, ORDER_TYPE_MARKET, &params)
	if err != nil {
		t.Fatal(err)
	}
	resp, ok := res.(*RespMarginOrderFull)
	if !ok {
		t.Fatal("convert response err")
	}
	fmt.Println(resp)
}

func TestLimitTrade(t *testing.T) {
	var (
		params = url.Values{}
	)
	params.Add("timeInForce", "GTC")
	params.Add("quantity", "11")
	params.Add("price", "1")

	res, err := apiClient.Order("BUSDUSDT", SIDE_TYPE_BUY, ORDER_TYPE_LIMIT, &params)
	if err != nil {
		t.Fatal(err)
	}
	resp, ok := res.(*RespMarginOrderFull)
	if !ok {
		t.Fatal("convert response err")
	}
	fmt.Println(resp)
}

func TestGetBalance(t *testing.T) {
	res, err := apiClient.Account()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
}

func TestGetOrder(t *testing.T) {
	params := url.Values{}
	params.Add("orderId", "9292172565")
	res, err := apiClient.GetOrder("ETHUSDT", &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", res)
}

func TestMarginLoanHistory(t *testing.T) {
	res, err := apiClient.MarginLoanHistory("USDT", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range res.Rows {
		fmt.Printf("%#v\n", item)
	}
}

func TestUserDataStream(t *testing.T) {
	res, err := apiClient.UserDataStream()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	_, err = apiClient.PutUserDataStream(res.ListenKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = apiClient.DELETEUserDataStream(res.ListenKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("done")
}

func TestMarginUserDataStream(t *testing.T) {
	res, err := apiClient.MarginUserDataStream()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	_, err = apiClient.PutMarginUserDataStream(res.ListenKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = apiClient.DELETEMarginUserDataStream(res.ListenKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("done")
}

func TestMarginIsolatedUserDataStream(t *testing.T) {
	symbol := "BTCBNB"
	res, err := apiClient.MarginIsolatedUserDataStream(symbol)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	_, err = apiClient.PutMarginIsolatedUserDataStream(symbol, res.ListenKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = apiClient.DELETEMarginIsolatedUserDataStream(symbol, res.ListenKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("done")
}

func TestIsolatedAllPairs(t *testing.T) {
	res, err := apiClient.MarginIsolatedAllPairs()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range res {
		fmt.Printf("%#v\n", item)
	}
}

func TestAssetTransfer(t *testing.T) {
	var (
		asset          = "USDT" // asset	STRING	YES
		amount float64 = 0.9    // amount	DECIMAL	YES
		params         = url.Values{}
	)

	res, err := apiClient.AssetTransfer(MOVE_TYPE_MARGIN_FUNDING, asset, amount, &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res:", res)
}

func TestLoan(t *testing.T) {
	var (
		asset          = "USDT" //asset	STRING	YES
		amount float64 = 10     //amount	DECIMAL	YES
		params         = url.Values{}
	)

	res, err := apiClient.MarginLoan(asset, amount, &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res:", res)
}

func TestRepay(t *testing.T) {
	var (
		asset          = "USDT" //asset	STRING	YES
		amount float64 = 15     //amount	DECIMAL	YES
		params         = url.Values{}
	)

	res, err := apiClient.MarginRepay(asset, amount, &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("res:", res)
}

func TestWeight(t *testing.T) {
	var count int64
	for {
		fmt.Println(apiClient.Ping())
		fmt.Println(apiClient.GetDepth("BTCUSDT", 1000))
		type_ := "ip_1m"
		cap_, remain := BinanceSpotWeight.GetInstance(base.RateLimitType(type_)).Instance.Remain()
		fmt.Println(count, "weight:", type_, BinanceSpotWeight.GetInstance(base.RateLimitType(type_)), cap_, remain)
		time.Sleep(time.Second)
		count++
		if count > 10000 {
			break
		}
	}
	/*
		使用滑动窗口，在窗口滑动时会有50us的延迟，令牌桶则不会，比较稳定
		平均时间：滑动窗口14us，令牌桶7us
	*/
}
