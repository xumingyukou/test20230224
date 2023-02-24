package spot_api

import (
	"bytes"
	"clients/conn"
	"clients/crypto"
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SpotApiClient To do CHECK
type SpotApiClient struct {
	base.APIConf
	//info *RespExchangeInfo
	HttpClient    *http.Client
	EndPoint      string
	ReqUrl        *ReqUrl
	GetSymbolName func(string) string
}

// NewApiClient To do CHECK
func NewApiClient(conf base.APIConf) *SpotApiClient {
	var (
		a = &SpotApiClient{
			APIConf:  conf,
			EndPoint: SPOT_API_BASE_URL,
		}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)
	// fmt.Println("proxy:", conf.ProxyUrl)
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
	a.ReqUrl = NewSpotReqUrl()

	return a
}

func (client *SpotApiClient) GetUrl(uri string, method string, params url.Values) string {
	signature := params.Encode()
	sign, err := crypto.GetParamHmacSHA256Base64Sign(client.SecretKey, method+"\n"+client.EndPoint+"\n"+uri+"\n"+signature)
	url := fmt.Sprintf("https://%s%s?%s&Signature=%s", client.EndPoint, uri, signature, transform.Quote(sign))
	if err != nil {
		return ""
	}
	return url
}

func (client *SpotApiClient) DoRequest(uri, method string, params url.Values, post_params interface{}, result interface{}) error {
	// fmt.Println("Do Request")

	header := &http.Header{}
	header.Add("Content-Type", CONTANTTYPE)

	url_params := url.Values{}
	// url_params.Add("AccessKeyId", ACCESSKEYID)
	url_params.Add("AccessKeyId", client.AccessKey)
	url_params.Add("SignatureMethod", SIGNATUREMETHOT)
	url_params.Add("SignatureVersion", "2")
	url_params.Add("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	if method == "GET" {
		ParseOptions(&params, &url_params)
	}

	// fmt.Println(uri, method, params)
	// get和post分别处理
	// get请求参数要合并在url中同时进行签名
	// post请求参数要放在body中不进行签名
	// rsp, err := conn.Request(client.HttpClient, client.GetUrl(uri, method, url_params), method, header, params)
	rsp, err := conn.HuoBiRequest(client.HttpClient, client.GetUrl(uri, method, url_params), method, header, post_params)
	fmt.Println("response", client.GetUrl(uri, method, url_params), string(rsp), err)

	if err == nil {
		if bytes.HasPrefix(rsp, []byte("\"code\"")) {
			var re RespError
			json.Unmarshal(rsp, &re)
			if re.Code != 0 && re.Msg != "" {
				unknown := false
				//-1006 UNEXPECTED_RESP 从消息总线收到意外的响应。 执行状态未知。
				//-1007 TIMEOUT 等待后端服务器响应超时。 发送状态未知； 执行状态未知。
				if re.Code == -1006 || re.Code == -1007 {
					unknown = true
				}
				logger.Logger.Debug(uri, params, re)
				return &base.ApiError{Code: 200, BizCode: re.Code, ErrMsg: re.Msg, UnknownStatus: unknown}
			}
		}
		err = json.Unmarshal(rsp, result)
		logger.Logger.Debug(uri, params, err, result)
		return err
	}
	logger.Logger.Debug(uri, params, err)
	// err not nil
	if v, ok := err.(*conn.HttpError); ok {
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
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

func (client *SpotApiClient) GetMarketStatus() (*RespGetMarketStatus, error) {
	params := url.Values{}
	post_params := EmptyPostParams{}
	res := &RespGetMarketStatus{}
	err := client.DoRequest(client.ReqUrl.MARKET_STATUS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *SpotApiClient) GetAllSymbols(options *url.Values) (*RespGetAllSymbols, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetAllSymbols{}
	err := client.DoRequest(client.ReqUrl.ALL_SYMBOLS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *SpotApiClient) GetAllCurrencies(options *url.Values) (*RespGetAllCurrencies, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetAllCurrencies{}
	err := client.DoRequest(client.ReqUrl.ALL_CURRENCIES_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取币种配置
func (client *SpotApiClient) GetCurrencysSettings(options *url.Values) (*RespGetCurrencysSettings, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetCurrencysSettings{}
	err := client.DoRequest(client.ReqUrl.CURRENCYS_SETTINGS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取交易对配置
func (client *SpotApiClient) GetSymbolsSettings(options *url.Values) (*RespGetSymbolsSettings, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetSymbolsSettings{}
	err := client.DoRequest(client.ReqUrl.SYMBOLS_SETTINGS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取市场交易对配置
func (client *SpotApiClient) GetMarketSymbolsSettings(options *url.Values) (*RespGetMarketSymbolsSettings, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetMarketSymbolsSettings{}
	err := client.DoRequest(client.ReqUrl.MARKET_SYMBOLS_SETTINGS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取链信息
func (client *SpotApiClient) GetChainsSettings(options *url.Values) (*RespGetChainsSettings, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetChainsSettings{}
	err := client.DoRequest(client.ReqUrl.CHAINS_SETTINGS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *SpotApiClient) GetCurrenciesChains(options *url.Values) (*RespGetCurrenciesChains, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetCurrenciesChains{}
	err := client.DoRequest(client.ReqUrl.CURRENCIES_CHAINS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *SpotApiClient) GetKline(symbol string, period string, options *url.Values) (*RespGetKline, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("period", period)
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetKline{}
	err := client.DoRequest(client.ReqUrl.KLINE_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 聚合行情(Ticker)
func (client *SpotApiClient) GetMerged(symbol string) (*RespGetMerged, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	post_params := EmptyPostParams{}
	res := &RespGetMerged{}
	err := client.DoRequest(client.ReqUrl.MERGED_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 所有交易对的最新Tickers
func (client *SpotApiClient) GetAllTickers() (*RespGetAllTickers, error) {
	params := url.Values{}
	res := &RespGetAllTickers{}
	post_params := EmptyPostParams{}
	err := client.DoRequest(client.ReqUrl.All_TICKERS_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取近期交易记录
func (client *SpotApiClient) GetHistoryTrade(symbol string, options *url.Values) (*RespGetHistoryTrade, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	post_params := EmptyPostParams{}
	res := &RespGetHistoryTrade{}
	err := client.DoRequest(client.ReqUrl.HISTORY_TRADE_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 最近24小时行情数据
func (client *SpotApiClient) GetDetail(symbol string) (*RespGetDetail, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	post_params := EmptyPostParams{}
	res := &RespGetDetail{}
	err := client.DoRequest(client.ReqUrl.DETAIL_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *SpotApiClient) GetTrade(symbol string) (*RespGetNewestTrade, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	post_params := EmptyPostParams{}
	res := &RespGetNewestTrade{}
	err := client.DoRequest(client.ReqUrl.NEWEST_TRADE_URL, "GET", params, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *SpotApiClient) Account() (*RespGetAccount, error) {
	params := url.Values{}
	res := &RespGetAccount{}
	err := client.DoRequest(client.ReqUrl.ACCOUNT_URL, "GET", params, EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *SpotApiClient) ExchangeInfo() (*RespExchangeInfo, error) {
	res := &RespExchangeInfo{}
	err := client.DoRequest(client.ReqUrl.EXCHANGEINFO_URL, "GET", url.Values{}, EmptyPostParams{}, res)
	return res, err
}

func (client *SpotApiClient) GetDepth(symbol string, limit int) (*RespDepth, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("depth", strconv.Itoa(limit))
	params.Add("type", "step0")
	res := &RespDepth{}
	err := client.DoRequest(client.ReqUrl.DEPTH_URL, "GET", params, EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (client *SpotApiClient) ServerTime() (*RespTime, error) {
	res := &RespTime{}
	err := client.DoRequest(client.ReqUrl.TIME_URL, "GET", url.Values{}, EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获取用户当前手续费率
func (client *SpotApiClient) AssetTradeFee(symbols string) (*RespAssetTradeFee, error) {
	params := url.Values{}
	params.Add("symbols", symbols)
	res := &RespAssetTradeFee{}
	err := client.DoRequest(client.ReqUrl.ASSET_TRADEFEE_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// APIv2币链参考信息
func (client *SpotApiClient) GetReferenceCurrencies(options *url.Values) (*RespGetReferenceCurrencies, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespGetReferenceCurrencies{}
	err := client.DoRequest(client.ReqUrl.GET_REFERENCE_CURRENCIES_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

func (client *SpotApiClient) AccountBalance(accountType string) (*RespAccountBalance, error) {
	params := url.Values{}
	res := &RespAccountBalance{}
	var err error
	if accountType == "spot" {
		err = client.DoRequest(client.ReqUrl.ACCOUNTBALANCE_URL+SPOT_ACCOUNT_ID+"/balance", "GET", params, EmptyPostParams{}, res)
	} else if accountType == "super-margin" {
		err = client.DoRequest(client.ReqUrl.ACCOUNTBALANCE_URL+SUPER_MARGIN_ACCOUNT_ID+"/balance", "GET", params, EmptyPostParams{}, res)
	} else {
		err = client.DoRequest(client.ReqUrl.ACCOUNTBALANCE_URL+MARGIN_ACCOUNT_ID+"/balance", "GET", params, EmptyPostParams{}, res)
	}

	return res, err
}

func (client *SpotApiClient) CapitalDepositHisRec(options *url.Values) (*RespCapitalDepositHisRec, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	params.Add("type", "deposit")
	res := &RespCapitalDepositHisRec{}
	err := client.DoRequest(client.ReqUrl.CAPITAL_DEPOSIT_HISREC_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

func (client *SpotApiClient) PostOrder(post_params ReqPostOrder) (*RespPostOrder, error) {
	params := url.Values{}
	res := &RespPostOrder{}
	err := client.DoRequest(client.ReqUrl.ORDER_URL, "POST", params, post_params, res)
	return res, err

}

func GetSymbolName(symbol string) string {
	symbol = strings.Replace(symbol, "/", "", 1)
	//symbol = strings.Replace(symbol, "-", "", 1)
	//symbol = strings.Replace(symbol, "_", "", 1)
	symbol = strings.ToLower(symbol)
	return symbol
}

// 获取平台资产总估值
func (client *SpotApiClient) GetValuation(options *url.Values) (*RespGetValuation, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespGetValuation{}
	err := client.DoRequest(client.ReqUrl.VALUATION_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 获取指定账户资产估值
func (client *SpotApiClient) GetAssetValuation(accountType string, options *url.Values) (*RespGetAssetValuation, error) {
	params := url.Values{}
	params.Add("accountType", accountType)
	ParseOptions(options, &params)
	res := &RespGetAssetValuation{}
	err := client.DoRequest(client.ReqUrl.ASSET_VALUATION_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 账户流水
func (client *SpotApiClient) GetAccountHistory(account_id string, options *url.Values) (*RespGetAccountHistory, error) {
	params := url.Values{}
	params.Add("account-id", account_id)
	ParseOptions(options, &params)
	res := &RespGetAccountHistory{}
	err := client.DoRequest(client.ReqUrl.ACCOUNT_HISTORY_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 财务流水
func (client *SpotApiClient) GetAccountLedger(account_id string, options *url.Values) (*RespGetAccountLedger, error) {
	params := url.Values{}
	params.Add("account-id", account_id)
	ParseOptions(options, &params)
	res := &RespGetAccountLedger{}
	err := client.DoRequest(client.ReqUrl.ACCOUNT_LEDGER_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 点卡余额查询
func (client *SpotApiClient) GetPointAccount(options *url.Values) (*RespGetPointAccount, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespGetPointAccount{}
	err := client.DoRequest(client.ReqUrl.POINT_ACCOUNT_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 点卡划转
func (client *SpotApiClient) PostPointTransfer(post_params ReqPostPointTransfer) (*RespPostPointTransfer, error) {
	params := url.Values{}
	res := &RespPostPointTransfer{}
	err := client.DoRequest(client.ReqUrl.POINT_TRANSFER_URL, "POST", params, post_params, res)
	return res, err
}

// 充币地址查询
func (client *SpotApiClient) GetDepositAddress(currency string) (*RespGetDepositAddress, error) {
	params := url.Values{}
	params.Add("currency", currency)
	res := &RespGetDepositAddress{}
	err := client.DoRequest(client.ReqUrl.DEPOSIT_ADDRESS_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 提币额度查询
func (client *SpotApiClient) GetWithdrawQuota(currency string) (*RespGetWithdrawQuota, error) {
	params := url.Values{}
	params.Add("currency", currency)
	res := &RespGetWithdrawQuota{}
	err := client.DoRequest(client.ReqUrl.WITHDRAW_QUOTA_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 提币地址查询，用户需要先通过Web端添加提币地址，才可以通过接口查询到
func (client *SpotApiClient) GetWithdrawAddress(currency string, options *url.Values) (*RespGetWithdrawAddress, error) {
	params := url.Values{}
	params.Add("currency", currency)
	ParseOptions(options, &params)
	res := &RespGetWithdrawAddress{}
	err := client.DoRequest(client.ReqUrl.WITHDRAW_ADDRESS_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 虚拟币提币
func (client *SpotApiClient) PostWithdrawCreate(post_params ReqPostWithdrawCreate) (*RespPostWithdrawCreate, error) {
	params := url.Values{}
	res := &RespPostWithdrawCreate{}
	err := client.DoRequest(client.ReqUrl.WITHDRAW_CREATE_URL, "POST", params, post_params, res)
	return res, err
}

// 通过clientOrderId查询提币订单
func (client *SpotApiClient) GetWithdrawClientOrderId(clientOrderId string) (*RespGetWithdrawClientOrderId, error) {
	params := url.Values{}
	params.Add("clientOrderId", clientOrderId)
	res := &RespGetWithdrawClientOrderId{}
	err := client.DoRequest(client.ReqUrl.WITHDRAW_CLIENTORDERID_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 取消提币
func (client *SpotApiClient) PostCancelWithdrawCreate(withdraw_id int64) (*RespPostCancelWithdrawCreate, error) {
	params := url.Values{}
	res := &RespPostCancelWithdrawCreate{}
	err := client.DoRequest(client.ReqUrl.CANCEL_WITHDRAW_CREATE_URL+strconv.FormatInt(withdraw_id, 10)+"/cancel", "POST", params, EmptyPostParams{}, res)
	return res, err
}

// 设置子用户手续费抵扣模式
func (client *SpotApiClient) PostSubUserDeductMode(post_params ReqPostSubUserDeductMode) (*RespPostSubUserDeductMode, error) {
	params := url.Values{}
	res := &RespPostSubUserDeductMode{}
	err := client.DoRequest(client.ReqUrl.SUB_USER_DEDUCT_MODE_URL, "POST", params, post_params, res)
	return res, err
}

// 母子用户API信息查询
func (client *SpotApiClient) GetUserApiKey(uid int64, options *url.Values) (*RespGetUserApiKey, error) {
	params := url.Values{}
	params.Add("uid", strconv.FormatInt(uid, 10))
	ParseOptions(options, &params)
	res := &RespGetUserApiKey{}
	err := client.DoRequest(client.ReqUrl.USER_API_KEY_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 母子用户获取用户UID
func (client *SpotApiClient) GetUserUid() (*RespGetUserUid, error) {
	params := url.Values{}
	res := &RespGetUserUid{}
	err := client.DoRequest(client.ReqUrl.USER_UID_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

func ToJson(v interface{}) (string, error) {
	result, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// 子用户创建，结构体数组参数传参有点问题
func (client *SpotApiClient) PostSubUserCreation(userList ReqPostSubUserCreation) (*RespPostSubUserCreation, error) {
	params := url.Values{}
	res := &RespPostSubUserCreation{}
	err := client.DoRequest(client.ReqUrl.SUB_USER_CREATION_URL, "POST", params, userList, res)
	return res, err
}

// 获取子用户列表
func (client *SpotApiClient) GetSubUserList(options *url.Values) (*RespGetSubUserList, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespGetSubUserList{}
	err := client.DoRequest(client.ReqUrl.SUB_USER_LIST_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 冻结/解冻子用户
func (client *SpotApiClient) PostSubUserManagement(post_params ReqPostSubUserManagement) (*RespPostSubUserManagement, error) {
	params := url.Values{}
	res := &RespPostSubUserManagement{}
	err := client.DoRequest(client.ReqUrl.SUB_USER_MANAGEMENT_URL, "POST", params, post_params, res)
	return res, err
}

// 获取特定子用户的用户状态
func (client *SpotApiClient) GetSubUserState(subUid int64) (*RespGetSubUserState, error) {
	params := url.Values{}
	params.Add("subUid", strconv.FormatInt(subUid, 10))
	res := &RespGetSubUserState{}
	err := client.DoRequest(client.ReqUrl.SUB_USER_STATE_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 批量下单
func (client *SpotApiClient) PostBatchOrders(post_params ReqPostBatchOrders) (*RespPostBatchOrders, error) {
	res := &RespPostBatchOrders{}
	err := client.DoRequest(client.ReqUrl.BATCH_ORDERS_URL, "POST", url.Values{}, post_params, res)
	return res, err
}

// 撤销订单
func (client *SpotApiClient) PostSubmitOrder(order_id string, post_params ReqPostSubmitOrder) (*RespPostSubmitOrder, error) {
	res := &RespPostSubmitOrder{}
	err := client.DoRequest(client.ReqUrl.CANCEL_ORDER_URL+order_id+"/submitcancel", "POST", url.Values{}, post_params, res)
	return res, err
}

// 查询订单详情
func (client *SpotApiClient) GetOrder(order_id string) (*RespGetOrder, error) {
	res := &RespGetOrder{}
	err := client.DoRequest(client.ReqUrl.GET_ORDER_URL+order_id, "GET", url.Values{}, EmptyPostParams{}, res)
	return res, err
}

// 搜索历史订单
func (client *SpotApiClient) GetHistoryOrders(symbol string, states string, optins *url.Values) (*RespGetHistoryOrders, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("states", states)
	ParseOptions(optins, &params)
	res := &RespGetHistoryOrders{}
	err := client.DoRequest(client.ReqUrl.GET_HISTORY_ORDERS_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 成交明细
func (client *SpotApiClient) GetOrderMatchresults(order_id string) (*RespGetOrderMatchresults, error) {
	res := &RespGetOrderMatchresults{}
	err := client.DoRequest(client.ReqUrl.GET_ORDER_MATCHRESULTS_URL+order_id+"/matchresults", "GET", url.Values{}, EmptyPostParams{}, res)
	return res, err
}

// 申请借币(逐仓)
func (client *SpotApiClient) PostMarginOrders(post_params ReqPostMarginOrders) (*RespPostMarginOrders, error) {
	res := &RespPostMarginOrders{}
	err := client.DoRequest(client.ReqUrl.POST_MARGIN_ORDERS_URL, "POST", url.Values{}, post_params, res)
	return res, err
}

// 申请借币(全仓)
func (client *SpotApiClient) PostCrossMarginOrders(post_params ReqPostCrossMarginOrders) (*RespPostCrossMarginOrders, error) {
	res := &RespPostCrossMarginOrders{}
	err := client.DoRequest(client.ReqUrl.POST_CROSS_MARGIN_ORDERS_URL, "POST", url.Values{}, post_params, res)
	return res, err
}

// 查询借币订单(逐仓)
func (client *SpotApiClient) GetMarginLoanOrders(symbol string, options *url.Values) (*RespGetMarginLoanOrders, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	res := &RespGetMarginLoanOrders{}
	err := client.DoRequest(client.ReqUrl.GET_MARGIN_LOAN_ORDERS_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 查询借币订单(全仓)
func (client *SpotApiClient) GetCrossMarginLoanOrders(options *url.Values) (*RespGetCrossMarginLoanOrders, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespGetCrossMarginLoanOrders{}
	err := client.DoRequest(client.ReqUrl.GET_CROSS_MARGIN_LOAN_ORDERS_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 归还借币(逐仓)
func (client *SpotApiClient) PostMarginRepay(order_id string, post_params ReqPostMarginRepay) (*RespPostMarginRepay, error) {
	res := &RespPostMarginRepay{}
	err := client.DoRequest(client.ReqUrl.POST_MARGIN_REPAY_URL+order_id+"/repay", "POST", url.Values{}, post_params, res)
	return res, err
}

// 归还借币(全仓)
func (client *SpotApiClient) PostCrossMarginRepay(order_id string, post_params ReqPostCrossMarginRepay) (*RespPostCrossMarginRepay, error) {
	res := &RespPostCrossMarginRepay{}
	err := client.DoRequest(client.ReqUrl.POST_CROSS_MARGIN_REPAY_URL+order_id+"/repay", "POST", url.Values{}, post_params, res)
	return res, err
}

// 还币交易记录查询
func (client *SpotApiClient) GetRepayment(options *url.Values) (*RespGetRepayment, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespGetRepayment{}
	err := client.DoRequest(client.ReqUrl.GET_REPAYMENT_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 当前和历史成交
func (client *SpotApiClient) GetMatchresults(symbol string, options *url.Values) (*RespGetMatchresults, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	ParseOptions(options, &params)
	res := &RespGetMatchresults{}
	err := client.DoRequest(client.ReqUrl.GET_MATCHRESULTS_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 查询订单详情(基于client order ID)
func (client *SpotApiClient) GetClientOrder(clientOrder_id string) (*RespGetClientOrder, error) {
	params := url.Values{}
	params.Add("clientOrderId", clientOrder_id)
	res := &RespGetClientOrder{}
	err := client.DoRequest(client.ReqUrl.GET_CLIENT_ORDER_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 查询当前未成交订单
func (client *SpotApiClient) GetOpenOrders(options *url.Values) (*RespGetOpenOrders, error) {
	params := url.Values{}
	ParseOptions(options, &params)
	res := &RespGetOpenOrders{}
	err := client.DoRequest(client.ReqUrl.GET_OPEN_ORDERS_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 充提记录
func (client *SpotApiClient) GetDepositHistory(type_ string, options *url.Values) (*RespGetDepositHistory, error) {
	params := url.Values{}
	params.Add("type", type_)
	ParseOptions(options, &params)
	res := &RespGetDepositHistory{}
	err := client.DoRequest(client.ReqUrl.GET_DEPOSIT_HISTORY_URL, "GET", params, EmptyPostParams{}, res)
	return res, err
}

// 币本位交割合约资金划转
func (c *SpotApiClient) PostCFuturesTransfer(post_params ReqPostCFuturesTransfer) (*RespPostCFuturesTransfer, error) {
	res := &RespPostCFuturesTransfer{}
	err := c.DoRequest(c.ReqUrl.C_BASE_FUTURE_TRANSFER_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 币本位永续合约资金划转
func (c *SpotApiClient) PostCSwapTransfer(post_params ReqPostCSwapTransfer) (*RespPostCSwapTransfer, error) {
	res := &RespPostCSwapTransfer{}
	err := c.DoRequest(c.ReqUrl.C_BASE_SWAP_TRANSFER_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// U本位划转
func (c *SpotApiClient) PostUTransfer(post_params ReqPostUTransfer) (*RespPostUTransfer, error) {
	res := &RespPostUTransfer{}
	err := c.DoRequest(c.ReqUrl.U_TRANSFER_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}
