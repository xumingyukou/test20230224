package gemini

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/gemini/spot_api"
	"clients/logger"
	"fmt"
)

type ClientGemini struct {
	api *spot_api.ApiClient
}

func NewClientGemini(conf base.APIConf) *ClientGemini {
	c := &ClientGemini{
		api: spot_api.NewApiClient(conf),
	}
	return c
}

func (client *ClientGemini) GetSymbols() []string {
	symbols, err := client.api.GetSymbols()
	if err != nil {
		logger.Logger.Error("get symbols ", err.Error())
		return []string{}
	}

	var symbollist []string

	for _, symbol := range *symbols {
		res, err := client.api.GetSymbolDetails(symbol)
		if err != nil {
			logger.Logger.Error("get symbols ", err.Error())
		}
		symbollist = append(symbollist, res.BaseCurrency+"/"+res.QuoteCurrency)
		fmt.Println(symbol)
	}
	return symbollist
}
