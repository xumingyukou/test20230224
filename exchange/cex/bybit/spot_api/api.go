package spot_api

import (
	"bytes"
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type ApiClient struct {
	base.APIConf
	HttpClient    *http.Client
	EndPoint      string
	ReqUrl        *ReqUrl
	GetSymbolName func(string) string
}

var (
	// OkexWeight 定义全局的限制
	BybitWeight *base.RateLimitMgr
)

func ClearCapacity(consumeMap []*base.RateLimitConsume) {
	for _, rate := range consumeMap {
		tmp := BybitWeight.GetInstance(rate.LimitTypeName)
		tmp.Instance.Clear()
	}
}

func WeightAllow(consumeMap []*base.RateLimitConsume) error {
	for _, rate := range consumeMap {
		if err := BybitWeight.Consume(*rate); err != nil {
			return err
		}
	}
	return nil
}

func NewApiClient(conf base.APIConf) *ApiClient {
	var (
		a = &ApiClient{
			APIConf:  conf,
			EndPoint: SPOT_API_BASE_URL,
		}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)
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

	if BybitWeight == nil {
		BybitWeight = base.NewRateLimitMgr()
	}

	a.ReqUrl = NewSpotReqUrl()

	return a
}

func (client *ApiClient) GetUrl(uri string, method string, params url.Values) string {
	return client.EndPoint + uri
}

func (client *ApiClient) DoRequest(uri base.ReqUrlInfo, method string, params url.Values, result interface{}) error {
	var err error
	// fmt.Println("Do Request")

	if err = WeightAllow(client.ReqUrl.GetURLRateLimit(uri, params)); err != nil {
		return err
	}
	header := &http.Header{}
	header.Add("Content-Type", CONTANTTYPE)

	params.Add("api_key", client.APIConf.AccessKey)
	//params.Add("signatureVersion", "2")
	//params.Add("timestamp", time.Now().Add(-time.Hour*8).Format("2006-01-02T15:04:05"))
	params.Add("timestamp", transform.XToString(time.Now().UnixMilli()))
	params.Add("sign", client.getSigned(params))
	rsp, err := conn.Request(client.HttpClient, client.GetUrl(uri.Url, method, params), method, header, params)
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
		if v.Code == 403 {
			ClearCapacity(client.ReqUrl.GetURLRateLimit(uri, params))
		}
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

// getSigned
func (clicent *ApiClient) getSigned(params url.Values) string {
	var param string

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var p []string
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%v=%v", k, params.Get(k)))
	}
	param = strings.Join(p, "&")

	sig := hmac.New(sha256.New, []byte(clicent.APIConf.SecretKey))
	sig.Write([]byte(param))
	signature := hex.EncodeToString(sig.Sum(nil))
	return signature
}

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sha)
}

func (client *ApiClient) GetSymbols() (*RespSymbols, error) {
	params := url.Values{}
	res := &RespSymbols{}
	err := client.DoRequest(client.ReqUrl.SYMBOLS_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetWithdrawFee() (*RespWithdrawFee, error) {
	params := url.Values{}
	res := &RespWithdrawFee{}
	err := client.DoRequest(client.ReqUrl.WITHDRAWFEE_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderBook(symbol string) (*RespOrderBook, error) {
	params := url.Values{}
	params.Add("symbol", Trans2Send(symbol))
	res := &RespOrderBook{}
	err := client.DoRequest(client.ReqUrl.ORDERBOOK_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Trans2Send(symbol string) string {
	return strings.ReplaceAll(symbol, "/", "")

}

func (client *ApiClient) GetTime(options *url.Values) (*RespTime, error) {
	params := url.Values{}
	res := &RespTime{}
	err := client.DoRequest(client.ReqUrl.TIME_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// todo 接口没有返回数据
func (client *ApiClient) GetFundFee(symbol string) (*RespFundFee, error) {
	params := url.Values{}
	//params.Add("ItCode", Trans2Send(symbol))
	res := &RespFundFee{}
	err := client.DoRequest(client.ReqUrl.FUNDFEE_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetAccountBalance() (*RespAccountBalance, error) {
	params := url.Values{}
	res := &RespAccountBalance{}
	err := client.DoRequest(client.ReqUrl.ACCOUNT_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// todo 未完成
func (client *ApiClient) GetMarginBalance(symbol string) (*RespAccountBalance, error) {
	params := url.Values{}
	res := &RespAccountBalance{}
	params.Add("coin", symbol)
	err := client.DoRequest(client.ReqUrl.FUNDFEE_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) PostPlaceOrder(symbol, orderQty, side, orderType, orderLinkId string, options ...url.Values) (*RespPlaceOrder, error) {
	params := url.Values{}
	params = MergeParams(params, options...)
	res := &RespPlaceOrder{}
	params.Add("symbol", symbol)
	params.Add("orderQty", orderQty)
	params.Add("side", side)
	params.Add("orderType", orderType)
	params.Add("orderLinkId", orderLinkId)
	err := client.DoRequest(client.ReqUrl.PLACEORDER_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) PostCancleOrder(orderId string, param2 ...url.Values) (*RespCancleOrder, error) {
	params := url.Values{}
	res := &RespCancleOrder{}
	if len(param2) > 0 {
		params = MergeParams(params, param2[0])
	}
	params.Add("orderId", orderId)
	err := client.DoRequest(client.ReqUrl.CANCLEORDER_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderInfo(orderId string, param2 ...url.Values) (*RespOrderInfo, error) {
	params := url.Values{}

	res := &RespOrderInfo{}
	params = MergeParams(params, param2...)
	params.Add("orderId", orderId)
	err := client.DoRequest(client.ReqUrl.ORDERINFO_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 获得正在进行的顶坛通过status进行判断
func (client *ApiClient) GetOrderHistory(param ...url.Values) (*RespOrderHistory, error) {

	params := url.Values{}
	params = MergeParams(params, param...)
	res := &RespOrderHistory{}
	err := client.DoRequest(client.ReqUrl.ORDERHISTORY_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderList(param ...url.Values) (*RespOrderHistory, error) {
	params := url.Values{}
	params = MergeParams(params, param...)
	res := &RespOrderHistory{}
	err := client.DoRequest(client.ReqUrl.ORDERINFO2_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func MergeParams(param1 url.Values, params ...url.Values) url.Values {
	if len(params) == 1 {
		for k := range params[0] {
			param1.Add(k, params[0].Get(k))
		}
	}
	return param1
}

// 待测试
func (client *ApiClient) Withdraw(param ...url.Values) (*RespWithdraw, error) {
	params := url.Values{}
	params = MergeParams(params, param...)
	res := &RespWithdraw{}
	err := client.DoRequest(client.ReqUrl.WITHDRAW_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) WithdrawHistory(param ...url.Values) (*RespWithdrawHistory, error) {

	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	res := &RespWithdrawHistory{}
	err := client.DoRequest(client.ReqUrl.WITHDRAWHISTORY_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) Transfer(transfer_id, coin, amount, from_account_type, to_account_type string, param ...url.Values) (*RespTransfer, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	params.Add("transfer_id", transfer_id)
	params.Add("coin", coin)
	params.Add("amount", amount)
	params.Add("from_account_type", from_account_type)
	params.Add("to_account_type", to_account_type)
	res := &RespTransfer{}
	err := client.DoRequest(client.ReqUrl.TRANSFER_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) TransferM2S(transfer_id, coin, amount, sub_user_id, type_ string, param ...url.Values) (*RespTransferM2S, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	params.Add("transfer_id", transfer_id)
	params.Add("coin", coin)
	params.Add("amount", amount)
	params.Add("sub_user_id", sub_user_id)
	params.Add("type_", type_)
	res := &RespTransferM2S{}
	err := client.DoRequest(client.ReqUrl.TRANSFER_M2S_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) ALLTransfer(transfer_id, coin, amount, from_member_id, to_member_id, from_account_type, to_account_type string, param ...url.Values) (*RespTransfer, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	params.Add("transfer_id", transfer_id)
	params.Add("coin", coin)
	params.Add("amount", amount)
	params.Add("from_member_id", from_member_id)
	params.Add("to_member_id", to_member_id)
	params.Add("from_account_type", from_account_type)
	params.Add("to_account_type", to_account_type)
	res := &RespTransfer{}
	err := client.DoRequest(client.ReqUrl.TRANSFER_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) TransferHistory(param ...url.Values) (*RespTransferHistory, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	res := &RespTransferHistory{}
	err := client.DoRequest(client.ReqUrl.TRANSFERHISTORY_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) TransferHistoryS2M(param ...url.Values) (*RespTransferHistoryS2M, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	res := &RespTransferHistoryS2M{}
	err := client.DoRequest(client.ReqUrl.TRANSFERHISTORYM2S_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) RecordHistory(timestamp string, param ...url.Values) (*RespRecordHistory, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	params.Add("timestamp", timestamp)
	res := &RespRecordHistory{}
	err := client.DoRequest(client.ReqUrl.RECORDHISTORY_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) PostLoan(coin, qyt string, param ...url.Values) (*RespLoan, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	params.Add("coin", coin)
	params.Add("qyt", qyt)
	res := &RespLoan{}
	err := client.DoRequest(client.ReqUrl.LOAN_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetLoanHistory(param ...url.Values) (*RespLoanHistory, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	res := &RespLoanHistory{}
	err := client.DoRequest(client.ReqUrl.LOANHISTORY_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// todo 还款unkown问题
func (client *ApiClient) PostRepay(coin, qyt string, param ...url.Values) (*RespRepay, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	params.Add("coin", coin)
	params.Add("qyt", qyt)
	res := &RespRepay{}
	err := client.DoRequest(client.ReqUrl.REPAY_URL, "POST", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetRepayHistory(param ...url.Values) (*RespRepayHistory, error) {
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	res := &RespRepayHistory{}
	err := client.DoRequest(client.ReqUrl.REPAYHISTORY_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetFutureSymbols() (*RespFutureSymbols, error) {
	params := url.Values{}
	res := &RespFutureSymbols{}
	err := client.DoRequest(client.ReqUrl.FUTURE_SYMBOL_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetFutureOrderbook(symbol string) (*RespFutureOrderbook, error) {
	params := url.Values{}
	params.Add("symbol", symbol)
	res := &RespFutureOrderbook{}
	err := client.DoRequest(client.ReqUrl.FUTURE_ORDERBOOK_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetMarketPrice(symbol, interval, from string) (*RespMarketPrice, error) {
	res := &RespMarketPrice{}
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("interval", interval)
	params.Add("from", from)

	//err := client.DoRequest2(client.ReqUrl.FUTURE_MARKET_URL, "GET", false, params, res, nil)
	err := client.DoRequest(client.ReqUrl.FUTURE_MARKET_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetPosition(param ...url.Values) (*RespFuturePosition, error) {
	res := &RespFuturePosition{}
	params := url.Values{}
	if param != nil {
		params = MergeParams(params, param[0])
	}
	err := client.DoRequest(client.ReqUrl.FUTURE_POSITION_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type OrderInfo struct {
	Symbol, Side, Order_type, TimeInForce, ReduceOnly string
	PositionIdx                                       int
	Qty                                               float64
	CloseOnTrigger                                    bool
}

func (client *ApiClient) PlaceFutureCoinOrder(orderInfo OrderInfo, options ...url.Values) (*RespFutureCoinOrder, error) {
	res := &RespFutureCoinOrder{}
	params := map[string]interface{}{}
	params["symbol"] = orderInfo.Symbol

	params["side"] = orderInfo.Side
	params["order_type"] = orderInfo.Order_type
	params["qty"] = orderInfo.Qty
	params["time_in_force"] = orderInfo.TimeInForce
	params["reduce_only"] = false
	if options != nil {
		params = mergeUrlValues(params, options[0])
	}
	err := client.DoRequest2(client.ReqUrl.FUTURE_PLACE_COIN_DELIVERY_URL, "POST", true, params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) PlaceFutureOrder(orderInfo OrderInfo, options ...url.Values) (*RespFuturePlaceOrder, error) {
	res := &RespFuturePlaceOrder{}
	params := map[string]interface{}{}

	params["symbol"] = orderInfo.Symbol
	params["side"] = orderInfo.Side
	params["order_type"] = orderInfo.Order_type
	params["qty"] = orderInfo.Qty
	params["time_in_force"] = orderInfo.TimeInForce
	params["reduce_only"] = orderInfo.ReduceOnly
	params["close_on_trigger"] = orderInfo.CloseOnTrigger
	params["reduce_only"] = false

	if options != nil {
		params = mergeUrlValues(params, options[0])
	}
	err := client.DoRequest2(client.ReqUrl.FUTURE_PLACE_URL, "POST", true, params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) PlaceFutureCoinDeliveryOrder(orderInfo OrderInfo, options ...url.Values) (*RespFutureCoinDeliveryPlaceOrder, error) {
	res := &RespFutureCoinDeliveryPlaceOrder{}
	params := map[string]interface{}{}

	params["symbol"] = orderInfo.Symbol
	params["side"] = orderInfo.Side
	params["order_type"] = orderInfo.Order_type
	params["qty"] = orderInfo.Qty
	params["reduce_only"] = false
	params["time_in_force"] = orderInfo.TimeInForce
	if options != nil {
		params = mergeUrlValues(params, options[0])
	}
	err := client.DoRequest2(client.ReqUrl.FUTURE_PLACE_COIN_DELIVERY_URL, "POST", true, params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func mergeUrlValues(params map[string]interface{}, options ...url.Values) map[string]interface{} {
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}
	return params
}

// 需要传入bool值，如果使用url.Values因为是string会失败
func (client *ApiClient) DoRequest2(uri base.ReqUrlInfo, method string, requireSign bool, params map[string]interface{}, result interface{}) error {

	var err error
	if err = WeightAllow(client.ReqUrl.GetURLRateLimit(uri, map2urlValues(params))); err != nil {
		return err
	}
	body := make([]byte, 0)
	urlParam := ""
	var header map[string]string
	if params != nil {
		if len(params) != 0 && (method == "GET" || method == "DELETE") {
			urlParam = "?"
			i := 0
			for key, element := range params {
				if elements, ok := element.([]string); ok {
					element = strings.Join(elements, ",")
				}
				if i == 0 {
					urlParam += key + "=" + fmt.Sprintf("%v", element)
				} else {
					urlParam += "&" + key + "=" + fmt.Sprintf("%v", element)
				}
				i += 1
			}
		} else if method == "POST" || method == "PUT" {
			// 不能调位置，需要在body之中加入apikey timestamp
			params["api_key"] = client.APIConf.AccessKey
			params["timestamp"] = transform.XToString(time.Now().UnixMilli())
			header, err = client.getHeader(requireSign, params)
			params["sign"] = header["sign"]
			body, err = json.Marshal(params)
			if err != nil {
				logger.Logger.Error("json marshal", params, err.Error())
				return err
			}
		}
	}

	if err != nil {
		logger.Logger.Error("get header", err.Error())
		return err
	}
	// fmt.Println(client.GetUrl(uri)+urlParam, method, string(body), header)

	rsp, err := conn.NewHttpRequest(client.HttpClient, method, client.GetUrl2(uri.Url)+urlParam, string(body), header)
	if err == nil {
		err = json.Unmarshal(rsp, result)
		logger.Logger.Debug(uri, params, err, result)
		return err
	}
	if v, ok := err.(*conn.HttpError); ok {
		if v.Code == 403 {
			ClearCapacity(client.ReqUrl.GetURLRateLimit(uri, map2urlValues(params)))
		}
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

func map2urlValues(m map[string]interface{}) url.Values {
	params := url.Values{}
	for k, v := range m {
		params.Add(k, transform.XToString(v))
	}
	return params
}

func (client *ApiClient) getHeader(requireSign bool, params map[string]interface{}) (map[string]string, error) {
	header := make(map[string]string)
	if requireSign {
		header["api_key"] = client.APIConf.AccessKey
		header["Content-Type"] = CONTANTTYPE
		header["timestamp"] = transform.XToString(time.Now().UnixMilli())
		header["sign"] = client.getSigned2(params)
	}

	header["Content-Type"] = "application/json"

	return header, nil
}

func (clicent *ApiClient) getSigned2(params map[string]interface{}) string {
	var param string

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var p []string
	for _, k := range keys {
		p = append(p, fmt.Sprintf("%v=%v", k, params[k]))
	}
	param = strings.Join(p, "&")

	sig := hmac.New(sha256.New, []byte(clicent.APIConf.SecretKey))
	sig.Write([]byte(param))
	signature := hex.EncodeToString(sig.Sum(nil))
	return signature
}

func (client *ApiClient) GetUrl2(url string) string {
	return client.EndPoint + url
}
