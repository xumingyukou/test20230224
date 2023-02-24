package spot_api

import (
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"strings"
)

func GetSpotSymbolName(symbol string) string {
	return strings.ToLower(strings.Replace(symbol, "/", "", -1))
}
func GetFutureSymbolName(symbol string) string {
	return strings.ToUpper(strings.Replace(symbol, "/", "", -1))
}

func GetNatsMarket(symbolType common.SymbolType) common.Market {
	switch symbolType {
	case common.SymbolType_SPOT_NORMAL, common.SymbolType_MARGIN_NORMAL, common.SymbolType_MARGIN_ISOLATED:
		return common.Market_SPOT
	case common.SymbolType_FUTURE_THIS_WEEK, common.SymbolType_FUTURE_NEXT_WEEK, common.SymbolType_FUTURE_THIS_MONTH, common.SymbolType_FUTURE_NEXT_MONTH, common.SymbolType_FUTURE_THIS_QUARTER, common.SymbolType_FUTURE_NEXT_QUARTER:
		return common.Market_FUTURE
	case common.SymbolType_FUTURE_COIN_THIS_WEEK, common.SymbolType_FUTURE_COIN_NEXT_WEEK, common.SymbolType_FUTURE_COIN_THIS_MONTH, common.SymbolType_FUTURE_COIN_NEXT_MONTH, common.SymbolType_FUTURE_COIN_THIS_QUARTER, common.SymbolType_FUTURE_COIN_NEXT_QUARTER:
		return common.Market_FUTURE_COIN
	case common.SymbolType_SWAP_COIN_FOREVER:
		return common.Market_SWAP_COIN
	case common.SymbolType_SWAP_FOREVER:
		return common.Market_SWAP
	default:
		return common.Market_INVALID_MARKET

	}
}

func IsUBaseMarket(market common.Market) bool {
	return market == common.Market_FUTURE || market == common.Market_SWAP
}
func IsCBaseMarket(market common.Market) bool {
	return market == common.Market_FUTURE_COIN || market == common.Market_SWAP_COIN
}

func IsUBaseSymbolType(symbolType common.SymbolType) bool {
	return symbolType == common.SymbolType_FUTURE_THIS_WEEK || symbolType == common.SymbolType_FUTURE_NEXT_WEEK || symbolType == common.SymbolType_FUTURE_THIS_MONTH || symbolType == common.SymbolType_FUTURE_NEXT_MONTH || symbolType == common.SymbolType_FUTURE_THIS_QUARTER || symbolType == common.SymbolType_FUTURE_NEXT_QUARTER || symbolType == common.SymbolType_SWAP_FOREVER
}

func IsCBaseSymbolType(symbolType common.SymbolType) bool {
	return symbolType == common.SymbolType_FUTURE_COIN_THIS_WEEK || symbolType == common.SymbolType_FUTURE_COIN_NEXT_WEEK || symbolType == common.SymbolType_FUTURE_COIN_THIS_MONTH || symbolType == common.SymbolType_FUTURE_COIN_NEXT_MONTH || symbolType == common.SymbolType_FUTURE_COIN_THIS_QUARTER || symbolType == common.SymbolType_FUTURE_COIN_NEXT_QUARTER || symbolType == common.SymbolType_SWAP_COIN_FOREVER
}

func GetSymbol(symbol string) string {
	symbol = strings.Replace(symbol, "/", "", 1)
	//symbol = strings.Replace(symbol, "-", "", 1)
	//symbol = strings.Replace(symbol, "_", "", 1)
	return strings.ToUpper(symbol)
}

func ParseSymbolName(symbol string) string {
	symbol = strings.ToUpper(strings.Split(symbol, "_")[0])
	//USD结尾特殊处理
	if strings.HasSuffix(symbol, "USD") { //"BUSD", "TUSD", "USD"
		if symbol == "BNBUSD" || symbol == "DOTUSD" || symbol == "VETUSD" || symbol == "GMTUSD" || symbol == "BATUSD" || symbol == "LPTUSD" || symbol == "GRTUSD" || symbol == "HNTUSD" || symbol == "ANTUSD" || symbol == "KSHIBUSD" || symbol == "ONTUSD" || symbol == "FETUSD" || symbol == "USDTUSD" || symbol == "OXTUSD" || symbol == "BTRSTUSD" || symbol == "DGBUSD" || symbol == "BNTUSD" || symbol == "TUSD" {
			return fmt.Sprint(symbol[:len(symbol)-3], "/", "USD")
		}
		if strings.HasSuffix(symbol, "BUSD") {
			return fmt.Sprint(symbol[:len(symbol)-4], "/", "BUSD")
		} else if strings.HasSuffix(symbol, "TUSD") {
			return fmt.Sprint(symbol[:len(symbol)-4], "/", "TUSD")
		}
		return fmt.Sprint(symbol[:len(symbol)-3], "/", "USD")
	}
	for _, quote := range []string{"UST", "BKRW", "AUD", "DOT", "NGN", "BVND", "XRP", "BTC", "TRX", "UAH", "BIDR", "VAI", "DAI", "DOGE", "GBP", "BRL", "USDP", "USDS", "USDT", "IDRT", "EUR", "ETH", "PAX", "TRY", "BNB", "RUB", "ZAR", "USDC"} {
		if strings.HasSuffix(symbol, quote) {
			baseCoin := symbol[:len(symbol)-len(quote)]
			return fmt.Sprint(baseCoin, "/", quote)
		}
	}
	return symbol
}
