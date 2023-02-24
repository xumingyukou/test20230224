package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warmplanet/proto/go/common"
)

// balance
func TestGetBalance(t *testing.T) {
	// 获取币现货balance数据
	// 要求：Total = Free + Frozen
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test GetBalance", t, func() {
			balances, err := api.GetBalance()
			So(err, ShouldBeNil)
			t.Logf("%#v\n", balances)
			if balances != nil {
				for _, balance := range balances.BalanceList {
					So(balance.Asset, ShouldNotEqual, "")
					So(balance.Free+balance.Frozen, ShouldEqual, balance.Total)
				}
			} else {
				t.Log(api.GetExchange(), "spot balance is nil")
			}
		})
	}
}
func TestGetMarginBalance(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test GetMarginBalance", t, func() {
			balances, err := api.GetMarginBalance()
			t.Logf("%#v\n", balances)
			So(err, ShouldBeNil)
			So(balances, ShouldNotBeNil)
			So(balances.QuoteAsset, ShouldNotEqual, "")
			So(balances.TotalAsset, ShouldBeGreaterThanOrEqualTo, 0)
			So(balances.TotalNetAsset, ShouldBeGreaterThanOrEqualTo, 0)
			So(balances.TotalLiabilityAsset, ShouldBeGreaterThanOrEqualTo, 0)
			So(balances.TotalNetAsset+balances.TotalLiabilityAsset, ShouldEqual, balances.TotalAsset) //总资产=用户总资产+已借
			for _, balance := range balances.MarginBalanceList {
				So(balance.Asset, ShouldNotEqual, "")
				So(balance.Total-(balance.Free+balance.Frozen), ShouldBeLessThan, 1e-6) //总资产=已借+用户资产+锁定
				So(balance.NetAsset+balance.Borrowed, ShouldEqual, balance.Free)        //可用=已借+用户资产
			}
		})
	}
}

// U本位
func TestUBaseFutureBalance(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestFutureBalance", t, func() {
			balances, err := api.GetFutureBalance(common.Market_FUTURE)
			t.Logf("%#v\n", balances)
			So(err, ShouldBeNil)
			if balances == nil {
				t.Log("UBase balance is nil")
				return
			}
			So(balances.Rights, ShouldEqual, balances.Balance+balances.Unprofit)
			//todo
			//So(balances.Used+balances.Available, ShouldBeGreaterThanOrEqualTo, balances.Balance)
			for _, balance := range balances.UBaseBalanceList {
				So(balance.Asset, ShouldNotEqual, "")
				So(balance.Rights, ShouldEqual, balance.Balance+balance.Unprofit) //总权益=用户资产+未实现盈亏
				//So(balance.Rights, ShouldEqual, balance.Available+balance.Used)   //总权益=可用+已用
			}
			for _, position := range balances.UBasePositionList {
				So(position.Price, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.Position, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.MaintainMargin, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.InitialMargin, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.Notional, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.Leverage, ShouldBeGreaterThanOrEqualTo, 0)
			}
		})
	}
}

// 币本位
func TestCBaseFutureBalance(t *testing.T) {
	// 获取杠杆balance数据
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestFutureBalance", t, func() {
			balances, err := api.GetFutureBalance(common.Market_FUTURE_COIN)
			t.Logf("%#v\n", balances)
			t.Log(balances)
			So(err, ShouldBeNil)
			if balances == nil {
				t.Log("CBase balance is nil")
				return
			}
			So(balances.Rights, ShouldEqual, balances.Balance+balances.Unprofit)
			So(balances.Used+balances.Available, ShouldBeGreaterThanOrEqualTo, balances.Balance)
			for _, balance := range balances.UBaseBalanceList {
				So(balance.Asset, ShouldNotEqual, "")
				So(balance.Rights, ShouldEqual, balance.Balance+balance.Unprofit) //总权益=用户资产+未实现盈亏
				//So(balance.Rights, ShouldEqual, balance.Available+balance.Used)   //总权益=可用+已用
			}
			for _, position := range balances.UBasePositionList {
				So(position.Price, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.Position, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.MaintainMargin, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.InitialMargin, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.Notional, ShouldBeGreaterThanOrEqualTo, 0)
				So(position.Leverage, ShouldBeGreaterThanOrEqualTo, 0)
			}
		})
	}
}
