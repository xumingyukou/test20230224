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

func TestUBaseTradeGroup(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range []string{"BTC/USDT"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_FUTURE,
			Type:   common.SymbolType_FUTURE_THIS_QUARTER,
		}
		chMap[info] = make(chan *client.WsTradeRsp, 10)
		//info = &client.SymbolInfo{
		//	Symbol: symbol,
		//	Type:   common.SymbolType_SWAP_FOREVER,
		//}
		//chMap[info] = make(chan *client.WsTradeRsp, 10)
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

func TestUBaseDepthIncrement(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	for _, symbol := range []string{"BTC/USDT", "ETH/USDT"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP,
			Type:   common.SymbolType_SWAP_FOREVER,
		}
		chMap[info] = make(chan *client.WsDepthRsp, 10)
	}
	err := wsClient.DepthIncrementGroup(ctx, 0, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range chMap {
		go PrintBookTicker(ch, "DepthIncrementGroup")
	}
	select {}
}

func TestUBaseDepthLimit(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"BTC/USDT", "BTC/USDC"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP,
			Type:   common.SymbolType_SWAP_FOREVER,
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

func TestUBaseBookTicker(t *testing.T) {
	//wss://fstream.binance.com/stream?streams=btcusdt@bookTicker
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	for _, symbol := range []string{"BTC/USDT"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP,
			Type:   common.SymbolType_SWAP_FOREVER,
		}
		chMap[info] = make(chan *client.WsBookTickerRsp, 10)
		info = &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_FUTURE,
			Type:   common.SymbolType_FUTURE_THIS_QUARTER,
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

func TestUBaseBinanceSpotWebsocket_Account(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl:    proxyUrl,
			ReadTimeout: 3000,
			AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	data := &client.WsAccountReq{
		Exchange: common.Exchange_BINANCE,
		Market:   common.Market_FUTURE,
		Type:     common.SymbolType_SPOT_NORMAL,
	}
	ch, err := wsClient.Account(ctx, data)
	fmt.Println("start")
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
		}
		break
	}
}

func TestUBaseBinanceSpotWebsocket_Order(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl:    proxyUrl,
			ReadTimeout: 3000,
			AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	data := &client.WsAccountReq{
		Exchange: common.Exchange_BINANCE,
		Market:   common.Market_FUTURE,
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
		break
	}
}

func TestUBaseFundingRate(t *testing.T) {
	//wss://fstream.binance.com/stream?streams=btcusdt@bookTicker
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	for _, symbol := range []string{"BTC/USDT", "BTC/USDC"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP,
			Type:   common.SymbolType_SWAP_FOREVER,
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

func TestUBaseDepthIncrementSnapshotGroup(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceUBaseWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"BTC/USDT"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Market: common.Market_SWAP,
			Type:   common.SymbolType_SWAP_FOREVER,
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
