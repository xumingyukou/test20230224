package base

import (
	"github.com/warmplanet/proto/go/depth"
)

func mergeDuplicate(deltaDepth DepthItemSlice) DepthItemSlice {
	var mergeIdx []int
	for i := 0; i < len(deltaDepth)-1; i++ {
		if deltaDepth[i+1].Price == deltaDepth[i].Price {
			mergeIdx = append(mergeIdx, i)
		}
	}
	if len(mergeIdx) == 0 {
		return deltaDepth
	}

	mergedDeltaDepth := []*depth.DepthLevel{}
	mergeIdx = append(mergeIdx, len(deltaDepth))
	for i := 0; i < len(mergeIdx); i++ {
		if i == 0 {
			mergedDeltaDepth = append(mergedDeltaDepth, deltaDepth[:mergeIdx[i]]...)
		} else {
			mergedDeltaDepth = append(mergedDeltaDepth, deltaDepth[mergeIdx[i-1]+1:mergeIdx[i]]...)
		}
	}
	return *DepthItemSlice(mergedDeltaDepth).Copy()
}

func MergeDepth(side SIDE, fullDepth *DepthItemSlice, deltaDepth DepthItemSlice) {
	deltaDepth = mergeDuplicate(deltaDepth) // 预先合并重复

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
		} else if (side == Ask && (*fullDepth)[j].Price < deltaDepth[i].Price) || (side == Bid && (*fullDepth)[j].Price > deltaDepth[i].Price) {
			// ask卖盘：全量的price小于增量的price，i不变，j增加
			j++
		} else {
			// ask卖盘：全量的price大于增量的price，需要在全量的j索引前面插入这个amount
			if deltaDepth[i].Amount > 0 {
				if j == 0 {
					*fullDepth = append(DepthItemSlice{
						&depth.DepthLevel{
							Price:  deltaDepth[i].Price,
							Amount: deltaDepth[i].Amount,
						},
					}, *fullDepth...)
				} else {
					tmp := append(DepthItemSlice{}, (*fullDepth)[:j]...)
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

func UpdateBidsAndAsks(deltaDepth *DeltaDepthUpdate, DepthCache *OrderBook, depthLevel int, depthData *depth.Depth) {
	MergeDepth(Ask, &DepthCache.Asks, deltaDepth.Asks)
	MergeDepth(Bid, &DepthCache.Bids, deltaDepth.Bids)
	// 调整最大容量
	if len(DepthCache.Asks) > depthLevel {
		DepthCache.Asks = DepthCache.Asks[:depthLevel]
	}
	if len(DepthCache.Bids) > depthLevel {
		DepthCache.Bids = DepthCache.Bids[:depthLevel]
	}
	DepthCache.Symbol = deltaDepth.Symbol
	DepthCache.Market = deltaDepth.Market
	DepthCache.Type = deltaDepth.Type
	DepthCache.TimeExchange = uint64(deltaDepth.TimeExchange)
	DepthCache.TimeReceive = uint64(deltaDepth.TimeReceive)
	DepthCache.UpdateId = deltaDepth.UpdateEndId
	if depthData != nil {
		DepthCache.Transfer2Depth(depthLevel, depthData)
	}
}
