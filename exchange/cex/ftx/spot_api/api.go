package spot_api

import (
	"bytes"
	"clients/config"
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/warmplanet/proto/go/client"
)

// ApiClient To do CHECK
type ApiClient struct {
	base.APIConf
	//info *RespExchangeInfo
	HttpClient *http.Client
	WeightInfo map[client.WeightType]*client.WeightInfo
	ReqUrl     *ReqUrl
}

var (
	FTXWeight *base.RateLimitMgr
)

func WeightAllow(consumeMap []*base.RateLimitConsume) error {
	for _, rate := range consumeMap {
		if err := FTXWeight.Consume(*rate); err != nil {
			return err
		}
	}
	return nil
}

// NewApiClient To do CHECK
func NewApiClient(conf base.APIConf, maps ...interface{}) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
		}
		proxyUrl      *url.URL
		transport     http.Transport
		err           error
		weightInfoMap map[string]int64
	)

	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
	}

	a.WeightInfo = make(map[client.WeightType]*client.WeightInfo)

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
	if FTXWeight == nil {
		FTXWeight = base.NewRateLimitMgr()
	}
	a.ReqUrl = NewSpotReqUrl(weightInfoMap, a.SubAccount)

	return a
}

func NewApiClient2(conf base.APIConf, cli *http.Client, maps ...interface{}) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
		}
		weightInfoMap map[string]int64
	)

	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
	}
	a.WeightInfo = make(map[client.WeightType]*client.WeightInfo)
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
	if FTXWeight == nil {
		FTXWeight = base.NewRateLimitMgr()
	}
	a.ReqUrl = NewSpotReqUrl(weightInfoMap, a.SubAccount)

	return a
}

func (client *ApiClient) sign(signaturePayload string) string {
	mac := hmac.New(sha256.New, []byte(client.SecretKey))
	mac.Write([]byte(signaturePayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (client *ApiClient) signRequest(method string, path string, body []byte, subAccount string) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := client.sign(signaturePayload)
	req, _ := http.NewRequest(method, client.GetUrl(path), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTX-KEY", client.AccessKey)
	req.Header.Set("FTX-SIGN", signature)
	req.Header.Set("FTX-TS", ts)
	if subAccount != "" {
		req.Header.Set("FTX-SUBACCOUNT", transform.Quote(subAccount))
	}
	return req
}

func ClearCount(consumeMap []*base.RateLimitConsume) {
	for _, rate := range consumeMap {
		tmp := FTXWeight.GetInstance(rate.LimitTypeName)
		tmp.Instance.Clear()
	}
}

func (client *ApiClient) DoRequest(uri base.ReqUrlInfo, method string, params map[string]interface{}, result interface{}) error {
	urlParam := ""
	parameters := make(map[string]interface{})
	if err := WeightAllow(client.ReqUrl.GetURLRateLimit(uri, params)); err != nil {
		return err
	}

	subAccount := ""
	// 如果配置的SubAccount形式为xxx.yyy，则yyy为子账户名称，其余形式不是子账户
	idx := strings.Index(client.SubAccount, ".")
	if idx >= 0 {
		subAccount = client.SubAccount[idx+1:]
	}

	// 如果用户设置了子账户，则覆盖默认配置的子账户
	if v, ok := params["_subAccount"]; ok {
		subAccount = v.(string)
		delete(params, "_subAccount")
	}

	if len(params) != 0 && (method == "GET" || method == "DELETE") {
		urlParam = "?"
		for key, element := range params {
			parameters[key] = element
			urlParam += key + "=" + fmt.Sprintf("%v", element) + "&"
		}

	} else if method == "POST" {
		for key, element := range params {
			parameters[key] = element
		}
	}

	body, _ := json.Marshal(parameters)
	preparedRequest := client.signRequest(method, uri.Url+urlParam, body, subAccount)
	rsp, err := conn.DoFTXRequest(client.HttpClient, preparedRequest)
	url_ := preparedRequest.URL

	if err == nil {
		if strings.Contains(string(rsp), `"success":false`) {
			var errResp RespError
			err = json.Unmarshal(rsp, &errResp)
			if err == nil {
				err = errors.New(errResp.Error)
			}
		} else {
			err = json.Unmarshal(rsp, result)
		}
		logger.Logger.Debug(url_, params, err, string(rsp))
		return err
	}
	logger.Logger.Debug(url_, params, err, string(rsp))
	// err not nil
	if v, ok := err.(*conn.HttpError); ok {
		if v.Code == 429 {
			ClearCount(client.ReqUrl.GetURLRateLimit(uri, params))
		}
		return &base.ApiError{Code: v.Code, UnknownStatus: v.Unknown, ErrMsg: ""}
	} else {
		return &base.ApiError{Code: 500, BizCode: 0, ErrMsg: err.Error(), UnknownStatus: true}
	}
}

func (client *ApiClient) GetUrl(url string) string {
	return client.EndPoint + url
}

func (client *ApiClient) GetMarkets() (*RespGetMarkets, error) {
	params := make(map[string]interface{})
	res := &RespGetMarkets{}
	err := client.DoRequest(client.ReqUrl.MARKETS_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetAccountInfo() (*RespGetAccountInfo, error) {
	params := make(map[string]interface{})
	res := &RespGetAccountInfo{}
	err := client.DoRequest(client.ReqUrl.ACCOUNT_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetPositions() (*RespGetPositions, error) {
	params := make(map[string]interface{})
	res := &RespGetPositions{}
	err := client.DoRequest(client.ReqUrl.POSITIONS_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderbook(marketName string, depth int) (*RespGetOrderbook, error) {
	params := make(map[string]interface{})
	params["depth"] = depth

	updatedURL := base.ReqUrlInfo{
		Url:                 client.ReqUrl.MARKETS_URL.Url + "/" + marketName + "/orderbook",
		RateLimitConsumeMap: client.ReqUrl.MARKETS_URL.RateLimitConsumeMap,
	}
	res := &RespGetOrderbook{}
	err := client.DoRequest(updatedURL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetBalances() (*RespGetBalances, error) {
	params := make(map[string]interface{})
	res := &RespGetBalances{}
	err := client.DoRequest(client.ReqUrl.WALLET_BALANCES_URL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) PlaceOrder(market string, side SideType, type_ OrderType, size float64, options *map[string]interface{}) (*RespPlaceOrder, error) {
	params := make(map[string]interface{})
	params["market"] = market
	params["side"] = string(side)
	params["type"] = string(type_)
	params["size"] = size
	//Optional: reduceOnly, ioc (immediate or cancel), postOnly, clientId, rejectOnPriceBand
	for key, element := range *options {
		params[key] = element
	}
	res := &RespPlaceOrder{}
	err := client.DoRequest(client.ReqUrl.ORDERS_URL, "POST", params, res)
	return res, err
}

func (client *ApiClient) GetOrderHistory(options *map[string]interface{}) (*RespGetOrderHistory, error) {
	res := &RespGetOrderHistory{}
	//options are market (string), side (string), orderType (string), start_time (int unix), end_time (int unix)
	err := client.DoRequest(client.ReqUrl.ORDERS_HISTORY_URL, "GET", *options, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOpenOrders(options *map[string]interface{}) (*RespGetOpenOrders, error) {
	res := &RespGetOpenOrders{}
	//option is market
	err := client.DoRequest(client.ReqUrl.ORDERS_URL, "GET", *options, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderStatus(orderId int) (*RespGetOrderStatus, error) {
	params := make(map[string]interface{})
	res := &RespGetOrderStatus{}
	updatedURL := base.ReqUrlInfo{
		Url:                 client.ReqUrl.ORDERS_URL.Url + "/" + strconv.Itoa(orderId),
		RateLimitConsumeMap: client.ReqUrl.MARKETS_URL.RateLimitConsumeMap,
	}
	err := client.DoRequest(updatedURL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderStatusByClientID(clientId string) (*RespGetOrderStatus, error) {
	params := make(map[string]interface{})
	res := &RespGetOrderStatus{}
	updatedURL := base.ReqUrlInfo{
		Url:                 client.ReqUrl.ORDERS_BY_CLIENT_ID_URL.Url + "/" + clientId,
		RateLimitConsumeMap: client.ReqUrl.MARKETS_URL.RateLimitConsumeMap,
	}
	err := client.DoRequest(updatedURL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) CancelOrder(orderId int) (*RespCancelOrder, error) {
	params := make(map[string]interface{})
	updatedURL := base.ReqUrlInfo{
		Url:                 client.ReqUrl.ORDERS_URL.Url + "/" + strconv.Itoa(orderId),
		RateLimitConsumeMap: client.ReqUrl.MARKETS_URL.RateLimitConsumeMap,
	}
	res := &RespCancelOrder{}
	err := client.DoRequest(updatedURL, "DELETE", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) CancelOrderByClientID(clientId string) (*RespCancelOrder, error) {
	params := make(map[string]interface{})
	updatedURL := base.ReqUrlInfo{
		Url:                 client.ReqUrl.ORDERS_BY_CLIENT_ID_URL.Url + "/" + clientId,
		RateLimitConsumeMap: client.ReqUrl.MARKETS_URL.RateLimitConsumeMap,
	}
	res := &RespCancelOrder{}
	err := client.DoRequest(updatedURL, "DELETE", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) CancelAllOrders(options *map[string]interface{}) (*RespCancelOrder, error) {
	params := make(map[string]interface{})
	//market, side, conditionalOrdersOnly, limitOrdersOnly
	for key, element := range *options {
		params[key] = element
	}
	res := &RespCancelOrder{}
	err := client.DoRequest(client.ReqUrl.ORDERS_URL, "DELETE", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Need to check functionality
func (client *ApiClient) RequestWithdrawal(coin string, size float64, address string, options *map[string]interface{}) (*RespRequestWithdrawal, error) {
	params := make(map[string]interface{})
	params["coin"] = coin
	params["size"] = size
	params["address"] = address
	//Options: tag, password, code
	for key, element := range *options {
		params[key] = element
	}
	res := &RespRequestWithdrawal{}
	err := client.DoRequest(client.ReqUrl.WALLET_WITHDRAWALS_URL, "POST", params, res)
	return res, err
}

func (client *ApiClient) GetWithdrawalFee(coin string, size float64, address string, options *map[string]interface{}) (*RespGetWithdrawalFee, error) {
	params := make(map[string]interface{})
	params["coin"] = coin
	params["size"] = size
	params["address"] = address
	//Options: tag, method
	for key, element := range *options {
		params[key] = element
	}
	res := &RespGetWithdrawalFee{}
	err := client.DoRequest(client.ReqUrl.WALLET_WITHDRAWAL_FEE_URL, "GET", params, res)
	return res, err
}

func (client *ApiClient) GetWithdrawalHistory(options *map[string]interface{}) (*RespGetWithdrawalHistory, error) {
	params := make(map[string]interface{})
	//Options: start_time, end_time
	for key, element := range *options {
		params[key] = element
	}
	res := &RespGetWithdrawalHistory{}
	err := client.DoRequest(client.ReqUrl.WALLET_WITHDRAWALS_URL, "GET", params, res)
	return res, err
}

func (client *ApiClient) GetFills(options *map[string]interface{}) (*RespFills, error) {
	res := &RespFills{}
	//option is market
	err := client.DoRequest(client.ReqUrl.FILLS_URL, "GET", *options, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetDepositHistory(options *map[string]interface{}) (*RespGetDepositHistory, error) {
	params := make(map[string]interface{})
	//Options: start_time, end_time
	for key, element := range *options {
		params[key] = element
	}
	res := &RespGetDepositHistory{}
	err := client.DoRequest(client.ReqUrl.WALLET_DEPOSITS_URL, "GET", params, res)
	return res, err
}

// Haven't tested yet
func (client *ApiClient) GetDepositAddress(coin string, options *map[string]interface{}) (*RespGetDepositAddress, error) {
	params := make(map[string]interface{})
	params["coin"] = coin
	//options: method
	if options != nil {
		for key, element := range *options {
			params[key] = element
		}
	}

	res := &RespGetDepositAddress{}
	err := client.DoRequest(client.ReqUrl.WALLET_DEPOSITS_URL, "GET", params, res)
	return res, err
}

func (client *ApiClient) GetLendingHistory(options *map[string]interface{}) (*RespGetLendingHistory, error) {
	params := make(map[string]interface{})
	//Options: start_time, end_time
	for key, element := range *options {
		params[key] = element
	}
	res := &RespGetLendingHistory{}
	err := client.DoRequest(client.ReqUrl.SPOT_MARGIN_HISTORY, "GET", params, res)
	return res, err
}

func (client *ApiClient) GetBorrowRates() (*RespGetBorrowRates, error) {
	params := make(map[string]interface{})
	res := &RespGetBorrowRates{}
	err := client.DoRequest(client.ReqUrl.SPOT_MARGIN_BORROW_RATES, "GET", params, res)
	return res, err
}

func (client *ApiClient) GetLendingRates() (*RespGetLendingRates, error) {
	params := make(map[string]interface{})
	res := &RespGetLendingRates{}
	err := client.DoRequest(client.ReqUrl.SPOT_MARGIN_LENDING_RATES, "GET", params, res)
	return res, err
}

func (client *ApiClient) GetMarketInfo(options *map[string]interface{}) (*RespGetMarketInfo, error) {
	res := &RespGetMarketInfo{}
	//option is market
	err := client.DoRequest(client.ReqUrl.SPOT_MARGIN_MARKET_INFO, "GET", *options, res)
	if err != nil {
		return nil, err
	}
	return res, nil
} //For Spot Margins

func (client *ApiClient) GetLendingOffers() (*RespGetLendingOffers, error) {
	params := make(map[string]interface{})
	res := &RespGetLendingOffers{}
	err := client.DoRequest(client.ReqUrl.SPOT_MARGIN_OFFERS, "GET", params, res)
	return res, err
}

func (client *ApiClient) SubmitLendingOffer(coin string, size float64, rate float64) (*RespSubmitLendingOffer, error) {
	params := make(map[string]interface{})
	params["coin"] = coin
	params["size"] = size
	params["rate"] = rate
	res := &RespSubmitLendingOffer{}
	err := client.DoRequest(client.ReqUrl.SPOT_MARGIN_OFFERS, "POST", params, res)
	return res, err
}

func (client *ApiClient) TransferSubAccount(coin string, size float64, source string, destination string) (*RespTransferSubaccount, error) {
	params := make(map[string]interface{})
	params["coin"] = coin
	params["size"] = size
	params["source"] = source
	params["destination"] = destination
	res := &RespTransferSubaccount{}
	err := client.DoRequest(client.ReqUrl.SUBACCOUNTS_TRANSFER, "POST", params, res)
	return res, err
}
