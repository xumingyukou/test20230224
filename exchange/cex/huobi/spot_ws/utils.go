package spot_ws

import (
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"strconv"
	"strings"
)

const (
	THIS_WEEK_SUB    = "CW"
	NEXT_WEEK_SUB    = "NW"
	THIS_QUARTER_SUB = "CQ"
	NEXT_QUARTER_SUB = "NQ"
)

func GetFutureSub(contractType string) string {
	var sub string
	switch contractType {
	case "this_week":
		sub = THIS_WEEK_SUB
	case "next_week":
		sub = NEXT_WEEK_SUB
	case "quarter":
		sub = THIS_QUARTER_SUB
	case "next_quarter":
		sub = NEXT_QUARTER_SUB
	}
	return sub
}

func GetFutureSub2(symbol *client.SymbolInfo) string {
	var s string
	switch symbol.Market {
	case common.Market_FUTURE:
		s = strings.ReplaceAll(symbol.Symbol, "/", "-") + "-"
	case common.Market_FUTURE_COIN:
		s = strings.ToUpper(strings.Split(symbol.Symbol, "/")[0]) + "_"
	}
	switch symbol.Type {
	case common.SymbolType_FUTURE_COIN_THIS_WEEK, common.SymbolType_FUTURE_THIS_WEEK:
		s = s + THIS_WEEK_SUB
	case common.SymbolType_FUTURE_COIN_NEXT_WEEK, common.SymbolType_FUTURE_NEXT_WEEK:
		s = s + NEXT_WEEK_SUB
	case common.SymbolType_FUTURE_COIN_THIS_QUARTER, common.SymbolType_FUTURE_THIS_QUARTER:
		s = s + THIS_QUARTER_SUB
	case common.SymbolType_FUTURE_COIN_NEXT_QUARTER, common.SymbolType_FUTURE_NEXT_QUARTER:
		s = s + NEXT_QUARTER_SUB
	}
	return s
}

func GetDepthFullSubscription(symbol *client.SymbolInfo, level, id int) Sub {
	var sub Sub
	switch symbol.Market {
	case common.Market_SPOT:
		sub = Sub{
			Req: "market." + GetSymbolName(symbol) + ".mbp." + strconv.Itoa(level),
			ID:  strconv.Itoa(id),
		}
	case common.Market_SWAP, common.Market_SWAP_COIN:
		if level > 150 {
			level = 150 // 合约不支持400的档位
		}
		sub = Sub{
			Sub:      "market." + GetSymbolName(symbol) + ".depth.size_" + strconv.Itoa(level) + ".high_freq",
			ID:       strconv.Itoa(id),
			DataType: "snapshot",
		}
	case common.Market_FUTURE, common.Market_FUTURE_COIN:
		if level > 150 {
			level = 150 // 合约不支持400的档位
		}
		sub = Sub{
			Sub:      "market." + GetFutureSub2(symbol) + ".depth.size_" + strconv.Itoa(level) + ".high_freq",
			ID:       strconv.Itoa(id),
			DataType: "snapshot",
		}
	}
	return sub
}

func GetDepthIncrementSubscription(symbol *client.SymbolInfo, level, id int) Sub {
	var sub Sub
	switch symbol.Market {
	case common.Market_SPOT:
		sub = Sub{
			Sub: "market." + GetSymbolName(symbol) + ".mbp." + strconv.Itoa(level),
			ID:  strconv.Itoa(id),
		}
	case common.Market_SWAP, common.Market_SWAP_COIN:
		if level > 150 {
			level = 150 // 合约不支持400的档位
		}
		sub = Sub{
			Sub:      "market." + GetSymbolName(symbol) + ".depth.size_" + strconv.Itoa(level) + ".high_freq",
			ID:       strconv.Itoa(id),
			DataType: "incremental",
		}
	case common.Market_FUTURE, common.Market_FUTURE_COIN:
		if level > 150 {
			level = 150 // 合约不支持400的档位
		}
		sub = Sub{
			Sub:      "market." + GetFutureSub2(symbol) + ".depth.size_" + strconv.Itoa(level) + ".high_freq",
			ID:       strconv.Itoa(id),
			DataType: "incremental",
		}
	}
	return sub
}

func GetDepthLimitSubscription(symbol *client.SymbolInfo, level, id int) Sub {
	var sub Sub
	switch symbol.Market {
	case common.Market_SPOT:
		sub = Sub{
			Sub: "market." + GetSymbolName(symbol) + ".mbp.refresh." + strconv.Itoa(level),
			ID:  strconv.Itoa(id),
		}
	case common.Market_SWAP, common.Market_SWAP_COIN:
		sub = Sub{
			Sub: "market." + GetSymbolName(symbol) + ".depth.step0",
			ID:  strconv.Itoa(id),
		}
	case common.Market_FUTURE, common.Market_FUTURE_COIN:
		sub = Sub{
			Sub: "market." + GetFutureSub2(symbol) + ".depth.step0",
			ID:  strconv.Itoa(id),
		}
	}
	return sub
}

func GetTickerSubscription(symbol *client.SymbolInfo) string {
	var sub string
	switch symbol.Market {
	case common.Market_SPOT:
		sub = "market." + GetSymbolName(symbol) + ".bbo"
	case common.Market_SWAP, common.Market_SWAP_COIN:
		sub = "market." + GetSymbolName(symbol) + ".bbo"
	case common.Market_FUTURE, common.Market_FUTURE_COIN:
		sub = "market." + GetFutureSub2(symbol) + ".bbo"
	}
	return sub
}

func GetTradeSubscription(symbol *client.SymbolInfo) string {
	var sub string
	switch symbol.Market {
	case common.Market_SPOT:
		sub = "market." + GetSymbolName(symbol) + ".trade.detail"
	case common.Market_SWAP, common.Market_SWAP_COIN:
		sub = "market." + GetSymbolName(symbol) + ".trade.detail"
	case common.Market_FUTURE, common.Market_FUTURE_COIN:
		sub = "market." + GetFutureSub2(symbol) + ".trade.detail"
	}
	return sub
}

func GetSymbolName(symbol *client.SymbolInfo) string {
	var symbolStr string
	switch symbol.Market {
	case common.Market_SPOT:
		symbolStr = strings.Replace(symbol.Symbol, "/", "", 1)
		symbolStr = strings.ToLower(symbolStr)
	case common.Market_SWAP, common.Market_SWAP_COIN:
		symbolStr = strings.Replace(symbol.Symbol, "/", "-", 1)
	}

	return symbolStr
}

func GetFundSubbscription(symbol *client.SymbolInfo) string {
	var sub string
	switch symbol.Market {
	case common.Market_SPOT:
		sub = "public." + GetSymbolName(symbol) + ".funding_rate"
	case common.Market_SWAP, common.Market_SWAP_COIN:
		sub = "public." + GetSymbolName(symbol) + ".funding_rate"
	case common.Market_FUTURE, common.Market_FUTURE_COIN:
		sub = "public." + GetFutureSub2(symbol) + ".funding_rate"
	}
	return sub
}
