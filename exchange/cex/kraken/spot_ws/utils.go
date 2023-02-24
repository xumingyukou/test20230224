package spot_ws

import (
	"clients/exchange/cex/base"
	"clients/logger"
	"clients/transform"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

func GetSide(side string) order.TradeSide {
	switch side {
	case "b":
		return order.TradeSide_BUY
	case "s":
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
	res, err := strconv.ParseFloat(x.(string), 64)
	res = res * 1000000

	if err != nil {
		logger.Logger.Error("transform failed", x)
	}
	return int64(res)
}

func ParseOrder(orders []*depth.DepthLevel, slice *base.DepthItemSlice) {
	for _, order := range orders {
		price, amount := order.Price, order.Amount
		*slice = append(*slice, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
}

func transferDiffDepth(r *client.WsDepthRsp, diff *base.DeltaDepthUpdate) {
	// 将binance返回的结构，解析为DeltaDepthUpdate，并将bids和ask进行排序
	diff.TimeExchange = int64(r.ExchangeTime) * 1000
	diff.TimeReceive = time.Now().UnixMicro()
	if diff.TimeReceive-diff.TimeExchange > 10000 {
		//fmt.Println(diff.Symbol, diff.TimeExchange, diff.TimeReceive, diff.TimeReceive-diff.TimeExchange)
	}
	diff.UpdateEndId = r.UpdateIdEnd
	//清空bids，asks
	diff.Bids = diff.Bids[:0]
	diff.Asks = diff.Asks[:0]
	ParseOrder(r.Bids, &diff.Bids)
	ParseOrder(r.Asks, &diff.Asks)
	sort.Stable(&diff.Asks)
	sort.Stable(sort.Reverse(&diff.Bids))
	return
}

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

func InsertPriceInBids(bids []Level, price decimal.Decimal, volume decimal.Decimal, priceF float64, volumF float64) []Level {
	bids = RemovePriceFromBids(bids, price)
	level := Level{Price: price, Volume: volume, PriceF: priceF, VolumeF: volumF}
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

func InsertPriceInAsks(asks []Level, price decimal.Decimal, volume decimal.Decimal, priceF float64, volumF float64) []Level {
	asks = RemovePriceFromAsks(asks, price)
	level := Level{Price: price, Volume: volume, PriceF: priceF, VolumeF: volumF}
	i := sort.Search(len(asks), func(i int) bool { return asks[i].Price.GreaterThan(price) })
	asks = append(asks, Level{})
	copy(asks[i+1:], asks[i:])
	asks[i] = level
	return asks
}

func getPriceAndVolume(ino interface{}) (decimal.Decimal, decimal.Decimal, int64, int, float64, float64) {
	el := ino.([]interface{})
	price, err := decimal.NewFromString(el[0].(string))
	if err != nil {
		log.Fatal(err)
	}
	volume, err := decimal.NewFromString(el[1].(string))
	if err != nil {
		log.Fatal(err)
	}
	priceF := transform.StringToX[float64](el[0].(string)).(float64)
	volumF := transform.StringToX[float64](el[1].(string)).(float64)
	exchangeTime := ParseI(el[2])
	return price, volume, exchangeTime, len(el), priceF, volumF
}

func kkd2orderbook(v KKDepth) (res *depth.Depth) {
	res = &depth.Depth{}
	asks, bids := base.DepthItemSlice{}, base.DepthItemSlice{}
	for _, vv := range v.Asks {
		price := vv.PriceF
		amount := vv.VolumeF
		asks = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	for _, vv := range v.Bids {
		price := vv.PriceF
		amount := vv.VolumeF
		bids = append(bids, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	res.TimeReceive = uint64(v.TimeReceive)
	res.TimeExchange = uint64(v.TimeExchange)
	res.Symbol = v.Symbol
	res.Asks = asks
	res.Bids = bids
	return
}

func kk2clientWs(v *KKDepth) (res *client.WsDepthRsp) {
	res = &client.WsDepthRsp{}
	asks, bids := base.DepthItemSlice{}, base.DepthItemSlice{}
	for _, vv := range v.Asks {
		price := vv.PriceF
		amount := vv.VolumeF
		asks = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	for _, vv := range v.Bids {
		price := vv.PriceF
		amount := vv.VolumeF
		bids = append(asks, &depth.DepthLevel{
			Price:  price,
			Amount: amount,
		})
	}
	res.ExchangeTime = v.TimeExchange
	res.ReceiveTime = v.TimeReceive
	res.Symbol = v.Symbol
	res.Asks = asks
	res.Bids = bids
	return
}

func GetChecksumInput(bids []Level, asks []Level, aPrecision, pPrecision int32) string {
	var str strings.Builder
	var num = 10
	if len(asks) < num {
		num = len(asks)
	}
	for _, level := range asks[:num] {
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
	num = 10
	if len(bids) < num {
		num = len(bids)
	}
	for _, level := range bids[:num] {
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
