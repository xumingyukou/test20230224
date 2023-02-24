package spot_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/spot_api"
	"clients/logger"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/warmplanet/proto/go/client"
)

func ParseOptions(options *url.Values, params *url.Values) {
	if options != nil {
		for key := range *options {
			if options.Get(key) != "" {
				params.Add(key, options.Get(key))
			}
		}
	}
}

var (
	// BinanceSpotWeight 定义全局的限制
	BinanceSpotWeight  *base.RateLimitMgr
	BinanceUBaseWeight *base.RateLimitMgr
	BinanceCBaseWeight *base.RateLimitMgr
	globalOnce         sync.Once
)

func NewApiClient(conf base.APIConf, maps ...interface{}) *spot_api.ApiClient {
	var (
		a = &spot_api.ApiClient{
			APIConf: conf,
		}
		proxyUrl      *url.URL
		transport     = http.Transport{}
		err           error
		weightInfoMap map[string]int64
	)
	if conf.EndPoint == "" {
		a.EndPoint = spot_api.SPOT_API_US_URL
	}
	a.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
	globalOnce.Do(func() {
		if BinanceSpotWeight == nil {
			BinanceSpotWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = BinanceSpotWeight
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
			weightInfoMap = t.Spot

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
	a.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}
	a.ReqUrl = spot_api.NewSpotReqUrl(weightInfoMap, a.SubAccount)
	a.GetSymbolName = spot_api.GetSpotSymbolName
	return a
}

func NewApiClient2(conf base.APIConf, cli *http.Client, maps ...interface{}) *spot_api.ApiClient {
	var (
		weightInfoMap map[string]int64
		a             = &spot_api.ApiClient{APIConf: conf}
	)

	if conf.EndPoint != "" {
		a.EndPoint = conf.EndPoint
	} else {
		a.EndPoint = spot_api.SPOT_API_US_URL
	}
	a.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
	globalOnce.Do(func() {
		if BinanceSpotWeight == nil {
			BinanceSpotWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = BinanceSpotWeight
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
			weightInfoMap = t.Spot
		}
	}
	a.HttpClient = cli
	a.ReqUrl = spot_api.NewSpotReqUrl(weightInfoMap, a.SubAccount)
	a.GetSymbolName = spot_api.GetSpotSymbolName
	return a
}
