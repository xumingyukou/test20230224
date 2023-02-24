package spot_api

import (
	"clients/exchange/cex/base"
	"fmt"
	"testing"
)

func TestGetSymbol(t *testing.T) {
	conf := base.APIConf{}
	a := NewApiClient(conf)
	res, err := a.GetSymbols()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res)
}

func TestSymbolDetails(t *testing.T) {
	conf := base.APIConf{}
	a := NewApiClient(conf)
	symbols, err := a.GetSymbols()
	if err != nil {
		t.Fatal(err)
	}
	for _, symbol := range *symbols {
		res, err := a.GetSymbolDetails(symbol)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(res.Symbol, res.BaseCurrency, res.QuoteCurrency)
	}
}
