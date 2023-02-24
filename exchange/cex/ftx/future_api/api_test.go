package future_api

import (
	"clients/exchange/cex/base"
	"encoding/json"
	"fmt"
	"github.com/warmplanet/proto/go/common"
	"strings"
	"testing"
)

func TestGetFuture(t *testing.T) {
	conf := base.APIConf{}
	a := NewUApiClient(conf)
	res, err := a.ListAllFutures()
	if err != nil {
		t.Fatal(err)
	}
	resSwap := make([]string, 0)
	resFuture := make([]string, 0)
	for _, i := range res.Result {
		if (i.Type == "perpetual" || i.Type == "future") && i.Enabled {
			symbolType := GetFutureTypeFromResp(i)
			var market common.Market
			if strings.Split(strings.ToLower(i.Name), "-")[1] == "perp" {
				market = common.Market_SWAP
				fmt.Println(i.Name, strings.Split(i.Name, "-")[0], market, symbolType, i.Type)
				resSwap = append(resSwap, fmt.Sprintf("%s_%d_%d", strings.Split(i.Name, "-")[0], market.Number(), symbolType.Number()))
			} else {
				market = common.Market_FUTURE
				fmt.Println(i.Name, strings.Split(i.Name, "-")[0], market, symbolType, i.Type)
				resFuture = append(resFuture, fmt.Sprintf("%s_%d_%d", strings.Split(i.Name, "-")[0], market.Number(), symbolType.Number()))
			}

		}
	}
	fmt.Println(len(resSwap))
	b, _ := json.Marshal(resSwap)
	fmt.Println(string(b))

	fmt.Println(len(resFuture))
	c, _ := json.Marshal(resFuture)
	fmt.Println(string(c))
}
