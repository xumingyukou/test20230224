package spot_api

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type ApiClient struct {
	base.APIConf
	HttpClient *http.Client
	ReqUrl     *ReqUrl
	lock       sync.Mutex
}

func NewApiClient(conf base.APIConf) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
		}
		proxyUrl  *url.URL
		transport = http.Transport{}
		err       error
	)
	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
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
	a.ReqUrl = NewSpotReqUrl()
	return a
}
func NewApiClient2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ApiClient {
	a := &ApiClient{
		APIConf: conf,
	}
	if conf.EndPoint != "" {
		a.EndPoint = conf.EndPoint
	} else {
		a.EndPoint = SPOT_API_BASE_URL
	}

	a.HttpClient = cli
	a.ReqUrl = NewSpotReqUrl()

	return a
}

func (a *ApiClient) GetUrl(url string) string {
	return a.EndPoint + url
}

func (a *ApiClient) DoRequest(uri, method string, params url.Values, result interface{}) error {
	header := &http.Header{}
	//header.Add("X-MBX-APIKEY", a.AccessKey)
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	url_ := a.GetUrl(uri)
	_, rsp, err := conn.RequestWithHeader(a.HttpClient, url_, method, header, params)

	if err == nil {
		err = json.Unmarshal(rsp, result)
		logger.Logger.Debug(url_, params, err, result, string(rsp))
		return err
	}
	logger.Logger.Debug(url_, params, err, string(rsp))
	// err not nil
	if v, ok := err.(*conn.HttpError); ok {
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

func (a *ApiClient) GetWebsocketsToken() (*RespWSToken, error) {
	params := url.Values{}
	res := &RespWSToken{}
	err := a.DoRequest(a.ReqUrl.WEBSOCKETS_TOKEN_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApiClient) GetTradingPairsInfo() (*RespTradingPairs, error) {
	params := url.Values{}
	res := &RespTradingPairs{}
	err := a.DoRequest(a.ReqUrl.TRADING_PAIRS_INFO_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApiClient) GetOrderbook(symbol string) (*RespOrderbook, error) {
	params := url.Values{}
	reqURL := fmt.Sprintf("%s%s/", a.ReqUrl.ORDER_BOOK_URL, ReformatSymbols(symbol))
	res := &RespOrderbook{}
	err := a.DoRequest(reqURL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ReformatSymbols(input string) string {
	return strings.ToLower(strings.Replace(input, "/", "", -1))
}
