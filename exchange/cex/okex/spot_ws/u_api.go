package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/okex/ok_api"
	"github.com/warmplanet/proto/go/sdk"

	"net/http"
)

type OkUBaseWebsocket struct {
	OkWebsocket
}

func NewBinanceUBaseWebsocket(conf base.WsConf) *OkUBaseWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &OkUBaseWebsocket{}
	d.WsConf = conf

	if conf.EndPoint == "" {
		d.EndPoint = WS_API_BASE_URL
	}
	// 未更改
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 3
	}

	d.symbolMap = sdk.NewCmapI()
	d.handler = NewWebSocketUBaseHandle(d.ChanCap)
	d.apiClient = ok_api.NewClientOkex(conf.APIConf)
	return d
}

func NewBinanceUBaseWebsocket2(conf base.WsConf, client *http.Client) *OkUBaseWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &OkUBaseWebsocket{}
	d.WsConf = conf

	if conf.EndPoint == "" {
		d.EndPoint = WS_API_BASE_URL
	}
	// 未更改
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 3
	}

	d.symbolMap = sdk.NewCmapI()
	d.handler = NewWebSocketUBaseHandle(d.ChanCap)
	d.apiClient = ok_api.NewClientOkex(conf.APIConf)
	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.apiClient = ok_api.NewClientOkex2(base.APIConf{
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		}, client)
	}
	return d
}
