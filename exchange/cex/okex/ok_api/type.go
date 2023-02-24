package ok_api

import (
	"encoding/json"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/order"
)

// 每个都继承了RespError表示都是实现了error。当code不为空时就返回error信息
type Resp_ServerTime struct {
	RespError
	Data []*struct {
		Ts string `json:"ts"`
	} `json:"data"`
}

type Resp_MarkPrice struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []*struct {
		InstType string `json:"instType"`
		InstId   string `json:"instId"`
		MarkPx   string `json:"markPx"`
		Ts       string `json:"ts"`
	} `json:"data"`
}

type RespError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func (r *RespError) Error() string {
	if r.Code == "0" {
		return ""
	} else {
		v, _ := json.Marshal(r)
		return string(v)
	}
}

type InstrumentsData struct {
	InstType     string `json:"instType"`
	InstId       string `json:"instId"`
	Uly          string `json:"uly"`
	Category     string `json:"category"`
	BaseCcy      string `json:"baseCcy"`
	QuoteCcy     string `json:"quoteCcy"`
	SettleCcy    string `json:"settleCcy"`
	CtVal        string `json:"ctVal"`
	CtMult       string `json:"ctMult"`
	CtValCcy     string `json:"ctValCcy"`
	OptType      string `json:"optType"`
	Stk          string `json:"stk"`
	ListTime     string `json:"listTime"`
	ExpTime      string `json:"expTime"`
	Lever        string `json:"lever"`
	TickSz       string `json:"tickSz"`
	LotSz        string `json:"lotSz"`
	MinSz        string `json:"minSz"`
	CtType       string `json:"ctType"`
	Alias        string `json:"alias"`
	State        string `json:"state"`
	MaxLmtSz     string `json:"maxLmtSz"`
	MaxMktSz     string `json:"maxMktSz"`
	MaxTwapSz    string `json:"maxTwapSz"`
	MaxIcebergSz string `json:"maxIcebergSz"`
	MaxTriggerSz string `json:"maxTriggerSz"`
	MaxStopSz    string `json:"maxStopSz"`
}

type RespInstruments struct {
	RespError
	Data []*InstrumentsData `json:"data"`
}

type Delivery_Exercise_HistoryData struct {
	Ts      string `json:"ts"`
	Details []*struct {
		Type   string `json:"type"`
		InstId string `json:"instId"`
		Px     string `json:"px"`
	} `json:"details"`
}

type Delivery_Exercise_History struct {
	RespError
	Data []*struct {
		Ts      string `json:"ts"`
		Details []*struct {
			Type   string `json:"type"`
			InstId string `json:"instId"`
			Px     string `json:"px"`
		} `json:"details"`
	} `json:"data"`
}

type Resp_Open_Interest struct {
	RespError
	Data []*struct {
		InstType string `json:"instType"`
		InstId   string `json:"instId"`
		Oi       string `json:"oi"`
		OiCcy    string `json:"oiCcy"`
		Ts       string `json:"ts"`
	} `json:"data"`
}

type DepthDate struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Ts   string     `json:"ts"`
}

type Resp_Market_Books struct {
	RespError
	Data []*DepthDate `json:"data"`
}

type Resp_Market_Trades struct {
	RespError
	Data []*struct {
		InstId  string `json:"instId"`
		Side    string `json:"side"`
		Sz      string `json:"sz"`
		Px      string `json:"px"`
		TradeId string `json:"tradeId"`
		Ts      string `json:"ts"`
	} `json:"data"`
}

type Resp_Market_IndexTickers struct {
	RespError
	Data []struct {
		InstId  string `json:"instId"`
		IdxPx   string `json:"idxPx"`
		High24H string `json:"high24h"`
		SodUtc0 string `json:"sodUtc0"`
		Open24H string `json:"open24h"`
		Low24H  string `json:"low24h"`
		SodUtc8 string `json:"sodUtc8"`
		Ts      string `json:"ts"`
	} `json:"data"`
}

type Resp_Market_HistoryTrades struct {
	RespError
	Data []*struct {
		InstId  string `json:"instId"`
		Side    string `json:"side"`
		Sz      string `json:"sz"`
		Px      string `json:"px"`
		TradeId string `json:"tradeId"`
		Ts      string `json:"ts"`
	} `json:"data"`
}

type Resp_Account_TradeFee struct {
	RespError
	Data []*struct {
		Category string `json:"category"`
		Delivery string `json:"delivery"`
		Exercise string `json:"exercise"`
		InstType string `json:"instType"`
		Level    string `json:"level"`
		Maker    string `json:"maker"`
		MakerU   string `json:"makerU"`
		Taker    string `json:"taker"`
		TakerU   string `json:"takerU"`
		Ts       string `json:"ts"`
	} `json:"data"`
}

type MarginAssetItem struct {
	AvailBal      string `json:"availBal"`
	AvailEq       string `json:"availEq"`
	CashBal       string `json:"cashBal"`
	Ccy           string `json:"ccy"`
	CrossLiab     string `json:"crossLiab"`
	DisEq         string `json:"disEq"`
	Eq            string `json:"eq"`
	EqUsd         string `json:"eqUsd"`
	FrozenBal     string `json:"frozenBal"`
	Interest      string `json:"interest"`
	IsoEq         string `json:"isoEq"`
	IsoLiab       string `json:"isoLiab"`
	IsoUpl        string `json:"isoUpl"`
	Liab          string `json:"liab"`
	MaxLoan       string `json:"maxLoan"`
	MgnRatio      string `json:"mgnRatio"`
	NotionalLever string `json:"notionalLever"`
	OrdFrozen     string `json:"ordFrozen"`
	Twap          string `json:"twap"`
	UTime         string `json:"uTime"`
	Upl           string `json:"upl"`
	UplLiab       string `json:"uplLiab"`
	StgyEq        string `json:"stgyEq"`
}

type Resp_Accout_Balance struct {
	RespError
	Data []*struct {
		AdjEq       string             `json:"adjEq"`
		Details     []*MarginAssetItem `json:"details"`
		Imr         string             `json:"imr"`
		IsoEq       string             `json:"isoEq"`
		MgnRatio    string             `json:"mgnRatio"`
		Mmr         string             `json:"mmr"`
		NotionalUsd string             `json:"notionalUsd"`
		OrdFroz     string             `json:"ordFroz"`
		TotalEq     string             `json:"totalEq"`
		UTime       string             `json:"uTime"`
	} `json:"data"`
}

type Resp_Asset_Balances struct {
	RespError
	Data []*struct {
		AvailBal  string `json:"availBal"`
		Bal       string `json:"bal"`
		Ccy       string `json:"ccy"`
		FrozenBal string `json:"frozenBal"`
	} `json:"data"`
}

type Resp_Account_Position struct {
	RespError
	Data []*struct {
		Adl         string `json:"adl"`
		AvailPos    string `json:"availPos"`
		AvgPx       string `json:"avgPx"`
		CTime       string `json:"cTime"`
		Ccy         string `json:"ccy"`
		DeltaBS     string `json:"deltaBS"`
		DeltaPA     string `json:"deltaPA"`
		GammaBS     string `json:"gammaBS"`
		GammaPA     string `json:"gammaPA"`
		Imr         string `json:"imr"`
		InstId      string `json:"instId"`
		InstType    string `json:"instType"`
		Interest    string `json:"interest"`
		Last        string `json:"last"`
		UsdPx       string `json:"usdPx"`
		Lever       string `json:"lever"`
		Liab        string `json:"liab"`
		LiabCcy     string `json:"liabCcy"`
		LiqPx       string `json:"liqPx"`
		MarkPx      string `json:"markPx"`
		Margin      string `json:"margin"`
		MgnMode     string `json:"mgnMode"`
		MgnRatio    string `json:"mgnRatio"`
		Mmr         string `json:"mmr"`
		NotionalUsd string `json:"notionalUsd"`
		OptVal      string `json:"optVal"`
		PTime       string `json:"pTime"`
		Pos         string `json:"pos"`
		PosCcy      string `json:"posCcy"`
		PosId       string `json:"posId"`
		PosSide     string `json:"posSide"`
		ThetaBS     string `json:"thetaBS"`
		ThetaPA     string `json:"thetaPA"`
		TradeId     string `json:"tradeId"`
		QuoteBal    string `json:"quoteBal"`
		BaseBal     string `json:"baseBal"`
		UTime       string `json:"uTime"`
		Upl         string `json:"upl"`
		UplRatio    string `json:"uplRatio"`
		VegaBS      string `json:"vegaBS"`
		VegaPA      string `json:"vegaPA"`
	} `json:"data"`
}

type Resp_Account_PositionsHistory struct {
	RespError
	Data []*struct {
		CTime         string `json:"cTime"`
		Ccy           string `json:"ccy"`
		CloseAvgPx    string `json:"closeAvgPx"`
		CloseTotalPos string `json:"closeTotalPos"`
		InstId        string `json:"instId"`
		InstType      string `json:"instType"`
		Lever         string `json:"lever"`
		MgnMode       string `json:"mgnMode"`
		OpenAvgPx     string `json:"openAvgPx"`
		OpenMaxPos    string `json:"openMaxPos"`
		Pnl           string `json:"pnl"`
		PnlRatio      string `json:"pnlRatio"`
		PosId         string `json:"posId"`
		PosSide       string `json:"posSide"`
		TriggerPx     string `json:"triggerPx"`
		Type          string `json:"type"`
		UTime         string `json:"uTime"`
		Uly           string `json:"uly"`
	} `json:"data"`
}

type Resp_Account_AccountPositionRisk struct {
	RespError
	Data []*struct {
		AdjEq   string `json:"adjEq"`
		BalData []*struct {
			Ccy   string `json:"ccy"`
			DisEq string `json:"disEq"`
			Eq    string `json:"eq"`
		} `json:"balData"`
		PosData []*struct {
			BaseBal     string `json:"baseBal"`
			Ccy         string `json:"ccy"`
			InstId      string `json:"instId"`
			InstType    string `json:"instType"`
			MgnMode     string `json:"mgnMode"`
			NotionalCcy string `json:"notionalCcy"`
			NotionalUsd string `json:"notionalUsd"`
			Pos         string `json:"pos"`
			PosCcy      string `json:"posCcy"`
			PosId       string `json:"posId"`
			PosSide     string `json:"posSide"`
			QuoteBal    string `json:"quoteBal"`
		} `json:"posData"`
		Ts string `json:"ts"`
	} `json:"data"`
}

type Resp_Account_Bills struct {
	RespError
	Data []*struct {
		Bal       string `json:"bal"`
		BalChg    string `json:"balChg"`
		BillId    string `json:"billId"`
		Ccy       string `json:"ccy"`
		ExecType  string `json:"execType"`
		Fee       string `json:"fee"`
		From      string `json:"from"`
		InstId    string `json:"instId"`
		InstType  string `json:"instType"`
		MgnMode   string `json:"mgnMode"`
		Notes     string `json:"notes"`
		OrdId     string `json:"ordId"`
		Pnl       string `json:"pnl"`
		PosBal    string `json:"posBal"`
		PosBalChg string `json:"posBalChg"`
		SubType   string `json:"subType"`
		Sz        string `json:"sz"`
		To        string `json:"to"`
		Ts        string `json:"ts"`
		Type      string `json:"type"`
	} `json:"data"`
}

type Resp_Account_BillsArchive struct {
	RespError
	Data []*struct {
		Bal       string `json:"bal"`
		BalChg    string `json:"balChg"`
		BillId    string `json:"billId"`
		Ccy       string `json:"ccy"`
		ExecType  string `json:"execType"`
		Fee       string `json:"fee"`
		From      string `json:"from"`
		InstId    string `json:"instId"`
		InstType  string `json:"instType"`
		MgnMode   string `json:"mgnMode"`
		Notes     string `json:"notes"`
		OrdId     string `json:"ordId"`
		Pnl       string `json:"pnl"`
		PosBal    string `json:"posBal"`
		PosBalChg string `json:"posBalChg"`
		SubType   string `json:"subType"`
		Sz        string `json:"sz"`
		To        string `json:"to"`
		Ts        string `json:"ts"`
		Type      string `json:"type"`
	} `json:"data"`
}

type Resp_Account_Config struct {
	RespError
	Data []*struct {
		AcctLv     string `json:"acctLv"`
		AutoLoan   bool   `json:"autoLoan"`
		CtIsoMode  string `json:"ctIsoMode"`
		GreeksType string `json:"greeksType"`
		Level      string `json:"level"`
		LevelTmp   string `json:"levelTmp"`
		MgnIsoMode string `json:"mgnIsoMode"`
		PosMode    string `json:"posMode"`
		Uid        string `json:"uid"`
	} `json:"data"`
}

type Resp_Account_SetPositionMode struct {
	RespError
	Data []*struct {
		PosMode string `json:"posMode"`
	} `json:"data"`
}

type Resp_Account_SetLeverage struct {
	RespError
	Data []*struct {
		Lever   string `json:"lever"`
		MgnMode string `json:"mgnMode"`
		InstId  string `json:"instId"`
		PosSide string `json:"posSide"`
	} `json:"data"`
}

type Resp_Account_MaxSize struct {
	RespError
	Data []*struct {
		Ccy     string `json:"ccy"`
		InstId  string `json:"instId"`
		MaxBuy  string `json:"maxBuy"`
		MaxSell string `json:"maxSell"`
	} `json:"data"`
}

type Resp_Account_MaxAvailSize struct {
	RespError
	Data []*struct {
		InstId    string `json:"instId"`
		AvailBuy  string `json:"availBuy"`
		AvailSell string `json:"availSell"`
	} `json:"data"`
}

type Resp_Account_PositionMarginBalance struct {
	RespError
	Data []*struct {
		Amt      string `json:"amt"`
		Ccy      string `json:"ccy"`
		InstId   string `json:"instId"`
		Leverage string `json:"leverage"`
		PosSide  string `json:"posSide"`
		Type     string `json:"type"`
	} `json:"data"`
}

type Resp_Account_LeverageInfo struct {
	RespError
	Data []*struct {
		InstId  string `json:"instId"`
		MgnMode string `json:"mgnMode"`
		PosSide string `json:"posSide"`
		Lever   string `json:"lever"`
	} `json:"data"`
}

type Resp_Account_MaxLoan struct {
	RespError
	Data []*struct {
		InstId  string `json:"instId"`
		MgnMode string `json:"mgnMode"`
		MgnCcy  string `json:"mgnCcy"`
		MaxLoan string `json:"maxLoan"`
		Ccy     string `json:"ccy"`
		Side    string `json:"side"`
	} `json:"data"`
}

type Resp_Account_InterestAccrued struct {
	RespError
	Data []*struct {
		Ccy          string `json:"ccy"`
		InstID       string `json:"instId"`
		Interest     string `json:"interest"`
		InterestRate string `json:"interestRate"`
		Liab         string `json:"liab"`
		MgnMode      string `json:"mgnMode"`
		Ts           string `json:"ts"`
		Type         string `json:"type"`
	} `json:"data"`
}

type Resp_Account_InterestRate struct {
	RespError
	Data []*struct {
		Ccy          string `json:"ccy"`
		InterestRate string `json:"interestRate"`
	} `json:"data"`
}

type Resp_Account_SetGreeks struct {
	RespError
	Data []*struct {
		GreeksType string `json:"greeksType"`
	} `json:"data"`
}

type Resp_Account_SetIsolatedMode struct {
	RespError
	Data []*struct {
		IsoMode string `json:"isoMode"`
	} `json:"data"`
}

type Resp_Account_MaxWithdrawal struct {
	RespError
	Data []*struct {
		Ccy     string `json:"ccy"`
		MaxWd   string `json:"maxWd"`
		MaxWdEx string `json:"maxWdEx"`
	} `json:"data"`
}

type Resp_Account_RiskState struct {
	RespError
	Data []*struct {
		AtRisk    bool          `json:"atRisk"`
		AtRiskIdx []interface{} `json:"atRiskIdx"`
		AtRiskMgn []interface{} `json:"atRiskMgn"`
		Ts        string        `json:"ts"`
	} `json:"data"`
}

type Resp_Account_BorrowRepay struct {
	RespError
	Data []*struct {
		Amt       string `json:"amt"`
		AvailLoan string `json:"availLoan"`
		Ccy       string `json:"ccy"`
		LoanQuota string `json:"loanQuota"`
		PosLoan   string `json:"posLoan"`
		Side      string `json:"side"`
		UsedLoan  string `json:"usedLoan"`
	} `json:"data"`
}

type Resp_Account_BorrowRepayHistory struct {
	RespError
	Data []*struct {
		Ccy        string `json:"ccy"`
		TradedLoan string `json:"tradedLoan"`
		Ts         string `json:"ts"`
		Type       string `json:"type"`
		UsedLoan   string `json:"usedLoan"`
	} `json:"data"`
}

type Resp_Account_InterestLimit struct {
	RespError
	Data []*struct {
		Debt             string `json:"debt"`
		Interest         string `json:"interest"`
		NextDiscountTime string `json:"nextDiscountTime"`
		NextInterestTime string `json:"nextInterestTime"`
		Records          []*struct {
			AvailLoan  string `json:"availLoan"`
			Ccy        string `json:"ccy"`
			Interest   string `json:"interest"`
			LoanQuota  string `json:"loanQuota"`
			PosLoan    string `json:"posLoan"`
			Rate       string `json:"rate"`
			SurplusLmt string `json:"surplusLmt"`
			UsedLmt    string `json:"usedLmt"`
			UsedLoan   string `json:"usedLoan"`
		} `json:"records"`
	} `json:"data"`
}

type Resp_Account_SimulatedMargin struct {
	RespError
	Data []*struct {
		Imr     string `json:"imr"`
		Mmr     string `json:"mmr"`
		Mr1     string `json:"mr1"`
		Mr2     string `json:"mr2"`
		Mr3     string `json:"mr3"`
		Mr4     string `json:"mr4"`
		Mr5     string `json:"mr5"`
		Mr6     string `json:"mr6"`
		Mr7     string `json:"mr7"`
		PosData []*struct {
			Delta       string `json:"delta"`
			Gamma       string `json:"gamma"`
			InstId      string `json:"instId"`
			InstType    string `json:"instType"`
			NotionalUsd string `json:"notionalUsd"`
			Pos         string `json:"pos"`
			Theta       string `json:"theta"`
			Vega        string `json:"vega"`
		} `json:"posData"`
		RiskUnit string `json:"riskUnit"`
		Ts       string `json:"ts"`
	} `json:"data"`
}

type Resp_Account_Greeks struct {
	RespError
	Data []*struct {
		ThetaBS string `json:"thetaBS"`
		ThetaPA string `json:"thetaPA"`
		DeltaBS string `json:"deltaBS"`
		DeltaPA string `json:"deltaPA"`
		GammaBS string `json:"gammaBS"`
		GammaPA string `json:"gammaPA"`
		VegaBS  string `json:"vegaBS"`
		VegaPA  string `json:"vegaPA"`
		Ccy     string `json:"ccy"`
		Ts      string `json:"ts"`
	} `json:"data"`
}

type Resp_Account_PositionTiers struct {
	RespError
	Data []*struct {
		MaxSz   string `json:"maxSz"`
		PosType string `json:"posType"`
		Uly     string `json:"uly"`
	} `json:"data"`
}

type Resp_Trde_Order struct {
	RespError
	Data []*struct {
		ClOrdId string `json:"clOrdId"`
		OrdId   string `json:"ordId"`
		Tag     string `json:"tag"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

func (r *Resp_Trde_Order) Error() string {
	if r.Code == "0" {
		return ""
	} else {
		v, _ := json.Marshal(r)
		return string(v)
	}
}

type Resp_Trde_BatchOrders struct {
	RespError
	Data []*struct {
		ClOrdId string `json:"clOrdId"`
		OrdId   string `json:"ordId"`
		Tag     string `json:"tag"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

func (r *Resp_Trde_BatchOrders) Error() string {
	if r.Code == "0" {
		return ""
	} else {
		v, _ := json.Marshal(r)
		return string(v)
	}
}

type Resp_Trade_CancelOrder struct {
	RespError
	Data []*struct {
		ClOrdId string `json:"clOrdId"`
		OrdId   string `json:"ordId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

func (r *Resp_Trade_CancelOrder) Error() string {
	if r.Code == "0" {
		return ""
	} else {
		v, _ := json.Marshal(r)
		return string(v)
	}
}

type Resp_Trade_CancelBatchOrders struct {
	RespError
	Data []*struct {
		ClOrdId string `json:"clOrdId"`
		OrdId   string `json:"ordId"`
		ReqId   string `json:"reqId"`
		SCode   string `json:"sCode"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

func (r *Resp_Trade_CancelBatchOrders) Error() string {
	if r.Code == "0" {
		return ""
	} else {
		v, _ := json.Marshal(r)
		return string(v)
	}
}

type Resp_Trade_ClosePosition struct {
	RespError
	Data []*struct {
		InstId  string `json:"instId"`
		PosSide string `json:"posSide"`
	} `json:"data"`
}

type Resp_TradeInfo_Info struct {
	RespError
	Data []*struct {
		InstType        string `json:"instType"`
		InstId          string `json:"instId"`
		Ccy             string `json:"ccy"`
		OrdId           string `json:"ordId"`
		ClOrdId         string `json:"clOrdId"`
		Tag             string `json:"tag"`
		Px              string `json:"px"`
		Sz              string `json:"sz"`
		Pnl             string `json:"pnl"`
		OrdType         string `json:"ordType"`
		Side            string `json:"side"`
		PosSide         string `json:"posSide"`
		TdMode          string `json:"tdMode"`
		AccFillSz       string `json:"accFillSz"`
		FillPx          string `json:"fillPx"`
		TradeId         string `json:"tradeId"`
		FillSz          string `json:"fillSz"`
		FillTime        string `json:"fillTime"`
		Source          string `json:"source"`
		State           string `json:"state"`
		AvgPx           string `json:"avgPx"`
		Lever           string `json:"lever"`
		TpTriggerPx     string `json:"tpTriggerPx"`
		TpTriggerPxType string `json:"tpTriggerPxType"`
		TpOrdPx         string `json:"tpOrdPx"`
		SlTriggerPx     string `json:"slTriggerPx"`
		SlTriggerPxType string `json:"slTriggerPxType"`
		SlOrdPx         string `json:"slOrdPx"`
		FeeCcy          string `json:"feeCcy"`
		Fee             string `json:"fee"`
		RebateCcy       string `json:"rebateCcy"`
		Rebate          string `json:"rebate"`
		TgtCcy          string `json:"tgtCcy"`
		Category        string `json:"category"`
		UTime           string `json:"uTime"`
		CTime           string `json:"cTime"`
	} `json:"data"`
}

type Resp_Trade_OrdersPending struct {
	RespError
	Data []*struct {
		AccFillSz       string `json:"accFillSz"`
		AvgPx           string `json:"avgPx"`
		CTime           string `json:"cTime"`
		Category        string `json:"category"`
		Ccy             string `json:"ccy"`
		ClOrdId         string `json:"clOrdId"`
		Fee             string `json:"fee"`
		FeeCcy          string `json:"feeCcy"`
		FillPx          string `json:"fillPx"`
		FillSz          string `json:"fillSz"`
		FillTime        string `json:"fillTime"`
		InstId          string `json:"instId"`
		InstType        string `json:"instType"`
		Lever           string `json:"lever"`
		OrdId           string `json:"ordId"`
		OrdType         string `json:"ordType"`
		Pnl             string `json:"pnl"`
		PosSide         string `json:"posSide"`
		Px              string `json:"px"`
		Rebate          string `json:"rebate"`
		RebateCcy       string `json:"rebateCcy"`
		Side            string `json:"side"`
		SlOrdPx         string `json:"slOrdPx"`
		SlTriggerPx     string `json:"slTriggerPx"`
		SlTriggerPxType string `json:"slTriggerPxType"`
		Source          string `json:"source"`
		State           string `json:"state"`
		Sz              string `json:"sz"`
		Tag             string `json:"tag"`
		TdMode          string `json:"tdMode"`
		TgtCcy          string `json:"tgtCcy"`
		TpOrdPx         string `json:"tpOrdPx"`
		TpTriggerPx     string `json:"tpTriggerPx"`
		TpTriggerPxType string `json:"tpTriggerPxType"`
		TradeId         string `json:"tradeId"`
		UTime           string `json:"uTime"`
	} `json:"data"`
}

type Resp_Trade_OrdersHistory struct {
	RespError
	Data []*struct {
		InstType        string `json:"instType"`
		InstId          string `json:"instId"`
		Ccy             string `json:"ccy"`
		OrdId           string `json:"ordId"`
		ClOrdId         string `json:"clOrdId"`
		Tag             string `json:"tag"`
		Px              string `json:"px"`
		Sz              string `json:"sz"`
		OrdType         string `json:"ordType"`
		Side            string `json:"side"`
		PosSide         string `json:"posSide"`
		TdMode          string `json:"tdMode"`
		AccFillSz       string `json:"accFillSz"`
		FillPx          string `json:"fillPx"`
		TradeId         string `json:"tradeId"`
		FillSz          string `json:"fillSz"`
		FillTime        string `json:"fillTime"`
		Source          string `json:"source"`
		State           string `json:"state"`
		AvgPx           string `json:"avgPx"`
		Lever           string `json:"lever"`
		TpTriggerPx     string `json:"tpTriggerPx"`
		TpTriggerPxType string `json:"tpTriggerPxType"`
		TpOrdPx         string `json:"tpOrdPx"`
		SlTriggerPx     string `json:"slTriggerPx"`
		SlTriggerPxType string `json:"slTriggerPxType"`
		SlOrdPx         string `json:"slOrdPx"`
		FeeCcy          string `json:"feeCcy"`
		Fee             string `json:"fee"`
		RebateCcy       string `json:"rebateCcy"`
		Rebate          string `json:"rebate"`
		TgtCcy          string `json:"tgtCcy"`
		Pnl             string `json:"pnl"`
		Category        string `json:"category"`
		UTime           string `json:"uTime"`
		CTime           string `json:"cTime"`
	} `json:"data"`
}

type Resp_Trade_FillsHistory struct {
	RespError
	Data []*struct {
		InstType string `json:"instType"`
		InstId   string `json:"instId"`
		TradeId  string `json:"tradeId"`
		OrdId    string `json:"ordId"`
		ClOrdId  string `json:"clOrdId"`
		BillId   string `json:"billId"`
		Tag      string `json:"tag"`
		FillPx   string `json:"fillPx"`
		FillSz   string `json:"fillSz"`
		Side     string `json:"side"`
		PosSide  string `json:"posSide"`
		ExecType string `json:"execType"`
		FeeCcy   string `json:"feeCcy"`
		Fee      string `json:"fee"`
		Ts       string `json:"ts"`
	} `json:"data"`
}

type Resp_Trade_OrderAlgo struct {
	RespError
	Data []*struct {
		AlgoId string `json:"algoId"`
		SCode  string `json:"sCode"`
		SMsg   string `json:"sMsg"`
	} `json:"data"`
}

func (r *Resp_Trade_OrderAlgo) Error() string {
	if r.Code == "0" {
		return ""
	} else {
		v, _ := json.Marshal(r)
		return string(v)
	}
}

type Resp_CancelAlgo struct {
	RespError
	Data []*struct {
		AlgoId string `json:"algoId"`
		SCode  string `json:"sCode"`
		SMsg   string `json:"sMsg"`
	} `json:"data"`
}

func (r *Resp_CancelAlgo) Error() string {
	if r.Code == "0" {
		return ""
	} else {
		v, _ := json.Marshal(r)
		return string(v)
	}
}

type Resp_OrdersAlgoPending struct {
	RespError
	Data []*struct {
		InstType        string `json:"instType"`
		InstId          string `json:"instId"`
		OrdId           string `json:"ordId"`
		Ccy             string `json:"ccy"`
		AlgoId          string `json:"algoId"`
		Sz              string `json:"sz"`
		OrdType         string `json:"ordType"`
		Side            string `json:"side"`
		PosSide         string `json:"posSide"`
		TdMode          string `json:"tdMode"`
		TgtCcy          string `json:"tgtCcy"`
		State           string `json:"state"`
		Lever           string `json:"lever"`
		TpTriggerPx     string `json:"tpTriggerPx"`
		TpTriggerPxType string `json:"tpTriggerPxType"`
		TpOrdPx         string `json:"tpOrdPx"`
		SlTriggerPx     string `json:"slTriggerPx"`
		SlTriggerPxType string `json:"slTriggerPxType"`
		TriggerPx       string `json:"triggerPx"`
		TriggerPxType   string `json:"triggerPxType"`
		OrdPx           string `json:"ordPx"`
		ActualSz        string `json:"actualSz"`
		ActualPx        string `json:"actualPx"`
		ActualSide      string `json:"actualSide"`
		PxVar           string `json:"pxVar"`
		PxSpread        string `json:"pxSpread"`
		PxLimit         string `json:"pxLimit"`
		SzLimit         string `json:"szLimit"`
		Tag             string `json:"tag"`
		TimeInterval    string `json:"timeInterval"`
		TriggerTime     string `json:"triggerTime"`
		CallbackRatio   string `json:"callbackRatio"`
		CallbackSpread  string `json:"callbackSpread"`
		ActivePx        string `json:"activePx"`
		MoveTriggerPx   string `json:"moveTriggerPx"`
		CTime           string `json:"cTime"`
	} `json:"data"`
}

type Resp_OrdersAlgoHistory struct {
	RespError
	Data []*struct {
		InstType        string `json:"instType"`
		InstId          string `json:"instId"`
		OrdId           string `json:"ordId"`
		Ccy             string `json:"ccy"`
		AlgoId          string `json:"algoId"`
		Sz              string `json:"sz"`
		OrdType         string `json:"ordType"`
		Side            string `json:"side"`
		PosSide         string `json:"posSide"`
		TdMode          string `json:"tdMode"`
		TgtCcy          string `json:"tgtCcy"`
		State           string `json:"state"`
		Lever           string `json:"lever"`
		TpTriggerPx     string `json:"tpTriggerPx"`
		TpTriggerPxType string `json:"tpTriggerPxType"`
		TpOrdPx         string `json:"tpOrdPx"`
		SlTriggerPx     string `json:"slTriggerPx"`
		SlTriggerPxType string `json:"slTriggerPxType"`
		TriggerPx       string `json:"triggerPx"`
		TriggerPxType   string `json:"triggerPxType"`
		OrdPx           string `json:"ordPx"`
		ActualSz        string `json:"actualSz"`
		ActualPx        string `json:"actualPx"`
		ActualSide      string `json:"actualSide"`
		PxVar           string `json:"pxVar"`
		PxSpread        string `json:"pxSpread"`
		PxLimit         string `json:"pxLimit"`
		SzLimit         string `json:"szLimit"`
		Tag             string `json:"tag"`
		TimeInterval    string `json:"timeInterval"`
		CallbackRatio   string `json:"callbackRatio"`
		CallbackSpread  string `json:"callbackSpread"`
		ActivePx        string `json:"activePx"`
		MoveTriggerPx   string `json:"moveTriggerPx"`
		TriggerTime     string `json:"triggerTime"`
		CTime           string `json:"cTime"`
	} `json:"data"`
}

type Resp_Rfq_CounterParties struct {
	RespError
	Data []*struct {
		TraderName string `json:"traderName"`
		TraderCode string `json:"traderCode"`
		Type       string `json:"type"`
	} `json:"data"`
}

type Resp_Asset_Currencies struct {
	RespError
	Data []*struct {
		CanDep               bool   `json:"canDep"`
		CanInternal          bool   `json:"canInternal"`
		CanWd                bool   `json:"canWd"`
		Ccy                  string `json:"ccy"`
		Chain                string `json:"chain"`
		LogoLink             string `json:"logoLink"`
		MainNet              bool   `json:"mainNet"`
		MaxFee               string `json:"maxFee"`
		MaxWd                string `json:"maxWd"`
		MinDep               string `json:"minDep"`
		MinDepArrivalConfirm string `json:"minDepArrivalConfirm"`
		MinFee               string `json:"minFee"`
		MinWd                string `json:"minWd"`
		MinWdUnlockConfirm   string `json:"minWdUnlockConfirm"`
		Name                 string `json:"name"`
		NeedTag              bool   `json:"needTag"`
		UsedWdQuota          string `json:"usedWdQuota"`
		WdQuota              string `json:"wdQuota"`
		WdTickSz             string `json:"wdTickSz"`
	} `json:"data"`
}

type Resp_Asset_Transfer struct {
	RespError
	Data []*struct {
		TransId  string `json:"transId"`
		Ccy      string `json:"ccy"`
		ClientId string `json:"clientId"`
		From     string `json:"from"`
		Amt      string `json:"amt"`
		To       string `json:"to"`
		State    string `json:"state"`
		Type     string `json:"type"`
	} `json:"data"`
}

type Resp_Asset_Transfer_State struct {
	RespError
	Data []*struct {
		Amt      string `json:"amt"`
		Ccy      string `json:"ccy"`
		ClientId string `json:"clientId"`
		From     string `json:"from"`
		InstId   string `json:"instId"`
		State    string `json:"state"`
		SubAcct  string `json:"subAcct"`
		To       string `json:"to"`
		ToInstId string `json:"toInstId"`
		TransId  string `json:"transId"`
		Type     string `json:"type"`
	} `json:"data"`
}

type Resp_Asset_Bills struct {
	RespError
	Data []*struct {
		BillId   string `json:"billId"`
		Ccy      string `json:"ccy"`
		ClientId string `json:"clientId"`
		BalChg   string `json:"balChg"`
		Bal      string `json:"bal"`
		Type     string `json:"type"`
		Ts       string `json:"ts"`
	} `json:"data"`
}

type Resp_User_SubAccountSetTransferOut struct {
	Code string `json:"code"`
	Data []*struct {
		SubAcct     string `json:"subAcct"`
		CanTransOut bool   `json:"canTransOut"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Resp_Asset_SubAccountBills struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []*struct {
		BillId  string `json:"billId"`
		Type    string `json:"type"`
		Ccy     string `json:"ccy"`
		Amt     string `json:"amt"`
		SubAcct string `json:"subAcct"`
		Ts      string `json:"ts"`
	} `json:"data"`
}

type Resp_Asset_Withdrawal struct {
	RespError
	Data []*struct {
		Amt      string `json:"amt"`
		WdId     string `json:"wdId"`
		Ccy      string `json:"ccy"`
		ClientId string `json:"clientId"`
		Chain    string `json:"chain"`
	} `json:"data"`
}

type Resp_Asset_WithdrawalLightning struct {
	RespError
	Data []*struct {
		WdId  string `json:"wdId"`
		CTime string `json:"cTime"`
	} `json:"data"`
}

type Resp_Asset_CancelWithdrawal struct {
	RespError
	Data []*struct {
		WdId string `json:"wdId"`
	} `json:"data"`
}

type Resp_Asset_WithdrawalHistory struct {
	RespError
	Data []*struct {
		Chain    string `json:"chain"`
		Fee      string `json:"fee"`
		Ccy      string `json:"ccy"`
		ClientId string `json:"clientId"`
		Amt      string `json:"amt"`
		TxId     string `json:"txId"`
		From     string `json:"from"`
		To       string `json:"to"`
		State    string `json:"state"`
		Ts       string `json:"ts"`
		WdId     string `json:"wdId"`
	} `json:"data"`
}

type Resp_Asset_DepositHistory struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []*struct {
		ActualDepBlkConfirm string `json:"actualDepBlkConfirm"`
		Amt                 string `json:"amt"`
		Ccy                 string `json:"ccy"`
		Chain               string `json:"chain"`
		DepId               string `json:"depId"`
		From                string `json:"from"`
		State               string `json:"state"`
		To                  string `json:"to"`
		Ts                  string `json:"ts"`
		TxId                string `json:"txId"`
	} `json:"data"`
}

func OkGetTransferTypeFromExchange(status OkexTransferStatus) client.TransferStatus {
	switch status {
	case OK_TRANSFER_TYPE_CREATED:
		return client.TransferStatus_TRANSFERSTATUS_CREATED
	case OK_TRANSFER_TYPE_WITHDRWAN:
		return client.TransferStatus_TRANSFERSTATUS_CANCELLED
		// 等待确认
	case OK_TRANSFER_TYPE_REVOCATION:
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	//	// 被拒绝
	case OK_TRANSFER_TYPE_CONFIRMING:
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	case OK_TRANSFER_TYPE_PROCESSING:
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	case OK_TRANSFER_TYPE_FAIL:
		return client.TransferStatus_TRANSFERSTATUS_FAILED
		// 已完成
	case OK_TRANSFER_RYPE_REMITTED:
		return client.TransferStatus_TRANSFERSTATUS_COMPLETE
	case OK_TRANSFER_TYPE_MANUALREVIEW:
		return client.TransferStatus_TRANSFERSTATUS_PROCESSING
	case OK_TRANSFER_TYPE_IDENTITYAUTHENTICATION:
		return client.TransferStatus_TRANSFERSTATUS_WAITCERTIFICATE
	default:
		return client.TransferStatus_TRANSFERSTATUS_INVALID
	}
}

func GetOrderTypeToExchange(ot interface{}) OrderType {
	switch ot {
	case order.OrderType_MARKET:
		return ORDER_TYPE_MARKET
	case order.OrderType_LIMIT:
		return ORDER_TYPE_LIMIT
	default:
		return ""
	}

}
