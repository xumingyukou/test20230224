package binance

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/spot_api"
	"clients/logger"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	GlobalInstance *ClientBinanceGlobal
)

type ClientBinanceGlobal struct {
	base.APIConf
	HttpClient *http.Client
	Symbols    map[string]interface{}
	Tokens     map[string]bool
	QuoteCoins map[string]bool
	VisitCount int
	VisitTotal int
	VisitTimer *time.Timer
	lock       sync.Mutex
}

func NewClientBinanceGlobal(url string, timeOffset int64) *ClientBinanceGlobal {
	c := &ClientBinanceGlobal{}
	c.ReadTimeout = timeOffset
	c.HttpClient = &http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Duration(timeOffset) * time.Second,
	}
	c.Symbols = make(map[string]interface{})
	c.Tokens = make(map[string]bool)
	c.QuoteCoins = make(map[string]bool)
	c.VisitTimer = time.NewTimer(time.Minute)

	if err := c.initTokenInfo(); err != nil {
		logger.Logger.Error("get token info err:", err)
	}
	go c.VisitHandle()
	return c
}

func (c *ClientBinanceGlobal) WeightCost() bool {
	if c.VisitCount > 1100 {
		return true
	} else {
		return false
	}
}

func (c *ClientBinanceGlobal) DoRequest(httpMethod, uri, reqBody string, response interface{}) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.WeightCost() {
		return errors.New(fmt.Sprint("api visit in high frequency", c.VisitCount))
	}
	url := spot_api.SPOT_API_BASE_URL + "/" + uri
	resp, header, err := conn.NewHttpRequestWithHeader(c.HttpClient, httpMethod, url, reqBody, map[string]string{})
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

func (c *ClientBinanceGlobal) GetSymbols() *map[string]interface{} {
	return &c.Symbols
}

func (c *ClientBinanceGlobal) GetTokens() *map[string]bool {
	return &c.Tokens
}

func (c *ClientBinanceGlobal) GetFullDepth(symbol string, limit int) (depth interface{}, err error) {
	var (
		url      = fmt.Sprintf("/api/v3/depth?symbol=%s&limit=%d", symbol, limit)
		resDepth spot_api.RespDepth
	)
	err = c.DoRequest("GET", url, "", &resDepth)
	if err != nil {
		return
	}
	depth = resDepth
	return
}

func (c *ClientBinanceGlobal) GetQuoteCoins() *map[string]bool {
	return &c.QuoteCoins
}

func (c *ClientBinanceGlobal) initTokenInfo() (err error) {
	var (
		url          = "/api/v3/exchangeInfo"
		exchangeInfo spot_api.RespExchangeInfo
	)
	err = c.DoRequest("GET", url, "", &exchangeInfo)
	if err != nil {
		logger.Logger.Error("get symbol list err:", err)
		return
	}
	for _, symbol := range exchangeInfo.Symbols {
		symbolStr := fmt.Sprint(symbol.BaseAsset, "/", symbol.QuoteAsset)
		if _, ok := c.Symbols[symbolStr]; !ok {
			c.Symbols[symbolStr] = symbol
		}
		if _, ok := c.Tokens[symbol.BaseAsset]; !ok {
			c.Tokens[symbol.BaseAsset] = true
		}
		if _, ok := c.Tokens[symbol.QuoteAsset]; !ok {
			c.Tokens[symbol.QuoteAsset] = true
		}
		if _, ok := c.QuoteCoins[symbol.QuoteAsset]; !ok {
			c.QuoteCoins[symbol.QuoteAsset] = true
		}
	}
	return
}

func (c *ClientBinanceGlobal) VisitHandle() {
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
