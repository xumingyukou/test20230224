package future_api

import (
	"clients/transform"
	"fmt"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/common"
)

type ContractType string

const (
	INVALID_CONTRACTTYPE ContractType = "move"
	PERPETUAL            ContractType = "perpetual" // 永续合约
	CURRENT_QUARTER      ContractType = "future"    // 当季交割合约
	NEXT_QUARTER         ContractType = ""          // 次季交割合约
)

func GetFutureTypeFromExchange(type_ ContractType) common.SymbolType {
	switch type_ {
	case PERPETUAL:
		return common.SymbolType_SWAP_FOREVER
	case CURRENT_QUARTER:
		return common.SymbolType_FUTURE_THIS_QUARTER
	case NEXT_QUARTER:
		return common.SymbolType_FUTURE_NEXT_QUARTER
	default:
		return common.SymbolType_INVALID_TYPE
	}
}

func GetFutureTypeFromResp(item *FutureItem) common.SymbolType {
	switch strings.ToLower(item.Type) {
	case "perpetual":
		return common.SymbolType_SWAP_FOREVER
	case "future":
		//FIXME: 修改算法
		if strings.Split(item.Name, "-")[1] == transform.GetThisQuarter(time.Now().UTC(), 5, 2).Format("0102") {
			return common.SymbolType_FUTURE_THIS_QUARTER
		} else {
			return common.SymbolType_FUTURE_NEXT_QUARTER
		}
	case "prediction":
		return common.SymbolType_INVALID_TYPE
	case "move":
		return common.SymbolType_INVALID_TYPE
	default:
		return common.SymbolType_INVALID_TYPE
	}
}

func GetBaseSymbol(symbol string, type_ common.SymbolType) string { //Converts BTC/USD future this quarter to BTC-0930
	//We ignore the Quote Currency on FTX, since it can only be USD, USDT, USD stablecoin values
	base := strings.Split(symbol, "/")[0]
	switch type_ {
	case common.SymbolType_SPOT_NORMAL:
		return symbol
	case common.SymbolType_SWAP_FOREVER:
		return fmt.Sprint(base, "-PERP")
	case common.SymbolType_FUTURE_THIS_QUARTER:
		return fmt.Sprintf("%v-%v", base, transform.GetThisQuarter(time.Now().UTC(), 5, 2).Format("0102"))
	case common.SymbolType_FUTURE_NEXT_QUARTER:
		return fmt.Sprintf("%v-%v", base, transform.GetNextQuarter(time.Now().UTC(), 5, 2).Format("0102"))
	default:
		return symbol //Did not parse into any known type
	}
} //Converts (BTC/USD, future_this_quarter) to (BTC-0930) "FTX Request format"
