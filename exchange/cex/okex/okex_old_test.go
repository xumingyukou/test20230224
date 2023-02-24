package okex

import (
	"fmt"
	"strings"
	"testing"
)

var (
	c = NewClientOkexOld(GLOBAL_API_BASE_URL, 2)
)

func TestClientOkex_GetSymbols(t *testing.T) {
	res := c.GetSymbols()
	var symbols []string
	for symbol, _ := range *res {
		symbols = append(symbols, symbol)
	}
	fmt.Println(len(*c.GetTokens()), c.GetTokens())
	fmt.Println(len(*c.GetQuoteCoins()), c.GetQuoteCoins())
	fmt.Println(strings.Join(symbols, ","))
}

func TestBalance(t *testing.T) {
	c.AccessKey = ""
	c.SecretKey = ""
	a, err := c.AccountBalance()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range a.Data {
		fmt.Printf("%#v\n", i)
		for _, j := range i.Details {
			fmt.Printf("%#v\n", j)
		}
	}
}

func TestTradeFee(t *testing.T) {
	a, err := c.AccountTradeFee()
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range a.Data {
		fmt.Printf("%#v\n", i)
	}
}
