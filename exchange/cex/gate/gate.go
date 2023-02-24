package gate

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/gate/gate_api"
	"clients/logger"
	"github.com/warmplanet/proto/go/client"
	"net/http"
	"strings"
)

type ClientGate struct {
	api *gate_api.ClientGate

	tradeFeeMap    map[string]*client.TradeFeeItem    //key: symbol
	transferFeeMap map[string]*client.TransferFeeItem //key: network+token
	precisionMap   map[string]*client.PrecisionItem   //key: symbol

	optionMap map[string]interface{}
}

func NewClientGate(conf base.APIConf, maps ...interface{}) *ClientGate {
	c := &ClientGate{
		api: gate_api.NewClientGate(conf),
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

func NewClientGate2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ClientGate {
	// 使用自定义http client
	c := &ClientGate{
		api: gate_api.NewClientGate2(conf, cli),
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

func (c *ClientGate) GetSymbols() []string {
	var symbols []string
	SymbolInfoRes, err := c.api.SymbolsInfo(nil)
	if err != nil {
		logger.Logger.Error("get exchange in error:", err)
		return symbols
	}
	for _, info := range *SymbolInfoRes {
		if info.TradeStatus == "tradable" {
			symbols = append(symbols, strings.ReplaceAll(info.Id, "_", "/"))
		}
	}
	return symbols
}
