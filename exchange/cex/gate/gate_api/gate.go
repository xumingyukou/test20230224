package gate_api

import (
	"clients/config"
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type ClientGate struct {
	base.APIConf
	ReqUrl        *ReqUrl
	HttpClient    *http.Client
	GetSymbolName func(string) string
}

func NewClientGateConf() *ClientGate {

	config.LoadExchangeConfig("conf/exchange.toml")

	conf := base.APIConf{}
	var (
		c = &ClientGate{
			APIConf: conf,
		}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)
	if conf.EndPoint == "" {
		c.EndPoint = GLOBAL_API_BASE_URL
	}
	if conf.ProxyUrl == "" {
		proxyUrl, err = url.Parse("http://127.0.0.1:7890")
		if err != nil {
			logger.Logger.Error("can not set proxy:", conf.ProxyUrl)
		}
		transport = http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	c.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}

	c.ReqUrl = NewSpotReqUrl()
	return c
}

func NewClientGate(conf base.APIConf) *ClientGate {
	var (
		c = &ClientGate{
			APIConf: conf,
		}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)
	if conf.EndPoint == "" {
		c.EndPoint = GLOBAL_API_BASE_URL
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
	c.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}

	c.ReqUrl = NewSpotReqUrl()
	return c
}

func NewClientGate2(conf base.APIConf, client *http.Client) *ClientGate {
	var (
		c = &ClientGate{
			APIConf: conf,
		}
	)
	if conf.EndPoint == "" {
		c.EndPoint = GLOBAL_API_BASE_URL
	}
	c.ReqUrl = NewSpotReqUrl()
	c.HttpClient = client
	return c
}

func (c *ClientGate) DoRequest(uri, method string, params url.Values, result interface{}, batch ...map[string]string) error {
	var (
		err error
		rsp []byte
	)

	header := &http.Header{}
	//loc, err := time.LoadLocation("English")

	if c.IsTest {
		header.Add("x-simulated-trading", "1")
	}
	rsp, err = conn.Request(c.HttpClient, c.GetUrl(uri), method, header, params)
	if err == nil {
		err = json.Unmarshal(rsp, result)
		if err != nil {
			return err
		}
		logger.Logger.Debug(uri, params, err, result)
		// 如果返回错误不为空
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

func (c *ClientGate) GetUrl(uri string) string {
	return c.EndPoint + uri
}

func (c *ClientGate) Market_Books_Info(pair string, options *url.Values) (*Resp_Market_Books, error) {
	params := url.Values{}
	params.Add("currency_pair", pair)
	params.Add("with_id", "true")
	ParseOptions(options, &params)
	res := &Resp_Market_Books{}
	err := c.DoRequest(c.ReqUrl.MARKET_BOOKS, "GET", params, res)
	return res, err
}

func ParseOptions(options *url.Values, params *url.Values) {
	if options != nil {
		for key := range *options {
			if options.Get(key) != "" {
				params.Add(key, options.Get(key))
			}
		}
	}
}

func (c *ClientGate) SymbolsInfo(options *url.Values) (*RespSymbols, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespSymbols{}
	err := c.DoRequest(c.ReqUrl.SYMBOLS, "GET", params, res)
	return res, err
}
