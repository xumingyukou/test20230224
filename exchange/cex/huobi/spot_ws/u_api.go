package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/huobi/spot_api"
	"clients/exchange/cex/huobi/u_api"
	"clients/logger"
	"github.com/warmplanet/proto/go/common"
	"strings"
)

var (
	getUBaseSwapSymbolOrNot   bool
	getUBaseFutureSymbolOrNot bool
)

type HuobiUBaseWebsocket struct {
	HuobiSpotWebsocket
}

func NewHuobiUBaseFutureWebsocket(conf base.WsConf, endPoint string) *HuobiCBaseWebsocket {
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
	d.handler.Market = common.Market_FUTURE

	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.UApiClient = u_api.NewUApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	} else {
		d.UApiClient = u_api.NewUApiClient(conf.APIConf)
	}
	d.WsReqUrl = spot_api.NewSpotWsUrl()

	InitUBaseFutureContractInfo(d.UApiClient)
	return d
}

func InitUBaseFutureContractInfo(uApiClient *u_api.UApiClient) {
	for !getUBaseFutureSymbolOrNot {
		getUBaseFutureSymbolOrNot = true
		logger.Logger.Info("uBase future map initializing......")
		data, err := uApiClient.GetContractInfo("futures")
		if err != nil {
			getUBaseFutureSymbolOrNot = false
			logger.Logger.Error("uBase future map initialize error:", err)
			continue
		}
		for _, pair := range data.Data {
			key := strings.ToUpper(pair.Pair) + "-" + GetFutureSub(pair.ContractType)
			symbolNameMap.Store(key, base.TransToSymbolInner(pair.Symbol+"/USDT", common.Market_FUTURE, u_api.TransStrToSymbolType(pair.ContractType)))
			contractSizeMap.Store(key, pair.ContractSize)
		}
	}
}

func NewHuoBiUBaseSwapWebsocket(conf base.WsConf, endPoint string) *HuobiUBaseWebsocket {
	if conf.ChanCap < 1 {
		conf.ChanCap = 1024
	}
	d := &HuobiUBaseWebsocket{}
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
	d.handler.Market = common.Market_SWAP

	if conf.AccessKey != "" && conf.SecretKey != "" {
		d.UApiClient = u_api.NewUApiClient(base.APIConf{
			ProxyUrl:    conf.ProxyUrl,
			ReadTimeout: conf.ReadTimeout,
			AccessKey:   conf.AccessKey,
			SecretKey:   conf.SecretKey,
		})
	} else {
		d.UApiClient = u_api.NewUApiClient(conf.APIConf)
	}
	d.WsReqUrl = spot_api.NewSpotWsUrl()

	InitUBaseSwapContractInfo(d.UApiClient)
	return d
}

func InitUBaseSwapContractInfo(uApiClient *u_api.UApiClient) {
	for !getUBaseSwapSymbolOrNot {
		getUBaseSwapSymbolOrNot = true
		logger.Logger.Info("uBase map initializing......")
		data, err := uApiClient.GetContractInfo("swap")
		if err != nil {
			getUBaseSwapSymbolOrNot = false
			logger.Logger.Error("uBase map initialize error:", err)
			continue
		}
		for _, pair := range data.Data {
			symbolNameMap.Store(pair.ContractCode, base.TransToSymbolInner(strings.Replace(pair.ContractCode, "-", "/", 1), common.Market_SWAP, common.SymbolType_SWAP_FOREVER))
			contractSizeMap.Store(pair.ContractCode, pair.ContractSize)
		}
	}
}
