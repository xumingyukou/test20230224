package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/huobi/c_api"
	"clients/exchange/cex/huobi/spot_api"
	"clients/logger"
	"net/url"
	"strings"

	"github.com/warmplanet/proto/go/common"
)

var (
	getCBaseSwapSymbolOrNot   bool
	getCBaseFutureSymbolOrNot bool
)

type HuobiCBaseWebsocket struct {
	HuobiSpotWebsocket
}

func NewHuobiCBaseFutureWebsocket(conf base.WsConf, endPoint string) *HuobiCBaseWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &HuobiCBaseWebsocket{}
	d.WsConf = conf
	if conf.EndPoint == "" {
		d.EndPoint = endPoint
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}

	d.handler = NewWebSocketSpotHandle(d.ChanCap)
	d.handler.Exchange = common.Exchange_HUOBI
	d.handler.Market = common.Market_FUTURE_COIN

	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.CApiClient = c_api.NewCApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	} else {
		d.CApiClient = c_api.NewCApiClient(conf.APIConf)
	}
	d.WsReqUrl = spot_api.NewSpotWsUrl()

	InitCBaseFutureContractInfo(d.CApiClient)
	return d
}

func InitCBaseFutureContractInfo(cApiClient *c_api.CApiClient) {
	for !getCBaseFutureSymbolOrNot {
		getCBaseFutureSymbolOrNot = true
		logger.Logger.Info("cBase future map initializing......")
		data, err := cApiClient.GetFutureContractInfo(&url.Values{})
		if err != nil {
			getCBaseFutureSymbolOrNot = false
			logger.Logger.Error("cBase future map initialize error:", err)
			continue
		}
		for _, pair := range data.Data {
			key := strings.ToUpper(pair.Symbol) + "_" + GetFutureSub(pair.ContractType)
			symbolNameMap.Store(key, base.TransToSymbolInner(pair.Symbol+"/USD", common.Market_FUTURE_COIN, c_api.TransStrToSymbolType(pair.ContractType)))
			contractSizeMap.Store(key, pair.ContractSize)
		}
	}
}

func NewHuobiCBaseSwapWebsocket(conf base.WsConf, endPoint string) *HuobiCBaseWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &HuobiCBaseWebsocket{}
	d.WsConf = conf
	if conf.EndPoint == "" {
		d.EndPoint = endPoint
	}
	if d.ReadTimeout == 0 {
		d.ReadTimeout = 300
	}
	if d.listenTimeout == 0 {
		d.listenTimeout = 1800
	}
	d.handler = NewWebSocketSpotHandle(d.ChanCap)
	d.handler.Exchange = common.Exchange_HUOBI
	d.handler.Market = common.Market_SWAP_COIN

	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.CApiClient = c_api.NewCApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	} else {
		d.CApiClient = c_api.NewCApiClient(conf.APIConf)
	}
	d.WsReqUrl = spot_api.NewSpotWsUrl()

	InitCBaseSwapContractInfo(d.CApiClient)
	return d
}

func InitCBaseSwapContractInfo(cApiClient *c_api.CApiClient) {
	for !getCBaseSwapSymbolOrNot {
		getCBaseSwapSymbolOrNot = true
		logger.Logger.Info("cBase swap map initializing......")
		data, err := cApiClient.GetSwapContractInfo(&url.Values{})
		if err != nil {
			getCBaseSwapSymbolOrNot = false
			logger.Logger.Error("cBase swap map initialize error:", err)
			continue
		}
		for _, pair := range data.Data {
			symbolNameMap.Store(pair.ContractCode, base.TransToSymbolInner(strings.Replace(pair.ContractCode, "-", "/", 1), common.Market_SWAP_COIN, common.SymbolType_SWAP_COIN_FOREVER))
			contractSizeMap.Store(pair.ContractCode, pair.ContractSize)
		}
	}
}
