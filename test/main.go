package main

import (
	// "clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance"
	"fmt"
)

func main() {
	conf := base.APIConf{
		// AccessKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.AccessKey,
		// SecretKey: config.ExchangeConfig.ExchangeList["binance"].ApiKeyConfig.SecretKey,
	}
	spotClient := binance.NewClientBinance(conf)
	// res, err := spotClient.GetBalance()
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(res)
	// }

	res, err := spotClient.GetPrecision()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
