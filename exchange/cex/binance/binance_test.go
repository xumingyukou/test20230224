package binance

import (
	"clients/config"
	"clients/exchange/cex/base"
	"fmt"
	"testing"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

var (
	ProxyUrl   = "http://127.0.0.1:1080"
	binance    *ClientBinance
	timeOffset int64 = 30
	conf       base.APIConf
)

func init() {
	config.LoadExchangeConfig("../../../conf/exchange.toml")
	conf = base.APIConf{
		ReadTimeout: timeOffset,
		ProxyUrl:    ProxyUrl,
		IsTest:      true,
		AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	}
	binance = NewClientBinance(conf)
	//testmove/binance/apikey  和  testmove/binance/secretkey
}

func TestNewClientBinance(t *testing.T) {
	res, err := binance.CancelOrder(&order.OrderCancelCEX{
		Base: &order.OrderBase{
			Id:       112321,
			Exchange: common.Exchange_BINANCE,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("BTC/USDT"),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetSymbols(t *testing.T) {
	res := binance.GetSymbols()
	fmt.Println(res)
	res2 := binance.GetFutureSymbols(common.Market_SPOT)
	for _, item := range res2 {
		fmt.Println(item)
	}
}

func TestOrder(t *testing.T) {
	precisionMap := make(map[string]*client.PrecisionItem)
	precisionMap["BUSD/USDT"] = &client.PrecisionItem{
		Symbol:    "",
		Type:      common.SymbolType_SPOT_NORMAL,
		Amount:    3,
		Price:     4,
		AmountMin: 1,
	}
	binance := NewClientBinance(conf, precisionMap)
	res, err := binance.PlaceOrder(&order.OrderTradeCEX{
		Hdr: &common.MsgHeader{},
		Base: &order.OrderBase{
			Id:       112321,
			Quote:    []byte("USDT"),
			Token:    []byte("BUSD"),
			Exchange: common.Exchange_BINANCE,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("BUSD/USDT"),
		},
		TradeType: order.TradeType_MAKER,
		OrderType: order.OrderType_LIMIT,
		Side:      order.TradeSide_SELL,
		Price:     1.02,
		Tif:       order.TimeInForce_GTC,
		Amount:    11,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetPrecision(t *testing.T) {
	res, err := binance.GetPrecision("BTC/USDT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetTradeFee(t *testing.T) {
	res, err := binance.GetTradeFee() //"BTC/USDT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	fmt.Println("")
	res, err = binance.GetFutureTradeFee(common.Market_SWAP) //"BTC/USDT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetTransferFee(t *testing.T) {
	res, err := binance.GetTransferFee(common.Chain_BSC, "BTC") //"BTC/USDT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestDepth(t *testing.T) {
	res, err := binance.GetDepth(&client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}, 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	res, err = binance.GetFutureDepth(&client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SWAP_FOREVER,
	}, 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestBalance(t *testing.T) {
	//res, err := binance.GetBalance()
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Printf("%#v\n", res)
	res3, err := binance.GetMarginBalance()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res3)
	//res2, err := binance.GetFutureBalance(common.Market_FUTURE)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(res2)
}

func TestBalance2(t *testing.T) {
	res2, err := binance.GetFutureBalance(common.Market_FUTURE)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)
}

func TestUBaseOrder(t *testing.T) {
	req := &order.OrderTradeCEX{
		Base: &order.OrderBase{
			Symbol: []byte("BTC/USDT"),
			Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
		},
		Side:      order.TradeSide_BUY,
		Tif:       order.TimeInForce_GTC,
		Amount:    5,
		OrderType: order.OrderType_LIMIT,
		Price:     1800,
	}
	res2, err := binance.PlaceFutureOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)
}

func TestCancelOrder(t *testing.T) {
	req := &order.OrderCancelCEX{
		Base: &order.OrderBase{
			IdEx:   "491140",
			Symbol: []byte("BTC/USDT"),
			Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
		},
		Side: order.TradeSide_BUY,
	}
	res2, err := binance.CancelOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)

}

func TestGetProcessingOrders(t *testing.T) {
	req := &client.OrderHistoryReq{
		Asset: "BTC/USDT",
		Type:  common.SymbolType_FUTURE_NEXT_QUARTER,
	}
	res2, err := binance.GetProcessingOrders(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)
}

func TestGetOrder(t *testing.T) {
	//https://testnet.binancefuture.com/fapi/v1/order GET &map[Content-Type:[application/x-www-form-urlencoded] X-Mbx-Apikey:[8fb5d52e655c58b79a9dc3a916e4e85ee917a4a44b6cd22176951449120ac411]] map[orderId:[491141] recvWindow:[60000] signature:[5cc6f7c415b50b415b8a482a7e872c9d5b23271d05f376e5a09d1fa03d574a33] symbol:[BTCUSDT_220930] timestamp:[1657785724184]]
	//https://testnet.binancefuture.com/fapi/v1/order GET &map[Content-Type:[application/x-www-form-urlencoded] X-Mbx-Apikey:[8fb5d52e655c58b79a9dc3a916e4e85ee917a4a44b6cd22176951449120ac411]] map[orderId:[491141] recvWindow:[60000] signature:[46547de381b8e4121281281cf2da7ae52a6e8f9d0982dd38f53dc951071b30cc] symbol:[BTCUSDT BTCUSDT_220930] timestamp:[1657786988055]]

	req := &order.OrderQueryReq{
		IdEx:   "491141",
		Symbol: []byte("BTC/USDT"),
		Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
	}
	res2, err := binance.GetOrder(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)
}

func TestFutureSymbols(t *testing.T) {

}

func TestMoveAsset(t *testing.T) {
	//普通划转
	req := &order.OrderMove{
		Asset:  "USDT",
		Amount: 400,
		Source: common.Market_SPOT,
		Target: common.Market_FUTURE,
	}
	res2, err := binance.MoveAsset(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)
	req1 := &client.MoveHistoryReq{
		Source: common.Market_SPOT,
		Target: common.Market_FUTURE,
	}
	res, err := binance.GetMoveHistory(req1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestSubAccountList(t *testing.T) {
	res2, err := binance.SubAccountList()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range res2.SubAccounts {

		fmt.Printf("%#v\n", i)
	}
}

func TestSubMoveAsset(t *testing.T) {
	//普通划转
	req := &order.OrderMove{
		Asset:  "BNB",
		Amount: 0.02,
		Source: common.Market_SPOT,
		Target: common.Market_FUTURE_COIN,
	}
	res2, err := binance.MoveAsset(req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res2)
	req1 := &client.MoveHistoryReq{
		Source: common.Market_SPOT,
		Target: common.Market_FUTURE_COIN,
	}
	res, err := binance.GetMoveHistory(req1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
