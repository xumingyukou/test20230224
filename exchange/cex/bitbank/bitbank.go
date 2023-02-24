package bitbank

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitbank/spot_api"
	"clients/logger"
	"strings"
)

func Canonical2Exchange(symbol string) string {
	return strings.Replace(strings.ToLower(symbol), "/", "_", 1)
}

func Exchange2Canonical(symbol string) string {
	return strings.Replace(strings.ToUpper(symbol), "_", "/", 1)
}

type ClientBitbank struct {
	api *spot_api.ApiClient
}

func NewClientBitbank(conf base.APIConf) *ClientBitbank {
	c := &ClientBitbank{
		api: spot_api.NewApiClient(conf),
	}
	return c
}

func (client *ClientBitbank) GetSymbols() []string {
	res, err := client.api.GetSymbols()
	var symbols []string
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return symbols
	}

	for _, symbol := range res {
		symbols = append(symbols, Exchange2Canonical(symbol))
	}
	return symbols
}
