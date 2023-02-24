package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/c_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/spot_ws"
	"clients/exchange/cex/binance/u_api"
	"net/http"
)

func NewBinanceCBaseWebsocket(conf base.WsConf) *spot_ws.BinanceCBaseWebsocket {
	if conf.EndPoint == "" {
		conf.EndPoint = spot_api.WS_CAPI_US_URL
	}
	d := &spot_ws.BinanceCBaseWebsocket{}
	d.WsConf = conf
	apiConf := base.APIConf{
		ProxyUrl:    conf.ProxyUrl,
		ReadTimeout: conf.ReadTimeout,
		AccessKey:   conf.AccessKey,
		SecretKey:   conf.SecretKey,
	}
	d.Init(spot_api.NewApiClient(apiConf), u_api.NewUApiClient(apiConf), c_api.NewCApiClient(apiConf), spot_ws.NewWebSocketCBaseHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	return d
}

func NewBinanceCBaseWebsocket2(conf base.WsConf, client *http.Client) *spot_ws.BinanceCBaseWebsocket {
	conf.EndPoint = spot_api.WS_CAPI_US_URL
	d := &spot_ws.BinanceCBaseWebsocket{}
	d.WsConf = conf
	apiConf := base.APIConf{
		ProxyUrl:    conf.ProxyUrl,
		ReadTimeout: conf.ReadTimeout,
		AccessKey:   conf.AccessKey,
		SecretKey:   conf.SecretKey,
	}
	d.Init(spot_api.NewApiClient2(apiConf, client), u_api.NewUApiClient2(apiConf, client), c_api.NewCApiClient2(apiConf, client), spot_ws.NewWebSocketCBaseHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	return d
}
