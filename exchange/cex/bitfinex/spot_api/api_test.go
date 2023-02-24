package spot_api

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

func TestGetExchange(t *testing.T) {
	conf := base.APIConf{}
	a := NewApiClient(conf)
	res, err := a.GetExchange()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestGetCurrency(t *testing.T) {
	conf := base.APIConf{}
	a := NewApiClient(conf)
	res, err := a.GetCurrency()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}
