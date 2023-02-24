package main

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/okex"
	"fmt"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
	"github.com/warmplanet/proto/go/sdk"
)

var (
	b              *okex.ClientOkex
	timeOffset     int64 = 5
	conf           base.APIConf
	apiKey, seckey = new(sdk.SecretConfig), new(sdk.SecretConfig)
)

func init() {
	if err := apiKey.UnmarshalText([]byte("awskms:/ap-northeast-1/testmove/okex/apikey")); err != nil {
		panic(err)
	}
	if err := seckey.UnmarshalText([]byte("awskms:/ap-northeast-1/testmove/binance/secretkey")); err != nil {
		panic(err)
	}
	conf = base.APIConf{
		ReadTimeout: timeOffset,
		AccessKey:   string(apiKey.V),
		SecretKey:   string(seckey.V),
	}
	b = okex.NewClientOkex(conf)

	//testmove/binance/apikey  和  testmove/binance/secretkey
}

func moveMaster2Sub() {
	req := &order.OrderMove{
		Asset:         "USDT",
		Amount:        10,
		Source:        common.Market_SPOT,
		Target:        common.Market_SPOT,
		AccountTarget: "testmove_virtual@b4xaf7pvnoemail.com",
	}
	res2, err := b.MoveAsset(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res2)
}
func getuniversalTransfer() {
	req1 := &client.MoveHistoryReq{
		Source:        common.Market_SPOT,
		Target:        common.Market_SPOT,
		AccountSource: "testmove_virtual@b4xaf7pvnoemail.com",
		ActionUser:    order.OrderMoveUserType_Master,
	}
	res, err := b.GetMoveHistory(req1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}

func main() {
	moveMaster2Sub()
}
