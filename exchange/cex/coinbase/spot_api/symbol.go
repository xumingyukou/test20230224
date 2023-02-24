package spot_api

import (
	"strings"
)

func GetSpotSymbolName(symbol string) string {
	return strings.ToLower(strings.Replace(strings.Replace(symbol, "/", "", -1), "-", "", -1))
}
