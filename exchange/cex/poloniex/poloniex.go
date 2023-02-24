package poloniex

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/poloniex/spot_api"
	"clients/logger"
	"strings"

	"github.com/warmplanet/proto/go/common"
)

func Exchange2Canonical(symbol string) string {
	return strings.Replace(symbol, "_", "/", 1)
}

func Canonical2Exchange(cs string) string {
	return strings.Replace(cs, "/", "_", 1)
}

type ClientPoloneix struct {
	api *spot_api.ApiClient
}

func NewClientPoloniex(conf base.APIConf) *ClientPoloneix {
	c := &ClientPoloneix{
		api: spot_api.NewApiClient(conf),
	}
	return c
}

func (client *ClientPoloneix) GetSymbols() []string {
	res, err := client.api.GetMarkets()
	var symbols []string
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return symbols
	}

	for _, symbol := range *res {
		symbols = append(symbols, Exchange2Canonical(symbol.Symbol))
	}
	return symbols
}

func (client *ClientPoloneix) GetExchange() common.Exchange {
	return common.Exchange_POLONIEX
}
