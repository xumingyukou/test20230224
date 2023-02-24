package test

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/huobi"
	"math/rand"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
)

/**
测试接口可用性
项目|测试范围|
-|-
01初始化|
02公共|获取交易所名称、获取所有现货币对、行情获取、交易所可用判断、交易手续费、提币手续费、交易精度
03余额|现货、杠杆、币本位合约、U本位合约
04下单|订单类型：limit、market、limit_maker，现货、U本位合约、币本位合约；订单取消、查询某一订单、查询历史订单
05划转|划转、历史查询
06提币|提币、提币历史
07子账户|划转、划转历史
08充值|充值历史
09订阅|depth、ticker、trade

todo 增加边界条件
*/

var (
	//trade config
	tradeSymbol         = "ETH/USDT"
	tradePrice  float64 = 1638
	tradeAmount         = 0.008

	//client config
	apiList      []base.CexApiInterface
	proxyUrl     = "http://127.0.0.1:1080"
	precisionMap = map[string]*client.PrecisionItem{
		"ETH/USDT": &client.PrecisionItem{
			Symbol:    "ETH/USDT",
			Type:      common.SymbolType_SPOT_NORMAL,
			Amount:    8,
			Price:     8,
			AmountMin: 0.001,
		},
	}
	clientInfoList = []*NewTestClient{
		// {
		// 	Name: "binance",
		// 	NewClientHandle: func(conf base.APIConf) base.CexApiInterface {
		// 		return binance.NewClientBinance(conf, precisionMap)
		// 	},
		// 	Exchange: common.Exchange_BINANCE,
		// },
		// {
		// 	Name: "okex",
		// 	NewClientHandle: func(conf base.APIConf) base.CexApiInterface {
		// 		return okex.NewClientOkex(conf, precisionMap)
		// 	},
		// 	Exchange: common.Exchange_OKEX,
		// },
		//{
		//	Name: "ftx",
		//	NewClientHandle: func(conf base.APIConf) base.CexApiInterface {
		//		return ftx.NewClientFTX(conf, precisionMap)
		//	},
		//	Exchange: common.Exchange_FTX,
		//},
		//{
		//	Name: "bybit",
		//	NewClientHandle: func(conf base.APIConf) base.CexApiInterface {
		//		return bybit.NewClientBybit(conf, precisionMap)
		//	},
		//	Exchange: common.Exchange_BYBIT,
		//},
		// {
		// 	Name: "kucoin",
		// 	NewClientHandle: func(conf base.APIConf) base.CexApiInterface {
		// 		return kucoin.NewClientKucoin(conf)
		// 	},
		// 	Exchange: common.Exchange_BYBIT,
		// },
		{
			Name: "huobi",
			NewClientHandle: func(conf base.APIConf) base.CexApiInterface {
				return huobi.NewClientHuobi(conf, precisionMap)
			},
			Exchange: common.Exchange_HUOBI,
		},
	}
)

type NewTestClient struct {
	Name            string
	NewClientHandle func(conf base.APIConf) base.CexApiInterface
	Exchange        common.Exchange
}

func init() {
	// 初始化配置
	rand.Seed(time.Now().UnixMilli())
	config.LoadExchangeConfig("../../../conf/exchange.toml")
	// 创建测试clients
	for _, clientInfo := range clientInfoList {
		conf := base.APIConf{
			ProxyUrl:    proxyUrl,
			ReadTimeout: 10,
			AccessKey:   config.ExchangeConfig.ExchangeList[clientInfo.Name].ApiKeyConfig.AccessKey,
			SecretKey:   config.ExchangeConfig.ExchangeList[clientInfo.Name].ApiKeyConfig.SecretKey,
			Passphrase:  config.ExchangeConfig.ExchangeList[clientInfo.Name].ApiKeyConfig.Passphrase,
		}
		apiList = append(apiList, clientInfo.NewClientHandle(conf))
	}
}
