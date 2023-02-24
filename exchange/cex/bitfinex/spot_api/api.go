package spot_api

import (
	"bytes"
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"net/http"
	"net/url"
	"time"

	"github.com/goccy/go-json"
)

type ApiClient struct {
	base.APIConf
	HttpClient *http.Client
	EndPoint   string
	ReqUrl     *ReqUrl
}

func NewApiClient(conf base.APIConf) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
		}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)
	if a.EndPoint == "" {
		a.EndPoint = SPOT_API_Public_URL
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

func (client *ApiClient) GetUrl(uri string, method string, params url.Values) string {
	return client.EndPoint + uri
}

func (client *ApiClient) DoRequest(uri, method string, params url.Values, result interface{}) error {
	// fmt.Println("Do Request")

	header := &http.Header{}
	header.Add("Content-Type", CONTANTTYPE)

	params.Add("SignatureVersion", "2")
	params.Add("Timestamp", time.Now().Add(-time.Hour*8).Format("2006-01-02T15:04:05"))

	// fmt.Println(uri, method, params)
	rsp, err := conn.Request(client.HttpClient, client.GetUrl(uri, method, params), method, header, nil)

	if err == nil {
		if bytes.HasPrefix(rsp, []byte("\"code\"")) {
			var re RespError
			json.Unmarshal(rsp, &re)
			if re.Code != 0 && re.Msg != "" {
				unknown := false
				//-1006 UNEXPECTED_RESP 从消息总线收到意外的响应。 执行状态未知。
				//-1007 TIMEOUT 等待后端服务器响应超时。 发送状态未知； 执行状态未知。
				if re.Code == -1006 || re.Code == -1007 {
					unknown = true
				}
				logger.Logger.Debug(uri, params, re)
				return &base.ApiError{Code: 200, BizCode: re.Code, ErrMsg: re.Msg, UnknownStatus: unknown}
			}
		}
		err = json.Unmarshal(rsp, result)
		logger.Logger.Debug(uri, params, err, result)
		return err
	}
	logger.Logger.Debug(uri, params, err)
	// err not nil
	if v, ok := err.(*conn.HttpError); ok {
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

func (client *ApiClient) GetExchange() (*ExchangeSymbols, error) {
	params := url.Values{}
	res := &ExchangeSymbols{}
	err := client.DoRequest(client.ReqUrl.EXCHANGE_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetCurrency() (*CurrencySymbols, error) {
	params := url.Values{}
	res := &CurrencySymbols{}
	err := client.DoRequest(client.ReqUrl.CURRENCY_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
