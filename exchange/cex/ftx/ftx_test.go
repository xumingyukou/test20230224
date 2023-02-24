package ftx

import (
	"clients/config"
	"clients/exchange/cex/base"
	"fmt"
	"testing"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/order"
)

var (
	ProxyUrl = "http://127.0.0.1:9999"
)

func TestOrder(t *testing.T) {
	var timeOffset int64 = 30
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		//ProxyUrl:    ProxyUrl,
		ReadTimeout: timeOffset,
		AccessKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	precisionMap := map[string]*client.PrecisionItem{
		"ETH/USDT": &client.PrecisionItem{
			Symbol:    "ETH/USDT",
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    8,
			Price:     8,
			AmountMin: 0.001,
		},
	}
	ftx := NewClientFTX(conf, precisionMap)
	res, err := ftx.PlaceOrder(&order.OrderTradeCEX{
		Hdr: &common.MsgHeader{
			Producer: []byte("1"),
		},
		Base: &order.OrderBase{
			Id:       112400,
			Quote:    []byte("USDT"),
			Token:    []byte("ETH"),
			Exchange: common.Exchange_FTX,
			Market:   common.Market_SPOT,
			Type:     common.SymbolType_SPOT_NORMAL,
			Symbol:   []byte("ETH/USDT"),
		},
		OrderType: order.OrderType_LIMIT,
		Side:      order.TradeSide_BUY,
		Price:     1500,
		Tif:       order.TimeInForce_GTC,
		Amount:    0.001,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetSymbols(t *testing.T) {
	var timeOffset int64 = 400
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		ProxyUrl:    ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	res := ftx.GetSymbols()
	fmt.Println(res)
}

func TestGetPrecision(t *testing.T) {
	conf := base.APIConf{
		ProxyUrl: ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	res, err := ftx.GetPrecision("BTC/USDT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

// Old Test function
func TestDepth(t *testing.T) {
	var timeOffset int64 = 30
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		ProxyUrl:    ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	res, err := ftx.GetDepth(nil, 10) //Try with BTC/USDT later
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestBalance(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl:  ProxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}

	ftx := NewClientFTX(conf)
	res, err := ftx.GetBalance()
	fmt.Println("Did we get here?")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
	res3, err := ftx.GetMarginBalance()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res3)
}

func TestGetTradeFee(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl:  ProxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	ftx := NewClientFTX(conf)
	res, err := ftx.GetTradeFee() //"BTC/USDT")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestClientFTX_GetDepositHistory(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		//ProxyUrl:  ProxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	ftx := NewClientFTX(conf)

	request := &client.DepositHistoryReq{
		StartTime: time.Date(2022, 04, 01, 01, 01, 01, 01, time.Local).Unix(),
		EndTime:   time.Now().Unix(),
	}
	res, err := ftx.GetDepositHistory(request)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestUDepth(t *testing.T) {
	var timeOffset int64 = 400
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		ProxyUrl:    ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	tester := &client.SymbolInfo{
		Symbol: "ATOM/USD",
		Name:   "ATOM-0930",
		Type:   common.SymbolType_FUTURE_THIS_QUARTER,
	}

	res1, err := ftx.GetDepth(tester, 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res1)
}

func TestGetFutureDepth(t *testing.T) {
	var timeOffset int64 = 400
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		//ProxyUrl:    ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	tester := &client.SymbolInfo{
		Symbol: "BTC/USD",
		Name:   "",
		Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
	}

	res1, err := ftx.GetFutureDepth(tester, 10)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res1)
}

func TestGetFutureSymbols(t *testing.T) {
	var timeOffset int64 = 400
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		//ProxyUrl:    ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	res := ftx.GetFutureSymbols(common.Market_FUTURE)
	fmt.Println(res)
}

func TestGetFutureMarkPrice(t *testing.T) {
	var timeOffset int64 = 400
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		//ProxyUrl:    ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	//tester := &client.SymbolInfo{
	//	Symbol: "BTC/USD",
	//	Name:   "",
	//	Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
	//}
	res, err := ftx.GetFutureMarkPrice(common.Market_FUTURE)
	if err != nil {
		fmt.Println("ERROR")
	}
	fmt.Println(res)
}

func TestGetFutureTradeFee(t *testing.T) {
	var timeOffset int64 = 400
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		AccessKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	ftx := NewClientFTX(conf)
	//tester := &client.SymbolInfo{
	//	Symbol: "BTC/USD",
	//	Name:   "",
	//	Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
	//}
	res, err := ftx.GetFutureTradeFee(common.Market_FUTURE)
	if err != nil {
		fmt.Println("ERROR")
	}
	fmt.Println(res)
}

func TestGetFuturePrecision(t *testing.T) {
	var timeOffset int64 = 400
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		//ProxyUrl:    ProxyUrl,
	}
	ftx := NewClientFTX(conf)
	tester := &client.SymbolInfo{
		Symbol: "BTC/USD",
		Name:   "",
		Type:   common.SymbolType_FUTURE_NEXT_QUARTER,
	}
	res, err := ftx.GetFuturePrecision(common.Market_FUTURE, tester)
	if err != nil {
		fmt.Println("ERROR")
	}
	fmt.Println(res)
}

func TestGetFutureBalance(t *testing.T) {
	var timeOffset int64 = 400
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		AccessKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	ftx := NewClientFTX(conf)
	res, err := ftx.GetFutureBalance(common.Market_FUTURE)
	if err != nil {
		fmt.Println("ERROR")
	}
	fmt.Println(res)
}

func TestWithdrawFee(t *testing.T) {
	var timeOffset int64 = 400
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl:    ProxyUrl,
		ReadTimeout: timeOffset,
		AccessKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	ftx := NewClientFTX(conf)
	res, err := ftx.GetTransferFee(common.Chain_SOLANA, "USDT", "BTC", "ETH", "SOL")
	if err != nil {
		fmt.Println(err)
	}
	print(res)
}

func Test_MoveAsset(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		//ProxyUrl:  ProxyUrl,
		AccessKey:  config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey:  config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
		SubAccount: "ftx_001",
	}
	ftx := NewClientFTX(conf)
	testOrder := &order.OrderMove{
		Asset:         "USD",
		Amount:        1,
		AccountSource: "ftx_001",
		AccountTarget: "ftx_002",
		ActionUser:    order.OrderMoveUserType_Sub,
	}
	res, err := ftx.MoveAsset(testOrder)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Println(res)
}

func Test_GetMoveHistory(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		//ProxyUrl:  ProxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	ftx := NewClientFTX(conf)
	testOrder := &client.MoveHistoryReq{
		AccountSource: "ftx_001",
		AccountTarget: "ftx_002",
		ActionUser:    order.OrderMoveUserType_Sub,
	}
	res, err := ftx.GetMoveHistory(testOrder)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Println(res)
}

func TestClientFTX_GetOrderHistory(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		//ProxyUrl:  ProxyUrl,
		AccessKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
		SecretKey: config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
	}
	ftx := NewClientFTX(conf)
	historyRequest := &client.OrderHistoryReq{
		Asset:  "ETH/USD",
		Market: common.Market_SWAP,
	}
	res, err := ftx.GetOrderHistory(historyRequest)
	if err != nil {
		fmt.Println("ERROR")
	}
	fmt.Println(res)
}
