package ok_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
)

var (
	c        *ClientOkex
	proxyUrl = "http://127.0.0.1:7890"
)

func init() {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl:   proxyUrl,
		AccessKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.AccessKey,
		SecretKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.SecretKey,
		Passphrase: config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.Passphrase,
		IsTest:     false,
	}
	c = NewClientOkex(conf)
}

func Test_ClientOkex_GetSymbols(t *testing.T) {
	var (
		res *RespInstruments
		err error
	)
	res, err = c.Instrument_Info("SPOT", nil)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
	var symbols []string

	for _, instrument := range res.Data {
		if instrument.State == "live" {
			words := strings.Split(instrument.InstId, "-")
			symbols = append(symbols, fmt.Sprint(words[0], "/", words[1]))
		}
	}
	fmt.Println(len(symbols))

	//res, err = c.Instrument_Info("MARGIN", nil)
	//if err != nil {
	//	fmt.Println(err)
	//	t.Fatal(err)
	//}
	//var symbols1 []string
	//for _, instrument := range res.Data {
	//	if instrument.State == "live" {
	//		words := strings.Split(instrument.InstId, "-")
	//		symbols1 = append(symbols1, fmt.Sprint(words[0], "/", words[1]))
	//	}
	//}
	//fmt.Println(len(symbols1))
	//{
	//	type void struct{}
	//	set := make(map[string]void)
	//	var member void
	//	for _, word := range symbols {
	//		set[word] = member
	//	}
	//	for _, word := range symbols1 {
	//		if _, ok := set[word]; !ok {
	//			fmt.Println(word)
	//			fmt.Println("false")
	//		}
	//	}
	//	fmt.Println("true")
	//	fmt.Println(len(set))
	//}
}

func Test_Delivery_Exercise_History(t *testing.T) {
	for i := 0; i < 10; i++ {
		go go1()
		go go2()
	}
	for {

	}
	res, err := c.Delivery_Exercise_History_Info("OPTION", "BTC-USD", nil)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func go1() {
	for i := 0; i < 100; i++ {
		res, err := c.Delivery_Exercise_History_Info("FUTURES", "BTC-USD", nil)
		if err != nil {
			fmt.Println("错误：", err.Error())
			err = nil
		} else {
			v, _ := json.Marshal(res)
			fmt.Println(string(v))
		}
	}
}

func go2() {
	for i := 0; i < 100; i++ {
		res, err := c.Delivery_Exercise_History_Info("FUTURES", "BTC-USD", nil)
		if err != nil {
			fmt.Println("错误：", err.Error())
			err = nil
		} else {
			v, _ := json.Marshal(res)
			fmt.Println(string(v))
		}
	}
}

func Test_User_SetSubAccountOut(t *testing.T) {
	param := url.Values{}
	param.Add("subAcct", "test003,test004")
	res, err := c.User_SetTransferOut(&param)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_ServerTime(t *testing.T) {
	res, err := c.Public_ServerTime_Info(nil)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_MarkPrice(t *testing.T) {
	for i := 0; i < 100; i++ {
		var options = &url.Values{}
		//if i%2 == 0 {
		options.Add("instId", "BTC-USD-SWAP")
		//}
		res, err := c.Public_MarkPrice_Info("SWAP", options)
		if err != nil {
			fmt.Println(err)
			t.Fatal(err)
		}
		fmt.Println(res)
		//for _, instrument := range res.Data {
		//	a, _ := json.Marshal(instrument)
		//	fmt.Println(string(a))
		//}
		//time.Sleep(time.Millisecond * 10)
	}
}

func Test_Open_Interest(t *testing.T) {
	res, err := c.Open_Interest_Info("SWAP", nil)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Market_Books(t *testing.T) {
	params := url.Values{}
	params.Add("sz", "2")
	res, err := c.Market_Books_Info("BTC-USDT", &params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Market_Trades(t *testing.T) {
	time.Sleep(2 * time.Second)
	res, err := c.Market_Trades_Info("BTC-USDT", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Market_IndexTickers(t *testing.T) {
	res, err := c.Market_IndexTickers_Info("BTC-USDT", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println(res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Market_HistoryTrade(t *testing.T) {
	for i := 0; i < 3; i++ {
		go s()
	}
	for {
	}
	//for _, instrument := range res.Data {
	//	a, _ := json.Marshal(instrument)
	//	fmt.Println(string(a))
	//}
}

func s() {
	for i := 0; i < 10; i++ {
		res, err := c.Market_HistoryTrades_Info("BTC-USDT", nil)
		if err != nil {
			fmt.Println("错误")
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}

func Test_Account_Balance(t *testing.T) {
	res, err := c.Account_Balance_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_Position(t *testing.T) {
	parapms := url.Values{}
	parapms.Add("insType", "FUTURES")
	res, err := c.Account_Positions_Info(&parapms)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_PositionHistory_Reflect(t *testing.T) {

	res, err := c.Account_PositionMarginBalance_Info("BTC-USDT-200626", "short", "add", "1", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}

	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_PositionHistory(t *testing.T) {
	res, err := c.Account_PositionsHistory_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_AccountPositionRisk(t *testing.T) {
	res, err := c.Account_AccountPositionRisk_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_Bills(t *testing.T) {
	res, err := c.Account_Bills_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_BillsArchive(t *testing.T) {
	params := url.Values{}
	//params.Add("begin", "1662566400000")
	//params.Add("end", "1662609600000")
	//params.Add("type", "1")
	res, err := c.Account_BillsArchive_Info(&params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_Config(t *testing.T) {
	res, err := c.Account_Config_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_SetPositionMode(t *testing.T) {
	res, err := c.Account_SetPositionMode_Info("long_short_mode", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_SetLeverage(t *testing.T) {
	params := url.Values{}
	res, err := c.Account_SetLeverage_Info("BTC-USDT", "5", "isolated", &params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_MaxSize(t *testing.T) {
	res, err := c.Account_MaxSize_Info("BTC-USDT", "isolated", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_MaxAvailSize(t *testing.T) {
	res, err := c.Account_MaxAvailSize_Info("BTC-USDT", "isolated", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_PositionMarginBalance(t *testing.T) {
	res, err := c.Account_PositionMarginBalance_Info("BTC-USDT-20062", "short", "add", "1", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_LeverageInfo(t *testing.T) {
	res, err := c.Account_LeverageInfo_Info("BTC-USDT-200626", "cross", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_MaxLoan(t *testing.T) {
	res, err := c.Account_MaxLoan_Info("BTC-USDT", "cross", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_TraddFee(t *testing.T) {
	res, err := c.Account_TradeFee_Info("FUTURES", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
	res, err = c.Account_TradeFee_Info("MARGIN", nil)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_InterestAccrued(t *testing.T) {
	res, err := c.Account_InterestAccrued_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_InterestRate(t *testing.T) {
	res, err := c.Account_InterestRate_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_SetGreeks(t *testing.T) {
	res, err := c.Account_SetGreeks_Info("PA", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_SetIsolatedMode(t *testing.T) {
	res, err := c.Account_SetIsolatedMode_Info("automatic", "MARGIN", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_MaxWithdrawal(t *testing.T) {
	res, err := c.Account_MaxWithdrawal_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_RiskState(t *testing.T) {
	res, err := c.Account_RiskState_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_BorrowRepay(t *testing.T) {
	res, err := c.Account_BorrowRepay_Info("USDT", "borrow", "1", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_BorrowRepayHistory(t *testing.T) {
	res, err := c.Account_BorrowRepayHistory_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_InterestLimits(t *testing.T) {
	res, err := c.Account_InterestLimits_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_SimulatedMargin(t *testing.T) {
	params := url.Values{}
	params.Add("inclRealPos", "true")
	//s := `{"pos":"10"","instId":"BTC-USDT_SWAP"}`
	//params.Add("simPos", s)
	res, err := c.Account_SimulatedMargin_Info(&params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_Greeks(t *testing.T) {
	res, err := c.Account_Greeks_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Account_PositionTiers(t *testing.T) {
	res, err := c.Account_PositionTiers_Info("SWAP", "BTC-USDT", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

// Operation is not supported under the current account mode
func Test_Trade_Order_Info(t *testing.T) {

	var f []Post_Trde_BatchOrders
	p := Post_Trde_BatchOrders{
		InstId:  "ETH/USDT",
		TdMode:  "cash",
		Side:    "buy",
		OrdType: "limit",
		Sz:      "2",
		Px:      "2.25",
		ClOrdId: "153",
	}

	f = append(f, p)
	//s := S2M(f, len(f))
	param := url.Values{}
	param.Add("px", "1")
	res, err := c.Trade_Order_Info(p.InstId, p.TdMode, p.Side, p.OrdType, p.Sz, &param)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_BatchOrders(t *testing.T) {
	var f []Post_Trde_BatchOrders
	p := Post_Trde_BatchOrders{
		InstId:  "BTC-USDT",
		TdMode:  "cash",
		Side:    "buy",
		OrdType: "limit",
		Sz:      "2",
		Px:      "2.25",
		ClOrdId: "h15",
	}
	p2 := Post_Trde_BatchOrders{
		InstId:   "BTC-USDT",
		TdMode:   "cash",
		Side:     "buy",
		OrdType:  "limit",
		Sz:       "2",
		Px:       "2.25",
		ClOrdId:  "f15",
		BanAmend: true,
	}
	f = append(f, p, p2)
	s := S2M(f, len(f))
	res, err := c.Trade_BatchOrders_Info(s)

	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_CancelOrder(t *testing.T) {
	p := Post_Trade_BatchCancelOrders{
		InstId: "BTC-USDT",
		OrdId:  "12312",
	}
	//f = append(f, p)
	//s := S2M(f, len(f))
	param := url.Values{}
	param.Add("ordId", "470256576257990657")
	res, err := c.Trade_CancelOrder_Info(p.InstId, &param)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func TestClientOkex_Trade_BatchCancelOrders(t *testing.T) {
	var f []Post_Trade_BatchCancelOrders
	p := Post_Trade_BatchCancelOrders{
		InstId: "BTC-USDT",
		OrdId:  "469535734538575877",
	}
	p2 := Post_Trade_BatchCancelOrders{
		InstId: "BTC-USDT",
		OrdId:  "469535482171498496",
	}
	f = append(f, p, p2)
	s := S2M(f, len(f))
	res, err := c.Trade_BatchCancelOrders_Info(s)

	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trde_ClosePosition(t *testing.T) {
	res, err := c.Trade_ClosePosition("BTC-USDT-SWAP", "cross", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_OrderInfo_Info(t *testing.T) {
	res, err := c.Trade_OrderInfo_Info("BTC-USDT", "471698501139832832", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_OrderHistoryWeek_Info(t *testing.T) {
	res, err := c.Trade_OrdersHistoryWeek_Info("SPOT", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_OrderHistoryArchive_Info(t *testing.T) {
	res, err := c.Trade_OrdersHistoryArchive_Info("SPOT", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_FillsHistoryThreeDays(t *testing.T) {
	res, err := c.Trade_FillsHistoryTreeDays_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_FillsHistoryArchive(t *testing.T) {
	res, err := c.Trade_FillsHistoryArchive_Info("SPOT", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_OrderAlgo(t *testing.T) {
	params := url.Values{}
	params.Add("tpTriggerPx", "15")
	params.Add("tpOrdPx", "18")
	params.Add("ccy", "USDT")
	res, err := c.Trade_OrderAlgo_Info("BTC-USDT", "cross", "buy", "conditional", "2", &params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_CancelAlgos(t *testing.T) {
	var f []Post_Trade_CancelAlgos
	p := Post_Trade_CancelAlgos{
		InstId: "BTC-USDT",
		AlgoId: "469558772399218691",
	}
	f = append(f, p)
	s := S2M(f, len(f))
	res, err := c.Trade_CancelAlgo_Info(s)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_CancelAdvanceAlgos(t *testing.T) {
	var f []Post_Trade_CancelAlgos
	p := Post_Trade_CancelAlgos{
		InstId: "BTC-USDT",
		AlgoId: "469558772399218691",
	}
	f = append(f, p)
	s := S2M(f, len(f))
	res, err := c.Trade_CancelAdvanceAlgo_Info(s)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_OrderAlgoPending(t *testing.T) {
	res, err := c.Trade_OrdersAlgoPending_Info("oco", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Trade_OrderAlgoHistory(t *testing.T) {
	params := url.Values{}
	params.Add("state", "effective")
	res, err := c.Trade_OrdersAlgoHistory_Info("conditional", &params)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_RFQ_CounterParties(t *testing.T) {
	res, err := c.Rfq_CounterParties_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_Currencies(t *testing.T) {
	res, err := c.Asset_Currencies_Info(nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_Transfer(t *testing.T) {
	param := url.Values{}
	param.Add("type", "2")
	param.Add("subAcct", "test003")
	res, err := c.Asset_Transfer_Info("USDT", "1", "6", "6", &param)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	a, _ := json.Marshal(res)
	fmt.Println("res: ", string(a))
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(":", string(a))
	}
}

func Test_Asset_TransferState(t *testing.T) {
	res, err := c.Asset_TransferState_Info("248155482", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_Bills(t *testing.T) {
	options := &url.Values{}
	options.Add("ccy", "USDT")
	//options.Add("type", "20")
	res, err := c.Asset_Bills_Info(options)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_SubAccountBills(t *testing.T) {
	options := &url.Values{}
	//options.Add("type", "20")
	res, err := c.Asset_SubAccountBills_Info(options)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_Withdrawal(t *testing.T) {
	res, err := c.Asset_Withdrawal_Info("BTC", "1", "4", "7DKe3kkkkiiiiTvAKKi2vMPbm1Bz3CMKw", "0.0005", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_WithdrawalLightning(t *testing.T) {
	res, err := c.Asset_WithdrawalLightning_Info("BTC", "lnbc100u1psnnvhtpp5yq2x3q5hhrzsuxpwx7ptphwzc4k4wk0j3stp0099968m44cyjg9sdqqcqzpgxqzjcsp5hz", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_CancleWithdrawal(t *testing.T) {
	res, err := c.Asset_CancleWithdrawal_Info("lnbc100u1psnnvhtpp5yq2x3q5hhrzsuxpwx7ptphwzc4k4wk0j3stp0099968m44cyjg9sdqqcqzpgxqzjcsp5hz", nil)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	fmt.Println("res: ", res)
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_WithdrawalHistory(t *testing.T) {
	res, err := c.Asset_WithdrawalHistory_Info(nil)
	fmt.Println("res: ", res)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Asset_DepositHistory(t *testing.T) {
	res, err := c.Asset_DepositHistory_Info(nil)
	fmt.Println("res: ", res)
	if err != nil {
		fmt.Println("错误")
		fmt.Println(err)
		t.Fatal(err)
	}
	for _, instrument := range res.Data {
		a, _ := json.Marshal(instrument)
		fmt.Println(string(a))
	}
}

func Test_Weight(t *testing.T) {
	for i := 0; i < 1000; i++ {
		res, err := c.Instrument_Info("SPOT", nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("res: ", res)
		//res2, err := c.Public_ServerTime_Info(nil)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//fmt.Println("res2: ", res2)
		time.Sleep(time.Millisecond * 500)
	}
}

func TestGetContractType(t *testing.T) {
	fmt.Println(GetContractType("BCH-USD-0203"))
}
