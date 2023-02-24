package spot_api

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
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
		a.EndPoint = REST_PUBLIC_BASE_URL
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

	// 判断unmarshal类型是否为列表的指针

	t := reflect.TypeOf(result).Elem().Kind()
	if t == reflect.Slice {
		// 列表直接解析
		err = json.Unmarshal(rsp, result)
		if err != nil {
			// 可能返回了错误码
			var re RespError
			err = json.Unmarshal(rsp, &re)
			if err != nil {
				logger.Logger.Error("parse code error", string(rsp))
				return err
			}

			if re.Code != 0 {
				return &base.ApiError{
					Code:          200,
					BizCode:       re.Code,
					ErrMsg:        re.Message,
					UnknownStatus: false,
				}
			} else {
				// 可能是其他格式
				err = errors.New("unknown status")
				return err
			}
		}
	} else if t == reflect.Struct {
		// 先解析错误码
		var re RespError
		err = json.Unmarshal(rsp, &re)
		if err != nil {
			logger.Logger.Error("parse code error", string(rsp))
			return err
		}

		if re.Code != 0 {
			return &base.ApiError{
				Code:          200,
				BizCode:       re.Code,
				ErrMsg:        re.Message,
				UnknownStatus: false,
			}
		}

		// 再解析业务数据
		err = json.Unmarshal(rsp, result)
		if err != nil {
			logger.Logger.Error("parse response body", string(rsp), err.Error())
			return err
		}

	} else {
		err = errors.New("unsupported json unmarshal target")
		return err
	}

	return nil
}

func (client *ApiClient) GetMarkets() (*RespGetMarkets, error) {
	res := &RespGetMarkets{}
	err := client.DoRequest(client.ReqUrl.MARKETS, "GET", false, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
