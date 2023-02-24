package bybit

import (
	"clients/exchange/cex/base"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
	"testing"
)

var (
	// ProxyUrl         = "http://127.0.0.1:1087"
	proxyUrl         = "http://127.0.0.1:7890"
	TimeOffset int64 = 30
	conf       base.APIConf
	bybit      *ClientBybit
	tradePrice float64 = 1700
)

func init() {
	precisionMap := map[string]*client.PrecisionItem{
		"ETH/USDT": &client.PrecisionItem{
			Symbol:    "ETH/USDT",
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    8,
			Price:     8,
			AmountMin: 0.001,
		},
	}

	conf := base.APIConf{
		ProxyUrl: proxyUrl,
		//AccessKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.AccessKey,
		AccessKey: "Ig0T9qTB4h8heN5aCo",
		SecretKey: "djSxj25DX8cHBAQiizRUdORq9BxIwML1NWZH",
		//Passphrase: config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.Passphrase,
		IsTest: false,
	}
	bybit = NewClientBybit(conf, precisionMap)
}

func TestGetSymbols(t *testing.T) {
	client := NewClientBybit(conf)
	symbols := client.GetSymbols()
	fmt.Println(symbols)
}

func TestGetFee(t *testing.T) {
	client := NewClientBybit(conf)
	symbols := client.GetSymbols()
	fmt.Println(symbols)
}

func TestDepth(t *testing.T) {
	symbol := client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}
	res, err := bybit.GetDepth(&symbol, 100)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestTime(t *testing.T) {
	res := bybit.IsExchangeEnable()
	fmt.Println(res)
}

func TestBalance(t *testing.T) {
	res, err := bybit.GetBalance()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestPlaceOrder(t *testing.T) {
	req := &order.OrderTradeCEX{
		//Side:      1,
		//OrderType: 1,
		//Base: &order.OrderBase{
		//	Market: common.Market_SPOT,
		//	Symbol: []byte("BTC/USDT"),
		//},
		//Amount: 1,
	}
	req = &order.OrderTradeCEX{
		Hdr: &common.MsgHeader{},
		Base: &order.OrderBase{
			Exchange: common.Exchange_BINANCE,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("ETH/USDT"),
		},
		Price:     tradePrice * 0.01,
		TradeType: order.TradeType_TAKER,
		OrderType: order.OrderType_LIMIT,
		Side:      order.TradeSide_BUY,
		Tif:       order.TimeInForce_GTC,
		Amount:    0.1,
	}
	res, err := bybit.PlaceOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCancle(t *testing.T) {
	o := &order.OrderCancelCEX{
		Base: &order.OrderBase{
			Symbol: []byte("ETH/USDT"),
			IdEx:   "1276079144750367488",
			Type:   common.SymbolType_FUTURE_THIS_WEEK,
		},
	}
	res, err := bybit.CancelOrder(o)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func TestCilentBybit_GetOrder(t *testing.T) {
	req := &order.OrderQueryReq{
		//Producer: []byte{123},
		Symbol: []byte("ETH/USDT"),
		Type:   common.SymbolType_SPOT_NORMAL,
		IdEx:   "1276079144750367488",
	}
	res, err := bybit.GetOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCilentBybit_GetOrderHistory(t *testing.T) {
	req := &client.OrderHistoryReq{
		Asset:  "ETH/USDT",
		Market: common.Market_SPOT,
		Type:   common.SymbolType_SPOT_NORMAL,
	}
	res, err := bybit.GetOrderHistory(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestCilentOkex_GetProcessingOrders(t *testing.T) {
	req := &client.OrderHistoryReq{}
	res, err := bybit.GetProcessingOrders(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestTransfer(t *testing.T) {
	o := &order.OrderTransfer{
		ExchangeToken: []byte("ETH"),
		Amount:        1,
		Chain:         2,
		//ExchangeTo:    []byte("Omni"),
	}
	res, _ := bybit.Transfer(o)
	fmt.Println(res)
}

func TestTransferHistory(t *testing.T) {
	req := &client.TransferHistoryReq{}
	res, _ := bybit.GetTransferHistory(req)
	fmt.Println("res:", res)
}

func TestDepositHistory(t *testing.T) {
	req := &client.DepositHistoryReq{}
	res, _ := bybit.GetDepositHistory(req)
	fmt.Println("res:", res)
}

func TestClientBybit_Loan(t *testing.T) {
	o := &order.OrderLoan{
		Asset:  "BTC",
		Amount: 0.001,
	}
	res, err := bybit.Loan(o)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestClientBybit_LoanHistory(t *testing.T) {
	o := &client.LoanHistoryReq{}
	res, err := bybit.GetLoanOrders(o)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func TestGetFutureSymbol(t *testing.T) {
	res := bybit.GetFutureSymbols(common.Market_FUTURE)
	fmt.Println(res)
}

func TestGetFutureDepth(t *testing.T) {
	symbol := client.SymbolInfo{
		Symbol: "ETH/USD",
		Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
		Market: common.Market_SWAP,
	}
	res, err := bybit.GetFutureDepth(&symbol, 100)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFutureMarketPrice(t *testing.T) {
	symbol := client.SymbolInfo{
		Symbol: "ETH/USDT",
		Type:   common.SymbolType_SWAP_FOREVER,
		Market: common.Market_SWAP,
	}
	res, err := bybit.GetFutureMarkPrice(common.Market_FUTURE, &symbol)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFutureTradeFee(t *testing.T) {
	//symbol := client.SymbolInfo{
	//	Symbol: "ETH/USDT",
	//	Type:   common.SymbolType_SPOT_NORMAL,
	//	Market: common.Market_SWAP,
	//}
	res, err := bybit.GetFutureTradeFee(common.Market_FUTURE)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFuturePrecision(t *testing.T) {
	symbol := client.SymbolInfo{
		Symbol: "ETH/USDT",
		Type:   common.SymbolType_FUTURE_COIN_NEXT_QUARTER,
		Market: common.Market_SWAP,
	}
	res, err := bybit.GetFuturePrecision(common.Market_FUTURE, &symbol)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFuturePositions(t *testing.T) {
	//symbol := client.SymbolInfo{
	//	Symbol: "ETH/USDT",
	//	Type:   common.SymbolType_SPOT_NORMAL,
	//	Market: common.Market_SWAP,
	//}
	res, err := bybit.GetFutureBalance(common.Market_FUTURE)
	fmt.Println(res)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPlaceFutureOrder(t *testing.T) {
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
			Type:     common.SymbolType_FUTURE_COIN_NEXT_QUARTER,
			Symbol:   []byte("ETH/USD"),
		},
		TradeType: order.TradeType_TAKER,
		OrderType: order.OrderType_LIMIT,
		Side:      order.TradeSide_SELL,
		Price:     tradePrice,
		Tif:       order.TimeInForce_GTC,
		Amount:    1,
	}
	orderRes, err := bybit.PlaceFutureOrder(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderRes)
}
