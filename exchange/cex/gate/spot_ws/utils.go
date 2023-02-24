package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"github.com/shopspring/decimal"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
	"log"
	"sort"
	"strconv"
	"strings"
)

func GetSide(side string) order.TradeSide {
	switch side {
	case "buy":
		return order.TradeSide_BUY
	case "sell":
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

func RemovePriceFromBids(bids []Level, price decimal.Decimal) []Level {
	i := sort.Search(len(bids), func(i int) bool { return bids[i].Price.LessThanOrEqual(price) })
	if i < len(bids) && bids[i].Price.Equals(price) {
		// price is present at bids[i]
		return remove(bids, i)
	} else {
		// price is not present in data, but i is the index where it would be inserted.
		// nothing to do here
		return bids
	}
}

func InsertPriceInBids(bids []Level, price decimal.Decimal, volume decimal.Decimal) []Level {
	bids = RemovePriceFromBids(bids, price)
	level := Level{Price: price, Volume: volume}
	i := sort.Search(len(bids), func(i int) bool { return bids[i].Price.LessThan(price) })
	bids = append(bids, Level{})
	copy(bids[i+1:], bids[i:])
	bids[i] = level
	return bids
}

func RemovePriceFromAsks(asks []Level, price decimal.Decimal) []Level {
	i := sort.Search(len(asks), func(i int) bool { return asks[i].Price.GreaterThanOrEqual(price) })
	if i < len(asks) && asks[i].Price.Equals(price) {
		// price is present at bids[i]
		return remove(asks, i)
	} else {
		// price is not present in data, but i is the index where it would be inserted.
		// nothing to do here
		return asks
	}
}

func InsertPriceInAsks(asks []Level, price decimal.Decimal, volume decimal.Decimal) []Level {
	asks = RemovePriceFromAsks(asks, price)
	level := Level{Price: price, Volume: volume}
	i := sort.Search(len(asks), func(i int) bool { return asks[i].Price.GreaterThan(price) })
	asks = append(asks, Level{})
	copy(asks[i+1:], asks[i:])
	asks[i] = level
	return asks
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

func kkd2orderbook(v KKDepth) (res *depth.Depth) {
	res = &depth.Depth{}
	asks, bids := base.DepthItemSlice{}, base.DepthItemSlice{}
	for _, vv := range v.Asks {
		price, _ := vv.Price.Float64()
		amount, _ := vv.Volume.Float64()
		asks = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	for _, vv := range v.Bids {
		price, _ := vv.Price.Float64()
		amount, _ := vv.Volume.Float64()
		bids = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	res.Symbol = v.Symbol
	res.Asks = asks
	res.Bids = bids
	return
}

func kk2clientWs(v KKDepth) (res *client.WsDepthRsp) {
	res = &client.WsDepthRsp{}
	asks, bids := base.DepthItemSlice{}, base.DepthItemSlice{}
	for _, vv := range v.Asks {
		price, _ := vv.Price.Float64()
		amount, _ := vv.Volume.Float64()
		asks = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	for _, vv := range v.Bids {
		price, _ := vv.Price.Float64()
		amount, _ := vv.Volume.Float64()
		bids = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	res.Symbol = v.Symbol
	res.Asks = asks
	res.Bids = bids
	return
}

func GetChecksumInput(bids []Level, asks []Level, aPrecision, pPrecision int32) string {
	var str strings.Builder
	for _, level := range asks[:10] {
		//price := level.Price.StringFixed(pPrecision)
		price := level.Price.StringFixed(pPrecision)
		price = strings.Replace(price, ".", "", 1)
		price = strings.TrimLeft(price, "0")
		str.WriteString(price)

		volume := level.Volume.StringFixed(aPrecision)
		volume = strings.Replace(volume, ".", "", 1)
		volume = strings.TrimLeft(volume, "0")
		str.WriteString(volume)
	}
	for _, level := range bids[:10] {
		price := level.Price.StringFixed(pPrecision)
		price = strings.Replace(price, ".", "", 1)
		price = strings.TrimLeft(price, "0")
		str.WriteString(price)

		volume := level.Volume.StringFixed(aPrecision)
		volume = strings.Replace(volume, ".", "", 1)
		volume = strings.TrimLeft(volume, "0")
		str.WriteString(volume)
	}
	return str.String()
}

func ParseSymbol(symbol string) string {
	return strings.ReplaceAll(symbol, "_", "/")
}

func ParseOrder(orders [][]string, slice *base.DepthItemSlice) {
	for _, order := range orders {
		price, amount, err := transform.ParsePriceAmountFloat(order)
		if err != nil {
			logger.Logger.Errorf("order float parse price error [%s] , response data = %s", err, order)
			continue
		}
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
}
