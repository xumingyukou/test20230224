package c_api

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/huobi/spot_api"
	"clients/logger"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type CApiClient struct {
	spot_api.SpotApiClient
}

// 币本位合约

func NewCApiClient(conf base.APIConf) *CApiClient {
	var (
		a         = &CApiClient{}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)
	a.APIConf = conf
	a.EndPoint = spot_api.C_BASE_API_URL
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
	a.ReqUrl = spot_api.NewSpotReqUrl()

	return a
}

// 交割合约系统状态
func (u *CApiClient) GetFutureContractState(options *url.Values) (*ResGetFutureContractState, error) {
	params := url.Values{}
	spot_api.ParseOptions(options, &params)
	res := &ResGetFutureContractState{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_STATE, "GET", params, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约信息，返回值中有合约价格最小变动精度，可能可以用到
func (u *CApiClient) GetFutureContractInfo(options *url.Values) (*ContractInfo, error) {
	params := url.Values{}
	spot_api.ParseOptions(options, &params)
	res := &ContractInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_INFO_URL, "GET", params, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约指数信息
func (u *CApiClient) GetFutureContractIndex(options *url.Values) (*ResGetFutureContractIndex, error) {
	params := url.Values{}
	spot_api.ParseOptions(options, &params)
	res := &ResGetFutureContractIndex{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_INDEX_URL, "GET", params, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约depth
func (u *CApiClient) GetFutureDepth(symbol, type_ string) (*DepthInfo, error) {
	options := url.Values{}
	options.Add("symbol", symbol)
	options.Add("type", type_)
	res := &DepthInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_DEPTH_URL, "GET", options, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约手续费率
func (u *CApiClient) PostFutureContractFee(post_params ReqPostFutureContractFee) (*ResPostFutureContractFee, error) {
	res := &ResPostFutureContractFee{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_FEE_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约账户总资产估值
func (u *CApiClient) PostFutureBalanceValuation(post_params ReqPostFutureBalanceValuation) (*ResPostFutureBalanceValuation, error) {
	res := &ResPostFutureBalanceValuation{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_BALANCE_VALUATION_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约账户信息
func (u *CApiClient) PostFutureContractAccountInfo(post_params ReqPostFutureContractAccountInfo) (*ResPostFutureContractAccountInfo, error) {
	res := &ResPostFutureContractAccountInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_ACCOUNT_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约持仓信息
func (u *CApiClient) PostFutureContractPositionInfo(post_params ReqPostFutureContractPositionInfo) (*ResPostFutureContractPositionInfo, error) {
	res := &ResPostFutureContractPositionInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_POSITION_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约下单
func (u *CApiClient) PostFutureOrder(post_params ReqPostFutureOrder) (*ResPostFutureOrder, error) {
	res := &ResPostFutureOrder{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_ORDER_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约取消订单
func (u *CApiClient) PostFutureCancelOrder(post_params ReqPostFutureCancelOrder) (*ResPostFutureCancelOrder, error) {
	res := &ResPostFutureCancelOrder{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CANCEL_ORDER_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约获取订单信息
func (u *CApiClient) PostFutureContractOrderInfo(post_params ReqPostFutureContractOrderInfo) (*ResPostFutureContractOrderInfo, error) {
	res := &ResPostFutureContractOrderInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_ORDER_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约获取历史成交记录
func (u *CApiClient) PostFutureContractMatchresults(post_params ReqPostFutureContractMatchresults) (*ResPostFutureContractMatchresults, error) {
	res := &ResPostFutureContractMatchresults{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_MATCHRESULTS_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 交割合约未成交委托
func (u *CApiClient) PostFutureContractOpenorders(post_params ReqPostFutureContractOpenorders) (*ResPostFutureContractOpenorders, error) {
	res := &ResPostFutureContractOpenorders{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_OPENORDERS_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约系统状态
func (u *CApiClient) GetSwapContractState(options *url.Values) (*ResGetSwapContractState, error) {
	params := url.Values{}
	spot_api.ParseOptions(options, &params)
	res := &ResGetSwapContractState{}
	err := u.DoRequest(u.ReqUrl.C_BASE_FUTURE_CONTRACT_STATE, "GET", params, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约信息，返回值中有合约价格最小变动精度，可能可以用到
func (u *CApiClient) GetSwapContractInfo(options *url.Values) (*ContractInfo, error) {
	params := url.Values{}
	spot_api.ParseOptions(options, &params)
	res := &ContractInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_CONTRACT_INFO_URL, "GET", params, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约用户账户信息
func (u *CApiClient) PostSwapAccountInfo(post_params ReqPostSwapAccountInfo) (*ResPostSwapAccountInfo, error) {
	res := &ResPostSwapAccountInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_ACCOUNT_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约用户持仓信息
func (u *CApiClient) PostSwapPositionInfo(post_params ReqPostSwapPositionInfo) (*ResPostSwapPositionInfo, error) {
	res := &ResPostSwapPositionInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_POSITION_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约指数信息
func (u *CApiClient) GetSwapContractIndex(symbol string) (*ResGetSwapContractIndex, error) {
	options := url.Values{}
	options.Add("symbol", symbol)
	res := &ResGetSwapContractIndex{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_CONTRACT_INDEX_URL, "GET", options, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约深度信息
func (u *CApiClient) GetSwapDepth(contract_code, type_ string) (*DepthInfo, error) {
	options := url.Values{}
	options.Add("contract_code", contract_code)
	options.Add("type", type_)
	res := &DepthInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_DEPTH_URL, "GET", options, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约资金费率
func (u *CApiClient) GetSwapFundingFeeRate() (*FundingFeeRateResponse, error) {
	res := &FundingFeeRateResponse{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_BATCH_FUNDING_RATE_URL, "GET", url.Values{}, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约账户总资产估值
func (u *CApiClient) PostSwapBalanceValuation(post_params ReqPostSwapBalanceValuation) (*ResPostSwapBalanceValuation, error) {
	res := &ResPostSwapBalanceValuation{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_BALANCE_VALUATION_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约下单
func (u *CApiClient) PostSwapOrder(post_params ReqPostSwapOrder) (*ResPostSwapOrder, error) {
	res := &ResPostSwapOrder{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_ORDER_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约取消订单
func (u *CApiClient) PostSwapCancelOrder(post_params ReqPostSwapCancelOrder) (*ResPostSwapCancelOrder, error) {
	res := &ResPostSwapCancelOrder{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_CANCEL_ORDER_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约获取订单信息
func (u *CApiClient) PostSwapContractOrderInfo(post_params ReqPostSwapContractOrderInfo) (*ResPostSwapContractOrderInfo, error) {
	res := &ResPostSwapContractOrderInfo{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_CONTRACT_ORDER_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约获取历史成交记录
func (u *CApiClient) PostSwapContractMatchresults(post_params ReqPostSwapContractMatchresults) (*ResPostSwapContractMatchresults, error) {
	res := &ResPostSwapContractMatchresults{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_CONTRACT_MATCHRESULTS_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约未成交委托
func (u *CApiClient) PostSwapContractOpenorders(post_params ReqPostSwapContractOpenorders) (*ResPostSwapContractOpenorders, error) {
	res := &ResPostSwapContractOpenorders{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_CONTRACT_OPENORDERS_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约手续费率
func (u *CApiClient) PostSwapContractFee(post_params ReqPostSwapContractFee) (*ResPostSwapContractFee, error) {
	res := &ResPostSwapContractFee{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_CONTRACT_FEE_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 永续合约标记价格的K线数据
func (u *CApiClient) GetSwapMarkPriceKline(contract_code string, period string, size int64) (*ResGetSwapMarkPriceKline, error) {
	options := url.Values{}
	options.Add("contract_code", contract_code)
	options.Add("period", period)
	options.Add("size", strconv.FormatInt(size, 10))
	res := &ResGetSwapMarkPriceKline{}
	err := u.DoRequest(u.ReqUrl.C_BASE_SWAP_MARK_PRICE_KLINE_URL, "GET", options, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}
