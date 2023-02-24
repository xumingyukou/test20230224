package spot_api

import (
	"clients/conn"
	"clients/conn/ratelimit"
	"clients/crypto"
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	tb = ratelimit.NewTokenBucket(6, 10)
	// KucoinWeight 定义全局的限制
	KucoinWeight *base.RateLimitMgr
)

type ApiClient struct {
	base.APIConf
	HttpClient *http.Client
	ReqUrl     *ReqUrl
}

func ClearCapacity(consumeMap []*base.RateLimitConsume) {
	for _, rate := range consumeMap {
		tmp := KucoinWeight.GetInstance(rate.LimitTypeName)
		tmp.Instance.Clear()
	}
}
func WeightAllow(consumeMap []*base.RateLimitConsume) error {
	for _, rate := range consumeMap {
		if err := KucoinWeight.Consume(*rate); err != nil {
			return err
		}
	}
	return nil
}

func NewApiClient(conf base.APIConf) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
			ReqUrl:  NewSpotReqUrl(conf.SubAccount),
		}
		proxyUrl  *url.URL
		transport = http.Transport{}
		err       error
	)
	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
	}

	if conf.ProxyUrl != "" {
		proxyUrl, err = url.Parse(conf.ProxyUrl)
		if err != nil {
			logger.Logger.Error("set proxy:", conf.ProxyUrl)
		}
		transport = http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	if KucoinWeight == nil {
		KucoinWeight = base.NewRateLimitMgr()
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
			ReqUrl:  NewSpotReqUrl(conf.SubAccount),
		}
	)
	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
	}

	if KucoinWeight == nil {
		KucoinWeight = base.NewRateLimitMgr()
	}

	a.HttpClient = cli
	return a
}

func (client *ApiClient) GetUrl(url string) string {
	return client.EndPoint + url
}

func (client *ApiClient) getHeader(method string, path string, requireSign bool, body []byte) (map[string]string, error) {
	header := make(map[string]string)
	if requireSign {
		ts := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
		signaturePayload := ts + method + path + string(body)
		signature, err := crypto.GetParamHmacSHA256Base64Sign(client.SecretKey, signaturePayload)

		if err != nil {
			logger.Logger.Error("sign payload with hmac-sha256", err.Error(), signaturePayload)
			return nil, err
		}

		passphrase, err := crypto.GetParamHmacSHA256Base64Sign(client.SecretKey, client.Passphrase)
		if err != nil {
			logger.Logger.Error("sign passphrase with hmac-sha256", err.Error())
			return nil, err
		}
		header["KC-API-PASSPHRASE"] = passphrase
		header["KC-API-KEY-VERSION"] = "2"
		header["KC-API-KEY"] = client.AccessKey
		header["KC-API-SIGN"] = signature
		header["KC-API-TIMESTAMP"] = ts
	}

	header["Content-Type"] = "application/json"

	return header, nil
}

func (client *ApiClient) DoRequest(uri base.ReqUrlInfo, method string, requireSign bool, params map[string]interface{}, result interface{}) error {
	// parameters := make(map[string]interface{})
	// TODO delete parameters, remove copy

	var err error
	if err = WeightAllow(client.ReqUrl.GetURLRateLimit(uri, map2urlValues(params))); err != nil {
		return err
	}
	body := make([]byte, 0)
	urlParam := ""
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
			body, err = json.Marshal(params)
			if err != nil {
				logger.Logger.Error("json marshal", params, err.Error())
				return err
			}
		}
	}

	header, err := client.getHeader(method, uri.Url+urlParam, requireSign, body)
	if err != nil {
		logger.Logger.Error("get header", err.Error())
		return err
	}
	// fmt.Println(client.GetUrl(uri)+urlParam, method, string(body), header)

	rsp, err := conn.NewHttpRequest(client.HttpClient, method, client.GetUrl(uri.Url)+urlParam, string(body), header)
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

func (client *ApiClient) GetSymbols() (*RespGetSymbols, error) {
	res := &RespGetSymbols{}
	err := client.DoRequest(client.ReqUrl.SYMBOLS_URL, "GET", false, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetTradeFees(symbols []string) ([]*TradeFeeInfo, error) {
	n := (len(symbols)-1)/10 + 1
	tradeFees := make([]*TradeFeeInfo, 0, len(symbols))

	for i := 0; i < n; i++ {
		res := &RespGetTradeFees{}
		params := make(map[string]interface{})
		if i == n-1 {
			params["symbols"] = strings.Join(symbols[i*10:], ",")
		} else {
			params["symbols"] = strings.Join(symbols[i*10:(i+1)*10], ",")
		}
		err := client.DoRequest(client.ReqUrl.TRADEFEE_URL, "GET", true, params, res)
		if err != nil {
			return tradeFees, err
		}
		tradeFees = append(tradeFees, res.Data...)
	}

	return tradeFees, nil
}

func (client *ApiClient) GetMarketsRestrict(symbolName string) (*RespGetMarkets, error) {
	var (
		endpoint base.ReqUrlInfo
		err      error
	)
	endpoint = client.ReqUrl.MARKET_LV2_URL
	params := make(map[string]interface{})
	params["symbol"] = symbolName
	res := &RespGetMarkets{}

	for !tb.Allow() {
		// 等待1s-3s不等
		dt := time.Duration(rand.Intn(2000)+1000) * time.Millisecond
		logger.Logger.Warnf("wait for %s: GetMarketsRestrict(%s) token bucket", dt, symbolName)
		time.Sleep(dt)
	}
	err = client.DoRequest(endpoint, "GET", true, params, res)
	if err != nil {
		if strings.Contains(err.Error(), "Too many requests") {
			logger.Logger.Warn("too many request, clear bucket")
			tb.Clear()
			// 等待1s-3s不等
			dt := time.Duration(rand.Intn(2000)+1000) * time.Millisecond
			logger.Logger.Warnf("wait for %s and retry", dt)
			time.Sleep(dt)
			return client.GetMarketsRestrict(symbolName)
		}
		return res, err
	} else {
		return res, err
	}
}

func (client *ApiClient) GetMarkets(symbolName string, limit int) (*RespGetMarkets, error) {
	var (
		endpoint base.ReqUrlInfo
		err      error
	)
	if limit <= 20 {
		endpoint = client.ReqUrl.MARKET_LV2_20_URL
	} else if limit <= 100 {
		endpoint = client.ReqUrl.MARKET_LV2_100_URL
	} else {
		return client.GetMarketsRestrict(symbolName)
	}
	params := make(map[string]interface{})
	params["symbol"] = symbolName
	res := &RespGetMarkets{}

	err = client.DoRequest(endpoint, "GET", false, params, res)
	return res, err
}

func (client *ApiClient) GetFutureMarkets(symbolName string, limit int) (*RespFutureMarkets, error) {
	var (
		endpoint base.ReqUrlInfo
		err      error
	)
	endpoint = client.ReqUrl.CONTRACT_DEPTH_URL
	if limit <= 20 {
		endpoint.Url += "20"
	} else {
		endpoint.Url += "100"
	}
	params := make(map[string]interface{})
	params["symbol"] = symbolName
	res := &RespFutureMarkets{}

	err = client.DoRequest(endpoint, "GET", false, params, res)
	return res, err
}

func (client *ApiClient) GetCurrencies() (*RespGetCurrencies, error) {
	res := &RespGetCurrencies{}
	err := client.DoRequest(client.ReqUrl.CURRENCIES_URL, "GET", false, nil, res)
	return res, err
}

func (client *ApiClient) GetAccounts(currency string, accountType string) (*RespGetAccounts, error) {
	params := make(map[string]interface{})
	if currency != "" {
		params["currency"] = currency
	}

	if accountType != "" {
		params["type"] = accountType
	}

	//client.ReqUrl.ACCOUNTS_URL += "/" + client.

	res := &RespGetAccounts{}
	err := client.DoRequest(client.ReqUrl.ACCOUNTS_URL, "GET", true, params, res)
	return res, err
}

func (client *ApiClient) GetMarginAcounts() (*RespGetMarginAccounts, error) {
	res := &RespGetMarginAccounts{}
	err := client.DoRequest(client.ReqUrl.MARGIN_ACCOUNTS_URL, "GET", true, nil, res)
	return res, err
}

func (client *ApiClient) GetIsolatedAccounts(balanceCurrency string) (*RespGetIsolatedAccounts, error) {
	res := &RespGetIsolatedAccounts{}
	params := make(map[string]interface{})
	if balanceCurrency != "" {
		if balanceCurrency != "USDT" && balanceCurrency != "KCS" && balanceCurrency != "BTC" {
			return res, errors.New("unrecognized currency")
		}
		params["balanceCurrency"] = balanceCurrency
	} else {
		// balanceCurrency = "BTC"
	}
	err := client.DoRequest(client.ReqUrl.ISOLATED_ACCOUNTS_URL, "GET", true, params, res)
	return res, err
}

func UpdateMapWith(dst map[string]interface{}, src map[string]interface{}) {
	for k, v := range src {
		if s, ok := v.(string); ok && s != "" {
			dst[k] = v
		}
	}
}

func (client *ApiClient) PlaceOrder(clientOid string, side string, symbol string, options map[string]interface{}) (*RespPlaceOrder, error) {
	res := &RespPlaceOrder{}
	params := make(map[string]interface{})
	params["clientOid"] = clientOid
	params["side"] = side
	params["symbol"] = symbol
	UpdateMapWith(params, options)
	err := client.DoRequest(client.ReqUrl.ORDERS_URL, "POST", true, params, res)
	return res, err
}

func (client *ApiClient) GetOrder(orderId string) (*RespGetOrder, error) {
	res := &RespGetOrder{}
	client.ReqUrl.ORDERS_URL.Url = client.ReqUrl.ORDERS_URL.Url + "/" + orderId
	err := client.DoRequest(client.ReqUrl.ORDERS_URL, "GET", true, nil, res)
	return res, err
}

func (client *ApiClient) GetOrderByClientOid(clientOid string) (*RespGetOrder, error) {
	res := &RespGetOrder{}
	client.ReqUrl.ORDER_CLIENT_URL.Url += "/" + clientOid
	err := client.DoRequest(client.ReqUrl.ORDER_CLIENT_URL, "GET", true, nil, res)
	return res, err
}

func (client *ApiClient) GetFills(all bool, options map[string]interface{}) ([]*FillInfo, error) {
	res := make([]*FillInfo, 0)
	resp := &RespGetFills{}

	if !all {
		err := client.DoRequest(client.ReqUrl.FILLS_URL, "GET", true, options, resp)
		if err != nil {
			return res, nil
		}
		res = append(res, resp.Data.Items...)
		return res, err
	}

	newParams := make(map[string]interface{})
	UpdateMapWith(newParams, options)
	if _, ok := newParams["pageSize"]; !ok {
		newParams["pageSize"] = 500
	}

	for {
		err := client.DoRequest(client.ReqUrl.FILLS_URL, "GET", true, newParams, resp)
		if err != nil {
			return res, err
		}
		res = append(res, resp.Data.Items...)
		if resp.Data.CurrentPage >= resp.Data.TotalPage {
			break
		}
		newParams["currentPage"] = resp.Data.CurrentPage + 1
	}

	return res, nil
}

func (client *ApiClient) PlaceMarginOrder(clientOid string, side string, symbol string, options map[string]interface{}) (*RespPlaceMarginOrder, error) {
	res := &RespPlaceMarginOrder{}
	params := make(map[string]interface{})
	params["clientOid"] = clientOid
	params["side"] = side
	params["symbol"] = symbol
	UpdateMapWith(params, options)
	err := client.DoRequest(client.ReqUrl.MARGIN_ORDER_URL, "POST", true, params, res)
	return res, err
}

func (client *ApiClient) GetStatus() (*RespGetStatus, error) {
	res := &RespGetStatus{}
	err := client.DoRequest(client.ReqUrl.STATUS_URL, "GET", false, nil, res)
	return res, err
}

func (client *ApiClient) InnerTransfer(clientOid string, currency string, from string, to string, amount string, options map[string]interface{}) (*RespInnerTransfer, error) {
	res := &RespInnerTransfer{}
	params := make(map[string]interface{})
	params["clientOid"] = clientOid
	params["currency"] = currency
	params["from"] = from
	params["to"] = to
	params["amount"] = amount
	UpdateMapWith(params, options)
	err := client.DoRequest(client.ReqUrl.ACCOUNTS_INNER_TRANSFER_URL, "POST", true, params, res)
	return res, err
}

func (client *ApiClient) SubTransfer(clientOid string, currency string, amount string, direction string, subUserId string, options map[string]interface{}) (*RespSubTransfer, error) {
	res := &RespSubTransfer{}
	params := make(map[string]interface{})
	params["clientOid"] = clientOid
	params["currency"] = currency
	params["amount"] = amount
	params["direction"] = direction
	params["subUserId"] = subUserId
	UpdateMapWith(params, options)
	err := client.DoRequest(client.ReqUrl.ACCOUNTS_SUB_TRANSFER_URL, "POST", true, params, res)
	return res, err
}

func (client *ApiClient) CancelOrderByClientID(clientOid string) (*RespCancelOrder, error) {
	res := &RespCancelOrder{}
	client.ReqUrl.ORDER_CLIENT_URL.Url += "/" + clientOid
	err := client.DoRequest(client.ReqUrl.ORDER_CLIENT_URL, "DELETE", true, nil, res)
	return res, err
}

// todo 返回为空
func (client *ApiClient) GetOrderHistory(params ...url.Values) (*RespOrderHistory, error) {
	res := &RespOrderHistory{}
	//client.ReqUrl.ORDER_HISTORY_URL.Url += parseUrl(client.ReqUrl.ORDER_HISTORY_URL.Url, "?", params...)
	//client.ReqUrl.ORDER_HISTORY_URL.Url = "/api/v1/limit/orders"
	err := client.DoRequest(client.ReqUrl.ORDER_HISTORY_URL, "GET", true, nil, res)
	return res, err
}

func parseUrl(url, sign string, params ...url.Values) string {
	if len(params) > 0 {
		if len(params[0]) > 0 {
			url += sign
			for k, v := range params[0] {
				url += (k + "=" + v[0])
			}
		}
	}
	return url
}

func (client *ApiClient) GetWsToken(isPrivate bool) (*RespGetWsToken, error) {
	res := &RespGetWsToken{}
	var url base.ReqUrlInfo
	if isPrivate {
		url = client.ReqUrl.TOKEN_PRIVATE_URL
	} else {
		url = client.ReqUrl.TOKEN_PUBLIC_URL
	}
	err := client.DoRequest(url, "POST", isPrivate, nil, res)
	return res, err
}

func (client *ApiClient) Withdraw(currency string, address string, amount float64, options ...url.Values) (*RespWithdraw, error) {
	res := &RespWithdraw{}
	params := make(map[string]interface{})
	params["currency"] = currency
	params["address"] = address
	params["amount"] = amount
	mergeUrlValues(params, options...)
	err := client.DoRequest(client.ReqUrl.WITHDRAW_URL, "POST", true, params, res)
	return res, err
}

func (client *ApiClient) WithdrawFee(currency string, options ...url.Values) (*RespWithdrawFee, error) {
	res := &RespWithdrawFee{}
	params := make(map[string]interface{})
	params["currency"] = currency
	err := client.DoRequest(client.ReqUrl.WITHDRAW_FEE_URL, "GET", true, params, res)
	return res, err
}

func mergeUrlValues(params map[string]interface{}, options ...url.Values) map[string]interface{} {
	if len(options) > 0 {
		for k, v := range options[0] {
			params[k] = v
		}
	}
	return params
}

func (client *ApiClient) WithdrawHistory(param ...url.Values) (*RespWithdrawHistory, error) {

	params := make(map[string]interface{})
	params = mergeUrlValues(params, param...)
	res := &RespWithdrawHistory{}
	err := client.DoRequest(client.ReqUrl.WITHDRAW_HIS_URL, "GET", true, params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ApiClient) MoveHistory(options ...url.Values) (*RespMoveHistory, error) {
	params := make(map[string]interface{})
	params = mergeUrlValues(params, options...)
	res := &RespMoveHistory{}
	err := c.DoRequest(c.ReqUrl.MOVE_HISTORY_URL, "GET", true, params, res)
	return res, err
}

func (c *ApiClient) Loan(currency, type_ string, size float64, options ...url.Values) (*RespLoan, error) {
	params := make(map[string]interface{})
	params["currency"] = currency
	params["type"] = type_
	params["size"] = size
	params = mergeUrlValues(params, options...)
	res := &RespLoan{}
	err := c.DoRequest(c.ReqUrl.LOAN_URL, "POST", true, params, res)
	return res, err
}

func (c *ApiClient) Repay(currency, tradId string, size float64, options ...url.Values) (*RespRepay, error) {
	params := make(map[string]interface{})
	params["currency"] = currency
	params["tradId"] = tradId
	params["size"] = size
	params = mergeUrlValues(params, options...)
	res := &RespRepay{}
	err := c.DoRequest(c.ReqUrl.REPAY_URL, "POST", true, params, res)
	return res, err
}

func (c *ApiClient) UnRepayHistory(options ...url.Values) (*RespUnrepayHistory, error) {
	params := make(map[string]interface{})
	params = mergeUrlValues(params, options...)
	res := &RespUnrepayHistory{}
	err := c.DoRequest(c.ReqUrl.UNREPAY_HISTORY_URL, "GET", true, params, res)
	return res, err
}

func (c *ApiClient) RepayHistory(options ...url.Values) (*RespRepayHistory, error) {
	params := make(map[string]interface{})
	params = mergeUrlValues(params, options...)
	res := &RespRepayHistory{}
	err := c.DoRequest(c.ReqUrl.REPAY_HISTORY_URL, "GET", true, params, res)
	return res, err
}

func (client *ApiClient) GetFutureSymbols() (*RespFutureSymbols, error) {
	res := &RespFutureSymbols{}
	err := client.DoRequest(client.ReqUrl.CONTRACT_SYMBOL_URL, "GET", false, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetFutureSpecifySymbol(symbol string) (*RespFutureSymbol, error) {
	res := &RespFutureSymbol{}
	client.ReqUrl.CONTRAC_SPECIFY_SYMBOL_URL.Url += symbol
	err := client.DoRequest(client.ReqUrl.CONTRAC_SPECIFY_SYMBOL_URL, "GET", false, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetFuturePositions() (*RespFuturePositions, error) {
	res := &RespFuturePositions{}
	err := client.DoRequest(client.ReqUrl.CONTRACT_POSITIONS_URL, "GET", true, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) PlaceFutureOrder(clientOid, side, symbol string, leverage int, options ...url.Values) (*RespFuturePlaceOrder, error) {
	res := &RespFuturePlaceOrder{}
	params := make(map[string]interface{})
	params["clientOid"] = clientOid
	params["side"] = side
	params["symbol"] = symbol
	params["leverage"] = leverage
	params = mergeUrlValues(params, options...)
	err := client.DoRequest(client.ReqUrl.CONTRACT_PALCE_ORDER_URL, "GET", true, params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
