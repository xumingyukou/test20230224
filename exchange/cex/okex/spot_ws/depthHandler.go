package spot_ws

import (
	"clients/exchange/cex/base"
	"fmt"
	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/depth"
	"hash/crc32"
	"strings"
	"sync"
	"time"
)

func init() {
}

const SaveType = "depth"

type DepthHandler struct {
	DepthChanCount int
	DepthChanMap   map[int]chan *depth.Depth
	pubSubject     string

	DepthCache       *base.OrderBook             // depth缓存
	depthPublishChan chan *depth.Depth           // 发送channel
	DepthUpdateChan  chan *base.DeltaDepthUpdate // 增量channel
	isStart          bool                        // 防止重复启动
	Lock             sync.Mutex                  // 停止、更新全局depth时加锁
	IsUseWsFunc      bool
	// 新增
	client.WsBookTickerRsp
	DepthCapLevel int
	DepthLevel    int
}

func NewDepthHandler(isUseWsFunc bool) *DepthHandler {
	d := &DepthHandler{}
	d.DepthChanMap = make(map[int]chan *depth.Depth)
	d.depthPublishChan = make(chan *depth.Depth, 1)
	//d.DepthUpdateChan = make(chan *base.DeltaDepthUpdate, conf.DepthChanCap)
	d.IsUseWsFunc = isUseWsFunc
	return d
}

func (d *DepthHandler) MergeDepth(side base.SIDE, fullDepth *base.DepthItemSlice, deltaDepth base.DepthItemSlice) {
	// 1. 与全量价格相等，增量amount不为0，则替换
	// 2. 与全量价格相等，增量amount为0，则删除
	// 3. 全量未找到，增量amount不为0，则插入
	var zeroIdx []int
	i, j := 0, 0
	for i < len(deltaDepth) && j < len(*fullDepth) {
		if (*fullDepth)[j].Price == deltaDepth[i].Price {
			if deltaDepth[i].Amount > 0 {
				(*fullDepth)[j].Amount = deltaDepth[i].Amount
			} else {
				zeroIdx = append(zeroIdx, j)
			}
			i++
			j++
		} else if (side == base.Ask && (*fullDepth)[j].Price < deltaDepth[i].Price) || (side == base.Bid && (*fullDepth)[j].Price > deltaDepth[i].Price) {
			// ask卖盘：全量的price小于增量的price，i不变，j增加
			j++
		} else {
			// ask卖盘：全量的price大于增量的price，需要在全量的j索引前面插入这个amount
			if deltaDepth[i].Amount > 0 {
				if j == 0 {
					*fullDepth = append(base.DepthItemSlice{
						&depth.DepthLevel{
							Price:  deltaDepth[i].Price,
							Amount: deltaDepth[i].Amount,
						},
					}, *fullDepth...)
				} else {
					tmp := append(base.DepthItemSlice{}, (*fullDepth)[:j]...)
					tmp = append(tmp, &depth.DepthLevel{
						Price:  deltaDepth[i].Price,
						Amount: deltaDepth[i].Amount,
					})
					*fullDepth = append(tmp, (*fullDepth)[j:]...)
				}
			} else {
				i++
				continue
			}
			i++
			j++
		}
	}
	//处理剩余的增量数据
	for k := 0; k < len(deltaDepth)-i; k++ {
		if deltaDepth[i+k].Amount > 0 {
			*fullDepth = append(*fullDepth, &depth.DepthLevel{
				Price:  deltaDepth[i+k].Price,
				Amount: deltaDepth[i+k].Amount,
			})
		}
	}
	for i1, idx := range zeroIdx {
		if idx == 0 {
			*fullDepth = append((*fullDepth)[0:0], (*fullDepth)[idx+1:]...)
		} else {
			*fullDepth = append((*fullDepth)[:idx-i1], (*fullDepth)[idx-i1+1:]...)
		}
	}
	return
}

func (d *DepthHandler) UpdateBidsAndAsks(deltaDepth *base.DeltaDepthUpdate) {

	d.MergeDepth(base.Ask, &d.DepthCache.Asks, deltaDepth.Asks)
	d.MergeDepth(base.Bid, &d.DepthCache.Bids, deltaDepth.Bids)
	// 调整最大容量
	if len(d.DepthCache.Asks) > d.DepthCapLevel {
		d.DepthCache.Asks = d.DepthCache.Asks[:d.DepthCapLevel]
	}
	if len(d.DepthCache.Bids) > d.DepthCapLevel {
		d.DepthCache.Bids = d.DepthCache.Bids[:d.DepthCapLevel]
	}
	d.DepthCache.TimeExchange = uint64(deltaDepth.TimeExchange)
	d.DepthCache.TimeReceive = uint64(deltaDepth.TimeReceive)
	d.DepthCache.UpdateId = deltaDepth.UpdateEndId
	depthData := depth.Depth{}
	d.DepthCache.Transfer2Depth(d.DepthLevel, &depthData)
	start := time.Now()
	d.Lock.Unlock()
	//if d.DataHandleAfterMerge != nil {
	//	if err := d.DataHandleAfterMerge(); err != nil {
	//		return
	//	}
	//}

	//fmt.Println(depthData.Symbol, len(d.DepthCache.Bids), len(d.DepthCache.Asks))
	if time.Now().UnixMicro()-start.UnixMicro() > 4000 {
		fmt.Println(depthData.Symbol, "send us time:", time.Now().UnixMicro()-start.UnixMicro(), time.Now().UnixMicro(), start.UnixMicro())
	}

	//base.S3Save(d.exchange, d.Symbol, deltaDepth.Transfer2S3Depth(d.DepthCapLevel))
}

func (d *DepthHandler) dataHandleAfterMerge() error {
	//if !d.checkFullDepth(d.DepthCache) {
	//	fmt.Println("check failed", d.exchange, d.Symbol)
	//	logger.Logger.Error("depth check failed, need to reset", d.exchange, d.Symbol)
	//	d.DepthHandler.Close()
	//	if err := d.DepthHandler.Start(); err != nil {
	//		utils.Logger.Error("restart okex full depth task error, sleep 30 seconds:", d.Symbol, err)
	//		time.Sleep(time.Second * 30)
	//		return err
	//	}
	//}
	return nil
}

func (d *DepthHandler) FullDepthCheck(dep *base.OrderBook) bool {
	/**
	rest api返回1000条数据，但是由于程序开始也是1000条，中间数据抛弃后，会有新的数据在rest api中与增量数据不同，所以不能对比所有，要限定数量
	*/
	lb, la := len(dep.Bids), len(dep.Asks)
	var fields []string
	for i := 0; i < 25; i++ {
		if i < lb {
			e := dep.Bids[i]
			fields = append(fields, ParseFloat(e.Price), ParseFloat(e.Amount))
		}
		if i < la {
			e := dep.Asks[i]
			fields = append(fields, ParseFloat(e.Price), ParseFloat(e.Amount))
		}
	}
	raw := strings.Join(fields, ":")
	cs := crc32.ChecksumIEEE([]byte(raw))
	return int32(dep.UpdateId) == int32(cs)
}
