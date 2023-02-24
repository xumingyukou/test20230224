package ok_api

type Post_Trde_BatchOrders struct {
	InstId     string `json:"instId"`
	TdMode     string `json:"tdMode"`
	Ccy        string `json:"ccy"`
	Tag        string `json:"tag"`
	PosSide    string `json:"pos_side"`
	ClOrdId    string `json:"clOrdId"`
	Side       string `json:"side"`
	OrdType    string `json:"ordType"`
	Px         string `json:"px"`
	Sz         string `json:"sz"`
	ReduceOnly bool   `json:"reduce_only"`
	TgtCcy     string `json:"tgtCcy"`
	BanAmend   bool   `json:"banAmend""`
}

type Post_Trade_BatchCancelOrders struct {
	InstId  string `json:"instId"`
	OrdId   string `json:"ordId"`
	ClOrdId string `json:"cl_ord_id"`
}

type Post_Trade_CancelAlgos struct {
	AlgoId string `json:"algoId"`
	InstId string `json:"instId"`
}
