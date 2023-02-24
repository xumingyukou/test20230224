package c_api

import (
	"clients/config"
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"github.com/warmplanet/proto/go/client"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CApiClient struct {
	spot_api.ApiClient
}

func NewCApiClient(conf base.APIConf, maps ...interface{}) *CApiClient {
	var (
		a             = &CApiClient{}
		proxyUrl      *url.URL
		transport     = http.Transport{}
		err           error
		weightInfoMap map[string]int64
	)
	a.APIConf = conf
	if conf.EndPoint == "" {
		a.EndPoint = spot_api.CBASE_API_BASE_URL
		if conf.IsTest {
			a.EndPoint = spot_api.UBASE_TEST_BASE_URL
		}
	}
	spot_api.GlobalCBaseOnce.Do(func() {
		if spot_api.BinanceCBaseWeight == nil {
			spot_api.BinanceCBaseWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = spot_api.BinanceCBaseWeight
	// 用户可以自定义限速规则
	for _, m := range maps {
		switch t := m.(type) {
		case map[client.WeightType]*client.WeightInfo:
			a.WeightInfo = t
			// 将传入的limit和interval归一化到number/minute
			for _, v := range a.WeightInfo {
				v.Limit = v.Limit * 60 / v.IntervalSec
				v.IntervalSec = 60
			}
		case config.ExchangeWeightInfo:
			weightInfoMap = t.CBase
		}
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
	a.ReqUrl = spot_api.NewCBaseReqUrl(weightInfoMap, a.SubAccount)
	a.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}
	a.GetSymbolName = spot_api.GetFutureSymbolName
	return a
}
func NewCApiClient2(conf base.APIConf, cli *http.Client, maps ...interface{}) *CApiClient {
	var (
		a             = &CApiClient{}
		weightInfoMap map[string]int64
	)
	a.APIConf = conf
	if conf.EndPoint != "" {
		a.EndPoint = conf.EndPoint
	} else {
		a.EndPoint = spot_api.CBASE_API_BASE_URL
	}
	spot_api.GlobalCBaseOnce.Do(func() {
		if spot_api.BinanceCBaseWeight == nil {
			spot_api.BinanceCBaseWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = spot_api.BinanceCBaseWeight
	// 用户可以自定义限速规则
	for _, m := range maps {
		switch t := m.(type) {
		case map[client.WeightType]*client.WeightInfo:
			a.WeightInfo = t
			// 将传入的limit和interval归一化到number/minute
			for _, v := range a.WeightInfo {
				v.Limit = v.Limit * 60 / v.IntervalSec
				v.IntervalSec = 60
			}
		case config.ExchangeWeightInfo:
			weightInfoMap = t.CBase
		}
	}
	a.ReqUrl = spot_api.NewCBaseReqUrl(weightInfoMap, a.SubAccount)
	a.HttpClient = cli
	a.GetSymbolName = spot_api.GetFutureSymbolName
	return a
}

func (a *CApiClient) PremiumIndex(options *url.Values) (*u_api.RespPremiumIndexList, error) {
	/*
		最新标记价格和资金费率:采集各大交易所数据加权平均
		symbol	STRING	NO	交易对
	*/
	var (
		params = url.Values{}
		res    *u_api.RespPremiumIndexList
		err    error
	)
	spot_api.ParseOptions(options, &params)
	if options == nil || options.Get("symbol") == "" {
		res = &u_api.RespPremiumIndexList{}
		err = a.DoRequest(a.ReqUrl.PREMIUMINDEX_URL, "GET", params, res)
	} else {
		resp := &u_api.RespPremiumIndexItem{}
		err = a.DoRequest(a.ReqUrl.PREMIUMINDEX_URL, "GET", params, resp)
		res = &u_api.RespPremiumIndexList{resp}
	}
	return res, err
}

func (a *CApiClient) PositionSideDual(dualSidePosition string) (*spot_api.RespError, error) {
	/*
		更改持仓模式(TRADE)
		变换用户在 所有symbol 合约上的持仓模式：双向持仓或单向持仓。
		dualSidePosition	STRING	YES	"true": 双向持仓模式；"false": 单向持仓模式
	*/
	params := url.Values{}
	params.Add("dualSidePosition", dualSidePosition)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &spot_api.RespError{}
	err = a.DoRequest(a.ReqUrl.POSITIONSIDE_DUAL_URL, "POST", params, res)
	return res, err
}

func (a *CApiClient) GetPositionSideDual() (*RespPositionSideDual, error) {
	/*
		更改持仓模式(TRADE)
		变换用户在 所有symbol 合约上的持仓模式：双向持仓或单向持仓。
		dualSidePosition	STRING	YES	"true": 双向持仓模式；"false": 单向持仓模式
	*/
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &RespPositionSideDual{}
	err = a.DoRequest(a.ReqUrl.POSITIONSIDE_DUAL_URL, "GET", params, res)
	return res, err
}

func (a *CApiClient) Order(symbol string, side spot_api.SideType, type_ spot_api.OrderType, options *url.Values) (*u_api.RespUBaseOrderResult, error) {
	/*
		symbol				STRING	YES	交易对
		side				ENUM	YES	买卖方向 SELL, BUY
		type				ENUM	YES	订单类型 LIMIT, MARKET, STOP, TAKE_PROFIT, STOP_MARKET, TAKE_PROFIT_MARKET, TRAILING_STOP_MARKET
		positionSide		ENUM	NO	持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
		reduceOnly			STRING	NO	true, false; 非双开模式下默认false；双开模式下不接受此参数； 使用closePosition不支持此参数。
		quantity			DECIMAL	NO	下单数量,使用closePosition不支持此参数。
		price				DECIMAL	NO	委托价格
		newClientOrderId	STRING	NO	用户自定义的订单号，不可以重复出现在挂单中。如空缺系统会自动赋值。必须满足正则规则 ^[\.A-Z\:/a-z0-9_-]{1,36}$
		stopPrice			DECIMAL	NO	触发价, 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
		closePosition		STRING	NO	true, false；触发后全部平仓，仅支持STOP_MARKET和TAKE_PROFIT_MARKET；不与quantity合用；自带只平仓效果，不与reduceOnly 合用
		activationPrice		DECIMAL	NO	追踪止损激活价格，仅TRAILING_STOP_MARKET 需要此参数, 默认为下单当前市场价格(支持不同workingType)
		callbackRate		DECIMAL	NO	追踪止损回调比例，可取值范围[0.1, 5],其中 1代表1% ,仅TRAILING_STOP_MARKET 需要此参数
		timeInForce			ENUM	NO	有效方法
		workingType			ENUM	NO	stopPrice 触发类型: MARK_PRICE(标记价格), CONTRACT_PRICE(合约最新价). 默认 CONTRACT_PRICE
		priceProtect		STRING	NO	条件单触发保护："TRUE","FALSE", 默认"FALSE". 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
		newOrderRespType	ENUM	NO	"ACK", "RESULT", 默认 "ACK"
		recvWindow			LONG	NO
		timestamp			LONG	YES
		根据 order type的不同，某些参数强制要求，具体如下:
		Type							强制要求的参数
		LIMIT							timeInForce, quantity, price
		MARKET							quantity
		STOP, TAKE_PROFIT				quantity, price, stopPrice
		STOP_MARKET, TAKE_PROFIT_MARKET	stopPrice
		TRAILING_STOP_MARKET			callbackRate
		条件单的触发必须:
		如果订单参数priceProtect为true:
		达到触发价时，MARK_PRICE(标记价格)与CONTRACT_PRICE(合约最新价)之间的价差不能超过改symbol触发保护阈值
		触发保护阈值请参考接口GET /fapi/v1/exchangeInfo 返回内容相应symbol中"triggerProtect"字段
		STOP, STOP_MARKET 止损单:
		买入: 最新合约价格/标记价格高于等于触发价stopPrice
		卖出: 最新合约价格/标记价格低于等于触发价stopPrice
		TAKE_PROFIT, TAKE_PROFIT_MARKET 止盈单:
		买入: 最新合约价格/标记价格低于等于触发价stopPrice
		卖出: 最新合约价格/标记价格高于等于触发价stopPrice
		TRAILING_STOP_MARKET 跟踪止损单:
		买入: 当合约价格/标记价格区间最低价格低于激活价格activationPrice,且最新合约价格/标记价高于等于最低价设定回调幅度。
		卖出: 当合约价格/标记价格区间最高价格高于激活价格activationPrice,且最新合约价格/标记价低于等于最高价设定回调幅度。
		TRAILING_STOP_MARKET 跟踪止损单如果遇到报错 {"code": -2021, "msg": "Order would immediately trigger."}
		表示订单不满足以下条件:
		买入: 指定的activationPrice 必须小于 latest price
		卖出: 指定的activationPrice 必须大于 latest price
		newOrderRespType 如果传 RESULT:
		MARKET 订单将直接返回成交结果；
		配合使用特殊 timeInForce 的 LIMIT 订单将直接返回成交或过期拒绝结果。
		STOP_MARKET, TAKE_PROFIT_MARKET 配合 closePosition=true:
		条件单触发依照上述条件单触发逻辑
		条件触发后，平掉当时持有所有多头仓位(若为卖单)或当时持有所有空头仓位(若为买单)
		不支持 quantity 参数
		自带只平仓属性，不支持reduceOnly参数
		双开模式下,LONG方向上不支持BUY; SHORT 方向上不支持SELL
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("side", string(side))
	params.Add("type", string(type_))
	spot_api.ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &u_api.RespUBaseOrderResult{}
	err = a.DoRequest(a.ReqUrl.ORDER_URL, "POST", params, res)
	return res, err
}

func (a *CApiClient) BatchOrder(batchOrders string) ([]*u_api.RespUBaseOrderResult, error) {
	/*
		batchOrders	list	YES	订单列表，最多支持5个订单
		其中batchOrders应以list of JSON格式填写订单参数 例子: /fapi/v1/batchOrders?batchOrders=[{"type":"LIMIT","timeInForce":"GTC","symbol":"BTCUSDT","side":"BUY","price":"10001","quantity":"0.001"}]
	*/
	params := url.Values{}
	params.Add("batchOrders", batchOrders)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res []*u_api.RespUBaseOrderResult
	err = a.DoRequest(a.ReqUrl.BATCHORDERS_URL, "POST", params, res)
	return res, err
}

func (a *CApiClient) GetOrder(symbol string, options *url.Values) (*u_api.RespUBaseOrderResult, error) {
	/*
		查询订单状态
		请注意，如果订单满足如下条件，不会被查询到：
		订单的最终状态为 CANCELED 或者 EXPIRED, 并且
		订单没有任何的成交记录, 并且
		订单生成时间 + 7天 < 当前时间

		symbol	STRING	YES	交易对
		orderId	LONG	NO	系统订单号
		origClientOrderId	STRING	NO	用户自定义的订单号
		recvWindow	LONG	NO
		timestamp	LONG	YES
		至少需要发送 orderId 与 origClientOrderId中的一个
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	spot_api.ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &u_api.RespUBaseOrderResult{}
	err = a.DoRequest(a.ReqUrl.ORDER_URL, "GET", params, res)
	return res, err
}

func (a *CApiClient) GetOpenOrder(symbol string, options *url.Values) (*u_api.RespUBaseOrderResult, error) {
	/*
		查询当前挂单
		symbol	STRING	YES	交易对
		orderId	LONG	NO	系统订单号
		origClientOrderId	STRING	NO	用户自定义的订单号
		recvWindow	LONG	NO
		timestamp	LONG	YES
		orderId 与 origClientOrderId 中的一个为必填参数
		查询的订单如果已经成交或取消，将返回报错 "Order does not exist."
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	spot_api.ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &u_api.RespUBaseOrderResult{}
	err = a.DoRequest(a.ReqUrl.OPENORDER_URL, "GET", params, res)
	return res, err
}

func (a *CApiClient) GetOpenOrders(options *url.Values) ([]*u_api.RespUBaseOrderResult, error) {
	/*
		查看当前全部挂单
		symbol	STRING	NO	交易对，权重: - 带symbol 1 - 不带 40
		不带symbol参数，会返回所有交易对的挂单
	*/
	params := url.Values{}
	spot_api.ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res []*u_api.RespUBaseOrderResult
	err = a.DoRequest(a.ReqUrl.OPENORDERS_URL, "GET", params, &res)
	return res, err
}

func (a *CApiClient) CancelOrder(symbol string, options *url.Values) (*u_api.RespUBaseOrderResult, error) {
	/*
		撤销订单 (TRADE)
		symbol	STRING	YES	交易对
		orderId	LONG	NO	系统订单号
		origClientOrderId	STRING	NO	用户自定义的订单号
		recvWindow	LONG	NO
		timestamp	LONG	YES
		orderId 与 origClientOrderId 必须至少发送一个
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	spot_api.ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &u_api.RespUBaseOrderResult{}
	err = a.DoRequest(a.ReqUrl.ORDER_URL, "DELETE", params, res)
	return res, err
}

func (a *CApiClient) CancelAllOpenOrders(symbol string, options *url.Values) (*spot_api.RespError, error) {
	/*
		撤销全部订单
		symbol	STRING	YES	交易对
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	spot_api.ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &spot_api.RespError{}
	err = a.DoRequest(a.ReqUrl.ALLOPENORDERS_URL, "DELETE", params, res)
	return res, err
}

func (a *CApiClient) CancelBatchOrders(symbol string, options *url.Values) ([]*u_api.RespUBaseOrderResult, error) {
	/*
		批量撤销订单
		symbol	STRING	YES	交易对
		orderIdList	LIST<LONG>	NO	系统订单号, 最多支持10个订单 比如[1234567,2345678]
		origClientOrderIdList	LIST<STRING>	NO	用户自定义的订单号, 最多支持10个订单 比如["my_id_1","my_id_2"] 需要encode双引号。逗号后面没有空格。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	spot_api.ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res []*u_api.RespUBaseOrderResult
	err = a.DoRequest(a.ReqUrl.BATCHORDERS_URL, "DELETE", params, res)
	return res, err
}

func (a *CApiClient) Balance() (*RespBalance, error) {
	/*
		账户余额
	*/
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &RespBalance{}
	err = a.DoRequest(a.ReqUrl.BALANCE_URL, "GET", params, res)
	return res, err
}

func (a *CApiClient) Account() (*u_api.RespAccount, error) {
	/*
		账户信息
		现有账户信息。 用户在单资产模式和多资产模式下会看到不同结果，响应部分的注释解释了两种模式下的不同。
	*/
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &u_api.RespAccount{}
	err = a.DoRequest(a.ReqUrl.ACCOUNT_URL, "GET", params, res)
	return res, err
}

func (a *CApiClient) Leverage(symbol string, leverage int64) (*spot_api.RespError, error) {
	/*
		调整用户在指定symbol合约的开仓杠杆。
		symbol	STRING	YES	交易对
		leverage	INT	YES	目标杠杆倍数：1 到 125 整数
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("leverage", strconv.FormatInt(leverage, 10))
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &spot_api.RespError{}
	err = a.DoRequest(a.ReqUrl.LEVERAGE_URL, "POST", params, res)
	return res, err
}

func (a *CApiClient) MarginType(symbol, marginType string) (*spot_api.RespError, error) {
	/*
		变换用户在指定symbol合约上的保证金模式：逐仓或全仓。
		symbol	STRING	YES	交易对
		marginType	ENUM	YES	保证金模式 ISOLATED(逐仓), CROSSED(全仓)
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("marginType", marginType)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &spot_api.RespError{}
	err = a.DoRequest(a.ReqUrl.MARGINTYPE_URL, "POST", params, res)
	return res, err
}

func (a *CApiClient) PositionMargin(symbol string, amount float64, type_ int64) (*spot_api.RespError, error) {
	/*
		针对逐仓模式下的仓位，调整其逐仓保证金资金。
		symbol	STRING	YES	交易对
		positionSide	ENUM	NO	持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
		amount	DECIMAL	YES	保证金资金
		type	INT	YES	调整方向 1: 增加逐仓保证金，2: 减少逐仓保证金
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	params.Add("type", strconv.FormatInt(type_, 10))
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &spot_api.RespError{}
	err = a.DoRequest(a.ReqUrl.MARGINTYPE_URL, "POST", params, res)
	return res, err
}

func (a *CApiClient) CommissionRate(symbol string) (*u_api.RespCommissionRate, error) {
	/*
		针对逐仓模式下的仓位，调整其逐仓保证金资金。
		symbol	STRING	YES	交易对
		positionSide	ENUM	NO	持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
		amount	DECIMAL	YES	保证金资金
		type	INT	YES	调整方向 1: 增加逐仓保证金，2: 减少逐仓保证金
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res = &u_api.RespCommissionRate{}
	err = a.DoRequest(a.ReqUrl.COMMISSIONRATE_URL, "GET", params, res)
	return res, err
}

func (a *CApiClient) UserDataStream() (*spot_api.RespUserDataStream, error) {
	/*
		开始一个新的数据流。除非发送 keepalive，否则数据流于60分钟后关闭。
		如果该帐户具有有效的listenKey，则将返回该listenKey并将其有效期延长60分钟。
		权重: 1
	*/
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &spot_api.RespUserDataStream{}
	err = a.DoRequest(a.ReqUrl.USERDATASTREAM_URL, "POST", params, res)
	return res, err
}

func (a *CApiClient) PutUserDataStream(listenKey string) (interface{}, error) {
	/*
		listenKey	STRING	YES
		有效期延长至本次调用后60分钟,建议每30分钟发送一个 ping 。
	*/
	params := url.Values{}
	params.Add("listenKey", listenKey)
	var res interface{}
	err := a.DoRequest(a.ReqUrl.USERDATASTREAM_URL, "PUT", params, &res)
	return res, err
}

func (a *CApiClient) DELETEUserDataStream(listenKey string) (interface{}, error) {
	/*
		listenKey	STRING	YES
		关闭用户数据流。
	*/
	params := url.Values{}
	params.Add("listenKey", listenKey)
	var res interface{}
	err := a.DoRequest(a.ReqUrl.USERDATASTREAM_URL, "DELETE", params, &res)
	return res, err
}
