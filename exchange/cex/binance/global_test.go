package binance

import (
	"clients/exchange/cex/binance/spot_api"
	"fmt"
	"strings"
	"testing"
)

func TestClientBinance_GetSymbols(t *testing.T) {
	if GlobalInstance == nil {
		GlobalInstance = NewClientBinanceGlobal(spot_api.SPOT_API_BASE_URL, 30)
	}
	res := GlobalInstance.GetSymbols()
	var symbols []string
	for symbol, _ := range *res {
		symbols = append(symbols, symbol)
	}
	fmt.Println(strings.Join(symbols, ","))
}

func TestClientBinance_GetFullDepth(t *testing.T) {
	if GlobalInstance == nil {
		GlobalInstance = NewClientBinanceGlobal(spot_api.SPOT_API_BASE_URL, 30)
	}
	res, _ := GlobalInstance.GetFullDepth("BTCUSDT", 1000)
	fmt.Println(res)
}
