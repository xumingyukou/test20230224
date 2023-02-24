package u_api

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/huobi/c_api"
	"clients/exchange/cex/huobi/spot_api"
	"clients/logger"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type UApiClient struct {
	spot_api.SpotApiClient
	Currency_size map[string]float64
}

func NewUApiClient(conf base.APIConf) *UApiClient {
	var (
		a         = &UApiClient{}
		proxyUrl  *url.URL
		transport http.Transport
		err       error
	)
	a.APIConf = conf
	a.EndPoint = spot_api.U_BASE_API_URL
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

	ResContractInfo, err1 := a.GetContractInfo("swap")
	if err1 != nil {
		fmt.Println("a.GetContractInfo err =", err1)
	}
	a.Currency_size = make(map[string]float64)
	for _, DataContractInfo := range ResContractInfo.Data {
		a.Currency_size[DataContractInfo.Symbol] = DataContractInfo.ContractSize
	}

	return a
}

func (u *UApiClient) GetFundingFeeRate() (*FundingFeeRateResponse, error) {
	res := &FundingFeeRateResponse{}
	err := u.DoRequest(u.ReqUrl.U_BASE_SWAP_BATCH_FUNDING_RATE_URL, "GET", url.Values{}, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u *UApiClient) GetContractInfo(businessType string) (*ContractInfo, error) {
	res := &ContractInfo{}
	param := url.Values{}
	param.Add("business_type", businessType)
	err := u.DoRequest(u.ReqUrl.U_BASE_CONTRACT_INFO_URL, "GET", param, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u *UApiClient) GetDepth(contract_code, type_ string) (*c_api.DepthInfo, error) {
	options := url.Values{}
	options.Add("contract_code", contract_code)
	options.Add("type", type_)
	res := &c_api.DepthInfo{}
	err := u.DoRequest(u.ReqUrl.U_DEPTH_URL, "GET", options, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 全仓持仓量限制的查询，用于查询标记价格
func (u *UApiClient) PostSwapCrossPositionLimit(post_params ReqPostSwapCrossPositionLimit) (*ResPostSwapCrossPositionLimit, error) {
	res := &ResPostSwapCrossPositionLimit{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CROSS_POSITION_LIMIT, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 用户当前的手续费费率
func (u *UApiClient) PostSwapFee(post_params ReqPostSwapFee) (*ResPostSwapFee, error) {
	res := &ResPostSwapFee{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_FEE, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u *UApiClient) PostSwapBalanceValuation(post_params ReqPostSwapBalanceValuation) (*ResPostSwapBalanceValuation, error) {
	res := &ResPostSwapBalanceValuation{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_BALANCE_VALUATION, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 全仓合约下单
func (u *UApiClient) PostSwapCrossOrder(post_params ReqPostSwapCrossOrder) (*ResPostSwapCrossOrder, error) {
	res := &ResPostSwapCrossOrder{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CROSS_ORDER, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 全仓撤销订单
func (u *UApiClient) PostSwapCrossCancelOrder(post_params ReqPostSwapCrossCancelOrder) (*ResPostSwapCrossCancelOrder, error) {
	res := &ResPostSwapCrossCancelOrder{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CROSS_CANCEL_ORDER, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 获取全仓订单信息
func (u *UApiClient) PostSwapCrossOrderInfo(post_params ReqPostSwapCrossOrderInfo) (*ResPostSwapCrossOrderInfo, error) {
	res := &ResPostSwapCrossOrderInfo{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CROSS_ORDER_INFO, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 获取历史成交记录
func (u *UApiClient) PostSwapCrossMatchresults(post_params ReqPostSwapCrossMatchresults) (*ResPostSwapCrossMatchresults, error) {
	res := &ResPostSwapCrossMatchresults{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CROSS_MATCHRESULTS, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 获取当前未成交委托
func (u *UApiClient) PostSwapOpenOrders(post_params ReqPostSwapOpenOrders) (*ResPostSwapOpenOrders, error) {
	res := &ResPostSwapOpenOrders{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_OPEN_ORDERS, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 全仓用户账户信息
func (u *UApiClient) PostSwapCrossAccountInfo(post_params ReqPostSwapCrossAccountInfo) (*ResPostSwapCrossAccountInfo, error) {
	res := &ResPostSwapCrossAccountInfo{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CROSS_ACCOUNT_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 全仓账户持仓信息
func (u *UApiClient) PostSwapCrossPositionInfo(post_params ReqPostSwapCrossPositionInfo) (*ResPostSwapCrossPositionInfo, error) {
	res := &ResPostSwapCrossPositionInfo{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CROSS_POSITION_INFO_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 获取合约信息
func (u *UApiClient) GetSwapContractInfo(business_type string) (*ResGetSwapContractInfo, error) {
	options := url.Values{}
	options.Add("business_type", business_type)
	res := &ResGetSwapContractInfo{}
	err := u.DoRequest(u.ReqUrl.U_SWAP_CONTRACT_INFO_URL, "GET", options, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

// 获取标记价格的K线数据
func (u *UApiClient) GetLinerSwapMarkPriceKline(contract_code string, period string, size int64) (*ResGetLinerSwapMarkPriceKline, error) {
	options := url.Values{}
	options.Add("contract_code", contract_code)
	options.Add("period", period)
	options.Add("size", strconv.FormatInt(size, 10))
	res := &ResGetLinerSwapMarkPriceKline{}
	err := u.DoRequest(u.ReqUrl.U_LINEAR_SWAP_MARK_PRICE_KLINE_URL, "GET", options, spot_api.EmptyPostParams{}, res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (u *UApiClient) PostSwapSwitchAccountType(post_params ReqPostSwapSwitchAccountType) (*ResPostSwapSwitchAccountType, error) {
	res := &ResPostSwapSwitchAccountType{}
	err := u.DoRequest(u.ReqUrl.U_BASE_SWAP_SWITCH_ACCOUNT_TYPE_URL, "POST", url.Values{}, post_params, res)
	if err != nil {
		return nil, err
	}
	return res, err
}
