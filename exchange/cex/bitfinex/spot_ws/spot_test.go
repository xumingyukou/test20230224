package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitfinex"
	"context"
	"fmt"
	"testing"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

var (
// proxyUrl = "http://127.0.0.1:9999"
)

func TestTradeGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl:  proxyUrl,
		},
	}
	wsClient := NewBitfinexSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range Exchange2Client {
		//for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
		info := &client.SymbolInfo{
			Name:   symbol,
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsTradeRsp, 100)
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
			if count%100 == 0 {
				fmt.Println(count, ok, res)
			}
		}
	}
}

func TestBookTickerGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl:  proxyUrl,
		},
	}
	wsClient := NewBitfinexSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	for _, symbol := range []string{"APENFT/UST", "AMP/USD"} {
		//for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
		info := &client.SymbolInfo{
			Name:   symbol,
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsBookTickerRsp, 100)
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

func TestIncrementGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl:  proxyUrl,
		},
	}
	wsClient := NewBitfinexSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	//for _, symbol := range symbolNameMap {
	for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
		info := &client.SymbolInfo{
			Name:   symbol,
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chMap[info] = make(chan *client.WsDepthRsp, 100)
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
			//fmt.Println(count, ok, title, res)
		}
	}
}

func TestDepthIncrementGroupSnapshot(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			// ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBitfinexSpotWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chDeltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	// for _, symbol := range symbolNameMap {
	for _, symbol := range []string{"APENFT/UST", "AMP/USD", "ALG/USD", "ALBT/USD", "AIX/UST", "ANT/ETH", "AMP/BTC", "ATO/BTC", "AAA/BBB", "ANT/USD", "APE/UST", "ADA/UST", "AMP/UST", "ANT/BTC"} {
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
			if count%1000 == 0 {
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
			if count%1000 == 0 {
				fmt.Println(count, ok, title, res)
			}
		}
	}
}

func TestAllSymbol(t *testing.T) {
	conf := base.APIConf{}
	a := bitfinex.NewClientBitfinex(conf)
	allSymbolList := a.GetSymbols()
	var num = 30
	quantity := (len(allSymbolList) / num) + 1
	var symbolLists = make([][]string, 0)
	var start, end, i int
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			symbolLists = append(symbolLists, allSymbolList[start:end])
		} else {
			symbolLists = append(symbolLists, allSymbolList[start:])
		}
		start = i * num
	}
	fmt.Println(symbolLists)
	for _, symbolList := range symbolLists {
		fmt.Println(symbolList)
		conf := base.WsConf{
			APIConf: base.APIConf{
				// ProxyUrl:  proxyUrl,
			},
		}
		wsClient := NewBitfinexSpotWebsocket(conf)
		ctx, _ := context.WithCancel(context.Background())
		chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
		for _, symbol := range symbolList {
			//for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
			info := &client.SymbolInfo{
				Name:   symbol,
				Symbol: symbol,
				Type:   common.SymbolType_SPOT_NORMAL,
			}
			chMap[info] = make(chan *client.WsTradeRsp, 100)
		}
		err := wsClient.TradeGroup(ctx, chMap)
		if err != nil {
			t.Fatal(err)
		}
		for _, ch := range chMap {
			go PrintTradeGroup(ch)
		}
		//select {}
	}
	for _, symbolList := range symbolLists {
		fmt.Println(symbolList)
		conf := base.WsConf{
			APIConf: base.APIConf{
				// ProxyUrl:  proxyUrl,
			},
		}
		wsClient := NewBitfinexSpotWebsocket(conf)
		ctx, _ := context.WithCancel(context.Background())
		chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
		for _, symbol := range symbolList {
			//for _, symbol := range []string{"BTC/USD", "ETH/USD"} {
			info := &client.SymbolInfo{
				Name:   symbol,
				Symbol: symbol,
				Type:   common.SymbolType_SPOT_NORMAL,
			}
			chMap[info] = make(chan *client.WsBookTickerRsp, 100)
		}
		err := wsClient.BookTickerGroup(ctx, chMap)
		if err != nil {
			t.Fatal(err)
		}
		for _, ch := range chMap {
			go PrintBookTickerGroup(ch)
		}
		//select {}
	}
	for _, symbolList := range symbolLists {
		fmt.Println(symbolList)
		conf := base.WsConf{
			APIConf: base.APIConf{
				// ProxyUrl: proxyUrl,
			},
		}
		wsClient := NewBitfinexSpotWebsocket(conf)
		ctx, _ := context.WithCancel(context.Background())
		chDeltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
		chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
		for _, symbol := range symbolList {
			//for _, symbol := range []string{"BTC/USD", "1INCH/USD", "ETHS/USD", "ETH/USD", "ETHW/USD", "ETHS/USD"} {
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
		//select {}
	}
}
