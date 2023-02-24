package spot_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"fmt"
	"testing"
	"time"
)

func TestGetMarkets(t *testing.T) {
	conf := base.APIConf{}
	a := NewApiClient(conf)
	res, err := a.GetMarkets()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range res.Result {
		fmt.Printf("%#v\n", i)
	}
}

func TestGetAccount(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	res, err := api.GetAccountInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetBalances(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	res, err := api.GetBalances()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPlaceOrder(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	params["price"] = "10"
	res, err := api.PlaceOrder("FTT/USD", "buy", "limit", 0.5, &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPlaceOrder2(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	params["price"] = nil
	res, err := api.PlaceOrder("ETH/USDT", "sell", "market", 0.005, &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestOrderHistory(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	//params["market"] = "ETH-PERP"
	//params["orderType"] = "limit"
	res, err := api.GetOrderHistory(&params)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%#v\n", res)
	for _, i := range res.Result {
		fmt.Println(i)
	}
}

func TestCancelOrder(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	orderId := 166471758605
	res, err := api.CancelOrder(orderId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestOrderbook(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})

	res, err := api.GetOrderbook("ETH-PERP", 35)
	if err != nil {
		t.Fatal(err)
	}
	for i, _ := range res.Result.Asks {
		fmt.Println("Ask - Price: ", res.Result.Asks[i][0], " || Size: ", res.Result.Asks[i][1])
	}
	for i, _ := range res.Result.Bids {
		fmt.Println("Bid - Price: ", res.Result.Asks[i][0], " || Size: ", res.Result.Asks[i][1])
	}
}

func TestApiClient_GetWithdrawalFee(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	res, err := api.GetWithdrawalFee("USDT", 35, "0x83a127952d266A6eA306c40Ac62A4a70668FE3BE", &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestOrderStatus(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	orderId := 166456041589
	res, err := api.GetOrderStatus(orderId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestApiClient_GetOpenOrders(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	res, err := api.GetOpenOrders(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestApiClient_GetFills(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	//params["market"] = "FTT/USD"
	res, err := api.GetFills(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestApiClient_GetWithdrawalHistory(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey:  config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey:  config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
		SubAccount: "ftx_001",
	})
	params := make(map[string]interface{})
	params["start_time"] = 1564146934
	res, err := api.GetWithdrawalHistory(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestApiClient_GetDepositHistory(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	params["start_time"] = 1564146934

	res, err := api.GetDepositHistory(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestApiClient_GetDepositAddress(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	params := make(map[string]interface{})
	//params["start_time"] = 1564146934
	res, err := api.GetDepositAddress("USDT", &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetPositions(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	res, err := api.GetPositions()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range res.Result {
		fmt.Println(item)
	}
}

func TestApiClient_TransferSubAccount(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	})
	res, err := api.TransferSubAccount("USD", 1, "ftx_001", "main")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Result)
}

func syncRun(api *ApiClient) {
	res, err := api.GetMarkets()
	fmt.Println(err, res)
}

func TestWeight(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
		ProxyUrl:  "http://127.0.0.1:7890",
	})
	for i := 0; i < 1000; i++ {
		go syncRun(api)
		time.Sleep(time.Millisecond * 10)
	}
}

func syncRun2(api *ApiClient) {
	var option map[string]interface{}
	res, err := api.PlaceOrder("SRM-USD", SIDE_TYPE_SELL, ORDER_TYPE_MARKET, 0.001, &option)
	fmt.Println(err, res)
}
func TestWeight2(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	api := NewApiClient(base.APIConf{
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
		ProxyUrl:  "http://127.0.0.1:7890",
	}, config.ExchangeConfig.ExchangeList["ftx"].Weight)
	for i := 0; i < 1000; i++ {
		go syncRun2(api)
		time.Sleep(time.Millisecond * 10)
	}
}
