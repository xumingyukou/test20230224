package spot_ws

import (
	"clients/exchange/cex/base"
	"context"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"testing"
)

func TestTradeGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBybitSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range []string{"BTC/USDT"} {
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
		go PrintTradeGroup(ch)
	}
	select {}
}

func PrintTradeGroup(ch chan *client.WsTradeRsp) {
	count := 0
	for {
		select {
		case res, ok := <-ch:
			count++
			if count%1 == 0 {
				fmt.Println(count, ok, res)
			}
		}
	}
}

func TestBookTickerGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBybitSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	for _, symbol := range symbolNameMap {
		info := &client.SymbolInfo{
			Name:   symbol,
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
		go PrintBookTickerGroup(ch)
	}
	select {}
}

func PrintBookTickerGroup(ch chan *client.WsBookTickerRsp) {
	count := 0
	for {
		select {
		case res, ok := <-ch:
			count++
			if count%100 == 0 {
				fmt.Println(count, ok, res)
			}
		}
	}
}

func TestDepthLimitGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBybitSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range symbolNameMap {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Name:   symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *depth.Depth, 10)
	}
	err := wsClient.DepthLimitGroup(ctx, 0, 10, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintDepthLimitGroup(ch)
	}
	select {}
}

func PrintDepthLimitGroup(ch chan *depth.Depth) {
	count := 0
	for {
		select {
		case res, ok := <-ch:
			count++
			if count%100 == 0 {
				fmt.Println(count, ok, res)
			}
		}
	}
}

func TestDepthIncrementGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBybitSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	for _, symbol := range symbolNameMap {
		info := &client.SymbolInfo{
			Name:   symbol,
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsDepthRsp, 10)
	}
	err := wsClient.DepthIncrementGroup(ctx, 0, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintIncrementGroup(ch)
	}
	select {}
}

func PrintIncrementGroup(ch chan *client.WsDepthRsp) {
	count := 0
	for {
		select {
		case res, ok := <-ch:
			count++
			if count%100 == 0 {
				fmt.Println(count, ok, res)
			}
		}
	}
}

func PrintChannel[T any](ch chan T, title string) {
	for {
		select {
		case res, ok := <-ch:
			if title == "" {
				continue
			}
			fmt.Println(title, ok, res)
		}
	}
}

func PrintDepthChannel(ch chan *depth.Depth, title string) {
	for {
		select {
		case res, ok := <-ch:
			if title == "" {
				continue
			}
			fmt.Println(title, ok, res)
			if res.Asks[0].Price <= res.Bids[0].Price {
				fmt.Println(223234e23)
			}
		}
	}
}
func TestDepthIncrementSnapshotGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: "http://127.0.0.1:7890",
		},
	}

	wsClient := NewBybitSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chDeltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	//for _, symbol := range symbolNameMap {
	for _, symbol := range []string{"BTC/USDT"} { //, "ETH/USDT"
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chDeltaMap[info] = make(chan *client.WsDepthRsp, 10)
		chFullMap[info] = make(chan *depth.Depth, 10)
		//go PrintDepthGroup(chDeltaMap[info], "delta")
		//go PrintDepthGroup2(chFullMap[info], "full")
	}
	err := wsClient.DepthIncrementSnapshotGroup(ctx, 10, 200, true, true, chDeltaMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chDeltaMap {
		go PrintChannel(ch, "")
	}
	for symbol, ch := range chFullMap {
		go PrintDepthChannel(ch, symbol.Symbol+"full depth")
	}
	select {}
}

func PrintDepthGroup2(ch chan *depth.Depth, title string) {
	count := 0
	for {
		select {
		case res, ok := <-ch:
			count++
			if count%100 == 0 {
				fmt.Println(count, ok, title, res)
			}
		}
	}
}

func PrintDepthGroup(ch chan *client.WsDepthRsp, title string) {
	count := 0
	for {
		select {
		case res, ok := <-ch:
			count++
			if count%100 == 0 {
				fmt.Println(count, ok, title, res)
			}
		}
	}
}
