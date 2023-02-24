package test

import (
	"math/rand"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warmplanet/proto/go/client"

	"testing"

	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

// 下单
func TestPlaceOrderLimit(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestPlaceOrderLimit", t, func() {
			balances, err := api.GetBalance()
			So(err, ShouldBeNil)
			haveAsset := false
			if balances != nil {
				for _, balance := range balances.BalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDT balance:", balance)
						haveAsset = true
						break
					}
				}
			} else {
				marginBalances, err := api.GetMarginBalance()
				So(err, ShouldBeNil)
				for _, balance := range marginBalances.MarginBalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDC balance:", balance)
						haveAsset = true
						break
					}
				}
			}

			if !haveAsset {
				t.Error("have no usdt to order")
			}
			// 限价单
			req := &order.OrderTradeCEX{
				Hdr: &common.MsgHeader{
					Producer: []byte("0493733"),
				},
				Base: &order.OrderBase{
					Id:       rand.Int63n(1000000),
					Exchange: common.Exchange_HUOBI,
					Market:   common.Market_SPOT,
					Type:     common.SymbolType_SPOT_NORMAL,
					Symbol:   []byte(tradeSymbol),
				},
				TradeType: order.TradeType_TAKER,
				OrderType: order.OrderType_LIMIT,
				Side:      order.TradeSide_SELL,
				Price:     tradePrice,
				Tif:       order.TimeInForce_GTC,
				Amount:    tradeAmount,
			}
			orderRes, err := api.PlaceOrder(req)
			So(err, ShouldBeNil)
			So(orderRes, ShouldNotBeNil)
			t.Log(orderRes)

			req2 := &order.OrderQueryReq{
				Symbol: []byte(tradeSymbol),
				IdEx:   orderRes.OrderId,
				Market: common.Market_SPOT,
				Type:   common.SymbolType_SPOT_NORMAL,
			}
			orderStatus, err := api.GetOrder(req2)
			So(err, ShouldBeNil)
			So(orderStatus.OrderId, ShouldEqual, orderRes.OrderId)
			So(orderStatus.Symbol, ShouldEqual, tradeSymbol)
			So(orderStatus.Price, ShouldBeGreaterThan, 0)
		})
	}
}
func TestPlaceOrderMarket(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestPlaceOrderMarket", t, func() {
			balances, err := api.GetBalance()
			So(err, ShouldBeNil)
			haveAsset := false
			if balances != nil {
				for _, balance := range balances.BalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDT balance:", balance)
						haveAsset = true
						break
					}
				}
			} else {
				marginBalances, err := api.GetMarginBalance()
				So(err, ShouldBeNil)
				for _, balance := range marginBalances.MarginBalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDT balance:", balance)
						haveAsset = true
						break
					}
				}
			}

			if !haveAsset {
				t.Error("have no usdt to order")
			}
			// 市价单
			req := &order.OrderTradeCEX{
				Hdr: &common.MsgHeader{},
				Base: &order.OrderBase{
					Exchange: common.Exchange_BINANCE,
					Market:   common.Market_SPOT,
					Type:     common.SymbolType_SPOT_NORMAL,
					Symbol:   []byte(tradeSymbol),
				},
				TradeType: order.TradeType_TAKER,
				OrderType: order.OrderType_MARKET,
				Side:      order.TradeSide_SELL,
				Tif:       order.TimeInForce_GTC,
				Amount:    tradeAmount,
			}
			orderRes, err := api.PlaceOrder(req)
			So(err, ShouldBeNil)
			So(orderRes, ShouldNotBeNil)
			t.Log(orderRes)

			req2 := &order.OrderQueryReq{
				IdEx:   orderRes.OrderId,
				Symbol: []byte(tradeSymbol),
				Market: common.Market_SPOT,
				Type:   common.SymbolType_SPOT_NORMAL,
			}
			orderStatus, err := api.GetOrder(req2)
			So(err, ShouldBeNil)
			So(orderStatus.OrderId, ShouldEqual, orderRes.OrderId)
			So(orderStatus.Symbol, ShouldEqual, tradeSymbol)
			So(orderStatus.Price, ShouldBeGreaterThanOrEqualTo, 0)
		})
	}
}
func TestPlaceOrderLimitMaker(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestPlaceOrderLimitMaker", t, func() {
			balances, err := api.GetBalance()
			So(err, ShouldBeNil)
			haveAsset := false
			if balances != nil {
				for _, balance := range balances.BalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDT balance:", balance)
						haveAsset = true
						break
					}
				}
			} else {
				marginBalances, err := api.GetMarginBalance()
				So(err, ShouldBeNil)
				for _, balance := range marginBalances.MarginBalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDT balance:", balance)
						haveAsset = true
						break
					}
				}
			}

			if !haveAsset {
				t.Error("have no usdt to order")
			}
			// 限价maker单
			req := &order.OrderTradeCEX{
				Hdr: &common.MsgHeader{},
				Base: &order.OrderBase{
					Exchange: common.Exchange_HUOBI,
					Market:   common.Market_SPOT,
					Type:     common.SymbolType_SPOT_NORMAL,
					Symbol:   []byte(tradeSymbol),
				},
				Price:     tradePrice,
				TradeType: order.TradeType_MAKER,
				OrderType: order.OrderType_LIMIT,
				Side:      order.TradeSide_BUY,
				Tif:       order.TimeInForce_GTC,
				Amount:    tradeAmount,
			}
			orderRes, err := api.PlaceOrder(req)
			So(err, ShouldBeNil)
			So(orderRes, ShouldNotBeNil)
			t.Log(orderRes)
			req2 := &order.OrderQueryReq{
				Symbol: []byte(tradeSymbol),
				IdEx:   orderRes.OrderId,
				Market: common.Market_SPOT,
				Type:   common.SymbolType_SPOT_NORMAL,
			}
			orderStatus, err := api.GetOrder(req2)
			So(err, ShouldBeNil)
			So(orderStatus.OrderId, ShouldEqual, orderRes.OrderId)
			So(orderStatus.Symbol, ShouldEqual, tradeSymbol)
		})
	}
}

// U本位合约下单
func TestUBaseFuturePlaceOrderMaker(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestFuturePlaceOrderMaker", t, func() {
			// 市价单
			req := &order.OrderTradeCEX{
				Hdr: &common.MsgHeader{},
				Base: &order.OrderBase{
					Id:       rand.Int63n(1000000),
					Exchange: common.Exchange_HUOBI,
					Market:   common.Market_SWAP,
					Type:     common.SymbolType_SWAP_FOREVER,
					Symbol:   []byte(tradeSymbol),
				},
				TradeType: order.TradeType_TAKER,
				OrderType: order.OrderType_LIMIT,
				Side:      order.TradeSide_SELL,
				Price:     tradePrice,
				Tif:       order.TimeInForce_GTC,
				Amount:    0.01,
			}
			orderRes, err := api.PlaceFutureOrder(req)
			So(err, ShouldBeNil)
			So(orderRes, ShouldNotBeNil)
			t.Log(orderRes)
			req2 := &order.OrderQueryReq{
				IdEx:   orderRes.OrderId,
				Market: common.Market_SWAP,
				Type:   common.SymbolType_SWAP_FOREVER,
				Symbol: []byte(tradeSymbol),
			}
			orderStatus, err := api.GetOrder(req2)
			t.Log(orderStatus)
			So(err, ShouldBeNil)
			So(orderStatus.OrderId, ShouldEqual, orderRes.OrderId)
			So(orderStatus.Symbol, ShouldEqual, tradeSymbol)
		})
	}
}

// 取消订单
func TestCancelOrder(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestCancelOrder", t, func() {
			balances, err := api.GetBalance()
			So(err, ShouldBeNil)
			haveAsset := false
			if balances != nil {
				for _, balance := range balances.BalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDT balance:", balance)
						haveAsset = true
						break
					}
				}
			} else {
				marginBalances, err := api.GetMarginBalance()
				So(err, ShouldBeNil)
				for _, balance := range marginBalances.MarginBalanceList {
					if balance.Asset == "USDT" {
						t.Log("asset USDT balance:", balance)
						haveAsset = true
						break
					}
				}
			}

			if !haveAsset {
				t.Error("have no usdt to order")
			}
			// 市价单
			req := &order.OrderTradeCEX{
				Hdr: &common.MsgHeader{},
				Base: &order.OrderBase{
					Exchange: common.Exchange_HUOBI,
					Market:   common.Market_SPOT,
					Type:     common.SymbolType_SPOT_NORMAL,
					Symbol:   []byte(tradeSymbol),
				},
				Price:     tradePrice * 0.8,
				TradeType: order.TradeType_TAKER,
				OrderType: order.OrderType_LIMIT,
				Side:      order.TradeSide_BUY,
				Tif:       order.TimeInForce_GTC,
				Amount:    tradeAmount,
			}
			t.Log(2222, req)
			orderRes, err := api.PlaceOrder(req)
			So(err, ShouldBeNil)
			So(orderRes, ShouldNotBeNil)
			t.Logf("%#v\n", orderRes)
			req2 := &order.OrderCancelCEX{
				Hdr: &common.MsgHeader{},
				Base: &order.OrderBase{
					Symbol: []byte(tradeSymbol),
					IdEx:   orderRes.OrderId,
					Type:   common.SymbolType_SPOT_NORMAL,
					Market: common.Market_SPOT,
				},
			}
			t.Log(333, req2)
			orderCancelRes, err := api.CancelOrder(req2)
			t.Logf("%#v\n", orderCancelRes)
			So(err, ShouldBeNil)
			So(orderCancelRes, ShouldNotBeNil)
		})
	}
}

// 订单历史
func TestGetOrderHistory(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestGetOrderHistory", t, func() {
			req := &client.OrderHistoryReq{
				Asset:  tradeSymbol,
				Market: common.Market_SWAP,
				Type:   common.SymbolType_SPOT_NORMAL,
			}
			orderRes, err := api.GetOrderHistory(req)
			So(err, ShouldBeNil)
			So(orderRes, ShouldNotBeNil)
			t.Log(orderRes)
		})
	}
}
