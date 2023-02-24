package spot_api

import (
	"bytes"
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/logger"
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
)

// ApiClient To do CHECK
type ApiClient struct {
	base.APIConf
	//info *RespExchangeInfo
	HttpClient *http.Client
	ReqUrl     *ReqUrl
}

// NewApiClient To do CHECK
func NewApiClient(conf base.APIConf) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
		}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)

	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
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
	a.ReqUrl = NewSpotReqUrl()

	return a
}

func NewApiClient2(conf base.APIConf, cli *http.Client) *ApiClient {
	var (
		a = &ApiClient{
			APIConf: conf,
		}
	)

	if conf.EndPoint == "" {
		a.EndPoint = SPOT_API_BASE_URL
	}
	a.HttpClient = cli
	a.ReqUrl = NewSpotReqUrl()

	return a
}

func (client *ApiClient) sign(signaturePayload string) string {
	mac := hmac.New(sha256.New, []byte(client.SecretKey))
	mac.Write([]byte(signaturePayload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (client *ApiClient) signRequest(method string, path string, body []byte) *http.Request {
	ts := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	signaturePayload := ts + method + "/api/" + path + string(body)
	signature := client.sign(signaturePayload)
	req, _ := http.NewRequest(method, client.GetUrl(path), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("FTXUS-KEY", client.AccessKey)
	req.Header.Set("FTXUS-SIGN", signature)
	req.Header.Set("FTXUS-TS", ts)
	//if client.Subaccount != "" {
	//	req.Header.Set("FTXUS-SUBACCOUNT", client.Subaccount)
	//}
	return req
}

func (client *ApiClient) DoRequest(uri, method string, params map[string]interface{}, result interface{}) error {
	urlParam := ""
	parameters := make(map[string]interface{})
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
	preparedRequest := client.signRequest(method, uri+urlParam, body)
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

func (client *ApiClient) GetOrderbook(marketName string, depth int) (*RespGetOrderbook, error) {
	params := make(map[string]interface{})
	params["depth"] = depth
	updatedURL := client.ReqUrl.MARKETS_URL + "/" + marketName + "/orderbook"
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
	updatedURL := client.ReqUrl.ORDERS_URL + "/" + strconv.Itoa(orderId)
	err := client.DoRequest(updatedURL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) GetOrderStatusByClientID(clientId string) (*RespGetOrderStatus, error) {
	params := make(map[string]interface{})
	res := &RespGetOrderStatus{}
	updatedURL := client.ReqUrl.ORDERS_BY_CLIENT_ID_URL + "/" + clientId
	err := client.DoRequest(updatedURL, "GET", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) CancelOrder(orderId int) (*RespCancelOrder, error) {
	params := make(map[string]interface{})
	updatedURL := client.ReqUrl.ORDERS_URL + "/" + strconv.Itoa(orderId)
	res := &RespCancelOrder{}
	err := client.DoRequest(updatedURL, "DELETE", params, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (client *ApiClient) CancelOrderByClientID(clientId string) (*RespCancelOrder, error) {
	params := make(map[string]interface{})
	updatedURL := client.ReqUrl.ORDERS_BY_CLIENT_ID_URL + "/" + clientId
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
	for key, element := range *options {
		params[key] = element
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
