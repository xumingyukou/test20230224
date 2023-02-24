package coinbase

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/coinbase/spot_api"
	"clients/logger"
	"github.com/warmplanet/proto/go/client"
)

type ClientCoinbase struct {
	api *spot_api.ApiClient

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol
}

func NewClientCoinbase(conf base.APIConf, maps ...interface{}) *ClientCoinbase {
	c := &ClientCoinbase{
		api: spot_api.NewApiClient(conf),
	}

	for _, m := range maps {
		switch t := m.(type) {
		case map[string]*client.TradeFeeItem:
			c.tradeFeeMap = t

		case map[string]*client.TransferFeeItem:
			c.transferFeeMap = t

		case map[string]*client.PrecisionItem:
			c.precisionMap = t
		}
	}
	return c
}

func (c *ClientCoinbase) GetSymbols() []string {
	var symbols []string
	exchangeInfoRes, err := c.api.GetProducts()
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, symbol := range exchangeInfoRes.ProductInfos {
		if symbol.TradingDisabled == false {
			symbols = append(symbols, symbol.DisplayName)
		}
	}
	return symbols
}
