package spot_api

import (
	"clients/exchange/cex/base"
	"clients/logger"
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
	// if conf.EndPoint == "" {
	// 	a.EndPoint = ""
	// }

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
	// if conf.EndPoint == "" {
	// 	a.EndPoint = SPOT_API_BASE_URL
	// }

	a.HttpClient = cli
	return a
}

func (cli *ApiClient) GetSymbols() ([]string, error) {
	symbols := []string{"btc_jpy", "xrp_jpy", "xrp_btc", "ltc_jpy", "ltc_btc", "eth_jpy", "eth_btc", "mona_jpy", "mona_btc", "bcc_jpy", "bcc_btc", "xlm_jpy", "xlm_btc", "qtum_jpy", "qtum_btc", "bat_jpy", "bat_btc", "omg_jpy", "omg_btc", "xym_jpy", "xym_btc", "link_jpy", "link_btc", "mkr_jpy", "mkr_btc", "boba_jpy", "boba_btc", "enj_jpy", "enj_btc", "matic_jpy", "matic_btc", "dot_jpy", "doge_jpy"}
	return symbols, nil
}
