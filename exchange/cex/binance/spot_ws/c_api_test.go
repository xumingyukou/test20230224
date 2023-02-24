package spot_ws

import (
	"clients/exchange/cex/base"
	"context"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"testing"
)

func TestCBaseTradeGroup(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceCBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP_COIN,
			Type:   common.SymbolType_SWAP_COIN_FOREVER,
		}
		chMap[info] = make(chan *client.WsTradeRsp, 10)
	}
	err := wsClient.TradeGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch, "TradeGroup")
	}
	select {}
}

func TestCBaseDepthIncrementGroup(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceCBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	for _, symbol := range []string{"BTC/USD"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP_COIN,
			Type:   common.SymbolType_SWAP_COIN_FOREVER,
		}
		chMap[info] = make(chan *client.WsDepthRsp, 10)
	}
	err := wsClient.DepthIncrementGroup(ctx, 100, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch, "DepthIncrementGroup")
	}
	select {}
}

func TestCBaseDepthLimitGroup(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceCBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP_COIN,
			Type:   common.SymbolType_SWAP_COIN_FOREVER,
		}
		chMap[info] = make(chan *depth.Depth, 10)
	}
	err := wsClient.DepthLimitGroup(ctx, 100, 10, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch, "DepthLimitGroup")
	}
	select {}
}

func TestCBaseBookTickerGroup(t *testing.T) {
	//wss://fstream.binance.com/stream?streams=btcusdt@bookTicker
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceCBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	for _, symbol := range []string{"BTC/USD"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_FUTURE_COIN,
			Type:   common.SymbolType_FUTURE_COIN_THIS_QUARTER,
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

func TestCBaseDepthIncrementSnapshotGroup(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceCBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"BTC/USD"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP_COIN,
			Type:   common.SymbolType_SWAP_COIN_FOREVER,
		}
		chMap[info] = make(chan *client.WsDepthRsp, 10)
		chFullMap[info] = make(chan *depth.Depth, 10)
	}
	err := wsClient.DepthIncrementSnapshotGroup(ctx, 0, 1000, true, true, chMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch, "DepthIncrementGroup")
	}
	select {}
}
