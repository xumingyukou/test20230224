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
	proxyUrl = "http://127.0.0.1:7890"
)

func TestTrade(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	for _, symbol := range []string{"BTC/USDT", "BTC/USDC"} {
		info := &client.SymbolInfo{
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
		go PrintBookTicker(ch, "TradeGroup")
	}
	select {}
}

func TestDepthIncrement(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	for _, symbol := range []string{"BTC/USDT", "BTC/USDC"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
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

func TestDepthLimit(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			//ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	for _, symbol := range []string{"BTC/USDT", "BTC/USDC"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
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

func PrintBookTicker[T any](ch chan T, title string) {
	for {
		select {
		case res, ok := <-ch:
			fmt.Sprint(title, ok, res)
		}
	}
}

func TestBookTickerGroup(t *testing.T) {
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}
	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chMap := make(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	for _, symbol := range []string{"BTC/USDT", "BTC/USDC"} {
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

func TestBinanceSpotWebsocket_Account(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl:    proxyUrl,
			ReadTimeout: 3000,
			AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
		},
	}
	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	data := &client.WsAccountReq{
		Exchange: common.Exchange_BINANCE,
		Market:   common.Market_SPOT,
		Type:     common.SymbolType_SPOT_NORMAL,
	}
	ch, err := wsClient.Account(ctx, data)
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

func TestBinanceSpotWebsocket_Balance(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl:    proxyUrl,
			ReadTimeout: 3000,
			AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
		},
	}
	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	data := &client.WsAccountReq{
		Exchange: common.Exchange_BINANCE,
		Market:   common.Market_SPOT,
		Type:     common.SymbolType_SPOT_NORMAL,
	}
	ch, err := wsClient.Balance(ctx, data)
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

func TestBinanceSpotWebsocket_Order(t *testing.T) {
	config.LoadExchangeConfig("./conf/exchange.toml")
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl:    proxyUrl,
			ReadTimeout: 3000,
			AccessKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
		},
	}
	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	data := &client.WsAccountReq{
		Exchange: common.Exchange_BINANCE,
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
		break
	}
}

func TestDepthIncrementGroupSnapshot(t *testing.T) {
	//wss://stream.binance.com:9443/ws/btcusdt@trade
	conf := base.WsConf{
		APIConf: base.APIConf{
			ProxyUrl: proxyUrl,
		},
	}

	wsClient := NewBinanceWebsocket(conf)
	ctx, _ := context.WithCancel(context.Background())
	chDeltaMap := make(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	chFullMap := make(map[*client.SymbolInfo]chan *depth.Depth)
	//for _, symbol := range []string{"YGG/BNB", "ATA/BTC", "CVX/BTC", "TRU/RUB", "BUSD/VAI"} { //{"YGG/BNB", "ATA/BTC", "CVX/BTC", "TRU/RUB", "ADA/RUB", "MATIC/RUB", "DIA/BUSD", "YFI/EUR", "PAXG/BNB", "BNX/BTC", "C98/BNB", "FTM/RUB", "ETH/USDP", "NEO/RUB", "BAND/BTC", "STMX/ETH", "TRXDOWN/USDT", "WRX/BNB", "CFX/BUSD", "ARK/BUSD", "GLM/BTC", "KNC/BNB", "FOR/USDT", "CLV/BTC", "BTG/BUSD", "BUSD/VAI", "BIFI/USDT", "ENJ/EUR", "XEM/BTC", "LRC/ETH", "SKL/BTC", "ORN/BUSD", "XRP/TUSD", "DNT/BUSD", "FIS/BTC", "DNT/USDT", "WOO/BUSD", "REN/BTC", "QKC/ETH", "FIRO/BTC", "AAVE/BNB", "CHESS/BTC", "ALPINE/EUR", "AKRO/BUSD", "DYDX/BTC", "BTC/TUSD", "CVC/BTC", "IOTA/ETH", "ATA/BNB", "BNB/USDP", "SNX/ETH", "STMX/BTC", "POLY/BUSD", "POLS/BNB", "MULTI/BUSD", "AUD/USDC", "ZEC/BNB", "ATOM/EUR", "FUN/USDT", "FOR/BTC", "GTC/BUSD", "APE/GBP", "BTS/BTC", "DATA/ETH", "REQ/BTC", "HARD/BUSD", "ENJ/GBP", "FLM/BTC", "PEOPLE/BNB", "AR/BNB", "LINA/BTC", "EGLD/BNB", "GAL/BNB", "ZIL/EUR", "MASK/BUSD", "BNT/BTC", "STX/BNB", "FRONT/BUSD", "ORN/BTC", "SOL/BRL", "MANA/BRL", "DREP/BUSD", "RIF/BTC", "ADA/GBP", "MLN/BUSD", "MKR/BTC", "ALGO/BNB", "QNT/BNB", "GALA/EUR", "WAXP/BNB", "BETA/ETH", "ACA/BTC", "DOGE/BRL", "CRV/ETH", "LTO/BUSD", "GLM/ETH", "TOMO/BUSD", "INJ/BNB", "MULTI/BTC", "XVS/BTC", "BAKE/BNB", "MDT/BTC", "DOGE/GBP", "JST/BUSD", "COS/BTC", "TRU/BTC", "COS/BNB", "LINK/BRL", "NEO/TRY", "ANT/BTC", "RNDR/BTC", "JASMY/ETH", "ENS/TRY", "LTC/UAH", "XRP/BRL", "MLN/BTC", "XVS/BUSD", "ICX/BUSD", "UNFI/BTC", "JASMY/BNB", "DEXE/ETH", "AUTO/BTC", "WBTC/BUSD", "REN/BUSD", "BNB/DAI", "BOND/BTC", "ALCX/USDT", "LSK/BUSD", "XRP/RUB", "ALPHA/BTC", "AXS/BIDR", "ALICE/BIDR", "XLM/TRY", "FIL/ETH", "PEOPLE/ETH", "AAVE/ETH", "XVS/BNB", "LPT/BNB", "IMX/BNB", "DOCK/BUSD", "MTL/ETH", "BSW/ETH", "MBL/BUSD", "MOB/BUSD", "NEO/BNB", "FIO/BTC", "MANA/ETH", "GAL/BTC", "ACA/BUSD", "HOT/EUR", "SOL/RUB", "POWR/BUSD", "ARPA/BTC", "VIDT/BTC", "XLM/EUR", "JASMY/EUR", "XMR/BNB", "COCOS/BUSD", "SAND/ETH", "TWT/BTC", "ICX/BTC", "ANC/BNB", "ZEN/BUSD", "CELR/ETH", "PLA/BUSD", "SCRT/BUSD", "MC/BNB", "CHZ/GBP", "CLV/BNB", "MANA/BNB", "XTZ/ETH", "REP/BTC", "BSW/BNB", "NKN/BUSD", "VET/GBP", "IRIS/USDT", "BNT/ETH", "DREP/BTC", "UTK/USDT", "APE/BRL", "BETA/BNB", "USDP/BUSD", "OOKI/BNB", "ALCX/BTC", "STRAX/BTC", "FET/BNB", "GAL/EUR", "IOTA/BNB", "UTK/BUSD", "SNX/BNB", "CELO/BTC", "T/BUSD", "REP/USDT"} {
	for _, symbol := range []string{"BTC/USD", "ATA/BTC", "CVX/BTC", "TRU/RUB", "BUSD/VAI"} { //{"YGG/BNB", "ATA/BTC", "CVX/BTC", "TRU/RUB", "ADA/RUB", "MATIC/RUB", "DIA/BUSD", "YFI/EUR", "PAXG/BNB", "BNX/BTC", "C98/BNB", "FTM/RUB", "ETH/USDP", "NEO/RUB", "BAND/BTC", "STMX/ETH", "TRXDOWN/USDT", "WRX/BNB", "CFX/BUSD", "ARK/BUSD", "GLM/BTC", "KNC/BNB", "FOR/USDT", "CLV/BTC", "BTG/BUSD", "BUSD/VAI", "BIFI/USDT", "ENJ/EUR", "XEM/BTC", "LRC/ETH", "SKL/BTC", "ORN/BUSD", "XRP/TUSD", "DNT/BUSD", "FIS/BTC", "DNT/USDT", "WOO/BUSD", "REN/BTC", "QKC/ETH", "FIRO/BTC", "AAVE/BNB", "CHESS/BTC", "ALPINE/EUR", "AKRO/BUSD", "DYDX/BTC", "BTC/TUSD", "CVC/BTC", "IOTA/ETH", "ATA/BNB", "BNB/USDP", "SNX/ETH", "STMX/BTC", "POLY/BUSD", "POLS/BNB", "MULTI/BUSD", "AUD/USDC", "ZEC/BNB", "ATOM/EUR", "FUN/USDT", "FOR/BTC", "GTC/BUSD", "APE/GBP", "BTS/BTC", "DATA/ETH", "REQ/BTC", "HARD/BUSD", "ENJ/GBP", "FLM/BTC", "PEOPLE/BNB", "AR/BNB", "LINA/BTC", "EGLD/BNB", "GAL/BNB", "ZIL/EUR", "MASK/BUSD", "BNT/BTC", "STX/BNB", "FRONT/BUSD", "ORN/BTC", "SOL/BRL", "MANA/BRL", "DREP/BUSD", "RIF/BTC", "ADA/GBP", "MLN/BUSD", "MKR/BTC", "ALGO/BNB", "QNT/BNB", "GALA/EUR", "WAXP/BNB", "BETA/ETH", "ACA/BTC", "DOGE/BRL", "CRV/ETH", "LTO/BUSD", "GLM/ETH", "TOMO/BUSD", "INJ/BNB", "MULTI/BTC", "XVS/BTC", "BAKE/BNB", "MDT/BTC", "DOGE/GBP", "JST/BUSD", "COS/BTC", "TRU/BTC", "COS/BNB", "LINK/BRL", "NEO/TRY", "ANT/BTC", "RNDR/BTC", "JASMY/ETH", "ENS/TRY", "LTC/UAH", "XRP/BRL", "MLN/BTC", "XVS/BUSD", "ICX/BUSD", "UNFI/BTC", "JASMY/BNB", "DEXE/ETH", "AUTO/BTC", "WBTC/BUSD", "REN/BUSD", "BNB/DAI", "BOND/BTC", "ALCX/USDT", "LSK/BUSD", "XRP/RUB", "ALPHA/BTC", "AXS/BIDR", "ALICE/BIDR", "XLM/TRY", "FIL/ETH", "PEOPLE/ETH", "AAVE/ETH", "XVS/BNB", "LPT/BNB", "IMX/BNB", "DOCK/BUSD", "MTL/ETH", "BSW/ETH", "MBL/BUSD", "MOB/BUSD", "NEO/BNB", "FIO/BTC", "MANA/ETH", "GAL/BTC", "ACA/BUSD", "HOT/EUR", "SOL/RUB", "POWR/BUSD", "ARPA/BTC", "VIDT/BTC", "XLM/EUR", "JASMY/EUR", "XMR/BNB", "COCOS/BUSD", "SAND/ETH", "TWT/BTC", "ICX/BTC", "ANC/BNB", "ZEN/BUSD", "CELR/ETH", "PLA/BUSD", "SCRT/BUSD", "MC/BNB", "CHZ/GBP", "CLV/BNB", "MANA/BNB", "XTZ/ETH", "REP/BTC", "BSW/BNB", "NKN/BUSD", "VET/GBP", "IRIS/USDT", "BNT/ETH", "DREP/BTC", "UTK/USDT", "APE/BRL", "BETA/BNB", "USDP/BUSD", "OOKI/BNB", "ALCX/BTC", "STRAX/BTC", "FET/BNB", "GAL/EUR", "IOTA/BNB", "UTK/BUSD", "SNX/BNB", "CELO/BTC", "T/BUSD", "REP/USDT"} {
		info := &client.SymbolInfo{
			Symbol: symbol,
			Type:   common.SymbolType_SPOT_NORMAL,
		}
		chDeltaMap[info] = make(chan *client.WsDepthRsp, 10)
		chFullMap[info] = make(chan *depth.Depth, 10)
		go PrintBookTicker(chDeltaMap[info], "delta")
		go PrintBookTicker(chFullMap[info], "full")
	}
	err := wsClient.DepthIncrementSnapshotGroup(ctx, 1000, 20, true, true, chDeltaMap, chFullMap)
	if err != nil {
		t.Fatal(err)
	}
	select {}
}
