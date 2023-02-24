package spot_api

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ApiClient struct {
	base.APIConf
	HttpClient    *http.Client
	ReqUrl        *ReqUrl
	GetSymbolName func(string) string
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
	a.GetSymbolName = GetSpotSymbolName
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
	a.GetSymbolName = GetSpotSymbolName
	return a
}

func (a *ApiClient) GetUrl(url string) string {
	return a.EndPoint + url
}

func (a *ApiClient) DoRequest(uri, method string, params url.Values, result interface{}) error {
	header := &http.Header{}
	//header.Add("X-MBX-APIKEY", a.AccessKey)
	//header.Add("Content-Type", "application/x-www-form-urlencoded")
	//fmt.Println(a.GetUrl(uri))
	rsp, err := conn.Request(a.HttpClient, a.GetUrl(uri), method, header, params)
	if err == nil {
		err = json.Unmarshal(rsp, result)
		if err != nil {
			return err
		}
		logger.Logger.Debug(uri, params, err, result)
		if v, ok := result.(error); ok && v.Error() != "" {
			return v
		}
		return err
	}
	logger.Logger.Debug(uri, params, err)
	if v, ok := err.(*conn.HttpError); ok {
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

func (a *ApiClient) GetBook(product string, level int) (*RespBook, error) {
	requestURL := fmt.Sprintf("/products/%s/book?level=%d", product, level)
	res := &RespBook{}
	err := a.DoRequest(requestURL, "GET", nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApiClient) GetProducts() (*RespProducts, error) {
	requestURL := fmt.Sprintf("/products")
	res := &RespProducts{}
	err := a.DoRequest(requestURL, "GET", nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
