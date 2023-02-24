package base

import (
	"clients/logger"
	"clients/transform"
	"fmt"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
)

type DepthItemSlice []*depth.DepthLevel

func (d DepthItemSlice) Len() int           { return len(d) }
func (d DepthItemSlice) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d DepthItemSlice) Less(i, j int) bool { return d[i].Price < d[j].Price }
func (d DepthItemSlice) Print() {
	msg := ""
	for _, item := range d {
		msg += fmt.Sprintf("(%.10f ,%.10f)", item.Price, item.Amount)
	}
	fmt.Println(msg)
	logger.Logger.Info(msg)
}
func (d *DepthItemSlice) Reverse() {
	for i := 0; i < d.Len()/2; i++ {
		(*d)[i], (*d)[d.Len()-i-1] = (*d)[d.Len()-i-1], (*d)[i]
	}
}
func (d DepthItemSlice) Copy() *DepthItemSlice {
	res := &DepthItemSlice{}
	for _, i := range d {
		*res = append(*res, &depth.DepthLevel{
			Price:  i.Price,
			Amount: i.Amount,
		})
	}
	return res
}

type DeltaDepthUpdate struct {
	UpdateStartId int64
	UpdateEndId   int64
	UpdateNextId  int64
	Market        common.Market
	Type          common.SymbolType
	Symbol        string
	TimeExchange  int64
	TimeReceive   int64
	IsFullDepth   bool

	Bids DepthItemSlice
	Asks DepthItemSlice
}

func (d *DeltaDepthUpdate) Transfer2Depth() (depthRes *client.WsDepthRsp) {
	depthRes = &client.WsDepthRsp{
		UpdateIdEnd:   d.UpdateEndId,
		UpdateIdStart: d.UpdateStartId,
		Symbol:        d.Symbol,
		ExchangeTime:  d.TimeExchange,
		ReceiveTime:   d.TimeReceive,
		Bids:          *d.Bids.Copy(),
		Asks:          *d.Asks.Copy(),
	}
	return
}

type OrderBook struct {
	Exchange     common.Exchange
	ExchangeAddr string
	Market       common.Market
	Type         common.SymbolType
	Symbol       string
	TimeExchange uint64
	TimeReceive  uint64
	Bids         DepthItemSlice
	Asks         DepthItemSlice
	UpdateId     int64
}

func (o *OrderBook) Limit(limit int) {
	bs := transform.Min(limit, len(o.Bids))
	as := transform.Min(limit, len(o.Asks))
	o.Bids = o.Bids[:bs]
	o.Asks = o.Asks[:as]
}

func (o *OrderBook) Copy() (oNew *OrderBook) {
	oNew = &OrderBook{
		Bids: make(DepthItemSlice, len(o.Bids), len(o.Bids)), // 长度和容量设置，开辟多余的空间数据会有很多none值
		Asks: make(DepthItemSlice, len(o.Asks), len(o.Asks)),
	}
	oNew.Exchange = o.Exchange
	//oNew.ExchangeAddr = o.ExchangeAddr
	oNew.Market = o.Market
	oNew.Type = o.Type
	oNew.Symbol = o.Symbol
	oNew.TimeExchange = o.TimeExchange
	oNew.TimeReceive = o.TimeReceive
	oNew.UpdateId = o.UpdateId
	for i, bid := range o.Bids {
		oNew.Bids[i] = &depth.DepthLevel{
			Price:  bid.Price,
			Amount: bid.Amount,
		}
	}
	for i, asks := range o.Asks {
		oNew.Asks[i] = &depth.DepthLevel{
			Price:  asks.Price,
			Amount: asks.Amount,
		}
	}
	return
}

func (o *OrderBook) Equal(oNew *OrderBook, depthLevel int) (isEqual bool) {
	if oNew.UpdateId == o.UpdateId {
		for i := 0; i < transform.Min(len(oNew.Asks), len(o.Asks), depthLevel); i++ {
			if oNew.Asks[i].Price != o.Asks[i].Price || oNew.Asks[i].Amount != o.Asks[i].Amount {
				return
			}
		}
		for i := 0; i < transform.Min(len(oNew.Bids), len(o.Bids), depthLevel); i++ {
			if oNew.Bids[i].Price != o.Bids[i].Price || oNew.Bids[i].Amount != o.Bids[i].Amount {
				return
			}
		}
		isEqual = true
	}
	return
}

func (o *OrderBook) Transfer2Depth(depthLimit int, dep *depth.Depth) {
	askCount, bidCount := transform.Min(len(o.Asks), depthLimit), transform.Min(len(o.Bids), depthLimit)
	dep.Hdr = &common.MsgHeader{
		Type: uint32(common.MsgType_DPETH),
	}
	dep.Exchange = o.Exchange
	dep.ExchangeAddr = o.ExchangeAddr
	dep.Market = o.Market
	dep.Type = o.Type
	dep.Symbol = o.Symbol
	dep.TimeExchange = o.TimeExchange
	dep.TimeReceive = o.TimeReceive
	dep.TimeOperate = uint64(time.Now().UnixMicro())
	dep.Asks = make(DepthItemSlice, 0, askCount) // 长度和容量设置，开辟多余的空间数据会有很多none值
	dep.Bids = make(DepthItemSlice, 0, bidCount)
	for i := 0; i < bidCount; i++ {
		dep.Bids = append(dep.Bids, &depth.DepthLevel{
			Price:  o.Bids[i].Price,
			Amount: o.Bids[i].Amount,
		})
	}
	for i := 0; i < askCount; i++ {
		dep.Asks = append(dep.Asks, &depth.DepthLevel{
			Price:  o.Asks[i].Price,
			Amount: o.Asks[i].Amount,
		})
	}
	return
}

func DepthDeepCopyByCustom(src *depth.Depth) *depth.Depth {
	dst := new(depth.Depth)
	hdr := src.Hdr
	dst.Hdr = hdr
	dst.Exchange = src.Exchange
	dst.ExchangeAddr = src.ExchangeAddr
	dst.Market = src.Market
	dst.Type = src.Type
	dst.Symbol = src.Symbol
	dst.TimeExchange = src.TimeExchange
	dst.TimeReceive = src.TimeReceive
	dst.TimeOperate = src.TimeOperate
	bids := make([]*depth.DepthLevel, len(src.Bids))
	for i, bid := range src.Bids {
		bids[i] = &depth.DepthLevel{
			Price:  bid.Price,
			Amount: bid.Amount,
		}
	}
	dst.Bids = bids
	asks := make([]*depth.DepthLevel, len(src.Asks))
	for i, ask := range src.Asks {
		asks[i] = &depth.DepthLevel{
			Price:  ask.Price,
			Amount: ask.Amount,
		}
	}
	dst.Asks = asks
	return dst
}

func MakeFirstDepthHdr() *common.MsgHeader {
	return &common.MsgHeader{
		Type: uint32(common.MsgType_FIRST_SUBSCRIBE),
	}
}
