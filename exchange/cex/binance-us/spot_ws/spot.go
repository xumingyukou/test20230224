package spot_ws

import (
	"clients/exchange/cex/base"
	usCBaseApi "clients/exchange/cex/binance-us/c_api"
	usSpotApi "clients/exchange/cex/binance-us/spot_api"
	usUBaseApi "clients/exchange/cex/binance-us/u_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/spot_ws"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
)

var (
	g          *spot_ws.GlobalDepthWight
	globalOnce sync.Once
)

func init() {
	rand.Seed(time.Now().Unix())
}

func NewBinanceWebsocket(conf base.WsConf) *spot_ws.BinanceSpotWebsocket {
	if conf.EndPoint == "" {
		conf.EndPoint = spot_api.WS_API_US_URL
	}
	d := &spot_ws.BinanceSpotWebsocket{
		WsConf: conf,
	}
	apiConf := base.APIConf{
		ProxyUrl:    d.ProxyUrl,
		ReadTimeout: d.ReadTimeout,
		AccessKey:   d.AccessKey,
		SecretKey:   d.SecretKey,
	}
	d.GlobalDW = g
	d.Init(usSpotApi.NewApiClient(apiConf), usUBaseApi.NewUApiClient(apiConf), usCBaseApi.NewCApiClient(apiConf), spot_ws.NewWebSocketSpotHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	globalOnce.Do(func() {
		if g == nil {
			g = new(spot_ws.GlobalDepthWight)
			d.GlobalDW = g
		}
		go func() {
			timer := time.NewTimer(time.Duration(d.ApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
			for {
				select {
				case <-timer.C:
					d.GlobalDW.GlobalLock.Lock()
					d.GlobalDW.GlobalWeight = 0
					d.GlobalDW.GlobalLock.Unlock()
					timer.Reset(time.Duration(d.ApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
				}
			}
		}()
	})
	return d
}

func NewBinanceWebsocket2(conf base.WsConf, cli *http.Client) *spot_ws.BinanceSpotWebsocket {
	d := &spot_ws.BinanceSpotWebsocket{
		WsConf: conf,
	}
	// ws socket目前固定写死
	d.EndPoint = spot_api.WS_API_US_URL
	apiConf := base.APIConf{
		ProxyUrl:    d.ProxyUrl,
		ReadTimeout: d.ReadTimeout,
		AccessKey:   d.AccessKey,
		SecretKey:   d.SecretKey,
	}
	d.Init(usSpotApi.NewApiClient2(apiConf, cli), usUBaseApi.NewUApiClient2(apiConf, cli), usCBaseApi.NewCApiClient2(apiConf, cli), spot_ws.NewWebSocketSpotHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	return d
}
