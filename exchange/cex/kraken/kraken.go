package kraken

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/kraken/kraken_api"
	"clients/logger"
	"net/http"
	"strings"

	"github.com/warmplanet/proto/go/client"
)

type ClientKraken struct {
	api *kraken_api.ClientKraken

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol

	optionMap map[string]interface{}
}

func NewClientKraken(conf base.APIConf, maps ...interface{}) *ClientKraken {
	c := &ClientKraken{
		api: kraken_api.NewClientKraken(conf),
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

func NewClientKraken2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ClientKraken {
	// 使用自定义http client
	c := &ClientKraken{
		api: kraken_api.NewClientKraken2(conf, cli),
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

func (c *ClientKraken) GetSymbols() []string {
	var symbols []string
	exchangeInfoRes, err := c.api.Asset_Pairs_Info(nil)
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, info := range exchangeInfoRes.Result {
		symbols = append(symbols, strings.ReplaceAll(info.Wsname, "XBT", "BTC"))
	}
	return symbols
}
