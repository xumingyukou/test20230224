package spot_ws

import (
	"clients/exchange/cex/base"
	"context"
	"fmt"
	"testing"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

var (
	proxyUrl = "http://127.0.0.1:1080"
)

func TestTradeGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewHuobiWebsocket(conf, "api.hbdm.com")
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	// for _, symbol := range symbolNameMap {
	// 	info := &client.SymbolInfo{
	// 		Name:   symbol,
	// 		Symbol: symbol,
	// 		Type:   common.SymbolType_SPOT_NORMAL,
	// 	}
	// 	chMap[info] = make(chan *client.WsTradeRsp, 100)
	// }
	// symbolNameMap.Range(func(k, symbol interface{}) bool {
	// 	info := &client.SymbolInfo{
	// 		Name:   symbol.(string),
	// 		Symbol: symbol.(string),
	// 		Type:   common.SymbolType_SPOT_NORMAL,
	// 	}
	// 	chMap[info] = make(chan *client.WsTradeRsp, 100)
	// 	return true
	// })
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
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
	}
}

func TestBookTickerGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewHuobiWebsocket(conf, "api.hbdm.com")
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	//for _, symbol := range symbolNameMap {
	//	info := &client.SymbolInfo{
	//		Symbol: symbol,
	//		Name:   symbol,
	//		Type:   common.SymbolType_SPOT_NORMAL,
	//	}
	//	chMap[info] = make(chan *client.WsBookTickerRsp, 100)
	//}
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
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
	}
}

func TestDepthLimitGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: "wss://127.0.0.1:1080",
		},
	}
	wsClient := NewHuobiWebsocket(conf, "api.hbdm.com")
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	// for _, symbol := range symbolNameMap {
	// 	info := &client.SymbolInfo{
	// 		Symbol: symbol,
	// 		Name:   symbol,
	// 		Type:   common.SymbolType_SPOT_NORMAL,
	// 	}
	// 	chMap[info] = make(chan *depth.Depth, 100)
	// }
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
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
	}
}

func TestDepthIncrementGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewHuobiWebsocket(conf, "api.hbdm.com")
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	//for _, symbol := range symbolNameMap {
	//	info := &client.SymbolInfo{
	//		Symbol: symbol,
	//		Name:   symbol,
	//		Type:   common.SymbolType_SPOT_NORMAL,
	//	}
	//	chMap[info] = make(chan *client.WsDepthRsp, 100)
	//}
	err := wsClient.DepthIncrementGroup(ctx, 0, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintDepthIncrementGroup(ch)
	}
	select {}
}

func PrintDepthIncrementGroup(ch chan *client.WsDepthRsp) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
	}
}

func TestDepthIncrementGroupSnapshot(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}

	wsClient := NewHuobiWebsocket(conf, "wss://api.hbdm.com/swap-notification")
	ctx, _ := context.WithCancel(context.Background())
	chDeltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"BTC/USDT"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chDeltaMap[info] = make(chan *client.WsDepthRsp, 100)
		chFullMap[info] = make(chan *depth.Depth, 100)
		go PrintDepthGroup(chDeltaMap[info], "delta")
		go PrintDepthGroup2(chFullMap[info], "full")
	}
	err := wsClient.DepthIncrementSnapshotGroup(ctx, 1000, 400, true, true, chDeltaMap, chFullMap)
	if err != nil {
		t.Fatal(err)
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
			//fmt.Println(count, ok, title, res)
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

func TestUBaseFundingRate(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: "http://127.0.0.1:7890",
		},
	}
	wsClient := NewHuobiWebsocket(conf, "wss://api.hbdm.com/linear-swap-notification")
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	for _, symbol := range []string{"ETH/USD", "BTC/USDT"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP,
			Type:   common.SymbolType_SWAP_COIN_FOREVER,
		}
		chMap[info] = make(chan *client.WsFundingRateRsp, 10)
	}
	err := wsClient.FundingRateGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch, "FundingRateGroup")
	}
	select {}
}

func PrintBookTicker[T any](ch chan T, title string) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(title, ok, res)
			//fmt.Sprint(title, ok, res)
		}
	}
}
