package okex

import (
	"clients/config"
	"clients/exchange/cex/base"
	"fmt"
	"testing"

	"github.com/warmplanet/proto/go/common"

	"github.com/warmplanet/proto/go/client"

	"github.com/warmplanet/proto/go/order"
)

var (
	okex *ClientOkex
)

func init() {
	config.LoadExchangeConfig("conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl:   "http://127.0.0.1:7890",
		IsTest:     false,
		AccessKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.AccessKey,
		SecretKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.SecretKey,
		Passphrase: config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.Passphrase,
	}
	precisionMap := map[string]*client.PrecisionItem{
		"ETH/USDT": &client.PrecisionItem{
			Symbol:    "ETH/USDT",
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    8,
			Price:     8,
			AmountMin: 0.001,
		},
	}
	okex = NewClientOkex(conf, precisionMap)
}

func TestGetSymbols(t *testing.T) {
	res := okex.GetSymbols()
	fmt.Println(res)
}

// [symbol:"ADA/USD/220805" symbol:"ADA/USD/220812" symbol:"ADA/USD/220930" symbol:"ADA/USD/221230" symbol:"BTC/USD/220805" symbol:"BTC/USD/220812" symbol:"BTC/USD/220930" symbol:"BTC/USD/221230" symbol:"DOT/USD/220805" symbol:"DOT/USD/220812" symbol:"DOT/USD/220930" symbol:"DOT/USD/221230" symbol:"ETH/USD/220805" symbol:"ETH/USD/220812" symbol:"ETH/USD/220930" symbol:"ETH/USD/221230" symbol:"LTC/USD/220805" symbol:"LTC/USD/220812" symbol:"LTC/USD/220930" symbol:"LTC/USD/221230" symbol:"AVAX/USD/220805" symbol:"AVAX/USD/220812" symbol:"AVAX/USD/220930" symbol:"AVAX/USD/221230" symbol:"BCH/USD/220805" symbol:"BCH/USD/220812" symbol:"BCH/USD/220930" symbol:"BCH/USD/221230" symbol:"BSV/USD/220930" symbol:"EOS/USD/220805" symbol:"EOS/USD/220812" symbol:"EOS/USD/220930" symbol:"EOS/USD/221230" symbol:"ETC/USD/220805" symbol:"ETC/USD/220812" symbol:"ETC/USD/220930" symbol:"ETC/USD/221230" symbol:"FIL/USD/220805" symbol:"FIL/USD/220812" symbol:"FIL/USD/220930" symbol:"FIL/USD/221230" symbol:"LINK/USD/220805" symbol:"LINK/USD/220812" symbol:"LINK/USD/220930" symbol:"LINK/USD/221230" symbol:"SOL/USD/220805" symbol:"SOL/USD/220812" symbol:"SOL/USD/220930" symbol:"SOL/USD/221230" symbol:"TRX/USD/220805" symbol:"TRX/USD/220812" symbol:"TRX/USD/220930" symbol:"TRX/USD/221230" symbol:"XRP/USD/220805" symbol:"XRP/USD/220812" symbol:"XRP/USD/220930" symbol:"XRP/USD/221230" symbol:"ADA/USDT/220805" symbol:"ADA/USDT/220812" symbol:"ADA/USDT/220930" symbol:"ADA/USDT/221230" symbol:"BTC/USDT/220805" symbol:"BTC/USDT/220812" symbol:"BTC/USDT/220930" symbol:"BTC/USDT/221230" symbol:"DOT/USDT/220805" symbol:"DOT/USDT/220812" symbol:"DOT/USDT/220930" symbol:"DOT/USDT/221230" symbol:"ETH/USDT/220805" symbol:"ETH/USDT/220812" symbol:"ETH/USDT/220930" symbol:"ETH/USDT/221230" symbol:"LTC/USDT/220805" symbol:"LTC/USDT/220812" symbol:"LTC/USDT/220930" symbol:"LTC/USDT/221230" symbol:"BCH/USDT/220805" symbol:"BCH/USDT/220812" symbol:"BCH/USDT/220930" symbol:"BCH/USDT/221230" symbol:"EOS/USDT/220805" symbol:"EOS/USDT/220812" symbol:"EOS/USDT/220930" symbol:"EOS/USDT/221230" symbol:"ETC/USDT/220805" symbol:"ETC/USDT/220812" symbol:"ETC/USDT/220930" symbol:"ETC/USDT/221230" symbol:"FIL/USDT/220805" symbol:"FIL/USDT/220812" symbol:"FIL/USDT/220930" symbol:"FIL/USDT/221230" symbol:"KNC/USDT/220805" symbol:"KNC/USDT/220812" symbol:"KNC/USDT/220930" symbol:"KNC/USDT/221230" symbol:"LINK/USDT/220805" symbol:"LINK/USDT/220812" symbol:"LINK/USDT/220930" symbol:"LINK/USDT/221230" symbol:"TRX/USDT/220805" symbol:"TRX/USDT/220812" symbol:"TRX/USDT/220930" symbol:"TRX/USDT/221230" symbol:"XRP/USDT/220805" symbol:"XRP/USDT/220812" symbol:"XRP/USDT/220930" symbol:"XRP/USDT/221230" symbol:"YFII/USDT/220805" symbol:"YFII/USDT/220812" symbol:"YFII/USDT/220930" symbol:"YFII/USDT/221230"]
func TestUGetSymbols(t *testing.T) {
	res := okex.GetUBaseSymbols()
	fmt.Println(res)
}

func TestUGetFutualSymbols(t *testing.T) {
	res := okex.GetFutureSymbols(common.Market(1))
	uBase := make([]string, 0)
	cBase := make([]string, 0)
	x := make([]string, 0)
	for _, r := range res {
		var tmp string
		if r.Type == common.SymbolType_SWAP_COIN_FOREVER {
			tmp = "\"" + r.Name + "\""
			cBase = append(cBase, tmp)
		} else if r.Type == common.SymbolType_SWAP_FOREVER {
			tmp = "\"" + r.Name + "\""
			uBase = append(uBase, tmp)
		} else {
			x = append(x, r.Name)
		}
	}
	fmt.Println(uBase)
	fmt.Println(cBase)
	fmt.Println(x)
}

func TestDepth(t *testing.T) {
	symbol := client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}
	res, err := okex.GetDepth(&symbol, 100)
	fmt.Println(res.Symbol)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUDepth(t *testing.T) {
	symbol := client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_FUTURE_THIS_WEEK,
	}
	res, err := okex.GetUBaseDepth(&symbol, 100)
	fmt.Println(res.Symbol)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTime(t *testing.T) {
	res := okex.IsExchangeEnable()
	fmt.Println(res)
}

func TestFee(t *testing.T) {
	res, _ := okex.GetTradeFee("ETC/BTC", "LTC/BTC")
	fmt.Println(res)
}

func TestCancle(t *testing.T) {
	o := &order.OrderCancelCEX{
		Base: &order.OrderBase{
			Symbol: []byte("BTC-USDT"),
			IdEx:   "4702665370081761",
			Type:   common.SymbolType_FUTURE_THIS_WEEK,
		},
	}
	res, err := okex.CancelOrder(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Status)
}

func TestSpotBalance(t *testing.T) {
	//{"adjEq":"","details":[{"availBal":"3","availEq":"","cashBal":"3",
	//"ccy":"BTC","crossLiab":"","disEq":"66417.9","eq":"3","eqUsd":"66417.9",
	//"frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"",
	//"maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096187","upl":"","uplLiab":"","stgyEq":"0"},
	//{"availBal":"15","availEq":"","cashBal":"15","ccy":"ETH","crossLiab":"","disEq":"23053.95",
	//"eq":"15","eqUsd":"23053.95","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"",
	//"liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"300","availEq":"",
	//"cashBal":"300","ccy":"JFI","crossLiab":"","disEq":"0","eq":"300","eqUsd":"15584.84415",
	//"frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"",
	//"mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096187","upl":"",
	//"uplLiab":"","stgyEq":"0"},{"availBal":"1500","availEq":"","cashBal":"1500","ccy":"UNI","crossLiab":"",
	//"disEq":"9097.125","eq":"1500","eqUsd":"10702.5","frozenBal":"0","interest":"","isoEq":"","isoLiab":"",
	//"isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096187","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000",
	//"ccy":"USDC","crossLiab":"","disEq":"9000","eq":"9000","eqUsd":"9000","frozenBal":"0","interest":"","isoEq":"",
	//"isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"8971","availEq":"","cashBal":"8997",
	//"ccy":"USDT","crossLiab":"","disEq":"8996.46018","eq":"8997","eqUsd":"8996.46018","frozenBal":"26","interest":"",
	//"isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"26","twap":"0","uTime":"1658286476214","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000","ccy":"PAX","crossLiab":"","disEq":"8910","eq":"9000","eqUsd":"8910","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096140","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000","ccy":"TUSD","crossLiab":"","disEq":"8910","eq":"9000","eqUsd":"8910","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000","ccy":"USDK","crossLiab":"","disEq":"8910","eq":"9000","eqUsd":"8910","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096092","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"300","availEq":"","cashBal":"300","ccy":"OKB","crossLiab":"","disEq":"3962.25","eq":"300","eqUsd":"4402.5","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096140","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"30000","availEq":"","cashBal":"30000","ccy":"TRX","crossLiab":"","disEq":"1657.5","eq":"30000","eqUsd":"1950","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"30","availEq":"","cashBal":"30","ccy":"LTC","crossLiab":"","disEq":"1604.835","eq":"30","eqUsd":"1689.3000000000002","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"3000","availEq":"","cashBal":"3000","ccy":"ADA","crossLiab":"","disEq":"1335.5010000000002","eq":"3000","eqUsd":"1483.89","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096282","upl":"","uplLiab":"","stgyEq":"0"}],"imr":"","isoEq":"","mgnRatio":"","mmr":"","notionalUsd":"","ordFroz":"","totalEq":"170011.34433","uTime":"1658713857034"}
	res, err := okex.GetBalance()
	fmt.Println(err)
	fmt.Println("res:", res)
}

func TestMarginBalance(t *testing.T) {
	//{"adjEq":"","details":[{"availBal":"3","availEq":"","cashBal":"3",
	//"ccy":"BTC","crossLiab":"","disEq":"66417.9","eq":"3","eqUsd":"66417.9",
	//"frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"",
	//"maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096187","upl":"","uplLiab":"","stgyEq":"0"},
	//{"availBal":"15","availEq":"","cashBal":"15","ccy":"ETH","crossLiab":"","disEq":"23053.95",
	//"eq":"15","eqUsd":"23053.95","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"",
	//"liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"300","availEq":"",
	//"cashBal":"300","ccy":"JFI","crossLiab":"","disEq":"0","eq":"300","eqUsd":"15584.84415",
	//"frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"",
	//"mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096187","upl":"",
	//"uplLiab":"","stgyEq":"0"},{"availBal":"1500","availEq":"","cashBal":"1500","ccy":"UNI","crossLiab":"",
	//"disEq":"9097.125","eq":"1500","eqUsd":"10702.5","frozenBal":"0","interest":"","isoEq":"","isoLiab":"",
	//"isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096187","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000",
	//"ccy":"USDC","crossLiab":"","disEq":"9000","eq":"9000","eqUsd":"9000","frozenBal":"0","interest":"","isoEq":"",
	//"isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0",
	//"uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"8971","availEq":"","cashBal":"8997",
	//"ccy":"USDT","crossLiab":"","disEq":"8996.46018","eq":"8997","eqUsd":"8996.46018","frozenBal":"26","interest":"",
	//"isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"26","twap":"0","uTime":"1658286476214","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000","ccy":"PAX","crossLiab":"","disEq":"8910","eq":"9000","eqUsd":"8910","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096140","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000","ccy":"TUSD","crossLiab":"","disEq":"8910","eq":"9000","eqUsd":"8910","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"9000","availEq":"","cashBal":"9000","ccy":"USDK","crossLiab":"","disEq":"8910","eq":"9000","eqUsd":"8910","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096092","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"300","availEq":"","cashBal":"300","ccy":"OKB","crossLiab":"","disEq":"3962.25","eq":"300","eqUsd":"4402.5","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096140","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"30000","availEq":"","cashBal":"30000","ccy":"TRX","crossLiab":"","disEq":"1657.5","eq":"30000","eqUsd":"1950","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"30","availEq":"","cashBal":"30","ccy":"LTC","crossLiab":"","disEq":"1604.835","eq":"30","eqUsd":"1689.3000000000002","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096235","upl":"","uplLiab":"","stgyEq":"0"},{"availBal":"3000","availEq":"","cashBal":"3000","ccy":"ADA","crossLiab":"","disEq":"1335.5010000000002","eq":"3000","eqUsd":"1483.89","frozenBal":"0","interest":"","isoEq":"","isoLiab":"","isoUpl":"","liab":"","maxLoan":"","mgnRatio":"","notionalLever":"","ordFrozen":"0","twap":"0","uTime":"1657794096282","upl":"","uplLiab":"","stgyEq":"0"}],"imr":"","isoEq":"","mgnRatio":"","mmr":"","notionalUsd":"","ordFroz":"","totalEq":"170011.34433","uTime":"1658713857034"}
	res, err := okex.GetMarginBalance()
	fmt.Println(err)
	fmt.Println("res:", res)
}

func TestFuturelance(t *testing.T) {
	res, _ := okex.GetFutureBalance(common.Market(1))
	fmt.Println("res:", res)
}

func TestFutureMarketPrice(t *testing.T) {
	sy := client.SymbolInfo{Symbol: "BZZ/USDT/SWAP"}
	res, _ := okex.GetFutureMarkPrice(common.Market_FUTURE, &sy)
	fmt.Println("res:", res)
}

func TestTransfer(t *testing.T) {
	o := &order.OrderTransfer{
		Hdr:             &common.MsgHeader{},
		Base:            &order.OrderBase{},
		ExchangeToken:   []byte("USDT"),
		Amount:          1,
		Chain:           common.Chain_ARBITRUM,
		TransferAddress: []byte("0x1cb3643Db2E039a4abdee676c544dfC13c9F1ea2"),
		//ExchangeTo:    []byte("Omni"),
	}
	res, _ := okex.Transfer(o)
	fmt.Println(res)
}

func TestTransferHistory(t *testing.T) {
	req := &client.TransferHistoryReq{}
	res, _ := okex.GetTransferHistory(req)
	fmt.Println("res:", res)
}

func TestAssetTransfer(t *testing.T) {
	req := &order.OrderMove{
		Asset:  "USDT",
		Amount: 10,
		Source: common.Market_SPOT,
		Target: common.Market_WALLET,
	}

	res, err := okex.MoveAsset(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res:", res)
}

func TestClientOkex_GetDepositHistory(t *testing.T) {
	//{"billId":"120600652","ccy":"USDT","clientId":"","balChg":"-1.5","bal":"3","type":"131","ts":"1658721901000"}
	//{"billId":"120600651","ccy":"USDT","clientId":"","balChg":"1.5","bal":"4.5","type":"130","ts":"1658721725000"}
	//{"billId":"120586805","ccy":"USDT","clientId":"","balChg":"1.5","bal":"3","type":"130","ts":"1658286476000"}
	//{"billId":"120586804","ccy":"USDT","clientId":"","balChg":"1.5","bal":"1.5","type":"130","ts":"1658286295000"}
	req := &client.DepositHistoryReq{}
	res, _ := okex.GetDepositHistory(req)
	for i, item := range res.DepositList {
		fmt.Println(i, item)
	}
}

func TestClientOkex_Loan(t *testing.T) {
	o := &order.OrderLoan{
		Asset:  "BTC",
		Amount: 0.001,
	}
	res, err := okex.Loan(o)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestClientOkex_LoanHistory(t *testing.T) {
	o := &client.LoanHistoryReq{}
	res, err := okex.GetLoanOrders(o)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCilentOkex_GetOrder(t *testing.T) {
	req := &order.OrderQueryReq{
		//Producer: []byte{123},
		Symbol: []byte("BTC/USDT"),
		Type:   common.SymbolType_FUTURE_THIS_WEEK,
		IdEx:   "471698501139832832",
	}
	res, err := okex.GetOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCilentOkex_GetOrderHistory(t *testing.T) {
	req := &client.OrderHistoryReq{
		Asset:  "ETH/USDT",
		Market: common.Market_SPOT,
		Type:   common.SymbolType_SPOT_NORMAL,
	}
	res, err := okex.GetOrderHistory(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCilentOkex_GetProcessingOrders(t *testing.T) {
	req := &client.OrderHistoryReq{}
	res, err := okex.GetProcessingOrders(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestClientOkex_GetOrderHistory(t *testing.T) {
	req := &client.OrderHistoryReq{}
	res, err := okex.GetOrderHistory(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestClientOkex_PlaceOrder(t *testing.T) {
	req := &order.OrderTradeCEX{
		//Side:      1,
		//OrderType: 1,
		//Base: &order.OrderBase{
		//	Market: common.Market_SPOT,
		//	Symbol: []byte("BTC/USDT"),
		//},
		//Amount: 1,
	}
	//req = &order.OrderTradeCEX{
	//	Hdr: &common.MsgHeader{},
	//	Base: &order.OrderBase{
	//		Exchange: common.Exchange_BINANCE,
	//		Market:   common.Market_SPOT,
	//		Type:     common.SymbolType_SPOT_NORMAL,
	//		Symbol:   []byte("ETH/USDT"),
	//	},
	//	Price:     1000,
	//	TradeType: order.TradeType_TAKER,
	//	OrderType: order.OrderType_LIMIT,
	//	Side:      order.TradeSide_BUY,
	//	Tif:       order.TimeInForce_GTC,
	//	Amount:    0.08,
	//}

	req = &order.OrderTradeCEX{
		Hdr: &common.MsgHeader{},
		Base: &order.OrderBase{
			Exchange: common.Exchange_BINANCE,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("ETH/USDT"),
		},
		TradeType: order.TradeType_TAKER,
		OrderType: order.OrderType_MARKET,
		Side:      order.TradeSide_BUY,
		Tif:       order.TimeInForce_GTC,
		Amount:    0.08,
	}
	res, err := okex.PlaceOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestClientUOkex_PlaceOrder(t *testing.T) {
	req := &order.OrderTradeCEX{
		Side:      1,
		OrderType: order.OrderType_MARKET,
		Base: &order.OrderBase{
			Market: common.Market_FUTURE,
			Symbol: []byte("BTC/USDT"),
			Type:   common.SymbolType_FUTURE_THIS_WEEK,
		},
		Amount: 1,
	}
	req = &order.OrderTradeCEX{
		Hdr: &common.MsgHeader{},
		Base: &order.OrderBase{
			Exchange: common.Exchange_BINANCE,
			Market:   common.Market_FUTURE,
			Type:     common.SymbolType_SWAP_FOREVER,
			Symbol:   []byte("ETH/USDT"),
		},
		TradeType: order.TradeType_TAKER,
		OrderType: order.OrderType_MARKET,
		Side:      order.TradeSide_SELL,
		Tif:       order.TimeInForce_GTC,
		Amount:    0.08,
	}
	res, err := okex.PlaceUBaseOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestClientUPrecision(t *testing.T) {
	//symbol := client.SymbolInfo{
	//	Symbol: "BTC/USDT",
	//	Type:   common.SymbolType_FUTURE_THIS_WEEK,
	//}
	res, err := okex.GetUBasePrecision()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestFuturePrecision(t *testing.T) {
	//symbol := client.SymbolInfo{
	//	Symbol: "BTC/USDT",
	//	Type:   common.SymbolType_FUTURE_THIS_WEEK,
	//}
	res, err := okex.GetFuturePrecision(common.Market_SWAP)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestClientUFee(t *testing.T) {
	//symbol := client.SymbolInfo{
	//	Symbol: "BTC/USDT",
	//	Type:   common.SymbolType_FUTURE_THIS_WEEK,
	//}
	sym := &client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_FUTURE_THIS_WEEK,
	}
	res, err := okex.GetUBaseTradeFee(sym)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
func TestFutureFee(t *testing.T) {
	//symbol := client.SymbolInfo{
	//	Symbol: "BTC/USDT",
	//	Type:   common.SymbolType_FUTURE_THIS_WEEK,
	//}
	sym := &client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_FUTURE_THIS_WEEK,
	}
	fmt.Println(sym)
	res, err := okex.GetFutureTradeFee(common.Market_FUTURE_COIN, sym)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestSubAssetTransfer(t *testing.T) {
	req := &order.OrderMove{
		Asset:         "USDT",
		Amount:        1,
		Source:        common.Market_SPOT,
		Target:        common.Market_SPOT,
		AccountSource: "test003",
		AccountTarget: "test004",
		ActionUser:    order.OrderMoveUserType_Sub,
	}
	res, err := okex.MoveAsset(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res:", res)
}

func TestMoveHistory(t *testing.T) {
	req1 := &client.MoveHistoryReq{
		Source:     common.Market_SPOT,
		Target:     common.Market_FUTURE_COIN,
		ActionUser: order.OrderMoveUserType_Master,
	}
	res, err := okex.GetMoveHistory(req1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestSubMoveHistory(t *testing.T) {
	req := &client.MoveHistoryReq{
		Source:        common.Market_SPOT,
		Target:        common.Market_SPOT,
		AccountSource: "test003",
		AccountTarget: "",
		ActionUser:    order.OrderMoveUserType_Master,
	}
	res, err := okex.GetMoveHistory(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("res:", res)
}
