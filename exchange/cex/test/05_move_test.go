package test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

var (
	moveAsset = "USDT"
)

//划转+划转历史

func TestMoveAsset(t *testing.T) {
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestMoveAsset", t, func() {
			req := &order.OrderMove{
				Hdr:    &common.MsgHeader{},
				Base:   &order.OrderBase{},
				Asset:  moveAsset,
				Amount: 10,
				Source: common.Market_FUTURE,
				Target: common.Market_SPOT,
			}
			moveRes, err := api.MoveAsset(req)
			So(err, ShouldBeNil)
			if moveRes == nil {
				t.Log("move res is nil")
				return
			}
			t.Log(moveRes)
			So(moveRes, ShouldNotBeNil)
		})
	}
}

func TestGetMoveHistory(t *testing.T) {
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test TestGetMoveHistory", t, func() {
			req := &client.MoveHistoryReq{
				Source:    common.Market_SPOT,
				Target:    common.Market_WALLET,
				StartTime: time.Date(2022, 11, 10, 1, 1, 1, 1, time.Local).UnixMilli(),
			}
			moveHistoryRes, err := api.GetMoveHistory(req)
			So(err, ShouldBeNil)
			if moveHistoryRes == nil {
				t.Log("moveHistoryRes is nil")
				return
			}
			for _, moveItem := range moveHistoryRes.MoveList {
				t.Log(moveItem)
			}
		})
	}
}
