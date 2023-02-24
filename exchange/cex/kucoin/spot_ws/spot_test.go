package spot_ws

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/kucoin/spot_api"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

var (
	ProxyUrl         = "http://127.0.0.1:1087"
	TimeOffset int64 = 30
	conf       base.WsConf
	apiConf    base.APIConf
)

func init() {
	// config.LoadExchangeConfig("./conf/exchange.toml")
	config.LoadExchangeConfig("/Users/yingtian/something/go/.conf/exchange.toml")
	apiConf := base.APIConf{
		ReadTimeout: TimeOffset,
		ProxyUrl:    ProxyUrl,
		EndPoint:    spot_api.SPOT_API_BASE_URL,
		AccessKey:   config.ExchangeConfig.ExchangeList["kucoin"].ApiKeyConfig.AccessKey,
		Passphrase:  config.ExchangeConfig.ExchangeList["kucoin"].ApiKeyConfig.Passphrase,
		SecretKey:   config.ExchangeConfig.ExchangeList["kucoin"].ApiKeyConfig.SecretKey,
		IsTest:      true,
	}

	// apiConf := base.APIConf{
	// 	ReadTimeout: TimeOffset,
	// 	ProxyUrl:    ProxyUrl,
	// 	EndPoint:    spot_api.SANDBOX_BASE_URL,
	// 	AccessKey:   config.ExchangeConfig.ExchangeList["kucoin_sandbox"].ApiKeyConfig.AccessKey,
	// 	Passphrase:  config.ExchangeConfig.ExchangeList["kucoin_sandbox"].ApiKeyConfig.Passphrase,
	// 	SecretKey:   config.ExchangeConfig.ExchangeList["kucoin_sandbox"].ApiKeyConfig.SecretKey,
	// 	IsTest:      true,
	// }

	conf = base.WsConf{
		APIConf: apiConf,
		ChanCap: 1024,
	}
}

func TestNewKucoinWebsocket2(t *testing.T) {
	proxyUrl, _ := url.Parse("http://127.0.0.1:1087")

	transport := http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	httpClient := &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}

	wsClient := NewKucoinWebsocket2(conf, httpClient)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	ch := make(chan *client.WsTradeRsp)
	chMap[&client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch

	err := wsClient.TradeGroup(ctx, chMap)
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

func TestTradeGroup(t *testing.T) {
	wsClient := NewKucoinWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	ch := make(chan *client.WsTradeRsp)
	chMap[&client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch

	err := wsClient.TradeGroup(ctx, chMap)
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

func TestBookTickerGroup(t *testing.T) {
	wsClient := NewKucoinWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	ch1 := make(chan *client.WsBookTickerRsp)
	chMap[&client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch1

	ch2 := make(chan *client.WsBookTickerRsp)
	chMap[&client.SymbolInfo{
		Symbol: "ETH/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch2

	err := wsClient.BookTickerGroup(ctx, chMap)
	if err != nil {
		t.Fatal(err)
	}

	rcv1 := false
	rcv2 := false

	for {
		select {
		case res, ok := <-ch1:
			fmt.Println(ok, res)
			rcv1 = ok
		case res, ok := <-ch2:
			fmt.Println(ok, res)
			rcv2 = ok
		}
		if rcv1 && rcv2 {
			break
		}
	}
}

func TestDepthLimitGroup(t *testing.T) {
	wsClient := NewKucoinWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	ch1 := make(chan *depth.Depth)
	chMap[&client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch1

	ch2 := make(chan *depth.Depth)
	chMap[&client.SymbolInfo{
		Symbol: "ETH/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch2

	err := wsClient.DepthLimitGroup(ctx, 0, 0, chMap)
	if err != nil {
		t.Fatal(err)
	}

	rcv1 := false
	rcv2 := false

	for {
		select {
		case res, ok := <-ch1:
			fmt.Println(ok, res)
			rcv1 = ok
		case res, ok := <-ch2:
			fmt.Println(ok, res)
			rcv2 = ok
		}
		if rcv1 && rcv2 {
			break
		}
	}
}

func TestDepthIncrementGroup(t *testing.T) {
	wsClient := NewKucoinWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	ch1 := make(chan *client.WsDepthRsp)
	chMap[&client.SymbolInfo{
		Symbol: "BTC/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch1

	ch2 := make(chan *client.WsDepthRsp)
	chMap[&client.SymbolInfo{
		Symbol: "ETH/USDT",
		Type:   common.SymbolType_SPOT_NORMAL,
	}] = ch2

	err := wsClient.DepthIncrementGroup(ctx, 0, chMap)
	if err != nil {
		t.Fatal(err)
	}

	rcv1 := false
	rcv2 := false

	for {
		select {
		case res, ok := <-ch1:
			fmt.Println(ok, res)
			rcv1 = ok
		case res, ok := <-ch2:
			fmt.Println(ok, res)
			rcv2 = ok
		}
		if rcv1 && rcv2 {
			break
		}
	}
}

func TestDepthIncrementSnapshotGroup(t *testing.T) {
	wsClient := NewKucoinWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	deltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	fullMap := make(map[*client.SymbolInfo]chan *depth.Depth)

	symbols := []string{"BTC/USDT", "ETH/USDT", "ETH/BTC", "BTC/USDC"}
	symbolInfos := make([]*client.SymbolInfo, 0, len(symbols))
	for _, symbol := range symbols {
		symbolInfos = append(symbolInfos, &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		})
	}

	deltaCh := make(chan *client.WsDepthRsp, 10*len(symbols))
	fullCh := make(chan *depth.Depth, 10*len(symbols))
	isDelta := true
	isFull := true
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

	err := wsClient.DepthIncrementSnapshotGroup(ctx, 0, 100, isDelta, isFull, deltaMap, fullMap)
	if err != nil {
		t.Fatal(err)
	}

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
