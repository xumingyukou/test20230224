package base

import (
	"clients/conn"
	"clients/logger"
	"clients/transform"
	"context"
	"fmt"
	"time"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/sdk"
)

type OptionalParameter map[string]interface{}

type Account struct{}
type Balance struct{}
type Order struct{}
type OrderResp struct{}
type TransferOrder struct{}
type DepositOrder struct{}

type SIDE int8
type SIDESTRING string

const (
	Err SIDE = iota
	Ask
	Bid
)

const (
	AskStr SIDESTRING = "sell"
	BidStr            = "buy"
)

func ConvertSide(side SIDE) SIDE {
	if side == Ask {
		return Bid
	} else {
		return Ask
	}
}

type APIConf struct {
	ProxyUrl    string
	EndPoint    string
	ReadTimeout int64 //second
	AccessKey   string
	SecretKey   string
	Passphrase  string
	IsTest      bool   //是否测试
	IsPrivate   bool   //是否使用私有地址
	SubAccount  string //如果是子账户，填入母账户名.子账户名，否则填入普通账户名
}

type WsConf struct {
	APIConf
	ChanCap int64
}

type IncrementDepthConf struct {
	Exchange               common.Exchange                       // 币对信息
	Market                 common.Market                         // 币对信息
	Type                   common.SymbolType                     // 币对信息
	IsPublishDelta         bool                                  // 是否发送增量depth
	IsPublishFull          bool                                  // 是否发送全量depth
	DepthCacheMap          sdk.ConcurrentMapI                    // string:*orderbook.OrderBook
	DepthCacheListMap      sdk.ConcurrentMapI                    // string:[]*orderbook.DeltaDepthUpdate
	DepthDeltaUpdateMap    sdk.ConcurrentMapI                    // string:chan *orderbook.DeltaDepthUpdate
	CheckDepthCacheChanMap sdk.ConcurrentMapI                    // string:chan *orderbook.OrderBook
	DepthNotMatchChanMap   map[*client.SymbolInfo]chan bool      // 发送校验错误
	CheckTimeSec           int                                   // 检查的时间间隔
	DepthCapLevel          int                                   // 缓存深度，防止数据不准
	DepthLevel             int                                   // 订阅深度
	DepthCheckLevel        int                                   // 检查depth深度
	GetFullDepth           func(string) (*OrderBook, error)      // 获取全量数据
	GetFullDepthLimit      func(string, int) (*OrderBook, error) // 获取有限档全量数据，用于检查
	Ctx                    context.Context                       // 上下文，用于关闭检查的goroutine
	CheckStates            sdk.ConcurrentMapI                    // string:bool 正在检查的标志位
}

type CheckDataSendStatus struct {
	CheckUpdateTimeMap sdk.ConcurrentMapI // string:updatetime
	CheckUpdateTimeSec int                // 检查数据更新的时间间隔
	CheckFutureDateMap map[string]string  // 期货绝对时间map
	GetDateFunc        func(transform.CONTRACTTYPE) string
	UpdateTimeoutChMap chan []*client.SymbolInfo          // 发送超时错误
	UpdateDateChMap    chan map[string]*client.SymbolInfo // okex期货重订阅
	SymbolInfoMap      map[string]*client.SymbolInfo      // 保存symbolinfo数据
	Ctx                context.Context                    // 上下文
	ReSubscribing      bool
	ReSubInfo          *ReSubInfo
}

func NewCheckDataSendStatus() *CheckDataSendStatus {
	c := &CheckDataSendStatus{}

	c.CheckUpdateTimeMap = sdk.NewCmapI()
	c.CheckFutureDateMap = make(map[string]string, 0)
	c.UpdateTimeoutChMap = make(chan []*client.SymbolInfo, 100)
	c.SymbolInfoMap = make(map[string]*client.SymbolInfo)
	c.ReSubInfo = &ReSubInfo{}
	return c
}

type ReSubInfo struct {
	ReSubFunc func(*conn.WsConn, []*client.SymbolInfo, string)
	WsConn    *conn.WsConn
	Event     string
}

// InitDecorator channel无法接受到数据，临时方案......
func (c *CheckDataSendStatus) InitDecorator(reSubFunc func(*conn.WsConn, []*client.SymbolInfo, string), checkUpdateTimeSec int, ctx context.Context, wsConn *conn.WsConn,
	event string, symbols ...*client.SymbolInfo) {
	c.ReSubInfo.WsConn = wsConn
	c.ReSubInfo.Event = event
	c.ReSubInfo.ReSubFunc = reSubFunc
	c.Init(checkUpdateTimeSec, ctx, symbols...)
}

func (c *CheckDataSendStatus) InitGetDateFunc(getDateFunc func(contracttype transform.CONTRACTTYPE) string) {
	c.GetDateFunc = getDateFunc
}

func (c *CheckDataSendStatus) Init(checkUpdateTimeSec int, ctx context.Context, symbols ...*client.SymbolInfo) {
	c.Ctx = ctx
	c.CheckUpdateTimeSec = checkUpdateTimeSec
	if checkUpdateTimeSec <= 0 {
		checkUpdateTimeSec = 600
	}
	for _, symbol := range symbols {
		formatSymbolName := SymInfoToString(symbol)
		c.CheckUpdateTimeMap.Set(formatSymbolName, time.Now().UnixMicro())
		c.SymbolInfoMap[formatSymbolName] = symbol
		if c.GetDateFunc != nil && (symbol.Market == common.Market_FUTURE || symbol.Market == common.Market_FUTURE_COIN) {
			c.CheckFutureDateMap[SymInfoToString(symbol)] = c.GetDateFunc(transform.CONTRACTTYPE(TransSymbolTypeToString(symbol.Type)))
		}
	}
	go c.Check()
	go c.CheckFutureReSub()
}

func (c *CheckDataSendStatus) Check() {
	/*
		定期检查数据更新情况
	*/
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Error("check data update time panic:", logger.PanicTrace(err))
			time.Sleep(time.Second)
			go c.Check()
		}
	}()
	heartTimer := time.NewTimer(time.Duration(600) * time.Second)
LOOP:
	for {
		select {
		case <-heartTimer.C:
			if c.ReSubscribing {
				continue
			}
			var symbolList []*client.SymbolInfo
			for _, symbol := range c.CheckUpdateTimeMap.Keys() {
				preTime, _ := c.CheckUpdateTimeMap.Get(symbol)
				if preTimeUs, ok := preTime.(int64); ok {
					if time.Now().Sub(time.UnixMicro(preTimeUs)) > time.Duration(c.CheckUpdateTimeSec)*time.Second {
						fmt.Println(symbol, time.Now(), time.UnixMicro(preTimeUs), time.Now().Sub(time.UnixMicro(preTimeUs)), time.Duration(c.CheckUpdateTimeSec)*time.Second)
						c.CheckUpdateTimeMap.Set(symbol, time.Now().UnixMicro())
						if symInfo, ok := c.SymbolInfoMap[symbol]; ok {
							symbolList = append(symbolList, symInfo)
						}
						logger.Logger.Error(symbol, " ", time.Now(), time.Now().Sub(time.UnixMicro(preTimeUs)), "need resubscribe")
					}
				}
			}
			if len(symbolList) > 0 {
				if c.ReSubInfo.ReSubFunc == nil {
					c.UpdateTimeoutChMap <- symbolList
				} else {
					c.ReSubscribing = true
					c.ReSubInfo.ReSubFunc(c.ReSubInfo.WsConn, symbolList, c.ReSubInfo.Event)
					c.ReSubscribing = false
				}
			}
			heartTimer.Reset(time.Duration(c.CheckUpdateTimeSec) * time.Second)
		case <-c.Ctx.Done():
			break LOOP
		}
	}
}

func (c *CheckDataSendStatus) CheckFutureReSub() {
	/*
		定期检查期货重订阅
	*/
	if len(c.CheckFutureDateMap) == 0 || c.GetDateFunc == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Error("check future resubscribe time panic:", logger.PanicTrace(err))
			time.Sleep(time.Second)
			go c.CheckFutureReSub()
		}
	}()

	heartTimer := time.NewTimer(time.Duration(600) * time.Second)
LOOP:
	for {
		select {
		case <-heartTimer.C:
			originDateSymbolMap := make(map[string]*client.SymbolInfo, 0)
			for formatSymbol, futureDate := range c.CheckFutureDateMap {
				symbolInfo := c.SymbolInfoMap[formatSymbol]
				date := c.GetDateFunc(transform.CONTRACTTYPE(TransSymbolTypeToString(symbolInfo.Type)))
				if futureDate == date {
					continue
				}
				c.CheckFutureDateMap[formatSymbol] = date

				originDateSymbolMap[futureDate] = symbolInfo
				logger.Logger.Info(symbolInfo.Symbol, " future date need resubscribe")
			}
			if len(originDateSymbolMap) > 0 {
				c.UpdateDateChMap <- originDateSymbolMap
			}
			heartTimer.Reset(time.Duration(c.CheckUpdateTimeSec) * time.Second)
		case <-c.Ctx.Done():
			break LOOP
		}
	}
}
