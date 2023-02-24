package binance_us

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance"
	"clients/exchange/cex/binance-us/c_api"
	"clients/exchange/cex/binance-us/spot_api"
	"clients/exchange/cex/binance-us/u_api"
	"net/http"
)

func NewClientBinance(conf base.APIConf, maps ...interface{}) *binance.ClientBinance {
	c := &binance.ClientBinance{
		Api:  spot_api.NewApiClient(conf, maps...),
		UApi: u_api.NewUApiClient(conf, maps...),
		CApi: c_api.NewCApiClient(conf, maps...),
	}
	c.InitConfMap(maps...)
	return c
}

func NewClientBinance2(conf base.APIConf, cli *http.Client, maps ...interface{}) *binance.ClientBinance {
	// 使用自定义http client
	c := &binance.ClientBinance{
		Api:  spot_api.NewApiClient2(conf, cli, maps...),
		UApi: u_api.NewUApiClient2(conf, cli, maps...),
		CApi: c_api.NewCApiClient2(conf, cli, maps...),
	}
	c.InitConfMap(maps...)
	return c
}
