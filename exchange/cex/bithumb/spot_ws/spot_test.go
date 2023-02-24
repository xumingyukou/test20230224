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
	ProxyUrl         = "http://127.0.0.1:1087"
	TimeOffset int64 = 30
	conf       base.WsConf
	apiConf    base.APIConf
	// symbols         = []string{"BTC/JPY", "ETH/BTC", "BAT/BTC", "DOGE/JPY"}
	// symbols = []string{"ETH/BTC", "BTC/JPY", "BAT/BTC", "DOGE/JPY", "ETH/JPY"}
	symbols = []string{"ETH/BTC"}
)

func init() {
	config.LoadExchangeConfig("/Users/yingtian/something/go/.conf/exchange.toml")
	apiConf = base.APIConf{
		ReadTimeout: TimeOffset,
		ProxyUrl:    ProxyUrl,
		EndPoint:    "",
		AccessKey:   "",
		Passphrase:  "",
		SecretKey:   "",
		IsTest:      true,
	}

	conf = base.WsConf{
		APIConf: apiConf,
		ChanCap: 1024,
	}
}

func TestTradeGroup(t *testing.T) {
	wsClient := NewBithumbWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	ch := make(chan *client.WsTradeRsp)
	for _, symbol := range symbols {
		chMap[&client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}] = ch
	}

	recv := make([]bool, len(chMap), len(chMap))

	err := wsClient.TradeGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
			for i, symbol := range symbols {
				if symbol == res.Symbol {
					recv[i] = ok
				}
			}
		}
		allRecv := true
		for _, r := range recv {
			allRecv = allRecv && r
		}
		if allRecv {
			break
		}
	}
}

func TestBookTickerGroup(t *testing.T) {
	wsClient := NewBithumbWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	ch := make(chan *client.WsBookTickerRsp)
	for _, symbol := range symbols {
		chMap[&client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}] = ch
	}

	recv := make([]bool, len(chMap), len(chMap))

	err := wsClient.BookTickerGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}
	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
			for i, symbol := range symbols {
				if symbol == res.Symbol {
					recv[i] = ok
				}
			}
		}
		allRecv := true
		for _, r := range recv {
			allRecv = allRecv && r
		}
		if allRecv {
			break
		}
	}
}

func TestDepthIncrementGroup(t *testing.T) {
	wsClient := NewBithumbWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	ch := make(chan *client.WsDepthRsp)
	for _, symbol := range symbols {
		chMap[&client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}] = ch
	}

	err := wsClient.DepthIncrementGroup(ctx, 0, chMap)
	if err != nil {
		t.Fatal(err)
	}

	recv := make([]bool, len(chMap), len(chMap))

	for {
		select {
		case res, ok := <-ch:
			fmt.Println(ok, res)
			for i, symbol := range symbols {
				if symbol == res.Symbol {
					recv[i] = ok
				}
			}
		}
		allRecv := true
		for _, r := range recv {
			allRecv = allRecv && r
		}
		if allRecv {
			break
		}
	}
}

func TestDepthIncrementSnapshotGroup(t *testing.T) {
	wsClient := NewBithumbWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	deltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	fullMap := make(map[*client.SymbolInfo]chan *depth.Depth)

	deltaCh := make(chan *client.WsDepthRsp, 10*len(symbols))
	fullCh := make(chan *depth.Depth, 10*len(symbols))
	isDelta := true
	isFull := true
	symbolInfos := make([]*client.SymbolInfo, 0, len(symbols))
	for _, symbol := range symbols {
		symbolInfos = append(symbolInfos, &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		})
	}
	if isDelta {
		for _, symbol := range symbolInfos {
			deltaMap[symbol] = deltaCh
		}
	}
	if isFull {
		for _, symbol := range symbolInfos {
			fullMap[symbol] = fullCh
		}
	}

	go func() {
		err := wsClient.DepthIncrementSnapshotGroup(ctx, 0, 100, isDelta, isFull, deltaMap, fullMap)
		if err != nil {
			panic(err)
		}
	}()

	fullRcv := make([]bool, len(symbols))
	deltaRcv := make([]bool, len(symbols))

	for {
		select {
		case res, ok := <-deltaCh:
			for i, symbol := range symbols {
				if res.Symbol == symbol {
					deltaRcv[i] = ok
					break
				}
			}
		case res, ok := <-fullCh:
			fmt.Println(res.Symbol, len(res.Asks), len(res.Bids))
			for i, symbol := range symbols {
				if res.Symbol == symbol {
					fullRcv[i] = ok
					break
				}
			}
		}
		fullAllRcv := true
		for _, rcv := range fullRcv {
			fullAllRcv = fullAllRcv && rcv
		}
		deltaAllRcv := true

		for _, rcv := range deltaRcv {
			deltaAllRcv = deltaAllRcv && rcv
		}
		if deltaAllRcv && fullAllRcv {
			break
		}
	}
}
