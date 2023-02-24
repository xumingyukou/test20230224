package spot_api

import (
	"clients/exchange/cex/base"
	"clients/transform"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

var (
	c        *ApiClient
	proxyUrl = "http://127.0.0.1:7890"
)

func init() {
	//config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl: proxyUrl,
		//AccessKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.AccessKey,
		AccessKey: "Ig0T9qTB4h8heN5aCo",
		SecretKey: "djSxj25DX8cHBAQiizRUdORq9BxIwML1NWZH",
		//Passphrase: config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.Passphrase,
		IsTest: false,
	}
	c = NewApiClient(conf)
}

func Test_Client_GetSymbols(t *testing.T) {
	for i := 0; i < 1; i++ {
		go RunCoroutine()
	}
	for {
	}
}

func RunCoroutine() {
	for i := 0; i < 50; i++ {
		res, err := c.GetSymbols()
		fmt.Println(res)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestApiClient_GetWithdrawFee(t *testing.T) {
	res, err := c.GetWithdrawFee()
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_OrderBook(t *testing.T) {
	res, err := c.GetOrderBook("BTC/USDT")
	if err != nil {
		fmt.Println(err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_ServerTime(t *testing.T) {
	res, err := c.GetTime(nil)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
}

func Test_FundFee(t *testing.T) {
	res, err := c.GetFundFee("BTCUSD")
	if err != nil {
		fmt.Println(err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_AccountBalance(t *testing.T) {
	res, err := c.GetAccountBalance()
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_PlaceOrder(t *testing.T) {
	param := url.Values{}
	param.Add("orderPrice", "1")
	res, err := c.PostPlaceOrder("ETHUSDT", "1", "Buy", "LIMIT", "test1", param)
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_CancleOrder(t *testing.T) {
	res, err := c.PostCancleOrder("test")
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_OrderInfo(t *testing.T) {
	res, err := c.GetOrderInfo("test")
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_OrderList(t *testing.T) {
	res, err := c.GetOrderList()
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_OrderHistory(t *testing.T) {
	res, err := c.GetOrderHistory()
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_Withdraw(t *testing.T) {
	res, err := c.Withdraw(nil)
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_WithdrawHistory(t *testing.T) {
	res, err := c.WithdrawHistory()
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_Transfer(t *testing.T) {
	res, err := c.Transfer("1", "USDT", "1", "SPOT", "CONTRACT")
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_AllTransfer(t *testing.T) {
	res, err := c.ALLTransfer("21ff1b44-2d5d-4293-913d-4597c5ad2611", "USDT", "1", "290118", "545366", "SPOT", "CONTRACT")
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_TransferHistory(t *testing.T) {
	res, err := c.TransferHistory(nil)
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_TransferS2MHistory(t *testing.T) {
	res, err := c.TransferHistoryS2M(nil)
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_RecordHistory(t *testing.T) {
	res, err := c.RecordHistory(transform.XToString(time.Now().UnixMilli()))
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_Loan(t *testing.T) {
	res, err := c.PostLoan("USDT", "1")
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_LoanHistory(t *testing.T) {
	res, err := c.GetLoanHistory()
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_Repay(t *testing.T) {
	res, err := c.PostRepay("USDT", "1")
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}

func Test_RepayHistory(t *testing.T) {
	res, err := c.GetRepayHistory()
	if err != nil {
		fmt.Println("error", err)
	}

	s, _ := json.Marshal(res)
	fmt.Println(string(s))
}
