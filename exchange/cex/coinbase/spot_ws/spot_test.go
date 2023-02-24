package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/coinbase"
	"context"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"testing"
)

var (
	proxyUrl = "http://127.0.0.1:7890"
)

func TestTradeGroup(t *testing.T) {
	//"wss://ws-feed.exchange.coinbase.com"
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewCoinbaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range []string{"ETH/BTC"} {
		info := &client.SymbolInfo{
			Name:   symbol,
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsTradeRsp, 10)
	}
	err := wsClient.TradeGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintTrades(ch)
	}
	select {}
}

func TestBookTickerGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewCoinbaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	for _, symbol := range []string{"ETH/BTC", "ETH/USD"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsBookTickerRsp, 10)
	}
	err := wsClient.BookTickerGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch, "BookTickerGroup")
	}
	select {}
}

func TestDepthLimitGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewCoinbaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"ETH/USD", "BTC/USD", "nonce"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Name:   symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *depth.Depth, 10)
	}
	err := wsClient.DepthLimitGroup(ctx, 10, 10, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintDepthGroup2(ch, "Limit")
	}
	select {}
}

func TestDepthIncrementGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewCoinbaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	for _, symbol := range []string{"ETH/USD", "BTC/USD", "kucasd"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Name:   symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsDepthRsp, 10)
	}
	err := wsClient.DepthIncrementGroup(ctx, 10, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintDepthGroup(ch, "increment")
	}
	select {}
}

func TestDepthIncrementGroupSnapshot(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}

	wsClient := NewCoinbaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chDeltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)

	coinbase := coinbase.NewClientCoinbase(conf.APIConf)
	symbols := coinbase.GetSymbols()
	symbols = []string{"LCX/USD", "ADA/USD"}
	for _, symbol := range symbols {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chDeltaMap[info] = make(chan *client.WsDepthRsp, 50)
		chFullMap[info] = make(chan *depth.Depth, 50)
		go PrintDepthGroup(chDeltaMap[info], "delta")
		go PrintDepthGroup2(chFullMap[info], "full")
	}
	err := wsClient.DepthIncrementSnapshotGroup(ctx, 1000, 20, true, true, chDeltaMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	select {}
}

/*Helper Functions*/
func PrintTrades(ch chan *client.WsTradeRsp) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
	}
}
func PrintBookTicker[T any](ch chan T, title string) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(title, ok, res)
		}
	}
}
func PrintDepthGroup(ch chan *client.WsDepthRsp, title string) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, title, res)
		}
	}
}
func PrintDepthGroup2(ch chan *depth.Depth, title string) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, title, res)
		}
	}
}

func TestErTest(t *testing.T) {
	//"wss://ws-feed.exchange.coinbase.com"
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewCoinbaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
		info := &client.SymbolInfo{
			Name:   symbol,
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsTradeRsp, 10)
	}
	err := wsClient.Tester(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintTrades(ch)
	}
	select {}
}

/*Deprecate Function
func TestGetFullDepth(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewCoinbaseWebsocket(conf)
	x, err := wsClient.GetFullDepth("ETH/USD")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Exchange time:", x.Asks)
}
*/
