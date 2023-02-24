package main

import (
	"clients/conn"
	"clients/exchange/cex/base"
	"fmt"
	"github.com/goccy/go-json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	depthCache *base.OrderBook
)

type req struct {
	Event        string       `json:"event"`
	Pair         []string     `json:"pair"`
	Subscription subscription `json:"subscription"`
}
type subscription struct {
	Name  string `json:"name"`
	Depth int    `json:"depth,omitempty"`
}

const WsSubscrible = "subscribe"

func websockt() {
	var (
		readTimeout = 100 * time.Second
		url         = "wss://ws.kraken.com"
		proxyUrl    = "http://127.0.0.1:7890"
		symbols_    = []string{"ZRX/BTC"}
		reqx        = req{
			Event: WsSubscrible,
			Pair:  symbols_,
			Subscription: subscription{
				Name:  "book",
				Depth: 10,
			},
		}
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler).ProxyUrl(proxyUrl)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		fmt.Println("create client err")
		return
	}
	err := Subscribe(wsClient, reqx)
	if err != nil {
		fmt.Println(err)
		return
	}
	select {}
}

type RespSymbols struct {
	Error  []interface{} `json:"error"`
	Result map[string]struct {
		Altname           string          `json:"altname"`
		Wsname            string          `json:"wsname"`
		AclassBase        string          `json:"aclass_base"`
		Base              string          `json:"base"`
		AclassQuote       string          `json:"aclass_quote"`
		Quote             string          `json:"quote"`
		Lot               string          `json:"lot"`
		PairDecimals      int             `json:"pair_decimals"`
		LotDecimals       int             `json:"lot_decimals"`
		LotMultiplier     int             `json:"lot_multiplier"`
		LeverageBuy       []int           `json:"leverage_buy"`
		LeverageSell      []int           `json:"leverage_sell"`
		Fees              [][]float64     `json:"fees"`
		FeesMaker         [][]interface{} `json:"fees_maker"`
		FeeVolumeCurrency string          `json:"fee_volume_currency"`
		MarginCall        int             `json:"margin_call"`
		MarginStop        int             `json:"margin_stop"`
		Ordermin          string          `json:"ordermin"`
	} `json:"result"`
}

func main() {
	var (
		uri         = "https://api.kraken.com/0/public/AssetPairs"
		proxyUrl, _ = url.Parse("http://127.0.0.1:7890")
		client      = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			Timeout: 10 * time.Second,
		}
	)

	resp, err := conn.Request(client, uri, "GET", &http.Header{}, nil)
	if err != nil {
		fmt.Println("get err:", err)
		return
	}
	var res RespSymbols
	err = json.Unmarshal(resp, &res)
	if err != nil {
		fmt.Println("parse err:", err)
		return
	}
	fmt.Println(res)
}

func Subscribe(client *conn.WsConn, reqx interface{}) error {
	err := client.Subscribe(reqx)
	if err != nil {
		fmt.Println("subscribe err:", err)
	}
	return err
}

func TransferFullDepth(data map[string][][]string) *base.OrderBook {
	return nil
}
func TransferDeltaDepth(data map[string]interface{}) *base.DeltaDepthUpdate {
	return nil
}

func handler(data []byte) error {
	var (
		content  []interface{}
		err      error
		symbol   string
		bookData map[string]interface{}
		ok       bool
		checkNum string
	)
	fmt.Sprint(symbol, checkNum)
	err = json.Unmarshal(data, &content)
	if err != nil {
		fmt.Println(err, string(data))
		return err
	}
	symbol, ok = content[len(content)-1].(string)
	if !ok {
		fmt.Println("get symbol name err")
		return err
	}
	if strings.Contains(string(data), "as") && strings.Contains(string(data), "bs") {
		var fullData map[string][][]string
		fullData, ok = content[1].(map[string][][]string)
		if !ok {
			fmt.Println("get book data err")
			return err
		}
		depthCache = TransferFullDepth(fullData)
	} else {
		bookData, ok = content[1].(map[string]interface{})
		if !ok {
			fmt.Println("get book data err")
			return err
		}
		update1 := TransferDeltaDepth(bookData)
		base.UpdateBidsAndAsks(update1, depthCache, 20, nil)
		bookData, ok = content[2].(map[string]interface{})
		if ok {
			update1 = TransferDeltaDepth(bookData)
			base.UpdateBidsAndAsks(update1, depthCache, 20, nil)
		}

	}

	return nil
}
