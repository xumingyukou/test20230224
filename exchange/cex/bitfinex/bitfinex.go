package bitfinex

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitfinex/spot_api"
	"clients/logger"
	"strings"
)

type ClientBitfinex struct {
	api *spot_api.ApiClient
}

func NewClientBitfinex(conf base.APIConf) *ClientBitfinex {
	c := &ClientBitfinex{
		api: spot_api.NewApiClient(conf),
	}
	return c
}

func (client *ClientBitfinex) GetSymbols() []string {
	exchangeCurrencys, err := client.api.GetExchange()
	if err != nil {
		logger.Logger.Error("get exchange ", err.Error())
		return nil
	}
	var fullSymbols []string
	for _, exchangeSymbols := range *exchangeCurrencys {
		for _, fullSymbol := range exchangeSymbols {
			fullSymbols = append(fullSymbols, fullSymbol)
		}
	}
	baseCurrencys, err := client.api.GetCurrency()
	if err != nil {
		logger.Logger.Error("get currency", err.Error())
		return nil
	}
	var preSymbols []string
	for _, baseSymbols := range *baseCurrencys {
		for _, baseSymbol := range baseSymbols {
			preSymbols = append(preSymbols, baseSymbol)
		}
	}
	var symbolsList []string
	for _, symbol := range fullSymbols {
		if find := strings.Contains(symbol, ":"); find {
			s := strings.Split(symbol, ":")
			symbolsList = append(symbolsList, s[0]+"/"+s[1])
			continue
		}
		for _, base := range preSymbols {
			if strings.HasPrefix(symbol, base) {
				quote := strings.Replace(symbol[len(base):], ":", "", -1)
				symbolsList = append(symbolsList, base+"/"+quote)
			}
		}
	}
	return symbolsList
}
