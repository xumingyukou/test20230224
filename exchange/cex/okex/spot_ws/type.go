package spot_ws

import (
	"encoding/json"
	"strconv"
)

type StringInt int64

// UnmarshalJSON create a custom unmarshal for the StringInt
/// this helps us check the type of our value before unmarshalling it

func (st *StringInt) UnmarshalJSON(b []byte) error {
	//convert the bytes into an interface
	//this will help us check the type of our value
	//if it is a string that can be converted into an int we convert it
	///otherwise we return an error
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		///here convert the string into
		///an integer
		i, err := strconv.Atoi(v)
		if err != nil {
			///the string might not be of integer type
			///so return an error
			return err

		}
		*st = StringInt(i)

	}
	return nil
}

type RespError struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
}

type req struct {
	Op   string      `json:"op"`
	Args []*argDates `json:"args"`
}
type argDates struct {
	Channel  DeepGear `json:"channel"`
	InstId   string   `json:"instId,omitempty"`
	InstType string   `json:"instType"`
}

type login struct {
	Op   string    `json:"op"`
	Args []account `json:"args"`
}

type account struct {
	ApiKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

//func (r *RespError) Error() string {
//	if r.Code == "0" {
//		return ""
//	} else {
//		v, _ := json.Marshal(r)
//		return string(v)
//	}
//}

type RespLimitDepthStreamp struct {
	Resp_Info
	Arg struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
	} `json:"arg"`
	Action string `json:"action"`
	Data   []struct {
		Asks     [][]string `json:"asks"`
		Bids     [][]string `json:"bids"`
		Ts       StringInt  `json:"ts, string, omitempty"`
		Checksum int        `json:"checksum"`
	} `json:"data"`
}

//type RespUserAccount struct {
//	Resp_Info
//	Arg struct {
//		Channel string `json:"channel"`
//		Uid     string `json:"uid"`
//	} `json:"arg"`
//	Data []struct {
//		PTime     StringInt `json:"pTime"`
//		EventType string    `json:"eventType"`
//		BalData   []struct {
//			Ccy     string `json:"ccy"`
//			CashBal string `json:"cashBal"`
//			UTime   string `json:"uTime"`
//		} `json:"balData"`
//		PosData []struct {
//			PosId    string `json:"posId"`
//			TradeId  string `json:"tradeId"`
//			InstId   string `json:"instId"`
//			InstType string `json:"instType"`
//			MgnMode  string `json:"mgnMode"`
//			PosSide  string `json:"posSide"`
//			Pos      string `json:"pos"`
//			Ccy      string `json:"ccy"`
//			PosCcy   string `json:"posCcy"`
//			AvgPx    string `json:"avgPx"`
//			UTime    string `json:"uTime"`
//		} `json:"posData"`
//	} `json:"data"`
//}

type RespPosition struct {
	Data []struct {
	} `json:"data"`
}

type RespUserAccount struct {
	Resp_Info
	Arg struct {
		Channel  string `json:"channel"`
		Ccy      string `json:"ccy"`
		InstType string `json:"instType"`
	} `json:"arg"`
	Data []struct {
		UTime       StringInt `json:"uTime"`
		TotalEq     string    `json:"totalEq"`
		IsoEq       string    `json:"isoEq"`
		AdjEq       string    `json:"adjEq"`
		OrdFroz     string    `json:"ordFroz"`
		Imr         string    `json:"imr"`
		Mmr         string    `json:"mmr"`
		NotionalUsd string    `json:"notionalUsd"`
		MgnRatio    string    `json:"mgnRatio"`
		Details     []struct {
			AvailBal      string    `json:"availBal"`
			AvailEq       string    `json:"availEq"`
			Ccy           string    `json:"ccy"`
			CashBal       string    `json:"cashBal"`
			UTime         StringInt `json:"uTime"`
			DisEq         string    `json:"disEq"`
			Eq            string    `json:"eq"`
			EqUsd         string    `json:"eqUsd"`
			FrozenBal     string    `json:"frozenBal"`
			Interest      string    `json:"interest"`
			IsoEq         string    `json:"isoEq"`
			Liab          string    `json:"liab"`
			MaxLoan       string    `json:"maxLoan"`
			MgnRatio      string    `json:"mgnRatio"`
			NotionalLever string    `json:"notionalLever"`
			OrdFrozen     string    `json:"ordFrozen"`
			Upl           string    `json:"upl"`
			UplLiab       string    `json:"uplLiab"`
			CrossLiab     string    `json:"crossLiab"`
			IsoLiab       string    `json:"isoLiab"`
			CoinUsdPrice  string    `json:"coinUsdPrice"`
			StgyEq        string    `json:"stgyEq"`
			IsoUpl        string    `json:"isoUpl"`
		} `json:"details"`
		// position数据
		Adl      string    `json:"adl"`
		AvailPos string    `json:"availPos"`
		AvgPx    string    `json:"avgPx"`
		CTime    StringInt `json:"cTime"`
		Ccy      string    `json:"ccy"`
		DeltaBS  string    `json:"deltaBS"`
		DeltaPA  string    `json:"deltaPA"`
		GammaBS  string    `json:"gammaBS"`
		GammaPA  string    `json:"gammaPA"`
		InstId   string    `json:"instId"`
		InstType string    `json:"instType"`
		Interest string    `json:"interest"`
		Last     string    `json:"last"`
		Lever    string    `json:"lever"`
		Liab     string    `json:"liab"`
		LiabCcy  string    `json:"liabCcy"`
		LiqPx    string    `json:"liqPx"`
		MarkPx   string    `json:"markPx"`
		Margin   string    `json:"margin"`
		MgnMode  string    `json:"mgnMode"`
		OptVal   string    `json:"optVal"`
		PTime    string    `json:"pTime"`
		Pos      string    `json:"pos"`
		PosCcy   string    `json:"posCcy"`
		PosId    string    `json:"posId"`
		PosSide  string    `json:"posSide"`
		ThetaBS  string    `json:"thetaBS"`
		ThetaPA  string    `json:"thetaPA"`
		TradeId  string    `json:"tradeId"`
		Upl      string    `json:"upl"`
		UplRatio string    `json:"uplRatio"`
		VegaBS   string    `json:"vegaBS"`
		VegaPA   string    `json:"vegaPA"`
	} `json:"data"`
}

type RespLogin struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
}

type RespTradeStream struct {
	Arg struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
	} `json:"arg"`
	Data []Data `json:"data"`
}
type Data struct {
	InstId  string    `json:"instId"`
	TradeId string    `json:"tradeId"`
	Px      string    `json:"px"`
	Sz      string    `json:"sz"`
	Side    string    `json:"side"`
	Ts      StringInt `json:"ts"`
}

type Resp_Info struct {
	RespError
	Event string `json:"event"`
	Arg   struct {
		Channel string `json:"channel"`
	} `json:"arg"`
}

type RespUserOrder struct {
	Resp_Info
	Arg struct {
		Channel  string `json:"channel"`
		InstType string `json:"instType"`
		InstId   string `json:"instId"`
		Uid      string `json:"uid"`
	} `json:"arg"`
	Data []struct {
		AccFillSz       string    `json:"accFillSz"`
		AmendResult     string    `json:"amendResult"`
		AvgPx           string    `json:"avgPx"`
		CTime           string    `json:"cTime"`
		Category        string    `json:"category"`
		Ccy             string    `json:"ccy"`
		ClOrdId         string    `json:"clOrdId"`
		Code            string    `json:"code"`
		ExecType        string    `json:"execType"`
		Fee             string    `json:"fee"`
		FeeCcy          string    `json:"feeCcy"`
		FillFee         string    `json:"fillFee"`
		FillFeeCcy      string    `json:"fillFeeCcy"`
		FillNotionalUsd string    `json:"fillNotionalUsd"`
		FillPx          string    `json:"fillPx"`
		FillSz          string    `json:"fillSz"`
		FillTime        string    `json:"fillTime"`
		InstId          string    `json:"instId"`
		InstType        string    `json:"instType"`
		Lever           string    `json:"lever"`
		Msg             string    `json:"msg"`
		NotionalUsd     string    `json:"notionalUsd"`
		OrdId           string    `json:"ordId"`
		OrdType         string    `json:"ordType"`
		Pnl             string    `json:"pnl"`
		PosSide         string    `json:"posSide"`
		Px              string    `json:"px"`
		Rebate          string    `json:"rebate"`
		RebateCcy       string    `json:"rebateCcy"`
		ReduceOnly      string    `json:"reduceOnly"`
		ReqId           string    `json:"reqId"`
		Side            string    `json:"side"`
		SlOrdPx         string    `json:"slOrdPx"`
		SlTriggerPx     string    `json:"slTriggerPx"`
		SlTriggerPxType string    `json:"slTriggerPxType"`
		Source          string    `json:"source"`
		State           string    `json:"state"`
		Sz              string    `json:"sz"`
		Tag             string    `json:"tag"`
		TdMode          string    `json:"tdMode"`
		TgtCcy          string    `json:"tgtCcy"`
		TpOrdPx         string    `json:"tpOrdPx"`
		TpTriggerPx     string    `json:"tpTriggerPx"`
		TpTriggerPxType string    `json:"tpTriggerPxType"`
		TradeId         string    `json:"tradeId"`
		UTime           StringInt `json:"uTime"`
	} `json:"data"`
}

type RespOrderS struct {
	Id string `json:"id"`
	Op string `json:"op"`
	RespError
}

type RespOrderItem struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	Tag     string `json:"tag"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

type RespOrder struct {
	RespOrderS
	RespLogin
	Data []*RespOrderItem `json:"data"`
}

type RespCancelOrderItem struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

type RespCancelOrder struct {
	RespOrderS
	RespLogin
	Data []*RespCancelOrderItem `json:"data"`
}

type ReqOrder struct {
	Id   string              `json:"id"`
	Op   string              `json:"op"`
	Args []map[string]string `json:"args"`
}

type RespFundingRate struct {
	Code string `json:"code"`
	Data []struct {
		FundingRate     string    `json:"fundingRate"`
		FundingTime     StringInt `json:"fundingTime"`
		InstId          string    `json:"instId"`
		InstType        string    `json:"instType"`
		NextFundingRate string    `json:"nextFundingRate"`
		NextFundingTime string    `json:"nextFundingTime"`
	} `json:"data"`
	Msg string `json:"msg"`
}
