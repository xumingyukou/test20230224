package base

import (
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
)

func SymInfoToString(symbolInfo *client.SymbolInfo) string {
	if symbolInfo.Type == common.SymbolType_SPOT_NORMAL {
		return fmt.Sprintf("%s", symbolInfo.Symbol)
	}
	return fmt.Sprintf("%s_%d_%d", symbolInfo.Symbol, symbolInfo.Market.Number(), symbolInfo.Type.Number())
}

func TransToSymbolInner(symbol string, market common.Market, symbolType common.SymbolType) string {
	return fmt.Sprintf("%s_%d_%d", symbol, market.Number(), symbolType.Number())
}

func TransSymbolTypeToString(symbolType common.SymbolType) string {
	switch symbolType {
	case common.SymbolType_FUTURE_COIN_THIS_WEEK, common.SymbolType_FUTURE_THIS_WEEK:
		return "this_week"
	case common.SymbolType_FUTURE_COIN_NEXT_WEEK, common.SymbolType_FUTURE_NEXT_WEEK:
		return "next_week"
	case common.SymbolType_FUTURE_COIN_THIS_MONTH, common.SymbolType_FUTURE_THIS_MONTH:
		return "this_month"
	case common.SymbolType_FUTURE_COIN_NEXT_MONTH, common.SymbolType_FUTURE_NEXT_MONTH:
		return "next_month"
	case common.SymbolType_FUTURE_COIN_THIS_QUARTER, common.SymbolType_FUTURE_THIS_QUARTER:
		return "this_quarter"
	case common.SymbolType_FUTURE_COIN_NEXT_QUARTER, common.SymbolType_FUTURE_NEXT_QUARTER:
		return "next_quarter"
	}
	return ""
}
