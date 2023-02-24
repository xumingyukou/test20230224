package spot_ws

import (
	"clients/config"
	"clients/exchange/cex/base"
	"context"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"testing"
)

var (
	proxyUrl = "http://127.0.0.1:9999"
)

func TestBookTickerGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewFTXWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	for _, symbol := range []string{"ETH/USD", "BTC/USD"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Name:   symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsBookTickerRsp, 10)
	}
	err := wsClient.BookTickerGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch)
	}
	select {}
}

func TestDepthLimitGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewFTXWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"ETH/USD", "BTC/USD"} {
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
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewFTXWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	for _, symbol := range []string{"ETH/USD", "BTC/USD"} {
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

func TestTradeGroup(t *testing.T) {
	//d1 := time.Date(2022, 8, 28, 1, 0, 0, 0, time.UTC)
	//fmt.Println(transform.GetThisQuarter(d1, 5, 2).Format("0102"))
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewFTXWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range []string{"USDTBEAR/USD", "ETH/USD", "BTC/USD"} { //"ETH/USD", "BTC/USD",
		info := &client.SymbolInfo{
			Symbol: symbol,
			Name:   symbol,
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

func TestDepthIncrementGroupSnapshot(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}

	wsClient := NewFTXWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chDeltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"ADABEAR/USD", "ETH/USDT", "kuco", "BTC/USD"} { //
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chDeltaMap[info] = make(chan *client.WsDepthRsp, 10)
		chFullMap[info] = make(chan *depth.Depth, 10)
		go PrintDepthGroup(chDeltaMap[info], "delta")
		go PrintDepthGroup2(chFullMap[info], "full")
	}
	err := wsClient.DepthIncrementSnapshotGroup(ctx, 1000, 20, true, true, chDeltaMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	select {}
}

func PrintBookTicker(ch chan *client.WsBookTickerRsp) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
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

func PrintTradeGroup(ch chan *client.WsTradeRsp) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
	}
}

func Test_Order(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl:    proxyUrl,
			ReadTimeout: 3000,
			AccessKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList["ftx"].ApiKeyConfig.SecretKey,
		},
	}
	wsClient := NewFTXWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	data := &client.WsAccountReq{
		Exchange: common.Exchange_FTX,
		Market:   common.Market_SPOT,
		Type:     common.SymbolType_SPOT_NORMAL,
	}
	ch, err := wsClient.Order(ctx, data)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
	}
}
