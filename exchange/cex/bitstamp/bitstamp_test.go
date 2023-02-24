package bitstamp

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

func TestGetSymbols(t *testing.T) {
	var timeOffset int64 = 30
	conf := base.APIConf{
		ReadTimeout: timeOffset,
		//ProxyUrl:    ProxyUrl,
	}
	bitstamp := NewClientBitstamp(conf)
	res := bitstamp.GetSymbols()
	fmt.Println(res)
}
