package okex

import (
	"clients/conn"
	"clients/crypto"
	"clients/logger"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	GLOBAL_API_BASE_URL = "https://www.okx.com/"
)

const (
	v5RestBaseUrl = "https://www.okex.com"
	v5WsBaseUrl   = "wss://ws.okex.com:8443/ws/v5"

	CONTENT_TYPE          = "Content-Type"
	ACCEPT                = "Accept"
	APPLICATION_JSON_UTF8 = "application/json; charset=UTF-8"
	APPLICATION_JSON      = "application/json"
	OK_ACCESS_KEY         = "OK-ACCESS-KEY"
	OK_ACCESS_SIGN        = "OK-ACCESS-SIGN"
	OK_ACCESS_TIMESTAMP   = "OK-ACCESS-TIMESTAMP"
	OK_ACCESS_PASSPHRASE  = "OK-ACCESS-PASSPHRASE"
	OK_SIMULATED_TRADING  = "x-simulated-trading"
)

var (
//Instance = NewClientOkex(GLOBAL_API_BASE_URL, 2)
)

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
	Code string             `json:"code"`
	Msg  string             `json:"msg"`
	Data []*InstrumentsData `json:"data"`
}

type BalanceDetail struct {
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

type BalanceItem struct {
	AdjEq       string           `json:"adjEq"`
	Details     []*BalanceDetail `json:"details"`
	Imr         string           `json:"imr"`
	IsoEq       string           `json:"isoEq"`
	MgnRatio    string           `json:"mgnRatio"`
	Mmr         string           `json:"mmr"`
	NotionalUsd string           `json:"notionalUsd"`
	OrdFroz     string           `json:"ordFroz"`
	TotalEq     string           `json:"totalEq"`
	UTime       string           `json:"uTime"`
}

type RespAccountBalance struct {
	Code string         `json:"code"`
	Data []*BalanceItem `json:"data"`
	Msg  string         `json:"msg"`
}

type TradeFee struct {
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
}

type RespTradeFee struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data []*TradeFee `json:"data"`
}

type ClientOkexOld struct {
	AccessKey, SecretKey string
	baseUrl              string
	httpClient           *http.Client
	timeOffset           int64 //nanosecond
	Symbols              map[string]interface{}
	Tokens               map[string]bool
	QuoteCoins           map[string]bool
	VisitCount           int
	VisitTotal           int
	VisitTimer           *time.Timer
	lock                 sync.Mutex
}

func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}

func (c *ClientOkexOld) doParamSign(httpMethod, uri, requestBody string) (string, string) {
	timestamp := IsoTime()
	preText := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), uri, requestBody)
	//log.Println("preHash", preText)
	sign, _ := crypto.GetParamHmacSHA256Base64Sign(c.SecretKey, preText)
	return sign, timestamp
}

func NewClientOkexOld(url_ string, timeOffset int64) *ClientOkexOld {
	c := new(ClientOkexOld)
	if strings.HasSuffix(url_, "/") {
		url_ = url_[:len(url_)-1]
	}
	proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	transport := http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	c.baseUrl = url_
	c.timeOffset = timeOffset
	c.httpClient = &http.Client{
		Transport: &transport,
		Timeout:   time.Duration(timeOffset) * time.Second,
	}
	c.Symbols = make(map[string]interface{})
	c.Tokens = make(map[string]bool)
	c.QuoteCoins = make(map[string]bool)
	c.VisitTimer = time.NewTimer(time.Minute)

	if err := c.initTokenInfo(); err != nil {
		fmt.Println("get token info err:", err)
		logger.Logger.Error("get token info err:", err)
	}
	go c.VisitHandle()
	return c
}

func (c *ClientOkexOld) WeightCost() bool {
	if c.VisitCount > 15 {
		return true
	} else {
		return false
	}
}

func (c *ClientOkexOld) DoRequest(httpMethod, uri, reqBody string, response interface{}) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	url := c.baseUrl + uri
	sign, timestamp := c.doParamSign(httpMethod, uri, reqBody)
	//logger.Log.Debug("timestamp=", timestamp, ", sign=", sign)
	resp, header, err := conn.NewHttpRequestWithHeader(c.httpClient, httpMethod, url, reqBody, map[string]string{
		CONTENT_TYPE: APPLICATION_JSON_UTF8,
		ACCEPT:       APPLICATION_JSON,
		//COOKIE:               LOCALE + "en_US",
		OK_ACCESS_KEY:        c.AccessKey,
		OK_ACCESS_PASSPHRASE: "",
		OK_ACCESS_SIGN:       sign,
		OK_ACCESS_TIMESTAMP:  fmt.Sprint(timestamp),
		OK_SIMULATED_TRADING: "1",
		//"Content-Type":       "application/x-www-form-urlencoded",
	})

	if c.WeightCost() {
		return errors.New(fmt.Sprint("api visit in high frequency", c.VisitCount))
	}

	//resp, header, err := conn.NewHttpRequestWithHeader(c.httpClient, httpMethod, url, reqBody, map[string]string{})
	weight := header.Get("X-Mbx-Used-Weight-1m")
	c.VisitCount, _ = strconv.Atoi(weight)
	if err != nil {
		logger.Logger.Error("DoRequest err:", httpMethod, uri, reqBody, err)
		return
	} else {
		logger.Logger.Debug(string(resp))
		return json.Unmarshal(resp, &response)
	}
}

func (c *ClientOkexOld) GetSymbols() *map[string]interface{} {
	return &c.Symbols
}

func (c *ClientOkexOld) GetTokens() *map[string]bool {
	return &c.Tokens
}
func (c *ClientOkexOld) GetQuoteCoins() *map[string]bool {
	return &c.QuoteCoins
}

func (c *ClientOkexOld) initTokenInfo() (err error) {
	var (
		url          = "/api/v5/public/instruments?instType=SPOT"
		exchangeInfo RespInstruments
	)
	err = c.DoRequest("GET", url, "", &exchangeInfo)
	if err != nil {
		logger.Logger.Error("get symbol list err:", err)
		return
	}
	for _, instrument := range exchangeInfo.Data {
		symbolStr := fmt.Sprint(instrument.BaseCcy, "/", instrument.QuoteCcy)
		if _, ok := c.Symbols[symbolStr]; !ok {
			c.Symbols[symbolStr] = true
		}
		if _, ok := c.Tokens[instrument.BaseCcy]; !ok {
			c.Tokens[instrument.BaseCcy] = true
		}
		if _, ok := c.Tokens[instrument.QuoteCcy]; !ok {
			c.Tokens[instrument.QuoteCcy] = true
		}
		if _, ok := c.QuoteCoins[instrument.QuoteCcy]; !ok {
			c.QuoteCoins[instrument.QuoteCcy] = true
		}
	}
	return
}

func (c *ClientOkexOld) GetFullDepth(symbol string, depthLevel int) (interface{}, error) {
	return nil, nil
}

func (c *ClientOkexOld) VisitHandle() {
	for {
		select {
		case <-c.VisitTimer.C:
			c.lock.Lock()
			c.VisitCount = 0
			c.lock.Unlock()
			c.VisitTimer.Reset(time.Minute)
		}
	}
}

func (c *ClientOkexOld) AccountBalance() (*RespAccountBalance, error) {
	var (
		url     = "/api/v5/account/balance"
		balance RespAccountBalance
	)
	err := c.DoRequest("GET", url, "", &balance)
	if err != nil {
		logger.Logger.Error("get symbol list err:", err)
		return nil, err
	}
	return &balance, nil
}

func (c *ClientOkexOld) AccountTradeFee() (*RespTradeFee, error) {
	var (
		url     = "/api/v5/account/trade-fee?instType=SPOT"
		balance RespTradeFee
	)
	err := c.DoRequest("GET", url, "", &balance)
	if err != nil {
		logger.Logger.Error("get symbol list err:", err)
		return nil, err
	}
	return &balance, nil
}
