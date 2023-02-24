package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/exchange/cex/binance/c_api"
	"clients/exchange/cex/binance/spot_api"
	"clients/exchange/cex/binance/u_api"
	"clients/logger"
	"context"
	"errors"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"net/http"
	"sort"
	"sync"
	"time"
)

var (
	cg          *GlobalDepthWight
	cglobalOnce sync.Once
)

type BinanceCBaseWebsocket struct {
	BinanceSpotWebsocket
}

func NewBinanceCBaseWebsocket(conf base.WsConf) *BinanceCBaseWebsocket {
	if conf.EndPoint == "" {
		conf.EndPoint = spot_api.WS_CAPI_BASE_URL
	}
	d := &BinanceCBaseWebsocket{}
	d.WsConf = conf
	apiConf := base.APIConf{
		ProxyUrl:    conf.ProxyUrl,
		ReadTimeout: conf.ReadTimeout,
		AccessKey:   conf.AccessKey,
		SecretKey:   conf.SecretKey,
	}
	d.Init(spot_api.NewApiClient(apiConf), u_api.NewUApiClient(apiConf), c_api.NewCApiClient(apiConf), NewWebSocketCBaseHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	d.GlobalDW = cg
	cglobalOnce.Do(func() {
		if cg == nil {
			cg = new(GlobalDepthWight)
			d.GlobalDW = cg
		}
		go func() {
			timer := time.NewTimer(time.Duration(d.CApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
			for {
				select {
				case <-timer.C:
					d.GlobalDW.GlobalLock.Lock()
					d.GlobalDW.GlobalWeight = 0
					d.GlobalDW.GlobalLock.Unlock()
					timer.Reset(time.Duration(d.CApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
				}
			}
		}()
	})
	return d
}

func NewBinanceCBaseWebsocket2(conf base.WsConf, cli *http.Client) *BinanceCBaseWebsocket {
	d := &BinanceCBaseWebsocket{}
	d.WsConf = conf
	apiConf := base.APIConf{
		ProxyUrl:    conf.ProxyUrl,
		ReadTimeout: conf.ReadTimeout,
		AccessKey:   conf.AccessKey,
		SecretKey:   conf.SecretKey,
	}
	d.Init(spot_api.NewApiClient2(apiConf, cli), u_api.NewUApiClient2(apiConf, cli), c_api.NewCApiClient2(apiConf, cli), NewWebSocketCBaseHandle(d.ChanCap), spot_api.NewSpotWsUrl(), d.GetFullDepth)
	d.GlobalDW = cg
	cglobalOnce.Do(func() {
		if cg == nil {
			cg = new(GlobalDepthWight)
			d.GlobalDW = cg
		}
		go func() {
			timer := time.NewTimer(time.Duration(d.CApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
			for {
				select {
				case <-timer.C:
					d.GlobalDW.GlobalLock.Lock()
					d.GlobalDW.GlobalWeight = 0
					d.GlobalDW.GlobalLock.Unlock()
					timer.Reset(time.Duration(d.CApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].IntervalSec) * time.Second)
				}
			}
		}()
	})
	return d
}

func (b *BinanceCBaseWebsocket) GetFullDepth(symbol string) (*base.OrderBook, error) {
	b.GlobalDW.GlobalLock.Lock()
	defer b.GlobalDW.GlobalLock.Unlock()
	if b.GlobalDW.GlobalWeight > 1000 {
		return nil, errors.New(fmt.Sprintf("api visit %s in high frequency %d", symbol, b.CApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].Value))
	}
	resp, err := b.CApiClient.GetDepth(spot_api.GetSymbol(symbol), 1000)
	b.GlobalDW.GlobalWeight = b.CApiClient.WeightInfo[client.WeightType_REQUEST_WEIGHT].Value
	if err != nil {
		return nil, err
	}
	symbol_, market, type_ := u_api.GetContractType(symbol)
	res := &base.OrderBook{
		Exchange: common.Exchange_BINANCE,
		Symbol:   symbol_,
		Market:   market,
		Type:     type_,
		UpdateId: resp.LastUpdateId,
	}
	ParseOrder(resp.Asks, &res.Asks)
	ParseOrder(resp.Bids, &res.Bids)
	// 币本位需要转换amount
	for bidIdx, _ := range res.Bids {
		res.Bids[bidIdx].Amount = getCbaseQty(res.Symbol, res.Bids[bidIdx].Amount, res.Bids[bidIdx].Price)
	}
	for askIdx, _ := range res.Asks {
		res.Asks[askIdx].Amount = getCbaseQty(res.Symbol, res.Asks[askIdx].Amount, res.Asks[askIdx].Price)
	}
	sort.Stable(res.Asks)
	sort.Stable(sort.Reverse(res.Bids))
	return res, nil
}

func (b *BinanceCBaseWebsocket) Account(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsAccountRsp, error) {
	key, err := b.AccountExecute(reqList[0], LISTEN, "")
	if err != nil {
		logger.Logger.Error("get listen key err:", err)
		return nil, err
	}
	var (
		url = b.GetUrl(spot_api.SINGLE_API_URL, key)
	)
	b.reqId++
	err = b.EstablishUserConn(b.reqId, reqList[0], key, url, b.handler.AccountHandle, ctx)
	return b.handler.GetChan("account").(chan *client.WsAccountRsp), err
}

func (b *BinanceCBaseWebsocket) Order(ctx context.Context, reqList ...*client.WsAccountReq) (<-chan *client.WsOrderRsp, error) {
	key, err := b.AccountExecute(reqList[0], LISTEN, "")
	if err != nil {
		logger.Logger.Error("get listen key err:", err)
		return nil, err
	}
	var (
		url = b.GetUrl(spot_api.SINGLE_API_URL, key)
	)
	b.reqId++
	err = b.EstablishUserConn(b.reqId, reqList[0], key, url, b.handler.OrderHandle, ctx)
	return b.handler.GetChan("order").(chan *client.WsOrderRsp), err
}

func (b *BinanceCBaseWebsocket) FundingRateGroup(ctx context.Context, chMap map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error {
	var symbols []*client.SymbolInfo
	for symbol := range chMap {
		symbols = append(symbols, symbol)
	}
	b.reqId++
	b.handler.SetFundingRateGroupChannel(chMap)
	err := b.EstablishConn(b.reqId, symbols, 0, 0, b.WsReqUrl.MARKPRICE_URL, spot_api.STREAM_API_URL, b.handler.FundingRateGroupHandle, ctx)
	return err
}
