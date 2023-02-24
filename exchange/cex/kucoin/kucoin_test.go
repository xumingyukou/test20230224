package kucoin

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/kucoin/spot_api"
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

var (
	ProxyUrl          = "http://127.0.0.1:7890"
	TimeOffset  int64 = 30
	conf        base.APIConf
	tradeSymbol         = "ETH/USDT"
	tradePrice  float64 = 1000
	tradeAmount         = 1.0
)

func init() {
	config.LoadExchangeConfig("/Users/yych/repositry/clients/conf/exchange.toml")
	conf = base.APIConf{
		ReadTimeout: TimeOffset,
		ProxyUrl:    ProxyUrl,
		EndPoint:    spot_api.FUTURUE_API_BASE_URL,
		//EndPoint:   spot_api.SPOT_API_BASE_URL,
		AccessKey:  config.ExchangeConfig.ExchangeList["kucoin"].ApiKeyConfig.AccessKey,
		Passphrase: config.ExchangeConfig.ExchangeList["kucoin"].ApiKeyConfig.Passphrase,
		SecretKey:  config.ExchangeConfig.ExchangeList["kucoin"].ApiKeyConfig.SecretKey,
		IsTest:     true,
	}

	//conf = base.APIConf{
	//	ReadTimeout: TimeOffset,
	//	ProxyUrl:    ProxyUrl,
	//	EndPoint:    spot_api.SANDBOX_BASE_URL,
	//	AccessKey:   config.ExchangeConfig.ExchangeList["kucoin_sandbox"].ApiKeyConfig.AccessKey,
	//	Passphrase:  config.ExchangeConfig.ExchangeList["kucoin_sandbox"].ApiKeyConfig.Passphrase,
	//	SecretKey:   config.ExchangeConfig.ExchangeList["kucoin_sandbox"].ApiKeyConfig.SecretKey,
	//	IsTest:      true,
	//}
}

func TestLoadConfig(t *testing.T) {
	fmt.Println(conf)
}

func TestGetExchange(t *testing.T) {
	client := NewClientKucoin(conf)
	fmt.Println(client.GetExchange())
}

func TestGetSymbols(t *testing.T) {
	client := NewClientKucoin(conf)
	symbols := client.GetSymbols()
	fmt.Println(symbols)
}

func TestGetTradeFee(t *testing.T) {
	cli := NewClientKucoin(conf)
	symbols := cli.GetSymbols()
	if len(symbols) >= 20 {
		symbols = symbols[:20]
	}
	tradeFee, err := cli.GetTradeFee(symbols...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tradeFee)
}

func TestGetDepth(t *testing.T) {
	cli := NewClientKucoin(conf)
	symbols := cli.GetSymbols()
	limit := 50

	for _, symbol := range symbols {

		symbolInfo := &client.SymbolInfo{
			Symbol: symbol,
		}
		dep, err := cli.GetDepth(symbolInfo, limit)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(dep)
		fmt.Println(len(dep.Bids), len(dep.Asks))
	}

	limit = 10
	for _, symbol := range symbols {
		symbolInfo := &client.SymbolInfo{
			Symbol: symbol,
		}
		dep, err := cli.GetDepth(symbolInfo, limit)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(dep)
		fmt.Println(len(dep.Bids), len(dep.Asks))
	}

	limit = 101
	for _, symbol := range symbols {
		symbolInfo := &client.SymbolInfo{
			Symbol: symbol,
		}
		_, err := cli.GetDepth(symbolInfo, limit)
		if err == nil {
			t.Fatal(errors.New("should exceed limit"))
		} else {
			fmt.Println(err.Error())
		}
	}
}

func TestGetTransferFee(t *testing.T) {
	cli := NewClientKucoin(conf)
	res, err := cli.GetTransferFee(common.Chain_ETH, "BTC", "ETH")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetBalance(t *testing.T) {
	cli := NewClientKucoin(conf)
	res, err := cli.GetBalance()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetTimestampNow(t *testing.T) {
	fmt.Println(GetCurrentTimestamp())
}

func TestGetMarginBalance(t *testing.T) {
	cli := NewClientKucoin(conf)
	res, err := cli.GetMarginBalance()
	if err != nil {
		t.Fatal(err)
	}
	for _, balanceItem := range res.MarginBalanceList {
		if balanceItem.NetAsset < 0 {
			t.Fatal(errors.New("user asset less than 0"))
		}
	}
	fmt.Println(res)
}

func TestParsePrecision(t *testing.T) {
	if ParsePrecision("0.0001") != 4 {
		t.Fatal(errors.New("parse error"))
	}

	if ParsePrecision("0.1") != 1 {
		t.Fatal(errors.New("parse error"))
	}

	if ParsePrecision("0.00000001") != 8 {
		t.Fatal(errors.New("parse error"))
	}
}

func TestGetPrecision(t *testing.T) {
	cli := NewClientKucoin(conf)
	symbols := cli.GetSymbols()
	if len(symbols) >= 20 {
		symbols = symbols[:20]
	}
	res, err := cli.GetPrecision(symbols...)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(res.PrecisionList), len(symbols))
	fmt.Println(symbols)
	fmt.Println(res)
}

func TestGetMarginIsolatedBalance(t *testing.T) {
	cli := NewClientKucoin(conf)
	symbols := cli.GetSymbols()
	if len(symbols) >= 20 {
		symbols = symbols[:20]
	}
	res, err := cli.GetMarginIsolatedBalance(symbols...)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(symbols)
	fmt.Println(res)
}

// func TestEndpoint(t *testing.T) {
// 	cli := NewClientKucoin(conf)
// 	res := &spot_api.RespError{}
// 	err := cli.api.DoRequest(cli.api.ReqUrl.ISOLATED_ACCOUNTS_URL+"/BTC-USDT", "GET", true, nil, res)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// producer:"0493733"  id:112321  order_id:"636a05c004b71f0001a9aab4"  timestamp:1667892672829  resp_type:FULL  symbol:"SOL/USDT"  status:FILLED  accum_qty:2.7556  accum_amount:0.1  fee_asset:"USDT"  fee:0.0027556
func TestPlaceOrder(t *testing.T) {
	cli := NewClientKucoin(conf)
	req := &order.OrderTradeCEX{
		Hdr: &common.MsgHeader{
			Producer: []byte("0493733"),
		},
		Base: &order.OrderBase{
			Id:       112321,
			Quote:    []byte("USDT"),
			Token:    []byte("SOL"),
			Exchange: common.Exchange_KUCOIN,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("ETH/USDT"),
		},
		OrderType: order.OrderType_MARKET,
		TradeType: order.TradeType_TAKER,
		Side:      order.TradeSide_BUY,
		Price:     27,
		Tif:       order.TimeInForce_GTC,
		Amount:    0.1,
	}
	req = &order.OrderTradeCEX{
		Hdr: &common.MsgHeader{
			Producer: []byte("0493733"),
		},
		Base: &order.OrderBase{
			Exchange: common.Exchange_KUCOIN,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("ETH/USDT"),
		},
		TradeType: order.TradeType_TAKER,
		OrderType: order.OrderType_MARKET,
		Side:      order.TradeSide_BUY,
		Tif:       order.TimeInForce_GTC,
		Price:     27,
		Amount:    tradeAmount,
	}

	req = &order.OrderTradeCEX{
		Hdr: &common.MsgHeader{
			Producer: []byte("0493733"),
		},
		Base: &order.OrderBase{
			Id:       rand.Int63n(1000000),
			Exchange: common.Exchange_BINANCE,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte(tradeSymbol),
		},
		TradeType: order.TradeType_TAKER,
		OrderType: order.OrderType_LIMIT,
		Side:      order.TradeSide_BUY,
		Price:     tradePrice,
		Tif:       order.TimeInForce_GTC,
		Amount:    tradeAmount,
	}

	res, err := cli.PlaceOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestCilentBybit_GetOrder(t *testing.T) {
	cli := NewClientKucoin(conf)
	req := &order.OrderQueryReq{
		//Producer: []byte{123},
		Symbol: []byte("ETH/USDT"),
		Type:   common.SymbolType_SPOT_NORMAL,
		IdEx:   "636a05c004b71f0001a9aab4",
	}
	res, err := cli.GetOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCancelOrder(t *testing.T) {
	cli := NewClientKucoin(conf)
	req := &order.OrderCancelCEX{
		Hdr: &common.MsgHeader{
			Producer: []byte("0493733"),
		},
		Base: &order.OrderBase{
			Id:       112321,
			Quote:    []byte("USDT"),
			Token:    []byte("ETH"),
			Exchange: common.Exchange_KUCOIN,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("ETH/USDT"),
		},
		CancelId: 1234,
	}

	res, err := cli.CancelOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestOrderHistory(t *testing.T) {
	cli := NewClientKucoin(conf)

	res, err := cli.GetOrderHistory(&client.OrderHistoryReq{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestProcessingOrderHistory(t *testing.T) {
	cli := NewClientKucoin(conf)

	res, err := cli.GetProcessingOrders(&client.OrderHistoryReq{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestTransfer(t *testing.T) {
	cli := NewClientKucoin(conf)
	o := &order.OrderTransfer{
		ExchangeToken: []byte("BTC"),
		Amount:        8,
		Chain:         2,
		//ExchangeTo:    []byte("Omni"),
	}
	res, _ := cli.Transfer(o)
	fmt.Println(res)
}

func TestTransferHistory(t *testing.T) {
	cli := NewClientKucoin(conf)
	//o := &order.OrderTransfer{
	//	ExchangeToken: []byte("BTC"),
	//	Amount:        8,
	//	Chain:         2,
	//	//ExchangeTo:    []byte("Omni"),
	//}
	req := &client.TransferHistoryReq{}
	res, _ := cli.GetTransferHistory(req)
	fmt.Println(res)
}

func TestMoveHistory(t *testing.T) {
	cli := NewClientKucoin(conf)
	req1 := &client.MoveHistoryReq{
		Source:     common.Market_SPOT,
		Target:     common.Market_FUTURE_COIN,
		ActionUser: order.OrderMoveUserType_Master,
	}
	res, err := cli.GetMoveHistory(req1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestDepositHistory(t *testing.T) {
	cli := NewClientKucoin(conf)
	req := &client.DepositHistoryReq{}
	res, _ := cli.GetDepositHistory(req)
	fmt.Println("res:", res)
}

func TestLoan(t *testing.T) {
	cli := NewClientKucoin(conf)
	o := &order.OrderLoan{
		Asset:  "USDT",
		Amount: 1,
	}
	res, err := cli.Loan(o)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestRepayHistory(t *testing.T) {
	cli := NewClientKucoin(conf)
	o := &client.LoanHistoryReq{}
	res, err := cli.GetLoanOrders(o)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("res", res, len(res.LoadList))
}

func TestGetFutureSymbol(t *testing.T) {
	cli := NewClientKucoin(conf)
	res := cli.GetFutureSymbols(common.Market_FUTURE)
	fmt.Println(res)
}

func TestGetFutureDepth(t *testing.T) {
	cli := NewClientKucoin(conf)
	symbol := client.SymbolInfo{
		Symbol: "ETH/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
		Market: common.Market_SWAP,
	}
	res, err := cli.GetFutureDepth(&symbol, 100)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFuturePrecision(t *testing.T) {
	cli := NewClientKucoin(conf)
	//symbol := client.SymbolInfo{
	//	Symbol: "ETH/USDT",
	//	Type:   common.SymbolType_SPOT_NORMAL,
	//	Market: common.Market_SWAP,
	//}
	res, err := cli.GetFuturePrecision(common.Market_FUTURE)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFutureMarketPrice(t *testing.T) {
	cli := NewClientKucoin(conf)
	//symbol := client.SymbolInfo{
	//	Symbol: "ETH/USDT",
	//	Type:   common.SymbolType_SPOT_NORMAL,
	//	Market: common.Market_SWAP,
	//}
	res, err := cli.GetFutureMarkPrice(common.Market_FUTURE)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFutureTradeFee(t *testing.T) {
	cli := NewClientKucoin(conf)
	//symbol := client.SymbolInfo{
	//	Symbol: "ETH/USDT",
	//	Type:   common.SymbolType_SPOT_NORMAL,
	//	Market: common.Market_SWAP,
	//}
	res, err := cli.GetFutureTradeFee(common.Market_FUTURE)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFuturePositions(t *testing.T) {
	cli := NewClientKucoin(conf)
	//symbol := client.SymbolInfo{
	//	Symbol: "ETH/USDT",
	//	Type:   common.SymbolType_SPOT_NORMAL,
	//	Market: common.Market_SWAP,
	//}
	res, err := cli.GetFutureBalance(common.Market_FUTURE)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlaceFutureOrder(t *testing.T) {
	cli := NewClientKucoin(conf)
	//symbol := client.SymbolInfo{
	//	Symbol: "ETH/USDT",
	//	Type:   common.SymbolType_SPOT_NORMAL,
	//	Market: common.Market_SWAP,
	//}
	req := &order.OrderTradeCEX{
		Hdr: &common.MsgHeader{},
		Base: &order.OrderBase{
			Id:       1,
			Exchange: common.Exchange_BINANCE,
			Market:   common.Market_SWAP,
			Type:     common.SymbolType_SWAP_FOREVER,
			Symbol:   []byte(tradeSymbol),
		},
		TradeType: order.TradeType_TAKER,
		OrderType: order.OrderType_LIMIT,
		Side:      order.TradeSide_SELL,
		Price:     tradePrice,
		Tif:       order.TimeInForce_GTC,
		Amount:    tradeAmount,
	}
	orderRes, err := cli.PlaceFutureOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderRes)
}
