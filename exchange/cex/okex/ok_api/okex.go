package ok_api

import (
	"clients/config"
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ClientOkex struct {
	base.APIConf
	ReqUrl        *ReqUrl
	HttpClient    *http.Client
	GetSymbolName func(string) string
}

var (
	// OkexWeight 定义全局的限制
	OkexWeight *base.RateLimitMgr
)

func WeightAllow(consumeMap []*base.RateLimitConsume) error {
	for _, rate := range consumeMap {
		if err := OkexWeight.Consume(*rate); err != nil {
			return err
		}
	}
	return nil
}

func ClearCapacity(consumeMap []*base.RateLimitConsume) {
	for _, rate := range consumeMap {
		tmp := OkexWeight.GetInstance(rate.LimitTypeName)
		tmp.Instance.Clear()
	}
}

func NewClientOkexConf() *ClientOkex {

	config.LoadExchangeConfig("conf/exchange.toml")

	conf := base.APIConf{
		AccessKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.AccessKey,
		SecretKey:  config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.SecretKey,
		Passphrase: config.ExchangeConfig.ExchangeList["okex"].ApiKeyConfig.Passphrase,
	}
	var (
		c = &ClientOkex{
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
	if OkexWeight == nil {
		OkexWeight = base.NewRateLimitMgr()
	}
	c.ReqUrl = NewSpotReqUrl(nil, conf.SubAccount)
	return c
}

func NewClientOkex(conf base.APIConf) *ClientOkex {
	var (
		c = &ClientOkex{
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

	if OkexWeight == nil {
		OkexWeight = base.NewRateLimitMgr()
	}
	c.ReqUrl = NewSpotReqUrl(nil, conf.SubAccount)
	return c
}

func NewClientOkex2(conf base.APIConf, client *http.Client) *ClientOkex {
	var (
		c = &ClientOkex{
			APIConf: conf,
		}
	)
	if conf.EndPoint == "" {
		c.EndPoint = GLOBAL_API_BASE_URL
	}
	if OkexWeight == nil {
		OkexWeight = base.NewRateLimitMgr()
	}
	c.ReqUrl = NewSpotReqUrl(nil, conf.SubAccount)
	c.HttpClient = client
	return c
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

func (c *ClientOkex) GetUrl(uri string) string {
	return c.EndPoint + uri
}

func (c *ClientOkex) GetUri(url string) string {
	return url[len(c.EndPoint):]
}

func (c *ClientOkex) GetUrlHeadPathBody(method string, uri string, params url.Values) string {
	if method == "GET" {
		// 如果没有参数，不能附加问号
		if len(params) == 0 {
			return uri
		} else {
			return uri + "?" + params.Encode()
		}
	} else if method == "POST" {
		params1 := make(map[string]string)
		for key, value := range params {
			params1[key] = value[0]
		}
		dataByte, err := json.Marshal(params1)
		if err != nil {
			panic(err)
		}
		body := string(dataByte)
		return uri + body
	} else {
		panic(errors.New("method not is post or get"))
	}
}

func (c *ClientOkex) GetUrlHeadPathBodyBatch(uri string, params []map[string]string) string {
	dataByte, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	body := string(dataByte)
	return uri + body
}

func (c *ClientOkex) DoRequest(uri base.ReqUrlInfo, method string, params url.Values, result interface{}, batch ...map[string]string) error {
	var (
		err  error
		rsp  []byte
		sign string
	)

	if err = WeightAllow(c.ReqUrl.GetURLRateLimit(uri, params)); err != nil {
		return err
	}
	header := &http.Header{}
	//loc, err := time.LoadLocation("English")

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	header.Add("OK-ACCESS-KEY", c.AccessKey)
	header.Add("OK-ACCESS-PASSPHRASE", c.Passphrase)
	header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	if c.IsTest {
		header.Add("x-simulated-trading", "1")
	}
	header.Add("Content-Type", "application/json")

	pathAndBody := c.GetUrlHeadPathBody(method, uri.Url, params)
	sign = ComputeHmacSha256(timestamp+method+pathAndBody, c.SecretKey)
	header.Add("OK-ACCESS-SIGN", sign)
	//fmt.Println(c.GetUrl(uri.Url), method, header, params)
	rsp, err = conn.Request(c.HttpClient, c.GetUrl(uri.Url), method, header, params)
	//fmt.Println(string(rsp), err)
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
		if v.Code == 429 {
			ClearCapacity(c.ReqUrl.GetURLRateLimit(uri, params))
		}
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

// BatchDoRequest
// @Description: ok批量处理api单独处理, 只能用于POST
// @receiver c
// @param uri
// @param method
// @param result
// @param batch
// @return error
func (c *ClientOkex) BatchDoRequest(uri base.ReqUrlInfo, method string, result interface{}, batch []map[string]string) error {
	var (
		err  error
		rsp  []byte
		sign string
	)

	header := &http.Header{}
	timestamp := time.Now().Add(-time.Hour * 8).Format("2006-01-02T15:04:05.000Z")

	pathAndBody := c.GetUrlHeadPathBodyBatch(uri.Url, batch)
	sign = ComputeHmacSha256(timestamp+method+pathAndBody, c.SecretKey)
	header.Add("OK-ACCESS-SIGN", sign)
	header.Add("OK-ACCESS-KEY", c.AccessKey)
	header.Add("OK-ACCESS-PASSPHRASE", c.Passphrase)
	header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	if c.IsTest {
		header.Add("x-simulated-trading", "1")
	}
	header.Add("Content-Type", "application/json")

	rsp, err = conn.BatchRequest(c.HttpClient, c.GetUrl(uri.Url), header, batch)
	if err == nil {
		err = json.Unmarshal(rsp, result)
		logger.Logger.Debug(uri, err, batch, result)
		return err
	}
	logger.Logger.Debug(uri, batch, err)
	if v, ok := err.(*conn.HttpError); ok {
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

func (c *ClientOkex) Instrument_Info(instType string, options *url.Values) (*RespInstruments, error) {
	params := url.Values{}
	params.Add("instType", instType)
	res := &RespInstruments{}
	err := c.DoRequest(c.ReqUrl.INSTRUMENT_URL, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Delivery_Exercise_History_Info(instType string, uly string, options *url.Values) (*Delivery_Exercise_History, error) {
	params := url.Values{}
	params.Add("instType", instType)
	params.Add("uly", uly)
	ParseOptions(options, &params)
	res := &Delivery_Exercise_History{}
	err := c.DoRequest(c.ReqUrl.DELIVERY_EXERCISE_HISTORY_URL, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Public_ServerTime_Info(options *url.Values) (*Resp_ServerTime, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_ServerTime{}
	err := c.DoRequest(c.ReqUrl.PUBLIC_SERVER_TIME, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Public_MarkPrice_Info(instType string, options *url.Values) (*Resp_MarkPrice, error) {
	params := url.Values{}
	params.Add("instType", instType)
	ParseOptions(options, &params)
	res := &Resp_MarkPrice{}
	err := c.DoRequest(c.ReqUrl.PUBLIC_MARK_PRICE, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Open_Interest_Info(instType string, options *url.Values) (*Resp_Open_Interest, error) {
	params := url.Values{}
	params.Add("instType", instType)
	ParseOptions(options, &params)
	res := &Resp_Open_Interest{}
	err := c.DoRequest(c.ReqUrl.OPEN_INTEREST_URL, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Market_Books_Info(instld string, options *url.Values) (*Resp_Market_Books, error) {
	params := url.Values{}
	params.Add("instId", instld)
	ParseOptions(options, &params)
	res := &Resp_Market_Books{}
	err := c.DoRequest(c.ReqUrl.MARKET_BOOKS, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Market_Trades_Info(instld string, options *url.Values) (*Resp_Market_Trades, error) {
	params := url.Values{}
	params.Add("instId", instld)
	ParseOptions(options, &params)
	res := &Resp_Market_Trades{}
	err := c.DoRequest(c.ReqUrl.MARKET_TRADES, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Market_IndexTickers_Info(instld string, options *url.Values) (*Resp_Market_IndexTickers, error) {
	params := url.Values{}
	params.Add("instId", instld)
	ParseOptions(options, &params)
	res := &Resp_Market_IndexTickers{}
	err := c.DoRequest(c.ReqUrl.MARKET_INDEX_TICKERS, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Market_HistoryTrades_Info(instld string, options *url.Values) (*Resp_Market_HistoryTrades, error) {
	params := url.Values{}
	params.Add("instId", instld)
	ParseOptions(options, &params)
	res := &Resp_Market_HistoryTrades{}
	err := c.DoRequest(c.ReqUrl.MARKET_HISTORY_TRADES, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Account_Balance_Info(options *url.Values) (*Resp_Accout_Balance, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Accout_Balance{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_BALANCE_INFO, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Asset_Balances_Info(options *url.Values) (*Resp_Asset_Balances, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Asset_Balances{}
	err := c.DoRequest(c.ReqUrl.ASSET_BALANCES_INFO, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Account_Positions_Info(options *url.Values) (*Resp_Account_Position, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_Position{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_POSITION, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Account_PositionsHistory_Info(options *url.Values) (*Resp_Account_PositionsHistory, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_PositionsHistory{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_POSITION, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Account_AccountPositionRisk_Info(options *url.Values) (*Resp_Account_AccountPositionRisk, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_AccountPositionRisk{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_POSITION, "GET", params, res)

	return res, err
}

// 账户流水 近7天
func (c *ClientOkex) Account_Bills_Info(options *url.Values) (*Resp_Account_Bills, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_Bills{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_ACCOUNT_BILLS, "GET", params, res)

	return res, err
}

// 账户流水 近三月
func (c *ClientOkex) Account_BillsArchive_Info(options *url.Values) (*Resp_Account_BillsArchive, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_BillsArchive{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_ACCOUNT_BILLS_ARCHIVE, "GET", params, res)

	return res, err
}

// 账户配置
// todo: Service temporarily unavailable, please try again later.
func (c *ClientOkex) Account_Config_Info(options *url.Values) (*Resp_Account_Config, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_Config{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_CONFIG, "GET", params, res)

	return res, err
}

// todo invalid
// 持仓模式
func (c *ClientOkex) Account_SetPositionMode_Info(posMode string, options *url.Values) (*Resp_Account_SetPositionMode, error) {
	params := url.Values{}
	params.Add("posMode", posMode)
	ParseOptions(options, &params)
	res := &Resp_Account_SetPositionMode{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_SET_POSITION_MODE, "POST", params, res)

	return res, err
}

func (c *ClientOkex) Account_SetLeverage_Info(instId string, lever string, mgnMode string, options *url.Values) (*Resp_Account_SetLeverage, error) {
	params := url.Values{}

	params.Add("instId", instId)
	params.Add("lever", lever)
	params.Add("mgnMode", mgnMode)
	ParseOptions(options, &params)
	res := &Resp_Account_SetLeverage{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_SET_LEVERAGE, "POST", params, res)

	return res, err
}

func (c *ClientOkex) Account_MaxSize_Info(instId string, tdMode string, options *url.Values) (*Resp_Account_MaxSize, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("tdMode", tdMode)
	ParseOptions(options, &params)
	res := &Resp_Account_MaxSize{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_MAX_SIZE, "GET", params, res)

	return res, err
}

func (c *ClientOkex) Account_MaxAvailSize_Info(instId string, tdMode string, options *url.Values) (*Resp_Account_MaxAvailSize, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("tdMode", tdMode)
	ParseOptions(options, &params)
	res := &Resp_Account_MaxAvailSize{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_MAX_AVAIL_SIZE, "GET", params, res)
	return res, err
}

// 模拟账户不支持
func (c *ClientOkex) Account_PositionMarginBalance_Info(instId, posSide, typew, amt string, options *url.Values) (*Resp_Account_PositionMarginBalance, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("posSide", posSide)
	params.Add("type", typew)
	params.Add("amt", amt)
	ParseOptions(options, &params)
	res := &Resp_Account_PositionMarginBalance{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_POSITION_MATGIN_BALANCE, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Account_LeverageInfo_Info(instId string, mgnMode string, options *url.Values) (*Resp_Account_LeverageInfo, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("mgnMode", mgnMode)
	ParseOptions(options, &params)
	res := &Resp_Account_LeverageInfo{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_MAX_AVAIL_SIZE, "GET", params, res)
	return res, err
}

// Operation is not supported under the current account mode
func (c *ClientOkex) Account_MaxLoan_Info(instId string, mgnMode string, options *url.Values) (*Resp_Account_MaxLoan, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("mgnMode", mgnMode)
	ParseOptions(options, &params)
	res := &Resp_Account_MaxLoan{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_MAX_LOAN, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Account_TradeFee_Info(instType string, options *url.Values) (*Resp_Account_TradeFee, error) {
	params := url.Values{}
	params.Add("instType", instType)
	ParseOptions(options, &params)
	res := &Resp_Account_TradeFee{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_TRADE_FEE, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Account_InterestAccrued_Info(options *url.Values) (*Resp_Account_InterestAccrued, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_InterestAccrued{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_INTEREST_ACCRUED, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Account_InterestRate_Info(options *url.Values) (*Resp_Account_InterestRate, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_InterestRate{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_INTEREST_RATE, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Account_SetGreeks_Info(greeksType string, options *url.Values) (*Resp_Account_SetGreeks, error) {
	params := url.Values{}
	params.Add("greeksType", greeksType)
	ParseOptions(options, &params)
	res := &Resp_Account_SetGreeks{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_SET_GREEKS, "POST", params, res)
	return res, err
}

// Operation is not supported under the current account mode
func (c *ClientOkex) Account_SetIsolatedMode_Info(isoMode, typew string, options *url.Values) (*Resp_Account_SetIsolatedMode, error) {
	params := url.Values{}
	params.Add("isoMode", isoMode)
	params.Add("type", typew)
	ParseOptions(options, &params)
	res := &Resp_Account_SetIsolatedMode{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_SET_ISOLATED_MODE, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Account_MaxWithdrawal_Info(options *url.Values) (*Resp_Account_MaxWithdrawal, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_MaxWithdrawal{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_MAX_WITHDRAWAL, "GET", params, res)
	return res, err
}

// Operation is not supported under the current account mode
func (c *ClientOkex) Account_RiskState_Info(options *url.Values) (*Resp_Account_RiskState, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_RiskState{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_RISK_STATE, "GET", params, res)
	return res, err
}

// Your account does not support VIP loan
// Parameter ccy can not be empty
func (c *ClientOkex) Account_BorrowRepay_Info(ccy, side, amt string, options *url.Values) (*Resp_Account_BorrowRepay, error) {
	params := url.Values{}
	params.Add("ccy", ccy)
	params.Add("side", side)
	params.Add("amt", amt)
	ParseOptions(options, &params)
	res := &Resp_Account_BorrowRepay{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_BORROW_REPAY, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Account_BorrowRepayHistory_Info(options *url.Values) (*Resp_Account_BorrowRepayHistory, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_BorrowRepayHistory{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_BORROW_REPAY_HISTORY, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Account_InterestLimits_Info(options *url.Values) (*Resp_Account_InterestLimit, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_InterestLimit{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_INTEREST_LIMITS, "GET", params, res)
	return res, err
}

// 一个post字段多个内容怎么做
func (c *ClientOkex) Account_SimulatedMargin_Info(options *url.Values) (*Resp_Account_SimulatedMargin, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Account_SimulatedMargin{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_SIMULATED_MARGIN, "POST", params, res)
	return res, err
}

// todo 错误
// unsupported operation
func (c *ClientOkex) Account_Greeks_Info(options *url.Values) (*Resp_Account_Greeks, error) {
	params := url.Values{}
	params.Add("ccy", "BTC")
	ParseOptions(options, &params)
	res := &Resp_Account_Greeks{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_GREEKS, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Account_PositionTiers_Info(instType, uly string, options *url.Values) (*Resp_Account_PositionTiers, error) {
	params := url.Values{}
	params.Add("instType", instType)
	params.Add("uly", uly)
	ParseOptions(options, &params)
	res := &Resp_Account_PositionTiers{}
	err := c.DoRequest(c.ReqUrl.ACCOUNT_POSITION_TIERS, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Trade_Order_Info(instId, tdMode, side, ordType, sz string, options *url.Values) (*Resp_Trde_Order, error) {
	params := url.Values{}
	instId = strings.ReplaceAll(instId, "/", "-")
	params.Add("tgtCcy", "base_ccy")
	params.Add("instId", instId)
	params.Add("tdMode", tdMode)
	params.Add("side", side)
	params.Add("ordType", strings.ToLower(ordType))
	params.Add("sz", sz)
	ParseOptions(options, &params)
	res := &Resp_Trde_Order{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDER, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Trade_BatchOrders_Info(batch []map[string]string) (*Resp_Trde_Order, error) {
	if len(batch) > 20 {
		err := errors.New("should be Should be less than 10 orders")
		panic(err)
	}
	res := &Resp_Trde_Order{}
	err := c.BatchDoRequest(c.ReqUrl.TRAD_BATCH_ORDERS, "POST", res, batch)
	return res, err
}

//	Trade_CancelOrder
//	@Description:撤销之前未完成的订单
//	@receiver c
//	@param instId
//	@param options : 订单id，ordId和clOrdId必须传一个，若传两个，以ordId为主
//	@return *Resp_Trade_CancelOrder
//	@return error
//
// Deprecated
func (c *ClientOkex) Trade_CancelOrder_Info(instId string, options *url.Values) (*Resp_Trade_CancelOrder, error) {
	params := url.Values{}
	params.Add("instId", instId)
	ParseOptions(options, &params)
	res := &Resp_Trade_CancelOrder{}
	err := c.DoRequest(c.ReqUrl.TRAD_CANCEL_ORDER, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Trade_BatchCancelOrders_Info(batch []map[string]string) (*Resp_Trade_CancelBatchOrders, error) {
	if len(batch) > 20 {
		err := errors.New("should be Should be less than 20 orders")
		panic(err)
	}
	res := &Resp_Trade_CancelBatchOrders{}
	err := c.BatchDoRequest(c.ReqUrl.TRAD_CANCEL_BATCH_ORDERS, "POST", res, batch)
	return res, err
}

func (c *ClientOkex) Trade_ClosePosition(instId string, mgnMode string, options *url.Values) (*Resp_Trade_ClosePosition, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("mgnMode", mgnMode)
	ParseOptions(options, &params)
	res := &Resp_Trade_ClosePosition{}
	err := c.DoRequest(c.ReqUrl.TRAD_CLOSE_POSITION, "POST", params, res)
	return res, err
}

// Trade_OrderInfo_Info
// @Description: 获取订单信息，和下单是同一个地址 method为GET
// @receiver c
// @param instId
// @param mgnMode
// @param options
// @return *Resp_Trade_ClosePosition
// @return error
func (c *ClientOkex) Trade_OrderInfo_Info(instId string, ordId string, options *url.Values) (*Resp_TradeInfo_Info, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("ordId", ordId)
	ParseOptions(options, &params)
	res := &Resp_TradeInfo_Info{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDER, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Trade_OrdersPending_Info(options *url.Values) (*Resp_Trade_OrdersPending, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Trade_OrdersPending{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDERS_PENDING, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Trade_OrdersHistoryWeek_Info(instType string, options *url.Values) (*Resp_Trade_OrdersHistory, error) {
	params := url.Values{}
	params.Add("instType", instType)
	ParseOptions(options, &params)
	res := &Resp_Trade_OrdersHistory{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDERS_HISTORY_WEEK, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Trade_OrdersHistoryArchive_Info(instType string, options *url.Values) (*Resp_Trade_OrdersHistory, error) {
	params := url.Values{}
	params.Add("instType", instType)
	ParseOptions(options, &params)
	res := &Resp_Trade_OrdersHistory{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDERS_HISTORY_ARCHIVE, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Trade_FillsHistoryTreeDays_Info(options *url.Values) (*Resp_Trade_FillsHistory, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Trade_FillsHistory{}
	err := c.DoRequest(c.ReqUrl.TRAD_FILLS_HISTORY_THREE_DAYS, "GET", params, res)
	return res, err
}

// 必须传入instType，和官网不同
func (c *ClientOkex) Trade_FillsHistoryArchive_Info(instType string, options *url.Values) (*Resp_Trade_FillsHistory, error) {
	params := url.Values{}
	params.Add("instType", instType)
	ParseOptions(options, &params)
	res := &Resp_Trade_FillsHistory{}
	err := c.DoRequest(c.ReqUrl.TRAD_FILLS_HISTORY_ARCHIVE, "GET", params, res)
	return res, err
}

// Either parameter slTriggerPx or tpTriggerPx is required
func (c *ClientOkex) Trade_OrderAlgo_Info(instId, tdMode, side, ordType, sz string, options *url.Values) (*Resp_Trade_OrderAlgo, error) {
	params := url.Values{}
	params.Add("instId", instId)
	params.Add("tdMode", tdMode)
	params.Add("side", side)
	params.Add("ordType", ordType)
	params.Add("sz", sz)
	ParseOptions(options, &params)
	res := &Resp_Trade_OrderAlgo{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDER_ALGO, "POST", params, res)
	return res, err
}

// 每次最多可以撤销10个策略委托单
// Service temporarily unavailable, please try again later
func (c *ClientOkex) Trade_CancelAlgo_Info(batch []map[string]string) (*Resp_CancelAlgo, error) {
	if len(batch) > 10 {
		err := errors.New("should be Should be less than 10 orders")
		panic(err)
	}
	res := &Resp_CancelAlgo{}
	err := c.BatchDoRequest(c.ReqUrl.TRAD_CANCEL_ALGO, "POST", res, batch)
	return res, err
}

// 每次最多可以撤销10个策略委托单
func (c *ClientOkex) Trade_CancelAdvanceAlgo_Info(batch []map[string]string) (*Resp_CancelAlgo, error) {
	if len(batch) > 10 {
		err := errors.New("should be Should be less than 10 orders")
		panic(err)
	}
	res := &Resp_CancelAlgo{}
	err := c.BatchDoRequest(c.ReqUrl.TRAD_CANCLE_ADVANCE_ALGOS, "POST", res, batch)
	return res, err
}

func (c *ClientOkex) Trade_OrdersAlgoPending_Info(orderType string, options *url.Values) (*Resp_OrdersAlgoPending, error) {
	params := url.Values{}
	params.Add("ordType", orderType)
	ParseOptions(options, &params)
	res := &Resp_OrdersAlgoPending{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDERS_ALGO_PENDING, "GET", params, res)
	return res, err
}

// state和algoId必填且只能填其一
func (c *ClientOkex) Trade_OrdersAlgoHistory_Info(ordType string, options *url.Values) (*Resp_OrdersAlgoHistory, error) {
	params := url.Values{}
	params.Add("ordType", ordType)
	ParseOptions(options, &params)
	res := &Resp_OrdersAlgoHistory{}
	err := c.DoRequest(c.ReqUrl.TRAD_ORDERS_ALGO_HISTORY, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Rfq_CounterParties_Info(options *url.Values) (*Resp_Rfq_CounterParties, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Rfq_CounterParties{}
	err := c.DoRequest(c.ReqUrl.RFQ_COUNTERPARTIES, "GET", params, res)
	return res, err
}

// System error, please try again later
func (c *ClientOkex) Asset_Currencies_Info(options *url.Values) (*Resp_Asset_Currencies, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Asset_Currencies{}
	err := c.DoRequest(c.ReqUrl.ASSET_CURRENCIES, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Asset_Transfer_Info(ccy, amt, from, to string, options *url.Values) (*Resp_Asset_Transfer, error) {
	params := url.Values{}
	params.Add("ccy", ccy)
	params.Add("amt", amt)
	params.Add("from", from)
	params.Add("to", to)
	ParseOptions(options, &params)
	res := &Resp_Asset_Transfer{}
	err := c.DoRequest(c.ReqUrl.ASSET_TRANSFER, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Asset_TransferState_Info(transId string, options *url.Values) (*Resp_Asset_Transfer, error) {
	params := url.Values{}
	params.Add("transId", transId)
	ParseOptions(options, &params)
	res := &Resp_Asset_Transfer{}
	err := c.DoRequest(c.ReqUrl.ASSET_TRANSFER_STATE, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Asset_Bills_Info(options *url.Values) (*Resp_Asset_Bills, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Asset_Bills{}
	err := c.DoRequest(c.ReqUrl.ASSET_BILLS, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Asset_SubAccountBills_Info(options *url.Values) (*Resp_Asset_SubAccountBills, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Asset_SubAccountBills{}
	err := c.DoRequest(c.ReqUrl.ASSET_SUBACCOUNT_BILLS, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Asset_Withdrawal_Info(ccy, amt, dest, toAddr, fee string, options *url.Values) (*Resp_Asset_Withdrawal, error) {
	params := url.Values{}
	params.Add("ccy", ccy)
	params.Add("amt", amt)
	params.Add("dest", dest)
	params.Add("toAddr", toAddr)
	params.Add("fee", fee)
	ParseOptions(options, &params)
	res := &Resp_Asset_Withdrawal{}
	err := c.DoRequest(c.ReqUrl.ASSET_WITHDRAWAL, "POST", params, res)
	return res, err
}

func (c *ClientOkex) User_SetTransferOut(options *url.Values) (*Resp_User_SubAccountSetTransferOut, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_User_SubAccountSetTransferOut{}
	err := c.DoRequest(c.ReqUrl.USER_SET_TRANSGER_OUT, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Asset_WithdrawalLightning_Info(ccy, invoice string, options *url.Values) (*Resp_Asset_WithdrawalLightning, error) {
	params := url.Values{}
	params.Add("ccy", ccy)
	params.Add("invoice", invoice)
	ParseOptions(options, &params)
	res := &Resp_Asset_WithdrawalLightning{}
	err := c.DoRequest(c.ReqUrl.ASSET_WITHDRAWAL_LIGHTNING, "POST", params, res)
	return res, err
}

func (c *ClientOkex) Asset_CancleWithdrawal_Info(wdId string, options *url.Values) (*Resp_Asset_CancelWithdrawal, error) {
	params := url.Values{}
	params.Add("wdId", wdId)
	ParseOptions(options, &params)
	res := &Resp_Asset_CancelWithdrawal{}
	err := c.DoRequest(c.ReqUrl.ASSET_CANCEL_WITHDRAWAL, "POST", params, res)
	return res, err
}

// Internal Server Error
func (c *ClientOkex) Asset_WithdrawalHistory_Info(options *url.Values) (*Resp_Asset_WithdrawalHistory, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Asset_WithdrawalHistory{}
	err := c.DoRequest(c.ReqUrl.ASSET_WITHDRAWAL_HISTORY, "GET", params, res)
	return res, err
}

func (c *ClientOkex) Asset_DepositHistory_Info(options *url.Values) (*Resp_Asset_DepositHistory, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &Resp_Asset_DepositHistory{}
	err := c.DoRequest(c.ReqUrl.ASSET_DEPOSIT_HISTORY, "GET", params, res)
	return res, err
}
