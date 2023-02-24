package main

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance"
	"clients/exchange/cex/okex"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
	"github.com/warmplanet/proto/go/sdk"
)

var (
	b          *binance.ClientBinance
	o          *okex.ClientOkex
	timeOffset int64 = 5
	conf       base.APIConf
)

func init() {
	apiKey, seckey := new(sdk.SecretConfig), new(sdk.SecretConfig)
	err := apiKey.UnmarshalText([]byte("awskms:/ap-northeast-1/trading/binance/binanmcemain/api/apikey"))
	if err != nil {
		panic(err)
	}
	err = seckey.UnmarshalText([]byte("awskms:/ap-northeast-1/trading/binance/binanmcemain/api/secretkey"))
	if err != nil {
		panic(err)
	}
	conf = base.APIConf{
		ReadTimeout: timeOffset,
		AccessKey:   string(apiKey.V),
		SecretKey:   string(seckey.V),
	}
	b = binance.NewClientBinance(conf)
	//o = okex.NewClientOkex(conf)

	//testmove/binance/apikey  和  testmove/binance/secretkey
}

func subaccountlist() {
	res2, err := b.SubAccountList()
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range res2.SubAccounts {

		fmt.Printf("%#v\n", i)
	}
}
func balance() {
	res, err := b.GetBalance()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range res.BalanceList {
		if item.Total == 0 {
			continue
		}
		fmt.Println(item)
	}
	for _, item := range res.WalletList {
		if item.Total == 0 {
			continue
		}
		fmt.Println(item)
	}
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
		AccountTarget: "testmove_virtual@b4xaf7pvnoemail.com",
		ActionUser:    order.OrderMoveUserType_Master,
		EndTime:       1672822975884,
	}
	res, err := b.GetMoveHistory(req1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}

func moveOkex() {
	reqOk := &order.OrderMove{}
	o.MoveAsset(reqOk)
}

func main() {
	//subaccountlist()
	//	balance()
	//moveMaster2Sub()
	balance()
}
