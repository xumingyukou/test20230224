package u_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	usSpotApi "clients/exchange/cex/binance-us/spot_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"github.com/warmplanet/proto/go/client"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	globalOnce sync.Once
)

func NewUApiClient(conf base.APIConf, maps ...interface{}) *u_api.UApiClient {
	var (
		a             = &u_api.UApiClient{}
		proxyUrl      *url.URL
		transport     = http.Transport{}
		err           error
		weightInfoMap map[string]int64
	)
	a.APIConf = conf
	if conf.EndPoint == "" {
		a.EndPoint = spot_api.UBASE_API_BASE_URL
	}
	// 用户可以自定义限速规则
	for _, m := range maps {
		switch t := m.(type) {
		case map[client.WeightType]*client.WeightInfo:
			a.WeightInfo = t
			// 将传入的limit和interval归一化到number/minute
			for _, v := range a.WeightInfo {
				v.Limit = v.Limit * 60 / v.IntervalSec
				v.IntervalSec = 60
			}
		case config.ExchangeWeightInfo:
			weightInfoMap = t.UBase
		}
	}
	if conf.ProxyUrl != "" {
		proxyUrl, err = url.Parse(conf.ProxyUrl)
		if err != nil {
			logger.Logger.Error("can not set proxy:", conf.ProxyUrl)
		}
		transport = http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	globalOnce.Do(func() {
		if usSpotApi.BinanceUBaseWeight == nil {
			usSpotApi.BinanceUBaseWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = usSpotApi.BinanceUBaseWeight
	a.ReqUrl = spot_api.NewUBaseReqUrl(weightInfoMap, a.SubAccount)
	a.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}
	a.GetSymbolName = spot_api.GetFutureSymbolName
	return a
}
func NewUApiClient2(conf base.APIConf, cli *http.Client, maps ...interface{}) *u_api.UApiClient {
	var (
		weightInfoMap map[string]int64
		a             = &u_api.UApiClient{}
	)
	a.APIConf = conf
	if conf.EndPoint != "" {
		a.EndPoint = conf.EndPoint
	} else {
		a.EndPoint = spot_api.UBASE_API_BASE_URL
	}
	// 用户可以自定义限速规则
	for _, m := range maps {
		switch t := m.(type) {
		case map[client.WeightType]*client.WeightInfo:
			a.WeightInfo = t
			// 将传入的limit和interval归一化到number/minute
			for _, v := range a.WeightInfo {
				v.Limit = v.Limit * 60 / v.IntervalSec
				v.IntervalSec = 60
			}
		case config.ExchangeWeightInfo:
			weightInfoMap = t.UBase
		}
	}
	globalOnce.Do(func() {
		if usSpotApi.BinanceUBaseWeight == nil {
			usSpotApi.BinanceUBaseWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = usSpotApi.BinanceUBaseWeight
	a.ReqUrl = spot_api.NewUBaseReqUrl(weightInfoMap, a.SubAccount)
	a.HttpClient = cli
	a.GetSymbolName = spot_api.GetFutureSymbolName
	return a
}
