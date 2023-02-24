package bithumb

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bithumb/spot_api"
	"clients/logger"
	"strings"

	"github.com/warmplanet/proto/go/common"
)

func Exchange2Canonical(symbol string) string {
	return strings.Replace(symbol, "-", "/", 1)
}

func Canonical2Exchange(cs string) string {
	return strings.Replace(cs, "/", "-", 1)
}

type ClientBithumb struct {
	api *spot_api.ApiClient
}

func NewClientBithumb(conf base.APIConf) *ClientBithumb {
	c := &ClientBithumb{
		api: spot_api.NewApiClient(conf),
	}
	return c
}

func (client *ClientBithumb) GetSymbols() []string {
	res, err := client.api.GetSpotConfig()
	var symbols []string
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return symbols
	}

	for _, symbol := range res.Data.SpotConfig {
		symbols = append(symbols, Exchange2Canonical(symbol.Symbol))
	}
	return symbols
}

func (client *ClientBithumb) GetExchange() common.Exchange {
	return common.Exchange_BITHUMB
}
