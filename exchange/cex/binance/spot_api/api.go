package spot_api

import (
	"bytes"
	"clients/config"
	"clients/conn"
	"clients/crypto"
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/warmplanet/proto/go/client"
)

func ParseOptions(options *url.Values, params *url.Values) {
	if options != nil {
		for key := range *options {
			if options.Get(key) != "" {
				params.Add(key, options.Get(key))
			}
		}
	}
}

var (
	// BinanceSpotWeight 定义全局的限制
	BinanceSpotWeight                                *base.RateLimitMgr
	BinanceUBaseWeight                               *base.RateLimitMgr
	BinanceCBaseWeight                               *base.RateLimitMgr
	globalSpotOnce, GlobalUBaseOnce, GlobalCBaseOnce sync.Once
)

type ApiClient struct {
	base.APIConf
	//info *RespExchangeInfo
	HttpClient    *http.Client
	ReqUrl        *ReqUrl
	GetSymbolName func(string) string
	WeightInfo    map[client.WeightType]*client.WeightInfo
	WeightMgr     *base.RateLimitMgr
	lock          sync.Mutex
}

func NewApiClient(conf base.APIConf, maps ...interface{}) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
		}
		proxyUrl      *url.URL
		transport     = http.Transport{}
		err           error
		weightInfoMap map[string]int64
	)
	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
		if conf.IsTest {
			a.EndPoint = SPOT_TEST_API_BASE_URL
		}
	}
	a.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
	globalSpotOnce.Do(func() {
		if BinanceSpotWeight == nil {
			BinanceSpotWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = BinanceSpotWeight
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
			weightInfoMap = t.Spot

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
	a.HttpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(conf.ReadTimeout) * time.Second,
	}
	a.ReqUrl = NewSpotReqUrl(weightInfoMap, a.SubAccount)
	a.GetSymbolName = GetSpotSymbolName
	return a
}

func NewApiClient2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ApiClient {
	var (
		weightInfoMap map[string]int64
		a             = &ApiClient{APIConf: conf}
	)

	if conf.EndPoint != "" {
		a.EndPoint = conf.EndPoint
	} else {
		a.EndPoint = SPOT_API_BASE_URL
	}
	a.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
	globalSpotOnce.Do(func() {
		if BinanceSpotWeight == nil {
			BinanceSpotWeight = base.NewRateLimitMgr()
		}
	})
	a.WeightMgr = BinanceSpotWeight
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
			weightInfoMap = t.Spot
		}
	}
	a.HttpClient = cli
	a.ReqUrl = NewSpotReqUrl(weightInfoMap, a.SubAccount)
	a.GetSymbolName = GetSpotSymbolName
	return a
}

func (a *ApiClient) GetUrl(url string) string {
	return a.EndPoint + url
}

func (a *ApiClient) WeightSet(header *http.Header) {
	if header == nil {
		return
	}
	weight := header.Get("X-MBX-USED-WEIGHT-1M")
	if v, ok := a.WeightInfo[client.WeightType_REQUEST_WEIGHT]; ok {
		atomic.StoreInt64(&v.Value, transform.StringToX[int64](weight).(int64))
	}
	count := header.Get("X-MBX-ORDER-COUNT-1M")
	if v, ok := a.WeightInfo[client.WeightType_ORDERS]; ok {
		atomic.StoreInt64(&v.Value, transform.StringToX[int64](count).(int64))
	}
}

func (a *ApiClient) WeightAllow(endpoint string, consumeMap []*base.RateLimitConsume) error {
	for _, rate := range consumeMap {
		if err := a.WeightMgr.Consume(*rate); err != nil {
			return err
		}
	}
	return nil
}

func (a *ApiClient) WeightUpdate(endpoint string, header *http.Header) {
	if header == nil {
		return
	}
	var (
		count string
		used  int64
	)
	for limitHeader, type_ := range BinanceWeightHeaderMap {
		count = header.Get(limitHeader)
		if count == "" {
			continue
		}
		used = transform.StringToX[int64](count).(int64)
		var ins *base.RateLimitInfo
		if strings.HasPrefix(type_, "account_") {
			ins = a.WeightMgr.GetInstance(base.RateLimitType(type_ + "_" + a.SubAccount))
		} else {
			ins = a.WeightMgr.GetInstance(base.RateLimitType(type_))
		}
		if ins != nil {
			total, _ := ins.Instance.Remain()
			ins.Instance.Set(transform.Min(used, total))
		} else {
			logger.Logger.Warnf("get weight instance err:%s %s %s %s", type_, a.SubAccount, count, limitHeader)
		}
	}
}

func (a *ApiClient) DoRequest(uri base.ReqUrlInfo, method string, params url.Values, result interface{}) error {
	header := &http.Header{}
	header.Add("X-MBX-APIKEY", a.AccessKey)
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	url_ := a.GetUrl(uri.Url)
	if err := a.WeightAllow(a.EndPoint, a.ReqUrl.GetURLRateLimit(uri, params)); err != nil {
		return err
	}
	//fmt.Println(url_, method, header, params)
	respHeader, rsp, err := conn.RequestWithHeader(a.HttpClient, url_, method, header, params)
	//fmt.Println(respHeader, string(rsp), err)
	a.WeightSet(respHeader)
	a.WeightUpdate(a.EndPoint, respHeader)
	//fmt.Println("header:", respHeader.Get("X-Mbx-Used-Weight"), respHeader.Get("X-Mbx-Used-Weight-1m"))
	if err == nil {
		// binance错误返回统一以{"code":开始
		if bytes.HasPrefix(rsp, []byte("{\"code\"")) {
			var re RespError
			err = json.Unmarshal(rsp, &re)
			if err != nil || (re.Code != 0 && re.Msg != "") {
				unknown := false
				//-1006 UNEXPECTED_RESP 从消息总线收到意外响应。 执行状态未知。
				//-1007 TIMEOUT 等待后端服务器响应超时。 发送状态未知； 执行状态未知。
				if re.Code == -1006 || re.Code == -1007 {
					unknown = true
				}
				logger.Logger.Debug(url_, params, re)
				return &base.ApiError{Code: 200, BizCode: re.Code, ErrMsg: re.Msg, UnknownStatus: unknown}
			}
		}
		err = json.Unmarshal(rsp, result)
		logger.Logger.Info("do request: ", url_, params, err)
		return err
	}
	logger.Logger.Debug(url_, params, err, string(rsp))
	// err not nil
	if v, ok := err.(*conn.HttpError); ok {
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: err.Error()}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

func (a *ApiClient) Signature(postForm *url.Values) error {
	postForm.Set("recvWindow", "60000")
	nonce := strconv.FormatInt(time.Now().UnixMilli()+100, 10)
	postForm.Set("timestamp", nonce)
	payload := transform.Unquote(postForm.Encode())
	sign, err := crypto.GetParamHmacSHA256Sign(a.SecretKey, payload)
	postForm.Set("signature", sign)
	return err
}

func (a *ApiClient) SignatureQuote(postForm *url.Values) error {
	postForm.Set("recvWindow", "60000")
	nonce := strconv.FormatInt(time.Now().UnixMilli()+100, 10)
	postForm.Set("timestamp", nonce)
	payload := postForm.Encode()
	sign, err := crypto.GetParamHmacSHA256Sign(a.SecretKey, payload)
	postForm.Set("signature", sign)
	return err
}

func (a *ApiClient) ServerTime() (*RespTime, error) {
	res := &RespTime{}
	err := a.DoRequest(a.ReqUrl.TIME_URL, "GET", url.Values{}, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApiClient) Ping() bool {
	res := &RespPing{}
	err := a.DoRequest(a.ReqUrl.PING_URL, "GET", url.Values{}, res)
	return err == nil
}

func (a *ApiClient) GetDepth(symbol string, limit int) (*RespDepth, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("limit", strconv.Itoa(limit))
	res := &RespDepth{}
	err := a.DoRequest(a.ReqUrl.DEPTH_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *ApiClient) TickerPrice(options *url.Values) (*RespTickerPrice, error) {
	/*
		最新价格
		symbol	STRING	NO	交易对
	*/
	var (
		params = url.Values{}
		res    *RespTickerPrice
		err    error
	)
	ParseOptions(options, &params)
	if options == nil || options.Get("symbol") == "" {
		res = &RespTickerPrice{}
		err = a.DoRequest(a.ReqUrl.PREMIUMINDEX_URL, "GET", params, res)
	} else {
		resp := &RespTickerPriceItem{}
		err = a.DoRequest(a.ReqUrl.PREMIUMINDEX_URL, "GET", params, resp)
		res = &RespTickerPrice{resp}
	}
	return res, err
}

func (a *ApiClient) GetTrades(symbol string, limit int) (*RespGetTrade, error) {
	params := url.Values{}
	params.Add("symbol", a.GetSymbolName(symbol))
	params.Add("limit", strconv.Itoa(limit))
	res := &RespGetTrade{}
	err := a.DoRequest(a.ReqUrl.TRADE_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApiClient) GetHistoricalTrades(symbol string, options *url.Values) (*RespHistoricalTrade, error) {
	/*
		symbol	STRING	YES
		limit	INT	NO	默认 500; 最大值 1000.
		fromId	LONG	NO	从哪一条成交id开始返回. 缺省返回最近的成交记录。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	res := &RespHistoricalTrade{}
	err := a.DoRequest(a.ReqUrl.HISTORICAL_TRADES_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *ApiClient) ExchangeInfo() (*RespExchangeInfo, error) {
	res := &RespExchangeInfo{}
	err := a.DoRequest(a.ReqUrl.EXCHANGEINFO_URL, "GET", url.Values{}, res)
	return res, err
}

func (a *ApiClient) CapitalConfigGetAll() (*RespCapitalConfigGetAll, error) {
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespCapitalConfigGetAll{}
	err = a.DoRequest(a.ReqUrl.CAPITAL_CONFIG_GETALL_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) AssetTradeFee() (*RespAssetTradeFee, error) {
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAssetTradeFee{}
	err = a.DoRequest(a.ReqUrl.ASSET_TRADEFEE_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) CapitalWithdrawApply(coin, address string, amount float64, options *url.Values) (*RespCapitalWithdrawApply, error) {
	/*
		coin	STRING	YES
		withdrawOrderId	STRING	NO	自定义提币ID
		network	STRING	NO	提币网络
		address	STRING	YES	提币地址
		addressTag	STRING	NO	某些币种例如 XRP,XMR 允许填写次级地址标签
		amount	DECIMAL	YES	数量
		transactionFeeFlag	BOOLEAN	NO	当站内转账时免手续费, true: 手续费归资金转入方; false: 手续费归资金转出方; . 默认 false.
		name	STRING	NO	地址的备注，填写该参数后会加入该币种的提现地址簿。地址簿上限为20，超出后会造成提现失败。地址中的空格需要encode成%20
		walletType	INTEGER	NO	表示出金使用的钱包，0为现货钱包，1为资金钱包，默认为现货钱包
	*/
	params := url.Values{}
	params.Add("coin", coin)
	params.Add("address", address)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespCapitalWithdrawApply{}
	err = a.DoRequest(a.ReqUrl.CAPITAL_WITHDRAW_APPLY_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) CapitalDepositHisRec(options *url.Values) (*RespCapitalDepositHisRec, error) {
	/*
		coin	STRING	NO
		status	INT	NO	0(0:pending,6: credited but cannot withdraw, 1:success)
		startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime	LONG	NO	默认当前时间戳
		offset	INT	NO	默认:0
		limit	INT	NO	默认：1000，最大1000
		recvWindow	LONG	NO
		timestamp	LONG	YES
		请注意startTime 与 endTime 的默认时间戳，保证请求时间间隔不超过90天.
		同时提交startTime 与 endTime间隔不得超过90天.
	*/
	params := url.Values{}
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespCapitalDepositHisRec{}
	err = a.DoRequest(a.ReqUrl.CAPITAL_DEPOSIT_HISREC_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) CapitalWithdrawHistory(options *url.Values) (*RespCapitalWithdrawHistory, error) {
	/*
		coin	STRING	NO
		withdrawOrderId	STRING	NO
		status	INT	NO	0(0:已发送确认Email,1:已被用户取消 2:等待确认 3:被拒绝 4:处理中 5:提现交易失败 6 提现完成)
		offset	INT	NO
		limit	INT	NO	默认：1000， 最大：1000
		startTime	LONG	NO	默认当前时间90天前的时间戳
		endTime	LONG	NO	默认当前时间戳
		recvWindow	LONG	NO
		timestamp	LONG	YES
		支持多网络提币前的历史记录可能不会返回network字段.
		请注意startTime 与 endTime 的默认时间戳，保证请求时间间隔不得超过90天.
		同时提交startTime 与 endTime间隔不得超过90天.
		若传了withdrawOrderId，则请求的startTime 与 endTime的时间间隔不得超过7天.
		若传了withdrawOrderId，没传startTime 与 endTime，则默认返回最近7天数据.
	*/
	params := url.Values{}
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespCapitalWithdrawHistory{}
	err = a.DoRequest(a.ReqUrl.CAPITAL_WITHDRAW_HISTORY_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) AssetTransfer(type_ MoveType, asset string, amount float64, options *url.Values) (*RespAssetTransfer, error) {
	/*
		type	ENUM	YES
		asset	STRING	YES
		amount	DECIMAL	YES
		fromSymbol	STRING	NO  //fromSymbol 必须要发送，当类型为 ISOLATEDMARGIN_MARGIN 和 ISOLATEDMARGIN_ISOLATEDMARGIN
		toSymbol	STRING	NO  //toSymbol 必须要发送，当类型为 MARGIN_ISOLATEDMARGIN 和 ISOLATEDMARGIN_ISOLATEDMARGIN
		recvWindow	LONG	NO
		timestamp	LONG	YES
	*/
	params := url.Values{}
	params.Add("type", string(type_))
	params.Add("asset", asset)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	ParseOptions(options, &params)
	if type_ == MOVE_TYPE_ISOLATEDMARGIN_MARGIN || type_ == MOVE_TYPE_ISOLATEDMARGIN_ISOLATEDMARGIN {
		if options.Get("fromSymbol") == "" {
			return nil, errors.New("fromSymbol 必须要发送")
		}
	}
	if type_ == MOVE_TYPE_MARGIN_ISOLATEDMARGIN || type_ == MOVE_TYPE_ISOLATEDMARGIN_ISOLATEDMARGIN {
		if options.Get("toSymbol") == "" {
			return nil, errors.New("toSymbol 必须要发送")
		}
	}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAssetTransfer{}
	err = a.DoRequest(a.ReqUrl.ASSET_TRANSFER_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) SubAccountUniversalTransfer(fromAccountType, toAccountType MoveAccountType, asset string, amount float64, options *url.Values) (*RespAssetTransfer, error) {
	/*
		fromEmail	STRING	NO
		toEmail	STRING	NO
		fromAccountType	STRING	YES	"SPOT","USDT_FUTURE","COIN_FUTURE","MARGIN"(Cross),"ISOLATED_MARGIN"
		toAccountType	STRING	YES	"SPOT","USDT_FUTURE","COIN_FUTURE","MARGIN"(Cross),"ISOLATED_MARGIN"
		clientTranId	STRING	NO	不可重复
		symbol	STRING	NO	仅在ISOLATED_MARGIN类型下使用
		asset	STRING	YES
		amount	DECIMAL	YES
		recvWindow	LONG	NO
		timestamp	LONG	YES

		需要开启母账户apikey“允许子母账户划转”权限。
		若 fromEmail 未传，默认从母账户转出。
		若 toEmail 未传，默认转入母账户。
		该接口支持的划转操作有：
		现货账户划转到现货账户、U本位合约账户、币本位合约账户（无论母账户或子账户）
		现货账户、U本位合约账户、币本位合约账户划转到现货账户（无论母账户或子账户）
		母账户现货账户划转到子账户杠杆全仓账户、杠杆逐仓账户
		子账户杠杆全仓账户、杠杆逐仓账户划转到母账户现货账户
	*/
	params := url.Values{}
	params.Add("fromAccountType", string(fromAccountType))
	params.Add("toAccountType", string(toAccountType))
	params.Add("asset", asset)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAssetTransfer{}
	err = a.DoRequest(a.ReqUrl.SUBACCOUNT_UNIVERSALTRANSFER_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) GetSubAccountUniversalTransfer(options *url.Values) (*RespSubAccountUniversalTransfer, error) {
	/*
		fromEmail	STRING	NO
		toEmail	STRING	NO
		clientTranId	STRING	NO
		startTime	LONG	NO
		endTime	LONG	NO
		page	INT	NO	默认 1
		limit	INT	NO	默认 500, 最大 500
		recvWindow	LONG	NO
		timestamp	LONG	YES

		本查询接口只可以单边查询，fromEmail 和 toEmail 不能同时传入。
		若 fromEmail 和 toEmail 都未传，默认返回 fromEmail 为母账户的划转记录。
		若 startTime 和 endTime 都未传，则只可查询最近30天的记录。
		查询时间范围最大不得超过30天。
	*/
	params := url.Values{}
	ParseOptions(options, &params)
	err := a.SignatureQuote(&params)
	if err != nil {
		return nil, err
	}
	res := &RespSubAccountUniversalTransfer{}
	err = a.DoRequest(a.ReqUrl.SUBACCOUNT_UNIVERSALTRANSFER_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) SubAccountTransferSubToMaster(asset string, amount float64) (*RespAssetTransfer, error) {
	/*
		asset	STRING	YES
		amount	DECIMAL	YES
		recvWindow	LONG	NO
		timestamp	LONG	YES
	*/
	params := url.Values{}
	params.Add("asset", asset)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAssetTransfer{}
	err = a.DoRequest(a.ReqUrl.SUBACCOUNT_TRANSFER_SUBTOMASTER_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) SubAccountTransferSubToSub(toEmail, asset string, amount float64) (*RespAssetTransfer, error) {
	/*
		toEmail	STRING	YES	接收者子邮箱地址 备注
		asset	STRING	YES
		amount	DECIMAL	YES
		recvWindow	LONG	NO
		timestamp	LONG	YES
	*/
	params := url.Values{}
	params.Add("toEmail", toEmail)
	params.Add("asset", asset)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAssetTransfer{}
	err = a.DoRequest(a.ReqUrl.SUBACCOUNT_TRANSFER_SUBTOSUB_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) SubAccountTransferSubUserHistory(options *url.Values) (*RespSubAccountTransferSubUserHistory, error) {
	/*
		asset	STRING	NO	如不提供，返回所有asset 划转记录
		type	INT	NO	1: transfer in, 2: transfer out; 如不提供，返回transfer out方向划转记录
		startTime	LONG	NO
		endTime	LONG	NO
		limit	INT	NO	默认值: 500
		recvWindow	LONG	NO
		timestamp	LONG	YES
	*/
	params := url.Values{}
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespSubAccountTransferSubUserHistory{}
	err = a.DoRequest(a.ReqUrl.SUBACCOUNT_TRANSFER_SUBUSERHISTORY_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) Account() (*RespAccount, error) {
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAccount{}
	err = a.DoRequest(a.ReqUrl.ACCOUNT_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) Asset(options *url.Values) (*RespAsset, error) {
	/*
		asset	STRING	NO
		needBtcValuation	STRING	NO	true or false
		recvWindow	LONG	NO
		timestamp	LONG	YES
	*/
	params := url.Values{}
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAsset{}
	err = a.DoRequest(a.ReqUrl.ASSET_URL, "POST", params, res)
	if err != nil {
		res2 := &SpotAssetItem{}
		if err = a.DoRequest(a.ReqUrl.ASSET_URL, "POST", params, res2); err == nil {
			res = &RespAsset{res2}
		}
	}
	return res, err
}

func (a *ApiClient) MarginAccount() (*RespMarginAccount, error) {
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginAccount{}
	err = a.DoRequest(a.ReqUrl.MARGIN_ACCOUNT_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) MarginIsolatedAccount(symbols ...string) (*RespMarginIsolatedAccount, error) {
	/*
		symbols	STRING	NO	最多可以传5个symbol; 由","分隔的字符串表示. e.g. "BTCUSDT,BNBUSDT,ADAUSDT"
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG	YES
	*/
	params := url.Values{}
	if len(symbols) > 0 {
		params.Add("symbols", strings.Join(symbols, ","))
	}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginIsolatedAccount{}
	err = a.DoRequest(a.ReqUrl.MARGIN_ISOLATED_ACCOUNT_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) Order(symbol string, side SideType, type_ OrderType, options *url.Values) (interface{}, error) {
	/*
		symbol	STRING	YES
		side	ENUM	YES	详见枚举定义：订单方向
		type	ENUM	YES	详见枚举定义：订单类型
		timeInForce	ENUM	NO	详见枚举定义：有效方式
		quantity	DECIMAL	NO
		quoteOrderQty	DECIMAL	NO
		price	DECIMAL	NO
		newClientOrderId	STRING	NO	客户自定义的唯一订单ID。 如果未发送，则自动生成
		stopPrice	DECIMAL	NO	仅 STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, 和TAKE_PROFIT_LIMIT 需要此参数。
		trailingDelta	LONG	NO	用于 STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, 和 TAKE_PROFIT_LIMIT 类型的订单. 更多追踪止盈止损订单细节, 请参考 追踪止盈止损(Trailing Stop)订单常见问题
		icebergQty	DECIMAL	NO	仅使用 LIMIT, STOP_LOSS_LIMIT, 和 TAKE_PROFIT_LIMIT 创建新的 iceberg 订单时需要此参数
		newOrderRespType	ENUM	NO	设置响应JSON。 ACK，RESULT或FULL； "MARKET"和" LIMIT"订单类型默认为"FULL"，所有其他订单默认为"ACK"。
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG

		基于订单 type不同，强制要求某些参数:
		类型	强制要求的参数
		LIMIT	timeInForce, quantity, price
		MARKET	quantity or quoteOrderQty
		STOP_LOSS	quantity, stopPrice 或者 trailingDelta
		STOP_LOSS_LIMIT	timeInForce, quantity, price, stopPrice 或者 trailingDelta
		TAKE_PROFIT	quantity, stopPrice 或者 trailingDelta
		TAKE_PROFIT_LIMIT	timeInForce, quantity, price, stopPrice 或者 trailingDelta
		LIMIT_MAKER	quantity, price

	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("side", string(side))
	params.Add("type", string(type_))
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if params.Get("newOrderRespType") == string(ORDER_RESP_TYPE_ACK) {
		res = &RespMarginOrderAck{}
	} else if params.Get("newOrderRespType") == string(ORDER_RESP_TYPE_RESULT) {
		res = &RespMarginOrderResult{}
	} else if params.Get("newOrderRespType") == string(ORDER_RESP_TYPE_FULL) {
		res = &RespMarginOrderFull{}
	} else {
		if type_ == ORDER_TYPE_MARKET || type_ == ORDER_TYPE_LIMIT {
			res = &RespMarginOrderFull{}
		} else {
			res = &RespMarginOrderAck{}
		}
	}
	err = a.DoRequest(a.ReqUrl.ORDER_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) MarginOrder(symbol string, isIsolated MarginType, side SideType, type_ OrderType, options *url.Values) (interface{}, error) {
	/*
		symbol	STRING	YES
		isIsolated	STRING	NO	是否逐仓杠杆，TRUE, FALSE, 默认 "FALSE"
		side	ENUM	YES	BUY SELL
		type	ENUM	YES	详见枚举定义：订单类型
		quantity	DECIMAL	NO
		quoteOrderQty	DECIMAL	NO
		price	DECIMAL	NO
		stopPrice	DECIMAL	NO	与STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, 和 TAKE_PROFIT_LIMIT 订单一起使用.
		newClientOrderId	STRING	NO	客户自定义的唯一订单ID。若未发送自动生成。
		icebergQty	DECIMAL	NO	与 LIMIT, STOP_LOSS_LIMIT, 和 TAKE_PROFIT_LIMIT 一起使用创建 iceberg 订单.
		newOrderRespType	ENUM	NO	设置响应: JSON. ACK, RESULT, 或 FULL; MARKET 和 LIMIT 订单类型默认为 FULL, 所有其他订单默认为 ACK.
		sideEffectType	ENUM	NO	NO_SIDE_EFFECT, MARGIN_BUY, AUTO_REPAY;默认为 NO_SIDE_EFFECT.
		timeInForce	ENUM	NO	GTC,IOC,FOK
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG	YES
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("isIsolated", string(isIsolated))
	params.Add("side", string(side))
	params.Add("type", string(type_))
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res interface{}
	if params.Get("newOrderRespType") == string(ORDER_RESP_TYPE_ACK) {
		res = &RespMarginOrderAck{}
	} else if params.Get("newOrderRespType") == string(ORDER_RESP_TYPE_RESULT) {
		res = &RespMarginOrderResult{}
	} else if params.Get("newOrderRespType") == string(ORDER_RESP_TYPE_FULL) {
		res = &RespMarginOrderFull{}
	} else {
		if type_ == ORDER_TYPE_MARKET || type_ == ORDER_TYPE_LIMIT {
			res = &RespMarginOrderFull{}
		} else {
			res = &RespMarginOrderAck{}
		}
	}
	err = a.DoRequest(a.ReqUrl.MARGIN_ORDER_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) CancelOrder(symbol string, options *url.Values) (*RespCancelOrder, error) {
	/*
		symbol	STRING	YES
		orderId	LONG	NO
		origClientOrderId	STRING	NO
		newClientOrderId	STRING	NO	用户自定义的本次撤销操作的ID(注意不是被撤销的订单的自定义ID)。如无指定会自动赋值。
		recvWindow	LONG	NO	赋值不得大于 60000
		timestamp	LONG	YES
		orderId 或 origClientOrderId 必须至少发送一个
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespCancelOrder{}
	err = a.DoRequest(a.ReqUrl.ORDER_URL, "DELETE", params, res)
	return res, err
}

func (a *ApiClient) CancelMarginOrder(symbol string, options *url.Values) (*RespCancelOrder, error) {
	/*
		symbol	STRING	YES
		isIsolated	STRING	NO	是否逐仓杠杆，"TRUE", "FALSE", 默认 "FALSE"
		orderId	LONG	NO
		origClientOrderId	STRING	NO
		newClientOrderId	STRING	NO	用于唯一识别此撤销订单，默认自动生成。
		recvWindow	LONG	NO	T赋值不能大于 60000
		timestamp	LONG	YES
		必须发送 orderId 或 origClientOrderId 其中一个。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespCancelOrder{}
	err = a.DoRequest(a.ReqUrl.MARGIN_ORDER_URL, "DELETE", params, res)
	return res, err
}

func (a *ApiClient) GetOrder(symbol string, options *url.Values) (*RespGetOrder, error) {
	/*
		symbol	STRING	YES
		orderId	LONG	NO
		origClientOrderId	STRING	NO
		recvWindow	LONG	NO	赋值不得大于 60000
		timestamp	LONG	YES
		注意:
		至少需要发送 orderId 与 origClientOrderId中的一个
		某些订单中cummulativeQuoteQty<0，是由于这些订单是cummulativeQuoteQty功能上线之前的订单。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	if options.Get("orderId") == "" && options.Get("origClientOrderId") == "" {
		return nil, errors.New("must have one of parameters: orderId or origClientOrderId")
	}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespGetOrder{}
	err = a.DoRequest(a.ReqUrl.ORDER_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) GetMarginOrder(symbol string, options *url.Values) (*RespGetMarginOrder, error) {
	/*
		symbol	STRING	YES
		isIsolated	STRING	NO	是否逐仓杠杆，TRUE, FALSE, 默认 "FALSE"
		orderId	LONG	NO
		origClientOrderId	STRING	NO
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG	YES
		必须发送 orderId 或 origClientOrderId 其中一个。
		一些历史订单的 cummulativeQuoteQty < 0, 是指当前数据不存在。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	if options.Get("orderId") == "" && options.Get("origClientOrderId") == "" {
		return nil, errors.New("must have one of parameters: orderId or origClientOrderId")
	}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespGetMarginOrder{}
	err = a.DoRequest(a.ReqUrl.MARGIN_ORDER_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) AssetTransferHistory(type_ MoveType, options *url.Values) (*RespAssetTransferHistory, error) {
	/*
		type	ENUM	YES
		startTime	LONG	NO
		endTime	LONG	NO
		current	INT	NO	默认 1  当前页面. 起始计数为1. 默认值1
		size	INT	NO	默认 10, 最大 100
		fromSymbol	STRING	NO    fromSymbol 必须要发送，当类型为 ISOLATEDMARGIN_MARGIN 和 ISOLATEDMARGIN_ISOLATEDMARGIN
		toSymbol	STRING	NO    toSymbol 必须要发送，当类型为 MARGIN_ISOLATEDMARGIN 和 ISOLATEDMARGIN_ISOLATEDMARGIN
		recvWindow	LONG	NO
		timestamp	LONG	YES

		仅支持查询最近半年（6个月）数据
		若startTime和endTime没传，则默认返回最近7天数据
	*/
	params := url.Values{}
	params.Add("type", string(type_))
	ParseOptions(options, &params)
	if type_ == MOVE_TYPE_ISOLATEDMARGIN_MARGIN || type_ == MOVE_TYPE_ISOLATEDMARGIN_ISOLATEDMARGIN {
		if options.Get("fromSymbol") == "" {
			return nil, errors.New("fromSymbol 必须要发送")
		}
	}
	if type_ == MOVE_TYPE_MARGIN_ISOLATEDMARGIN || type_ == MOVE_TYPE_ISOLATEDMARGIN_ISOLATEDMARGIN {
		if options.Get("toSymbol") == "" {
			return nil, errors.New("toSymbol 必须要发送")
		}
	}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAssetTransferHistory{}
	err = a.DoRequest(a.ReqUrl.ASSET_TRANSFER_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) MarginLoan(asset string, amount float64, options *url.Values) (*RespMarginLoan, error) {
	/*
		asset	STRING	YES
		isIsolated	STRING	NO	是否逐仓杠杆，TRUE, FALSE, 默认 "FALSE"
		symbol	STRING	NO	逐仓交易对，配合逐仓使用
		amount	DECIMAL	YES
		recvWindow	LONG	NO	赋值不能超过 60000
		timestamp	LONG	YES
		如果 isIsolated = “TRUE”, 表示逐仓借贷，此时 symbol 必填
		如果isIsolated = “FALSE” 表示全仓借贷
	*/
	params := url.Values{}
	params.Add("asset", asset)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	ParseOptions(options, &params)
	if options.Get("isIsolated") == "true" && options.Get("symbol") == "" {
		return nil, errors.New("isolated margin must have symbol parameter")
	}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginLoan{}
	err = a.DoRequest(a.ReqUrl.MARGIN_LOAN_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) MarginRepay(asset string, amount float64, options *url.Values) (*RespMarginRepay, error) {
	/*
		asset	STRING	YES
		isIsolated	STRING	NO	是否逐仓杠杆，TRUE, FALSE, 默认 "FALSE"
		symbol	STRING	NO	逐仓交易对，配合逐仓使用
		amount	DECIMAL	YES
		recvWindow	LONG	NO	赋值不能超过 60000
		timestamp	LONG	YES
	*/
	params := url.Values{}
	params.Add("asset", asset)
	params.Add("amount", strconv.FormatFloat(amount, 'f', -1, 64))
	ParseOptions(options, &params)
	if options.Get("isIsolated") == "true" && options.Get("symbol") == "" {
		return nil, errors.New("isolated margin must have symbol parameter")
	}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginRepay{}
	err = a.DoRequest(a.ReqUrl.MARGIN_REPAY_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) MarginLoanHistory(asset string, options *url.Values) (*RespMarginLoanHistory, error) {
	/*
		asset	STRING	YES
		isolatedSymbol	STRING	NO	逐仓symbol
		txId	LONG	NO	tranId in POST /sapi/v1/margin/loan
		startTime	LONG	NO
		endTime	LONG	NO
		current	LONG	NO	当前查询页。 开始值 1。 默认:1
		size	LONG	NO	默认:10 最大:100
		archived	STRING	NO	默认: false. 查询6个月以前的数据，需要设为 true
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG	YES
		必须发送txId 或 startTime，txId 优先。
		响应返回为降序排列。
		如果发送isolatedSymbol，返回指定逐仓symbol指定asset的借贷记录。
		查询时间范围最大不得超过30天。
		若startTime和endTime没传，则默认返回最近7天数据
		如果想查询6个月以前数据，设置 archived 为 true。
	*/
	params := url.Values{}
	params.Add("asset", asset)
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginLoanHistory{}
	err = a.DoRequest(a.ReqUrl.MARGIN_LOAN_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) MarginRepayHistory(asset string, options *url.Values) (*RespMarginRepayHistory, error) {
	/*
		asset	STRING	YES
		isolatedSymbol	STRING	NO	逐仓symbol
		txId	LONG	NO	返回 /sapi/v1/margin/repay
		startTime	LONG	NO
		endTime	LONG	NO
		current	LONG	NO	当前查询页。开始值 1. 默认:1
		size	LONG	NO	默认:10 最大:100
		archived	STRING	NO	默认: false. 查询6个月以前的数据，需要设为 true
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG	YES
		必须发送txId 或 startTime，txId 优先。
		响应返回为降序排列。
		如果发送isolatedSymbol，返回指定逐仓symbol指定asset的还贷记录。
		查询时间范围最大不得超过30天。
		若startTime和endTime没传，则默认返回最近7天数据
		如果想查询6个月以前数据，设置 archived 为 true。
	*/
	params := url.Values{}
	params.Add("asset", asset)
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginRepayHistory{}
	err = a.DoRequest(a.ReqUrl.MARGIN_REPAY_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) MyTrades(symbol string, options *url.Values) (*RespMyTrades, error) {
	/*
		symbol	STRING	YES
		orderId	LONG	NO	必须要和参数symbol一起使用.
		startTime	LONG	NO
		endTime	LONG	NO
		fromId	LONG	NO	起始Trade id。 默认获取最新交易。
		limit	INT	NO	默认 500; 最大 1000.
		recvWindow	LONG	NO	赋值不能超过 60000
		timestamp	LONG	YES
		注意: 如果设定 fromId , 获取订单 >= fromId. 否则返回最近订单。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMyTrades{}
	err = a.DoRequest(a.ReqUrl.MYTRADES_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) MarginAllOrders(symbol string, options *url.Values) (*RespMarginAllOrders, error) {
	/*
		symbol	STRING	YES
		isIsolated	STRING	NO	是否逐仓杠杆，"TRUE", "FALSE", 默认 "FALSE"
		orderId	LONG	NO
		startTime	LONG	NO
		endTime	LONG	NO
		limit	INT	NO	默认 500;最大500.
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG	YES
		如果设置 orderId , 获取订单 >= orderId， 否则返回近期订单历史。
		一些历史订单的 cummulativeQuoteQty < 0, 是指当前数据不存在。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginAllOrders{}
	err = a.DoRequest(a.ReqUrl.MARGIN_ALLORDERS_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) AllOrders(symbol string, options *url.Values) (*RespAllOrders, error) {
	/*
		symbol	STRING	YES
		orderId	LONG	NO
		startTime	LONG	NO
		endTime	LONG	NO
		limit	INT	NO	默认 500; 最大 1000.
		recvWindow	LONG	NO	赋值不得大于 60000
		timestamp	LONG	YES
		注意:
		如设置 orderId , 订单量将 >= orderId。否则将返回最新订单。
		一些历史订单 cummulativeQuoteQty < 0, 是指数据此时不存在。
		如果设置 startTime 和 endTime, orderId 就不需要设置。
		数据源: 数据库
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespAllOrders{}
	err = a.DoRequest(a.ReqUrl.ALLORDERS_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) OpenOrders(options *url.Values) (*RespOpenOrders, error) {
	/*
		symbol	STRING	NO
		recvWindow	LONG	NO	赋值不得大于 60000
		timestamp	LONG	YES
		不带symbol参数，会返回所有交易对的挂单
	*/
	params := url.Values{}
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespOpenOrders{}
	err = a.DoRequest(a.ReqUrl.OPENORDERS_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) MarginOpenOrders(options *url.Values) (*RespMarginOpenOrders, error) {
	/*
		symbol	STRING	NO
		isIsolated	STRING	NO	是否逐仓杠杆，"TRUE", "FALSE", 默认 "FALSE"
		recvWindow	LONG	NO	赋值不能大于 60000
		timestamp	LONG	YES
		如未发送symbol，返回所有 symbols 订单记录。
		当返回所有symbols时，针对限速器计数的请求数量等于当前在交易所交易的symbols数量。
		如果 isIsolated = "TRUE", symbol 为必填
	*/
	params := url.Values{}
	ParseOptions(options, &params)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	res := &RespMarginOpenOrders{}
	err = a.DoRequest(a.ReqUrl.MARGIN_OPENORDERS_URL, "GET", params, res)
	return res, err
}

func (a *ApiClient) UserDataStream() (*RespUserDataStream, error) {
	/*
		开始一个新的数据流。除非发送 keepalive，否则数据流于60分钟后关闭。
		如果该帐户具有有效的listenKey，则将返回该listenKey并将其有效期延长60分钟。
		权重: 1
	*/
	params := url.Values{}
	//err := a.Signature(&params)
	//if err != nil {
	//	return nil, err
	//}
	res := &RespUserDataStream{}
	err := a.DoRequest(a.ReqUrl.USERDATASTREAM_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) PutUserDataStream(listenKey string) (interface{}, error) {
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

func (a *ApiClient) DELETEUserDataStream(listenKey string) (interface{}, error) {
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

func (a *ApiClient) MarginUserDataStream() (*RespUserDataStream, error) {
	/*
		开始一个新的数据流。除非发送 keepalive，否则数据流于60分钟后关闭。
		如果该帐户具有有效的listenKey，则将返回该listenKey并将其有效期延长60分钟。
		权重: 1
	*/
	params := url.Values{}
	ParseOptions(nil, &params)
	res := &RespUserDataStream{}
	err := a.DoRequest(a.ReqUrl.MARGINUSERDATASTREAM_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) PutMarginUserDataStream(listenKey string) (interface{}, error) {
	/*
		listenKey	STRING	YES
		有效期延长至本次调用后60分钟,建议每30分钟发送一个 ping 。
	*/
	params := url.Values{}
	params.Add("listenKey", listenKey)
	var res interface{}
	err := a.DoRequest(a.ReqUrl.MARGINUSERDATASTREAM_URL, "PUT", params, &res)
	return res, err
}

func (a *ApiClient) DELETEMarginUserDataStream(listenKey string) (interface{}, error) {
	/*
		listenKey	STRING	YES
		关闭用户数据流。
	*/
	params := url.Values{}
	params.Add("listenKey", listenKey)
	var res interface{}
	err := a.DoRequest(a.ReqUrl.MARGINUSERDATASTREAM_URL, "DELETE", params, &res)
	return res, err
}

func (a *ApiClient) MarginIsolatedUserDataStream(symbol string) (*RespUserDataStream, error) {
	/*
		开始一个新的数据流。除非发送 keepalive，否则数据流于60分钟后关闭。
		如果该帐户具有有效的listenKey，则将返回该listenKey并将其有效期延长60分钟。
		权重: 1
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	res := &RespUserDataStream{}
	err := a.DoRequest(a.ReqUrl.MARGINISOLATEDUSERDATASTREAM_URL, "POST", params, res)
	return res, err
}

func (a *ApiClient) PutMarginIsolatedUserDataStream(symbol, listenKey string) (interface{}, error) {
	/*
		listenKey	STRING	YES
		有效期延长至本次调用后60分钟,建议每30分钟发送一个 ping 。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("listenKey", listenKey)
	var res interface{}
	err := a.DoRequest(a.ReqUrl.MARGINISOLATEDUSERDATASTREAM_URL, "PUT", params, &res)
	return res, err
}

func (a *ApiClient) DELETEMarginIsolatedUserDataStream(symbol, listenKey string) (interface{}, error) {
	/*
		listenKey	STRING	YES
		关闭用户数据流。
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("listenKey", listenKey)
	var res interface{}
	err := a.DoRequest(a.ReqUrl.MARGINISOLATEDUSERDATASTREAM_URL, "DELETE", params, &res)
	return res, err
}

func (a *ApiClient) MarginIsolatedPair(symbol string) (*RespIsolatedPair, error) {
	/*
		symbol	STRING	YES
		查询逐仓杠杆交易对 (USER_DATA)
	*/
	params := url.Values{}
	params.Add("symbol", symbol)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res RespIsolatedPair
	err = a.DoRequest(a.ReqUrl.MARGIN_ISOLATED_PAIR_URL, "GET", params, &res)
	return &res, err
}

func (a *ApiClient) MarginIsolatedAllPairs() (RespIsolatedAllPairs, error) {
	/*
		获取所有逐仓杠杆交易对(USER_DATA)
	*/
	params := url.Values{}
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res RespIsolatedAllPairs
	err = a.DoRequest(a.ReqUrl.MARGIN_ISOLATED_ALLPAIRS_URL, "GET", params, &res)
	return res, err
}

func (a *ApiClient) CreateSubAccount(subAccountEmail string) (RespIsolatedAllPairs, error) {
	/*
		subAccountString	STRING	YES	请输入字符串，我们将为您创建一个虚拟邮箱进行注册
	*/
	params := url.Values{}
	params.Add("subAccountEmail", subAccountEmail)
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res RespIsolatedAllPairs
	err = a.DoRequest(a.ReqUrl.SUBACCOUNT_VIRTUALSUBACCOUNT_URL, "GET", params, &res)
	return res, err
}

func (a *ApiClient) SubAccountList(options *url.Values) (*RespSubAccountList, error) {
	/*
		email	STRING	NO	Sub-account email
		isFreeze	STRING	NO	true or false
		page	INT	NO	默认: 1
		limit	INT	NO	默认: 1, 最大: 200
		recvWindow	LONG	NO
		timestamp	LONG	YES
	*/
	params := url.Values{}
	params.Add("limit", "200")
	err := a.Signature(&params)
	if err != nil {
		return nil, err
	}
	var res RespSubAccountList
	err = a.DoRequest(a.ReqUrl.SUBACCOUNT_LIST_URL, "GET", params, &res)
	return &res, err
}
