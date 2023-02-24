package spot_api

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ApiClient struct {
	base.APIConf
	HttpClient *http.Client
	ReqUrl     *ReqUrl
}

func NewApiClient(conf base.APIConf) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
			ReqUrl:  NewSpotReqUrl(),
		}
		proxyUrl  *url.URL
		transport = http.Transport{}
		err       error
	)

	if conf.EndPoint == "" {
		a.EndPoint = API_BASE_URL
	}

	if conf.ProxyUrl != "" {
		proxyUrl, err = url.Parse(conf.ProxyUrl)
		if err != nil {
			logger.Logger.Error("set proxy:", conf.ProxyUrl)
		}
		transport = http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	a.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}
	return a
}

func NewApiClient2(conf base.APIConf, cli *http.Client) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
			ReqUrl:  NewSpotReqUrl(),
		}
	)

	a.HttpClient = cli
	return a
}

func (client *ApiClient) getHeader(method string, path string, requireSign bool, body []byte) (map[string]string, error) {
	header := make(map[string]string)
	// TODO: authentication

	return header, nil
}

func (client *ApiClient) GetUrl(url string) string {
	return client.EndPoint + url
}

func (client *ApiClient) DoRequest(uri string, method string, requireSign bool, params map[string]interface{}, result interface{}) error {

	var err error
	body := make([]byte, 0)
	urlParam := ""
	if params != nil {
		if len(params) != 0 && (method == "GET" || method == "DELETE") {
			urlParam = "?"
			i := 0
			for key, element := range params {
				if i == 0 {
					urlParam += key + "=" + fmt.Sprintf("%v", element)
				} else {
					urlParam += "&" + key + "=" + fmt.Sprintf("%v", element)
				}
				i += 1
			}
		} else if method == "POST" || method == "PUT" {
			body, err = json.Marshal(params)
			if err != nil {
				logger.Logger.Error("json marshal", params, err.Error())
				return err
			}
		}
	}

	header, err := client.getHeader(method, uri+urlParam, requireSign, body)
	if err != nil {
		logger.Logger.Error("get header", err.Error())
		return err
	}
	fmt.Println(client.GetUrl(uri)+urlParam, method, string(body), header)

	rsp, err := conn.NewHttpRequest(client.HttpClient, method, client.GetUrl(uri)+urlParam, string(body), header)
	if err != nil {
		if v, ok := err.(*conn.HttpError); ok {
			return &base.ApiError{
				Code:          v.Code,
				UnknownStatus: v.Unknown,
			}
		} else {
			return &base.ApiError{
				Code:          500,
				ErrMsg:        err.Error(),
				UnknownStatus: true,
			}
		}
	}

	var re RespError
	err = json.Unmarshal(rsp, &re)
	if err != nil {
		logger.Logger.Error("parse response error", string(rsp))
		return err
	}

	code, err := transform.Str2Int32(re.Code)
	if err != nil {
		logger.Logger.Error("parser response error", re.Code)
		return err
	}

	if code != 0 {
		return &base.ApiError{
			Code:          200,
			BizCode:       code,
			ErrMsg:        re.Msg,
			UnknownStatus: false,
		}
	}

	err = json.Unmarshal(rsp, result)

	if err != nil {
		logger.Logger.Error("parse response body", string(rsp), err.Error())
		return err
	}

	return nil
}

func (client *ApiClient) GetSpotConfig() (*RespSpotConfig, error) {
	res := &RespSpotConfig{}
	err := client.DoRequest(client.ReqUrl.SPOT_CONFIG, "GET", false, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderBook(symbol string) (*RespOrderBook, error) {
	res := &RespOrderBook{}
	params := make(map[string]interface{})
	params["symbol"] = symbol
	err := client.DoRequest(client.ReqUrl.SPOT_ORDERBOOK, "GET", false, params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
