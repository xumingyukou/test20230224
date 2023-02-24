package spot_api

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestExchangeInfo(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.ExchangeInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetDepth(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetDepth("ustcusdt", 5)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetMarkets(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetMarketStatus()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetAllSymbols(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetAllSymbols(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetAllCurrencies(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetAllCurrencies(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetCurrencysSettings(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetCurrencysSettings(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetSymbolsSettings(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetSymbolsSettings(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetMarketSymbolsSettings(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetMarketSymbolsSettings(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetChainsSettings(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetChainsSettings(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetCurrenciesChains(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetCurrenciesChains(&params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetKline(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetKline("btcusdt", "1day", &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetMerged(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetMerged("btcusdt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetAllTickers(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetAllTickers()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetHistoryTrade(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	params := url.Values{}
	res, err := a.GetHistoryTrade("btcusdt", &params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetDetail(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetDetail("btcusdt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", res)
}

func TestGetTrade(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetTrade("btcusdt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetAccount(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.Account()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestAccountBalance(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.AccountBalance("spot")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestAssetTradeFee(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl:  "http://127.0.0.1:1080",
		AccessKey: "",
	}
	a := NewApiClient(conf)
	res, err := a.AssetTradeFee("btcusdt, ethusdt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestOrder(t *testing.T) { // account-frozen-account-inexistent-error
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	post_params := ReqPostOrder{
		AccountID:     SPOT_ACCOUNT_ID,
		Amount:        "10.1",
		Price:         "100.1",
		Source:        "spot-api",
		Symbol:        "ethusdt",
		Type:          "buy-limit",
		ClientOrderID: "a0001",
	}
	res, err := a.PostOrder(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetValuation(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := url.Values{}
	res, err := a.GetValuation(&options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetAssetValuation(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := url.Values{}
	res, err := a.GetAssetValuation("spot", &options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetAccountHistory(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := url.Values{}
	res, err := a.GetAccountHistory(SPOT_ACCOUNT_ID, &options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetAccountLedger(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := url.Values{}
	res, err := a.GetAccountLedger(SPOT_ACCOUNT_ID, &options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostFuturesTransfer(t *testing.T) { // Invalid request
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	post_params := ReqPostCFuturesTransfer{
		Currency: "btc",
		Amount:   0.001,
		Type:     "pro-to-futures",
	}
	res, err := a.PostCFuturesTransfer(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetPointAccount(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := url.Values{}
	res, err := a.GetPointAccount(&options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

// func TestPostPointTransfer(t *testing.T) { // 只有一个账户无法传参
// 	conf := base.APIConf{
// 		ProxyUrl: "http://127.0.0.1:1080",
// 	}
// 	a := NewApiClient(conf)
// 	res, err := a.PostPointTransfer()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(res)
// }

func TestGetDepositAddress(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetDepositAddress("btc")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetWithdrawAddress(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := url.Values{}
	res, err := a.GetWithdrawAddress("btc", &options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostWithdrawCreate(t *testing.T) { //API withdrawal does not support temporary addresses
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	post_params := ReqPostWithdrawCreate{
		Address:  "0xde709f2102306220921060314715629080e2fb77",
		Amount:   "0.05",
		Currency: "eth",
		Fee:      "0.01",
	}
	res, err := a.PostWithdrawCreate(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetWithdrawClientOrderId(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.GetWithdrawClientOrderId("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

// func TestPostCancelWithdrawCreate(t *testing.T) { // 没有提币ID，没法测
// 	conf := base.APIConf{
// 		ProxyUrl: "http://127.0.0.1:1080",
// 	}
// 	a := NewApiClient(conf)
// 	res, err := a.PostCancelWithdrawCreate("")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	fmt.Println(res)
// }

func TestPostSubUserCreation(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	userList := ReqPostSubUserCreation{
		[]struct {
			UserName string `json:"userName"`
			Note     string `json:"note"`
		}{
			{"subuser1412", "sub-user-1-note"},
			{"subuser1413", "sub-user-2-note"},
		},
	}
	res, err := a.PostSubUserCreation(userList)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostBatchOrders(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	post_params := ReqPostBatchOrders{
		{
			AccountID:     SPOT_ACCOUNT_ID,
			Amount:        "10.1",
			Price:         "100.1",
			Source:        "spot-api",
			Symbol:        "ethusdt",
			Type:          "buy-limit",
			ClientOrderID: "a0001",
		},
		{
			AccountID: SPOT_ACCOUNT_ID,
			Type:      "buy-limit",
			Source:    "spot-api",
			Symbol:    "btcusdt",
			Price:     "1.1",
			Amount:    "1",
		},
	}
	res, err := a.PostBatchOrders(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostSubmitOrder(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	post_params := ReqPostSubmitOrder{
		Symbol: "btcusdt",
	}
	res, err := a.PostSubmitOrder("1", post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostMarginOrders(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	post_params := ReqPostMarginOrders{
		Symbol:   "ethusdt",
		Currency: "eth",
		Amount:   "1.0",
	}
	res, err := a.PostMarginOrders(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostCrossMarginOrders(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	post_params := ReqPostCrossMarginOrders{
		Currency: "eth",
		Amount:   "1.0",
	}
	res, err := a.PostCrossMarginOrders(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetMarginLoanOrders(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := &url.Values{}
	res, err := a.GetMarginLoanOrders("ethusdt", options)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetReferenceCurrencies(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	options := &url.Values{}
	res, _ := a.GetReferenceCurrencies(options)

	Chain_value := map[string]int32{
		"INVALID_CAHIN": 0,
		"BTC":           1,
		"ETH":           2,
		"BSC":           3,
		"AVALANCHE":     4,
		"SOLANA":        5,
		"FANTOM":        6,
		"TRON":          7,
		"POLYGON":       8,
		"ARBITRUM":      9,
		"CRONOS":        10,
		"HECO":          11,
		"OKC":           12,
		"BNB":           13,
		"OPTIMISM":      14,
	}

	result := []RespTestGetReferenceCurrencies{}
	for _, data := range res.Data {
		result_tmp := RespTestGetReferenceCurrencies{}
		result_tmp.Currency = data.Currency
		for _, chain := range data.Chains {
			if _, ok := Chain_value[chain.BaseChain]; ok {
				result_tmp_tmp := Chain{
					Chain:       chain.Chain,
					DisplayName: chain.DisplayName,
					BaseChain:   chain.BaseChain,
				}
				result_tmp.Chains = append(result_tmp.Chains, result_tmp_tmp)
			} else if _, ok := Chain_value[chain.DisplayName]; ok {
				result_tmp_tmp := Chain{
					Chain:       chain.Chain,
					DisplayName: chain.DisplayName,
					BaseChain:   chain.BaseChain,
				}
				result_tmp.Chains = append(result_tmp.Chains, result_tmp_tmp)
			}
		}
		if result_tmp.Chains != nil {
			result = append(result, result_tmp)
		}

	}

	for _, data := range result {
		fmt.Printf("%+v", data)
	}

	filename := "chains.json"
	filePtr, _ := os.Create(filename)
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	encoder.Encode(result)

}

func TestGetOrder(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	order_id := "718603211381256"
	res, err := a.GetOrder(order_id)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostTransfer(t *testing.T) {
	post_params := ReqPostCFuturesTransfer{
		Currency: "fil",
		Amount:   1,
		Type:     "futures-to-pro",
	}
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.PostCFuturesTransfer(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPostCSwapTransfer(t *testing.T) {
	post_params := ReqPostCSwapTransfer{
		Currency: "fil",
		Amount:   1,
		From:     "swap",
		To:       "spot",
	}
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:1080",
	}
	a := NewApiClient(conf)
	res, err := a.PostCSwapTransfer(post_params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
