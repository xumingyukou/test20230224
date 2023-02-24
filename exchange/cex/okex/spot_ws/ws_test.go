package spot_ws

import (
	"clients/config"
	"clients/exchange/cex/base"
	"context"
	"fmt"
	"testing"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

var (
	b *OkWebsocket
)

func init() {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl:   "http://127.0.0.1:7890",
		EndPoint:   WS_API_BASE_URL,
		AccessKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.AccessKey,
		SecretKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.SecretKey,
		Passphrase: config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.Passphrase,
		//IsTest:   true,
	}

	if conf.IsTest {
		conf.EndPoint = TEST_WS_API_BASE_URL
	}
	con := base.WsConf{
		APIConf: conf,
	}
	b = NewOkWebsocket(con, common.Market_SPOT)
}

func TestDepthIncreGroup(t *testing.T) {
	ch1 := make(chan *client.WsDepthRsp)
	ch2 := make(chan *client.WsDepthRsp)
	bookMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	symbol1 := client.SymbolInfo{
		Symbol: "BTC/USDT",
	}
	symbol2 := client.SymbolInfo{
		Symbol: "ETH/USDT",
	}
	bookMap[&symbol1] = ch1
	bookMap[&symbol2] = ch2
	ctx, _ := context.WithCancel(context.Background())
	err := b.DepthIncrementGroup(ctx, 0, bookMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2 res", res)
		}
		//break
	}
}

func TestTradeGroup(t *testing.T) {
	ch1 := make(chan *client.WsTradeRsp)
	ch2 := make(chan *client.WsTradeRsp)
	bookMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	symbol1 := client.SymbolInfo{
		Symbol: "BTC/USDT",
	}
	symbol2 := client.SymbolInfo{
		Symbol: "ETH/USDT",
	}
	bookMap[&symbol1] = ch1
	bookMap[&symbol2] = ch2
	ctx, _ := context.WithCancel(context.Background())
	err := b.TradeGroup(ctx, bookMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2 res", res)
		}
		//break
	}
}

func TestBookTickerGroup(t *testing.T) {
	ch1 := make(chan *client.WsBookTickerRsp)
	ch2 := make(chan *client.WsBookTickerRsp)
	bookMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	symbol1 := client.SymbolInfo{
		Symbol: "BTC/USDT",
	}
	symbol2 := client.SymbolInfo{
		Symbol: "ETH/USDT",
	}
	bookMap[&symbol1] = ch1
	bookMap[&symbol2] = ch2
	ctx, _ := context.WithCancel(context.Background())
	err := b.BookTickerGroup(ctx, bookMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2 res", res)
		}
		//break
	}
}

func TestDepthLimitGroup(t *testing.T) {
	ch1 := make(chan *depth.Depth)
	ch2 := make(chan *depth.Depth)
	bookMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	symbol1 := client.SymbolInfo{
		Symbol: "BTC/USDT",
	}
	symbol2 := client.SymbolInfo{
		Symbol: "ETH/USDT",
	}
	bookMap[&symbol1] = ch1
	bookMap[&symbol2] = ch2
	ctx, _ := context.WithCancel(context.Background())
	err := b.DepthLimitGroup(ctx, 0, 0, bookMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2 res", res)
		}
		//break
	}
}

// 推送时间没有变化
func TestAccount(t *testing.T) {
	b.ProxyUrl = "127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	ch, err := b.Account(ctx, &client.WsAccountReq{
		Market: 4,
	})
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch:
			fmt.Println("chan res:", res)
		}
		//break

	}
}

func TestOrder(t *testing.T) {
	b.ProxyUrl = "127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	ch, err := b.Order(ctx, &client.WsAccountReq{Market: common.Market_FUTURE})
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch:
			fmt.Println("chan res:", res)
			//case <-time.After(time.Second):
		}
		//break

	}
}

func TestIncrementSnapshotGroup(t *testing.T) {
	ch1 := make(chan *client.WsDepthRsp)
	ch2 := make(chan *depth.Depth)
	chDelatMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	symbol1 := client.SymbolInfo{
		Symbol: "BTC/USDT",
	}
	ch3 := make(chan *client.WsDepthRsp)
	ch4 := make(chan *depth.Depth)
	symbol2 := client.SymbolInfo{
		Symbol: "LTC/USDT",
	}
	chDelatMap[&symbol1] = ch1
	chFullMap[&symbol1] = ch2
	chDelatMap[&symbol2] = ch3
	chFullMap[&symbol2] = ch4
	ctx, _ := context.WithCancel(context.Background())
	err := b.DepthIncrementSnapshotGroup(ctx, 0, 1000, true, true, chDelatMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2 i111212  res", res)
		case res, _ := <-ch3:
			fmt.Println("ch3 res", res)
		case res, _ := <-ch4:
			fmt.Println("ch4 i111212  res", res)
		}
		//break
	}
}
