package c_api

import (
	"github.com/warmplanet/proto/go/common"
)

func TransStrToSymbolType(contractType string) common.SymbolType {
	var res common.SymbolType
	switch contractType {
	case "this_week":
		res = common.SymbolType_FUTURE_COIN_THIS_WEEK
	case "next_week":
		res = common.SymbolType_FUTURE_COIN_NEXT_WEEK
	case "quarter":
		res = common.SymbolType_FUTURE_COIN_THIS_QUARTER
	case "next_quarter":
		res = common.SymbolType_FUTURE_COIN_NEXT_QUARTER
	}
	return res
}
