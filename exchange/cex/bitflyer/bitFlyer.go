package bitFlyer

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/bitflyer/bitFlyer_api"
	"clients/logger"
	"net/http"
	"strings"

	"github.com/warmplanet/proto/go/client"
)

type ClientBitFlyer struct {
	api *bitFlyer_api.ClientBitFlyer

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol

	optionMap map[string]interface{}
}

func NewClientBitFlyer(conf base.APIConf, maps ...interface{}) *ClientBitFlyer {
	c := &ClientBitFlyer{
		api: bitFlyer_api.NewClientBitFlyer(conf),
	}

	for _, m := range maps {
		switch t := m.(type) {
		case map[string]*client.TradeFeeItem:
			c.tradeFeeMap = t

		case map[string]*client.TransferFeeItem:
			c.transferFeeMap = t

		case map[string]*client.PrecisionItem:
			c.precisionMap = t

		case map[string]interface{}:
			c.optionMap = t
		}
	}

	return c
}

func NewClientBitFlyer2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ClientBitFlyer {
	// 使用自定义http client
	c := &ClientBitFlyer{
		api: bitFlyer_api.NewClientBitFlyer2(conf, cli),
	}

	for _, m := range maps {
		switch t := m.(type) {
		case map[string]*client.TradeFeeItem:
			c.tradeFeeMap = t

		case map[string]*client.TransferFeeItem:
			c.transferFeeMap = t

		case map[string]*client.PrecisionItem:
			c.precisionMap = t

		case map[string]interface{}:
			c.optionMap = t
		}
	}

	return c
}

func (c *ClientBitFlyer) GetSymbols() []string {
	var symbols []string
	SymbolInfoRes, err := c.api.SymbolsInfo(nil)
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, info := range *SymbolInfoRes {
		if info.MarketType == "Spot" {
			symbols = append(symbols, strings.ReplaceAll(info.ProductCode, "_", "/"))
		}
	}
	return symbols
}
