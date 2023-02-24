package bitMex

import (
	"clients/exchange/cex/base"
	bitMex_api "clients/exchange/cex/bitmex/bitMex_api"
	"clients/logger"
	"github.com/warmplanet/proto/go/client"
	"net/http"
	"strings"
)

type ClientBitMex struct {
	api *bitMex_api.ClientBitMex

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol

	optionMap map[string]interface{}
}

func NewClientBitMex(conf base.APIConf, maps ...interface{}) *ClientBitMex {
	c := &ClientBitMex{
		api: bitMex_api.NewClientBitMex(conf),
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

func NewClientBitMex2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ClientBitMex {
	// 使用自定义http client
	c := &ClientBitMex{
		api: bitMex_api.NewClientBitMex2(conf, cli),
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

func (c *ClientBitMex) GetSymbols() []string {
	var symbols []string
	var set = make(map[string]struct{})
	SymbolInfoRes, err := c.api.SymbolsInfo(nil)
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, info := range *SymbolInfoRes {
		if info.Typ != "IFXXXP" {
			continue
		}
		s := info.Underlying + "/" + info.QuoteCurrency
		if _, ok := set[s]; !ok {
			symbols = append(symbols, strings.ReplaceAll(s, "XBT", "BTC"))
			set[s] = struct{}{}
		}
	}
	return symbols
}
