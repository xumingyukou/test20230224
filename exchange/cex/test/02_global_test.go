package test

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
)

// 获取交易所名称
func TestGetExchange(t *testing.T) {
	// 测试交易所返回非空
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey("Test GetExchange", t, func() {
			So(api.GetExchange(), ShouldNotBeNil)
		})
	}
}

// 获取所有现货币对
func TestGetSymbols(t *testing.T) {
	// 测试币本位币对返回数量大于10
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey("Test GetSymbols", t, func() {
			symbols := api.GetSymbols()
			So(len(symbols), ShouldBeGreaterThan, 10)
			t.Logf("%#v\n", symbols)
		})
	}
}

// 行情获取
func TestGetDepth(t *testing.T) {
	// 测试随机选取2个币对，获取100档行情数据，单边行情大于10个
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey("Test GetSymbols", t, func() {
			symbolList := api.GetSymbols()
			symbols := []*client.SymbolInfo{
				{
					Symbol: symbolList[rand.Intn(len(symbolList))],
					Type:   common.SymbolType_SPOT_NORMAL,
				},
				{
					Symbol: symbolList[rand.Intn(len(symbolList))],
					Type:   common.SymbolType_SPOT_NORMAL,
				},
			}
			for _, symbol := range symbols {
				depthData, err := api.GetDepth(symbol, 10)
				So(err, ShouldBeNil)
				So(len(depthData.Asks), ShouldBeGreaterThan, 0)
				So(len(depthData.Bids), ShouldBeGreaterThan, 0)
				t.Logf("%#v\n", depthData)
			}
		})
	}
}

// 可用判断
func TestIsExchangeEnable(t *testing.T) {
	// 测试交易所是否可用
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey("Test IsExchangeEnable", t, func() {
			So(api.IsExchangeEnable(), ShouldBeTrue)
		})
	}
}

// 交易手续费
func TestGetTradeFee(t *testing.T) {
	// 随机选取10个币对，测试交易手续费，taker、maker费率大于-1
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test GetTradeFee", t, func() {
			symbolList := api.GetSymbols()
			var symbols []string
			for i := 0; i < 10; i++ {
				symbols = append(symbols, symbolList[rand.Intn(len(symbolList)-1)])
			}
			fees, err := api.GetTradeFee(symbols...)
			So(err, ShouldBeNil)
			for _, fee := range fees.TradeFeeList {
				So(fee.Symbol, ShouldNotEqual, "")
				So(fee.Taker, ShouldBeGreaterThan, -1)
				So(fee.Maker, ShouldBeGreaterThan, -1)
			}
			t.Logf("%#v\n", fees)
		})
	}
}

// 提币手续费
func TestGetTransferFee(t *testing.T) {
	// 随机选取10种币，测试转账手续费，使用ETH网络，手续费大于-1
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test GetTransferFee", t, func() {
			var (
				symbols = []string{"USDT", "ETH", "SOL", "BTC"}
			)
			fees, err := api.GetTransferFee(common.Chain_SOLANA, symbols...)
			So(err, ShouldBeNil)
			//So(len(fees.TransferFeeList), ShouldEqual, 10)
			for _, fee := range fees.TransferFeeList {
				if fee != nil {
					t.Log(fee)
					So(fee.Token, ShouldNotEqual, "")
					So(fee.Network, ShouldEqual, common.Chain_SOLANA)
					So(fee.Fee, ShouldBeGreaterThan, -1)
				}
			}
			t.Logf("%#v\n", fees)
		})
	}
}

// 精度
func TestGetPrecision(t *testing.T) {
	// 随机选取10个币对，测试交易精度，AmountMin大于0，Amount大于等于-10，Price大于等于-10，
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test GetPrecision", t, func() {
			symbolList := api.GetSymbols()
			var (
				symbols   []string
				symbolMap = make(map[string]bool)
			)
			for i := 0; i < 10; i++ {
				symbol := symbolList[rand.Intn(len(symbolList)-1)]
				if _, ok := symbolMap[symbol]; !ok {
					symbols = append(symbols, symbol)
					symbolMap[symbol] = true
				} else {
					i--
				}
			}
			precisions, err := api.GetPrecision(symbols...)
			So(err, ShouldBeNil)
			So(len(precisions.PrecisionList), ShouldEqual, 10)
			for _, fee := range precisions.PrecisionList {
				So(fee.Symbol, ShouldNotEqual, "")
				So(fee.AmountMin, ShouldBeGreaterThan, 0)
				So(fee.Amount, ShouldBeGreaterThanOrEqualTo, -10)
				So(fee.Price, ShouldBeGreaterThanOrEqualTo, -10)
			}
			t.Logf("%#v\n", precisions)
		})
	}
}
