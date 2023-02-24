package base

import (
	"context"

	"github.com/warmplanet/proto/go/client"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/order"
)

// CexSpotApiInterface SpotApi
type CexSpotApiInterface interface {
	//常规信息
	GetExchange() common.Exchange
	GetSymbols() []string                                   //所有交易对
	GetDepth(*client.SymbolInfo, int) (*depth.Depth, error) //获取行情
	IsExchangeEnable() bool                                 //交易所是否可用
	//Weight()                                  // 权值信息

	//交易手续费
	GetTradeFee(...string) (*client.TradeFee, error)                     //查询交易手续费
	GetTransferFee(common.Chain, ...string) (*client.TransferFee, error) //查询转账手续费
	GetPrecision(...string) (*client.Precision, error)                   //查询交易对精读信息
	//资产查询
	GetBalance() (*client.SpotBalance, error)         //获得现货的balance信息
	GetMarginBalance() (*client.MarginBalance, error) //获得全仓杠杆的balance信息
	//todo GetMaxMarginBorrow(string)() // 逐仓、全仓的最大可借量

	//交易相关
	PlaceOrder(*order.OrderTradeCEX) (*client.OrderRsp, error)                   //下单
	CancelOrder(*order.OrderCancelCEX) (*client.OrderRsp, error)                 //取消订单
	GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error)                 //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
	GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error)     //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
	GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
	//提币
	Transfer(*order.OrderTransfer) (*client.OrderRsp, error)                               //-转账订单，根据提币手续费配置划转金额，binance不扣除，ok手动添加手续费参数
	GetTransferHistory(req *client.TransferHistoryReq) (*client.TransferHistoryRsp, error) //查询提款记录
	//划转
	MoveAsset(*order.OrderMove) (*client.OrderRsp, error)                      //-资产划转订单，支持子账户
	GetMoveHistory(req *client.MoveHistoryReq) (*client.MoveHistoryRsp, error) //-查询划转记录，支持子账户
	//充值
	GetDepositHistory(req *client.DepositHistoryReq) (*client.DepositHistoryRsp, error) //-查询存款记录，支持子账户

	//借贷还款
	Loan(*order.OrderLoan) (*client.OrderRsp, error)                          //借贷
	GetLoanOrders(*client.LoanHistoryReq) (*client.LoanHistoryRsp, error)     //获取已放款订单
	Return(*order.OrderReturn) (*client.OrderRsp, error)                      //还币
	GetReturnOrders(*client.LoanHistoryReq) (*client.ReturnHistoryRsp, error) //获取已还币订单
}

type CexFutureApiInterface interface {
	//常规信息
	GetFutureSymbols(common.Market) []*client.SymbolInfo                                   //所有交易对
	GetFutureDepth(*client.SymbolInfo, int) (*depth.Depth, error)                          //获取行情
	GetFutureMarkPrice(common.Market, ...*client.SymbolInfo) (*client.RspMarkPrice, error) //标记价格
	// 币本位永续合约的标记价格
	// 标记价格 = 中位数* (价位1, 价位2, 合约价格)
	// 价位 1 = 价格指数* (1 + 资金费率 *(距离下次资金费率收取的时间（小时）/8))
	// 价位 2 = 价格指数+ 移动平均值(30分钟基础)*
	// 移动平均线（30分钟基础）=移动平均线（（Bid1 + Ask1）/ 2-价格指数），以30分钟为间隔，每分钟采样取值
	// 中位数: 价位1, 价位2, 合约价格三个数取中间那个, 例如价位1 < 价位2 < 合约价格，则标记价格取价位2
	//Weight()                                  // 权值信息

	//交易手续费
	GetFutureTradeFee(common.Market, ...*client.SymbolInfo) (*client.TradeFee, error)   //查询交易手续费
	GetFuturePrecision(common.Market, ...*client.SymbolInfo) (*client.Precision, error) //查询交易对精读信息

	//资产查询
	GetFutureBalance(common.Market) (*client.UBaseBalance, error) //todo 获得合约的balance信息

	//交易相关
	PlaceFutureOrder(*order.OrderTradeCEX) (*client.OrderRsp, error)             //下单
	CancelOrder(*order.OrderCancelCEX) (*client.OrderRsp, error)                 //取消订单
	GetOrder(req *order.OrderQueryReq) (*client.OrderRsp, error)                 //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
	GetOrderHistory(req *client.OrderHistoryReq) ([]*client.OrderRsp, error)     //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
	GetProcessingOrders(req *client.OrderHistoryReq) ([]*client.OrderRsp, error) //查询订单信息。不支持ws回报订单行情的交易所，如bitbank, 需要频发查单获知订单状态
}

type CexApiInterface interface {
	CexSpotApiInterface
	CexFutureApiInterface
}

// CexDepthInterface depth行情相关
type CexDepthInterface interface {
	KeepAliveDepth( /*info_precious*/ )           //保持ws连接的方法，有些是ping-pong，上层服务将data发送到对方服务
	GetConnectDepthUrl( /*market, symbol*/ )      // 返回ws接口订阅url作为connect-str，订阅特定的交易区的订单簿
	GetBooks( /*market, symbol*/ )                //rest 主动获得交易所的当前行情
	ParseDepth( /*ws.data*/ )                     //解析ws接口返回的数据，拆装箱，得到Tickpape统一格式的订单簿数据
	CheckDepthData( /*depth_all, depth_update*/ ) //判断当前的增量数据是否有断档、丢失。发现数据断档或丢失，主流程通常的做法是断开重连，或忽略继续，并额外调用get_books方法更新本地的全量数据，主流程根据ret来定
}

// CexOrderInterface 订单处理相关
type CexOrderInterface interface {
	KeepAliveOrder( /*info_precious*/ ) //保持ws连接的方法，有些是ping-pong，上层服务将data发送到对方服务
	GetConnectOrderUrl()                //返回ws接口订阅登录签名信息作为connect-str，订阅订单的状态变化
	ParseTradeInfo( /*ws.data*/ )       //解析ws接口返回的数据，拆装箱，解析获得ordersys统一格式的订单信息
}

type CexWebsocketPublicInterface interface {
	FundingRateGroup(context.Context, map[*client.SymbolInfo]chan *client.WsFundingRateRsp) error                                                                      //资金费率
	TradeGroup(context.Context, map[*client.SymbolInfo]chan *client.WsTradeRsp) error                                                                                  //交易数据
	BookTickerGroup(context.Context, map[*client.SymbolInfo]chan *client.WsBookTickerRsp) error                                                                        //最优挂单
	DepthLimitGroup(context.Context, int, int, map[*client.SymbolInfo]chan *depth.Depth) error                                                                         //有限档全量
	DepthIncrementGroup(context.Context, int, map[*client.SymbolInfo]chan *client.WsDepthRsp) error                                                                    //增量
	DepthIncrementSnapshotGroup(context.Context, int, int, bool, bool, map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth) error //增量合并为全量副本，并校验
}

type CexWebsocketInterface interface {
	Account(context.Context, ...*client.WsAccountReq) (<-chan *client.WsAccountRsp, error)
	Balance(context.Context, ...*client.WsAccountReq) (<-chan *client.WsBalanceRsp, error)
	Order(context.Context, ...*client.WsAccountReq) (<-chan *client.WsOrderRsp, error)
}

// DepthIncrementSnapshot 合并增量数据，并发送增量或者全量数据
type DepthIncrementSnapshot interface {
	FullDepthCheck() bool //检查主函数
}

type WebSocketPublicHandleInterface interface {
	FundingRateGroupHandle([]byte) error
	AggTradeGroupHandle([]byte) error
	TradeGroupHandle([]byte) error
	BookTickerGroupHandle([]byte) error
	DepthIncrementGroupHandle([]byte) error
	DepthLimitGroupHandle([]byte) error
	DepthIncrementSnapShotGroupHandle([]byte) error

	//设置chan map
	SetFundingRateGroupChannel(map[*client.SymbolInfo]chan *client.WsFundingRateRsp)
	SetTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsTradeRsp)
	SetBookTickerGroupChannel(map[*client.SymbolInfo]chan *client.WsBookTickerRsp)
	SetDepthLimitGroupChannel(map[*client.SymbolInfo]chan *depth.Depth)
	SetDepthIncrementGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp)
	SetDepthIncrementSnapshotGroupChannel(map[*client.SymbolInfo]chan *client.WsDepthRsp, map[*client.SymbolInfo]chan *depth.Depth)
	SetAggTradeGroupChannel(map[*client.SymbolInfo]chan *client.WsAggTradeRsp) // binance

	SetDepthIncrementSnapShotConf([]*client.SymbolInfo, *IncrementDepthConf)
}

type WebSocketPrivateHandleInterface interface {
	AccountHandle([]byte) error
	BalanceHandle([]byte) error
	MarginAccountHandle([]byte) error
	MarginBalanceHandle([]byte) error
	OrderHandle([]byte) error
	GetChan(chName string) interface{}
}

type WebSocketOrderHandleI interface {
	PlaceOrders([]*order.OrderTradeCEX) ([]*client.OrderRsp, error)   //批量下单
	CancelOrders([]*order.OrderCancelCEX) ([]*client.OrderRsp, error) //批量取消订单
	CreateConn() error                                                //创建连接
}

type WebSocketHandleInterface interface {
	WebSocketPublicHandleInterface
	WebSocketPrivateHandleInterface
}
