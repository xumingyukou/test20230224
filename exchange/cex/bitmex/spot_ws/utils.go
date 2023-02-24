package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"github.com/shopspring/decimal"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
	"log"
	"strconv"
	"strings"
)

func GetSide(side string) order.TradeSide {
	switch side {
	case "BUY":
		return order.TradeSide_BUY
	case "SELL":
		return order.TradeSide_SELL
	default:
		return order.TradeSide_InvalidSide
	}
}

// interface转float64
func ParseF(x interface{}) float64 {
	res, err := strconv.ParseFloat(x.(string), 64)
	if err != nil {
		logger.Logger.Error("transform failed", x)
	}
	return res
}

// interface转Int64
func ParseI(x interface{}) int64 {
	res, err := strconv.ParseInt(x.(string), 10, 64)
	if err != nil {
		logger.Logger.Error("transform failed", x)
	}
	return res
}

//func ParseOrder(orders []*depth.DepthLevel, slice *base.DepthItemSlice) {
//	for _, order := range orders {
//		price, amount := order.Price, order.Amount
//		*slice = append(*slice, &depth.DepthLevel{
//			Price:  price,
//			Amount: amount,
//		})
//	}
//}

//func transferDiffDepth(r *client.WsDepthRsp, diff *base.DeltaDepthUpdate) {
//	// 将binance返回的结构，解析为DeltaDepthUpdate，并将bids和ask进行排序
//	diff.TimeExchange = int64(r.ExchangeTime) * 1000
//	diff.TimeReceive = time.Now().UnixMicro()
//	if diff.TimeReceive-diff.TimeExchange > 10000 {
//		//fmt.Println(diff.Symbol, diff.TimeExchange, diff.TimeReceive, diff.TimeReceive-diff.TimeExchange)
//	}
//	diff.UpdateEndId = r.UpdateIdEnd
//	//清空bids，asks
//	diff.Bids = diff.Bids[:0]
//	diff.Asks = diff.Asks[:0]
//	ParseOrder(r.Bids, &diff.Bids)
//	ParseOrder(r.Asks, &diff.Asks)
//	sort.Sort(&diff.Asks)
//	sort.Sort(sort.Reverse(&diff.Bids))
//	return
//}

func transDepth2OrderBook(d *depth.Depth, res *base.OrderBook) {
	res.Asks = d.Asks
	res.Bids = d.Bids
	res.TimeReceive = d.TimeReceive
	res.TimeExchange = d.TimeExchange
}

func ParseFloat(f float64) string {
	res := strconv.FormatFloat(f, 'f', -1, 64)
	return res
}

func getPriceAndVolume(ino interface{}) (decimal.Decimal, decimal.Decimal, int) {
	el := ino.([]interface{})
	price, err := decimal.NewFromString(el[0].(string))
	if err != nil {
		log.Fatal(err)
	}
	volume, err := decimal.NewFromString(el[1].(string))
	if err != nil {
		log.Fatal(err)
	}
	return price, volume, len(el)
}

func ParseSymbol(symbol string) string {
	return strings.ReplaceAll(symbol, "_", "/")
}

func ParseOrder(orders []levelList, slice *base.DepthItemSlice) {
	for _, order := range orders {
		*slice = append(*slice, &depth.DepthLevel{
			Price:  order.Price,
			Amount: order.Size,
		})
	}
}
