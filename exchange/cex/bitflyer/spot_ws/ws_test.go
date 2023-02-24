package spot_ws

import (
	"clients/exchange/cex/base"
	"context"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"testing"
)

var (
	b *BitFlyerWebsocket
)

func init() {
	//config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.APIConf{
		ProxyUrl: "http://127.0.0.1:7890",
		EndPoint: WS_PUBLIC_BASE_URL,
	}

	con := base.WsConf{
		APIConf: conf,
	}
	b = NewBitFlyerWebsocket(con)
}

func TestDepthIncre(t *testing.T) {
	b.ProxyUrl = "http://127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	symbol := client.SymbolInfo{
		Symbol: "BTC/JPY",
		Type:   1,
	}
	symbol_ := client.SymbolInfo{
		Symbol: "LTC/USD",
		Type:   1,
	}
	depthMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	ch1 := make(chan *client.WsDepthRsp)
	ch2 := make(chan *client.WsDepthRsp)
	depthMap[&symbol] = ch1
	depthMap[&symbol_] = ch2

	b.DepthIncrementGroup(ctx, 10, depthMap)
	//if err != nil {
	//	t.Fatal(err)
	//}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2  res", res)
		}
	}
}

func TestTradeGroup(t *testing.T) {
	b.ProxyUrl = "http://127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	symbol := client.SymbolInfo{
		Symbol: "BTC/USD",
		Type:   1,
	}
	symbol_ := client.SymbolInfo{
		Symbol: "BTC/JPY",
		Type:   1,
	}
	tradeMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	ch1 := make(chan *client.WsTradeRsp)
	ch2 := make(chan *client.WsTradeRsp)
	tradeMap[&symbol] = ch1
	tradeMap[&symbol_] = ch2

	err := b.TradeGroup(ctx, tradeMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2  res", res)
		}
	}
}

func TestTickerGroup(t *testing.T) {
	b.ProxyUrl = "http://127.0.0.1:7890"
	ctx, _ := context.WithCancel(context.Background())
	symbol := client.SymbolInfo{
		Symbol: "BTC/JPY",
		Type:   1,
	}
	symbol_ := client.SymbolInfo{
		Symbol: "BTC/USD",
		Type:   1,
	}
	bookMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	ch1 := make(chan *client.WsBookTickerRsp)
	ch2 := make(chan *client.WsBookTickerRsp)
	bookMap[&symbol] = ch2
	bookMap[&symbol_] = ch1

	err := b.BookTickerGroup(ctx, bookMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, _ := <-ch1:
			fmt.Println("ch1 res", res)
		case res, _ := <-ch2:
			fmt.Println("ch2  res", res)
		}
	}
}

func TestSnapShotGroup(t *testing.T) {
	b.ProxyUrl = "http://127.0.0.1:7890"

	ch1 := make(chan *client.WsDepthRsp)
	ch2 := make(chan *depth.Depth)
	chDelatMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	symbol1 := client.SymbolInfo{
		Symbol: "BTC/USD",
		Type:   1,
	}

	chDelatMap[&symbol1] = ch1
	chFullMap[&symbol1] = ch2
	//
	symbol2 := client.SymbolInfo{
		Symbol: "BTC/JPY",
		Type:   1,
	}
	ch3 := make(chan *client.WsDepthRsp)
	ch4 := make(chan *depth.Depth)
	chDelatMap[&symbol2] = ch3
	chFullMap[&symbol2] = ch4
	//symbols := []string{"BTC/JPY", "XRP/JPY", "ETH/JPY", "XLM/JPY", "MONA/JPY", "ETH/BTC", "BCH/BTC"}
	//for i := 0; i < len(symbols); i += 10 {
	//	for j := i; j < i+10 && j < len(symbols); j++ {
	//
	//		symbol := client.SymbolInfo{Symbol: strings.ReplaceAll(symbols[j], "XBT", "BTC"), Type: 1}
	//		chDelatMap[&symbol] = nil
	//		chFullMap[&symbol] = nil
	//
	//	}
	//
	//}

	ctx, _ := context.WithCancel(context.Background())
	err := b.DepthIncrementSnapshotGroup(ctx, 0, 1000, true, true, chDelatMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res := <-ch1:
			fmt.Println("增量ch1 res", res)
		case res := <-ch2:
			fmt.Println("全量i111212  res", res)
		case res := <-ch3:
			fmt.Println("ch3 res", res)
		case res := <-ch4:
			fmt.Println("ch4 i111212  res", res)
		}
		//con()
		//break
	}
}
