package test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warmplanet/proto/go/client"
)

//充值历史

func TestGetDepositHistory(t *testing.T) {
	for _, api := range apiList {
		t.Log(api.GetExchange())
		Convey(api.GetExchange().String()+"Test GetDepositHistory", t, func() {
			req := &client.DepositHistoryReq{
				StartTime: time.Date(2022, 8, 3, 1, 1, 1, 1, time.Local).UnixMilli(),
				Status:    client.DepositStatus_DEPOSITSTATUS_CONFIRMED,
			}
			depositHistoryRes, err := api.GetDepositHistory(req)
			So(err, ShouldBeNil)
			t.Log(depositHistoryRes)
			So(len(depositHistoryRes.DepositList), ShouldBeGreaterThanOrEqualTo, 0)
		})
	}
}
