package spot_ws

import (
	"clients/exchange/cex/okex"
	"clients/exchange/cex/okex/ok_api"
	"clients/logger"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"strconv"
	"strings"
)

func ParseFloat(f float64) string {
	res := strconv.FormatFloat(f, 'f', -1, 64)
	return res
}

func GetInstId(symbols ...*client.SymbolInfo) string {
	var type_ string
	symbol := okex.GetSymbol(symbols[0].Symbol)
	if symbols[0].Type != common.SymbolType_SPOT_NORMAL {
		tmp := ok_api.GetFutureTypeFromNats(symbols[0].Type)
		type_ = ok_api.GetUBaseSymbol(symbol, tmp)
		symbol = strings.ReplaceAll(type_, "_", "-")
		//fmt.Println(type_)
	}
	return symbol
}

func GetOriginInstId(instId, originDate string) string {
	instIdList := strings.Split(instId, "-")
	if len(instIdList) < 3 {
		return ""
	}
	nowDate := instIdList[2]
	return strings.Replace(instId, nowDate, originDate, 1)
}

// transformSymbol
// @Description: 将返回的参数转换为统一的用斜杠划分的格式
// @param s
// @return string
func transSymbolToSymbolInfo(s string) *client.SymbolInfo {
	symbol, market, symbolType := ok_api.GetContractType(s)
	return &client.SymbolInfo{
		Symbol: symbol,
		Market: market,
		Type:   symbolType,
	}
}

func transformSymbol(s string) string {
	return strings.Replace(s, "-", "/", 1)
}

func TransContractSize(symbolInfoStr string, market common.Market, tradeRsp *client.WsTradeRsp, levels ...*depth.DepthLevel) {
	var transFunc func(price, amount, contractSize float64) float64
	switch market {
	case common.Market_FUTURE_COIN, common.Market_SWAP_COIN:
		transFunc = func(price, amount, contractSize float64) float64 {
			if amount == 0 {
				return amount
			}
			return amount * contractSize / price
		}
	case common.Market_FUTURE, common.Market_SWAP:
		transFunc = func(price, amount, contractSize float64) float64 {
			return amount * contractSize
		}
	case common.Market_SPOT:
		fallthrough
	default:
		return
	}
	contractSizeIf, ok := contractSizeMap.Load(symbolInfoStr)
	if !ok {
		logger.Logger.Error("contractSizeMap key error: ", symbolInfoStr)
		return
	}
	contractSize, okk := contractSizeIf.(float64)
	if !okk {
		logger.Logger.Error("contractSize assert error: ", symbolInfoStr)
		return
	}

	// depth和ticker数据
	for _, level := range levels {
		level.Amount = transFunc(level.Price, level.Amount, contractSize)
	}

	// trade数据
	if tradeRsp != nil {
		tradeRsp.Amount = transFunc(tradeRsp.Price, tradeRsp.Amount, contractSize)
	}
}
