package future_api

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/ftx/spot_api"
	"clients/logger"
	"net/http"
	"net/url"
	"time"
)

type UApiClient struct {
	spot_api.ApiClient
}

func NewUApiClient(conf base.APIConf, maps ...interface{}) *UApiClient {
	var (
		a         = &UApiClient{}
		proxyUrl  *url.URL
		transport = http.Transport{}
		err       error
	)
	a.APIConf = conf
	if conf.EndPoint == "" {
		a.EndPoint = spot_api.SPOT_API_BASE_URL
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
	if spot_api.FTXWeight == nil {
		spot_api.FTXWeight = base.NewRateLimitMgr()
	}
	a.ReqUrl = spot_api.NewUBaseReqUrl()
	a.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}
	return a
}
func NewUApiClient2(conf base.APIConf, client *http.Client, maps ...interface{}) *UApiClient {
	var (
		a = &UApiClient{}
	)
	a.APIConf = conf
	if conf.EndPoint == "" {
		a.EndPoint = spot_api.FUTURE_API_BASE_URL
	}
	if spot_api.FTXWeight == nil {
		spot_api.FTXWeight = base.NewRateLimitMgr()
	}
	a.ReqUrl = spot_api.NewUBaseReqUrl()
	a.HttpClient = client
	return a
}

func (client *UApiClient) ListAllFutures() (*RespListAllFutures, error) {
	params := make(map[string]interface{})
	res := &RespListAllFutures{}
	err := client.DoRequest(client.ReqUrl.FUTURES_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *UApiClient) GetFuture(futureName string) (*RespGetFuture, error) {
	params := make(map[string]interface{})
	res := &RespGetFuture{}
	updatedURL := base.ReqUrlInfo{
		Url:                 client.ReqUrl.FUTURES_URL.Url + "/" + futureName,
		RateLimitConsumeMap: client.ReqUrl.MARKETS_URL.RateLimitConsumeMap,
	}
	err := client.DoRequest(updatedURL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *UApiClient) GetFundingRates(startTime int64, endTime int64) (*RespFundingRates, error) {
	params := make(map[string]interface{})
	params["start_time"] = startTime
	params["end_time"] = endTime
	res := &RespFundingRates{}
	err := client.DoRequest(client.ReqUrl.FUNDING_RATES_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//func (a *UApiClient) Order(symbol string, side spot_api.SideType, type_ spot_api.OrderType, options *url.Values) (*RespUBaseOrderResult, error) {
