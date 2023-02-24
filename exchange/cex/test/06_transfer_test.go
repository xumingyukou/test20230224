package test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
	"testing"
)

//提币 + 提币历史
var (
	transferAsset   = "USDT"
	transferAddress = "0x1cb3643Db2E039a4abdee676c544dfC13c9F1ea2"
	transferChain   = common.Chain_ETH
)

func TestTransfer(t *testing.T) {
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test Transfer", t, func() {
			req := &order.OrderTransfer{
				Hdr: &common.MsgHeader{
					Producer: []byte("test"),
				},
				Base:            &order.OrderBase{},
				Chain:           transferChain,
				Amount:          10,
				TransferAddress: []byte(transferAddress),
				ExchangeToken:   []byte(transferAsset),
			}
			transferRes, err := api.Transfer(req)
			So(err, ShouldBeNil)
			So(transferRes.OrderId, ShouldNotEqual, "")
			So(transferRes.Id, ShouldNotEqual, 0)
		})
	}
}

func TestGetTransferHistory(t *testing.T) {
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test GetTransferHistory", t, func() {
			req := &client.TransferHistoryReq{}
			transferHistoryRes, err := api.GetTransferHistory(req)
			So(err, ShouldBeNil)
			So(transferHistoryRes, ShouldNotBeNil)
		})
	}
}
