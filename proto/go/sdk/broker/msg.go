package broker

import (
	"bytes"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tchap/go-patricia/v2/patricia"
	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/sdk"
	"github.com/warmplanet/proto/go/sdk/pond"
)

/**
Pub/Sub和Client/Server是基本的nats封装，进行不可靠消息传送可以使用
NakPub/NakSub实现了可靠消息广播，需要可靠消息传送的可以使用
API是线程安全的，回调函数不是线程安全的，需要调用者自己保护
*/

// protobuf结构体
type PbMsg interface {
	GetHdr() *common.MsgHeader
	MarshalVT() (dAtA []byte, err error)
	UnmarshalVT(dAtA []byte) error

	MarshalToVT(dAtA []byte) (int, error)
}

// 原始消息处理函数
type DataHandler func(subject string, data []byte) []byte

// Protobuf消息处理函数
type PbMsgHandler func(subject string, msg PbMsg) []byte

// Protobuf结构体创建方法
type PbMsgNew func() PbMsg

// 消费者回调函数集
type SubHandlers struct {
	Data DataHandler
	Pb   PbMsgHandler
	New  PbMsgNew
}

// 服务器回调函数集
type SrvHandlers struct {
	Data DataHandler
	Pb   PbMsgHandler
	New  PbMsgNew
}

/*
	一个pubsub实例使用一条nats连接，可以发起多个订阅，同一个订阅内部的消息是有序的，
	不同订阅之间的消息顺序是不确定的，使用中可以根据不同的业务使用多条连接，每条连接使用多个subject和固定的消息类型
 	nats会自动重连所有的订阅，但是nats不会对订阅去重，如果一个连接上的多个订阅匹配消息，
 	订阅者会接收到一条消息多次，为此在应用层进行去重处理，确保所有的订阅不重复且不重叠
*/

// server订阅subjects时，绑定到固定的group
const serverGroup string = "dt"

// 只支持形如a.b.c或者a.b.*或者a.b.>的通配符订阅，不支持通配符在中间
func isValidSub(subject string) (bool, string) {
	idx1 := strings.Index(subject, "*")
	idx2 := strings.Index(subject, ">")

	if (idx1 == len(subject)-1 && idx2 == -1) || (idx2 == len(subject)-1 && idx1 == -1) {
		return true, subject[:len(subject)-1]
	} else if idx1 == -1 && idx2 == -1 {
		return true, subject
	} else {
		return false, subject
	}
}

type Pub struct {
	sdk.PubConfig
	con  *nats.Conn
	Uid  []byte // 生产者唯一ID
	seqs sdk.ConcurrentMap
}

// 每个连接可以发起多个订阅，默认每个订阅会使用一个单独的go routine处理消息
// 如果想提升消息处理的并行度，可以初始化时设置workerPoolSize参数
type Sub struct {
	sdk.SubConfig
	con        *nats.Conn
	msgHandler nats.MsgHandler
	tree       *patricia.Trie
	treeMtx    sync.Mutex

	pool *pond.WorkerPool
}

type Client struct {
	con *nats.Conn
}

// Server和普通sub唯一的不同是要求msg必须携带reply-to字段，是它能在handler里调用msg.Respond
// 同时Server也不支持动态添加和删除subject
type Server struct {
	Sub
}

func connect(config *sdk.TtNatsConfig) *nats.Conn {
	options := make([]nats.Option, 0)

	if config.ConnectionName != "" {
		o := nats.Name(config.ConnectionName)
		options = append(options, o)
	}

	if config.NkeysSeed != "" {
		o, err := nats.NkeyOptionFromSeed(config.NkeysSeed)
		if err != nil {
			log.Fatal(err)
		}
		options = append(options, o)
	}

	if config.ConnectTimeout.Nanoseconds() != 0 {
		o := nats.Timeout(config.ConnectTimeout.Duration)
		options = append(options, o)
	}

	if config.PingInterval.Nanoseconds() != 0 {
		o := nats.PingInterval(config.PingInterval.Duration)
		options = append(options, o)
	}

	if config.MaxPingsOutstanding != 0 {
		o := nats.MaxPingsOutstanding(config.MaxPingsOutstanding)
		options = append(options, o)
	}

	if config.MaxReconnects != 0 {
		o := nats.MaxReconnects(config.MaxReconnects)
		options = append(options, o)
	}

	if config.ReconnectWait.Nanoseconds() != 0 {
		o := nats.ReconnectWait(config.ReconnectWait.Duration)
		options = append(options, o)
	}

	if config.ReconnectBufSize != 0 {
		o := nats.ReconnectBufSize(config.ReconnectBufSize)
		options = append(options, o)
	}

	options = append(options,
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			BrokerLogger.Warnf("Error while handle msg from broker: %v", err)
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			BrokerLogger.Warnf("Disconnect broker due to %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			BrokerLogger.Warnf("Reconnect broker done")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			BrokerLogger.Errorf("Connection to broker closed")
		}))

	con, err := nats.Connect(strings.Join(config.Servers, ","), options...)
	if err != nil {
		log.Fatal(err)
	}

	return con
}

func (ps *Sub) sub(tree *patricia.Trie, subject string, queueGroup string, handler nats.MsgHandler) {
	var key string
	var ok bool

	if ok, key = isValidSub(subject); !ok {
		log.Fatal("Invalid subject subscribe: " + subject)
	} else {
		overlapped := 0
		existed := 0
		kk := patricia.Prefix(key)
		// 寻找所有匹配key部分前缀的节点，如果找到，认为有overlap，报错
		tree.VisitPrefixes(kk, func(prefix patricia.Prefix, item patricia.Item) error {
			// 完全相等不算overlap
			if !bytes.Equal(kk, prefix) {
				overlapped = overlapped + 1
			} else {
				existed = 1
			}
			return nil
		})
		if overlapped > 0 {
			BrokerLogger.Errorf("Subject %s overlapped with previous", subject)
			return
		} else if existed > 0 {
			BrokerLogger.Infof("Subject %s already subscribed", subject)
			return
		}
	}

	var subptr *nats.Subscription
	var err error
	if queueGroup != "" {
		subptr, err = ps.con.QueueSubscribe(subject, queueGroup, handler)
	} else {
		subptr, err = ps.con.Subscribe(subject, handler)
	}

	if err != nil {
		BrokerLogger.Error(err)
	} else {
		tree.Insert(patricia.Prefix(key), subptr)
	}
}

func (ps *Sub) handlerConvert(handlers SubHandlers) nats.MsgHandler {
	var fn func(msg *nats.Msg)

	if handlers.Data != nil {
		handler := handlers.Data
		fn = func(msg *nats.Msg) {
			// 返回值忽略，正常应该返回nil
			handler(msg.Subject, msg.Data)
		}
	} else if handlers.Pb != nil && handlers.New != nil {
		msgNew := handlers.New
		handler := handlers.Pb
		fn = func(msg *nats.Msg) {
			pbMsg := msgNew()
			err := pbMsg.UnmarshalVT(msg.Data)
			if err == nil {
				handler(msg.Subject, pbMsg)
			} else {
				BrokerLogger.Errorf("Sub handle pb msg error: %v", err)
			}
		}
	} else {
		log.Fatal("Invalid SrvHandlers")
	}

	if ps.WorkerPoolSize > 0 {
		pool := ps.pool
		return func(msg *nats.Msg) {
			pool.Submit(func() {
				fn(msg)
			})
		}
	} else {
		return fn
	}
}

// 注册处理原始消息数据的回调函数
func (ps *Sub) Init(config *sdk.TtNatsConfig, handlers SubHandlers) {
	// 只能初始化一次
	if ps.msgHandler != nil {
		log.Fatal("Initing initialized Sub!")
	}
	ps.con = connect(config)
	tree := patricia.NewTrie()
	if ps.WorkerPoolSize > 0 {
		ps.pool = pond.New(ps.WorkerPoolSize, ps.QueueSize, pond.Strategy(pond.Eager()))
	}

	ps.msgHandler = ps.handlerConvert(handlers)
	ps.tree = tree
}

func (ps *Sub) Subscribe(subject string) {
	ps.treeMtx.Lock()
	defer ps.treeMtx.Unlock()

	ps.sub(ps.tree, subject, "", ps.msgHandler)
}

func (ps *Sub) SubscribeEx(subject string, handlers SubHandlers) {
	ps.treeMtx.Lock()
	defer ps.treeMtx.Unlock()

	ps.sub(ps.tree, subject, "", ps.handlerConvert(handlers))
}

// 重复订阅，使用者需要保存返回的handle用于取消订阅
func (ps *Sub) SubscribeDup(subject string) (subHandle interface{}) {
	var err error
	subHandle, err = ps.con.Subscribe(subject, ps.msgHandler)

	if err != nil {
		BrokerLogger.Error(err)
		return nil
	} else {
		return
	}
}

func (ps *Sub) SubscribeDupEx(subject string, handlers SubHandlers) (subHandle interface{}) {
	var err error
	subHandle, err = ps.con.Subscribe(subject, ps.handlerConvert(handlers))

	if err != nil {
		BrokerLogger.Error(err)
		return nil
	} else {
		return
	}
}

func (ps *Sub) Unsubscribe(subject string) bool {
	ps.treeMtx.Lock()
	defer ps.treeMtx.Unlock()

	var prefix string
	var ok bool
	if ok, prefix = isValidSub(subject); !ok {
		return false
	} else {
		item, ok := (ps.tree.Get(patricia.Prefix(prefix))).(*nats.Subscription)
		if ok && item != nil && item.Subject == subject {
			if item.Unsubscribe() == nil {
				ps.tree.Delete(patricia.Prefix(prefix))
				return true
			}
		}
	}

	return false
}

func (ps *Sub) UnsubscribeDup(subHandle interface{}) bool {
	if s, ok := subHandle.(*nats.Subscription); !ok {
		return false
	} else if s.Unsubscribe() == nil {
		return true
	} else {
		return false
	}
}

// 模拟多个消费者读取队列，每条队列里的消息只处理一次
func (ps *Sub) QueueSubscribe(subject string, queueGroup string) {
	ps.treeMtx.Lock()
	defer ps.treeMtx.Unlock()

	if queueGroup == "" {
		log.Fatal("queueGroup is null")
	}

	ps.sub(ps.tree, subject, queueGroup, ps.msgHandler)
}

func (ps *Sub) QueueSubscribeEx(subject string, queueGroup string, handlers SubHandlers) {
	ps.treeMtx.Lock()
	defer ps.treeMtx.Unlock()

	if queueGroup == "" {
		log.Fatal("queueGroup is null")
	}

	ps.sub(ps.tree, subject, queueGroup, ps.handlerConvert(handlers))
}

func (ps *Sub) QueueSubscribeDup(subject string, queueGroup string) (subHandle interface{}) {
	var err error
	subHandle, err = ps.con.QueueSubscribe(subject, queueGroup, ps.msgHandler)

	if err != nil {
		BrokerLogger.Error(err)
		return nil
	} else {
		return
	}
}

func (ps *Sub) QueueSubscribeDupEx(subject string, queueGroup string, handlers SubHandlers) (subHandle interface{}) {
	var err error
	subHandle, err = ps.con.QueueSubscribe(subject, queueGroup, ps.handlerConvert(handlers))

	if err != nil {
		BrokerLogger.Error(err)
		return nil
	} else {
		return
	}
}

func (ps *Pub) Init(config *sdk.TtNatsConfig) {
	ps.seqs = sdk.NewCmap()
	ps.Uid = []byte(ps.UniqueName)
	if len(ps.Uid) == 0 {
		ps.Uid = []byte(strconv.Itoa(int(rand.Int31())))
	}
	ps.con = connect(config)
}

// 发送自定义消息
func (ps *Pub) Publish(subject string, data []byte) error {
	return ps.con.Publish(subject, data)
}

// 每个pub对应每个subject一个自增序列
func upsert(exist bool, valueInMap int64, newValue int64) int64 {
	if !exist {
		return newValue
	}

	return valueInMap + 1
}

// 发送protobuf消息
func (ps *Pub) PublishMsg(subject string, msg PbMsg) error {
	hdr := msg.GetHdr()
	if hdr == nil {
		return errors.New("nil header message")
	}
	// 微秒级时间戳
	hdr.Timestamp = uint64(time.Now().UnixMicro())
	hdr.Producer = ps.Uid
	// 自增序列号
	nextSeq := ps.seqs.Upsert(subject, 1, upsert)
	hdr.Sequence = nextSeq

	data, err := msg.MarshalVT()
	if err != nil {
		BrokerLogger.Errorf("Pub send pb msg error: %v", err)
		return err
	}

	// TODO: 如果开启了持久化，需要写入消息
	//if ps.Persist {}

	return ps.con.Publish(subject, data)
}

func (ps *Pub) PublishMsg2(subject string, msg []byte) error {
	return ps.con.Publish(subject, msg)
}

func (ps *Pub) SerializeMsg(subject string, msg PbMsg) ([]byte, error) {
	hdr := msg.GetHdr()
	if hdr == nil {
		return nil, errors.New("nil header message")
	}
	// 微秒级时间戳
	hdr.Timestamp = uint64(time.Now().UnixMicro())
	hdr.Producer = ps.Uid
	// 自增序列号
	nextSeq := ps.seqs.Upsert(subject, 1, upsert)
	hdr.Sequence = nextSeq

	data, err := msg.MarshalVT()
	if err != nil {
		BrokerLogger.Errorf("serialize send pb msg error: %v", err)
		return nil, err
	}
	return data, nil
}

func (c *Client) Init(config *sdk.TtNatsConfig) {
	c.con = connect(config)
}

func (c *Client) Request(subject string, data []byte, timeout time.Duration) ([]byte, error) {
	rsp, err := c.con.Request(subject, data, timeout)
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (c *Client) RequestMsg(subject string, msg PbMsg, timeout time.Duration) ([]byte, error) {
	data, err := msg.MarshalVT()
	if err != nil {
		return nil, err
	}

	rsp, err := c.con.Request(subject, data, timeout)
	if err != nil {
		return nil, err
	}

	return rsp.Data, nil
}

func (s *Server) Init(config *sdk.TtNatsConfig, handlers []SrvHandlers, subjects []string) {
	s.con = connect(config)
	var fn func(msg *nats.Msg)

	for i, subject := range subjects {
		if handlers[i].Data != nil {
			handler := handlers[i].Data
			fn = func(msg *nats.Msg) {
				rsp := handler(msg.Subject, msg.Data)

				if rsp != nil {
					if err := msg.Respond(rsp); err != nil {
						BrokerLogger.Errorf("Send response by '%s' error %v", config.ConnectionName)
					}
				} else {
					BrokerLogger.Error("Response data from handler is nil")
				}
			}

		} else if handlers[i].Pb != nil && handlers[i].New != nil {
			msgNew := handlers[i].New
			handler := handlers[i].Pb
			fn = func(msg *nats.Msg) {
				pbMsg := msgNew()
				err := pbMsg.UnmarshalVT(msg.Data)

				var rsp []byte
				if err == nil {
					rsp = handler(msg.Subject, pbMsg)
				}

				if rsp != nil {
					if err := msg.Respond(rsp); err != nil {
						BrokerLogger.Errorf("Send response by '%s' error %v", config.ConnectionName)
					}
				} else {
					BrokerLogger.Error("Response data from handler is nil")
				}
			}
		} else {
			log.Fatal("Invalid SrvHandlers")
		}

		if s.WorkerPoolSize > 0 {
			pool := pond.New(s.WorkerPoolSize, s.QueueSize, pond.Strategy(pond.Eager()))
			s.pool = pool
			if _, err := s.con.QueueSubscribe(subject, serverGroup, func(msg *nats.Msg) {
				pool.Submit(func() {
					fn(msg)
				})
			}); err != nil {
				BrokerLogger.Error(err)
			}
		} else {
			if _, err := s.con.QueueSubscribe(subject, serverGroup, fn); err != nil {
				BrokerLogger.Error(err)
			}
		}
	}
}

func (s *Server) Subscribe(subject string) {
	log.Fatal("Unsupported 'Subscribe' on Server")
}

func (s *Server) SubscribeEx(subject string) {
	log.Fatal("Unsupported 'Subscribe' on Server")
}

func (s *Server) SubscribeDup(subject string) (subHandler interface{}) {
	log.Fatal("Unsupported 'Subscribe' on Server")
	return nil
}

func (s *Server) SubscribeDupEx(subject string) (subHandler interface{}) {
	log.Fatal("Unsupported 'Subscribe' on Server")
	return nil
}

func (s *Server) Unsubscribe(subject string) bool {
	log.Fatal("Unsupported 'Unsubscribe' on Server")
	return false
}

func (ps *Server) UnsubscribeDup(subHandle interface{}) bool {
	log.Fatal("Unsupported 'Unsubscribe' on Server")
	return false
}

func (s *Server) QueueSubscribe(subject string, queueGroup string) {
	log.Fatal("Unsupported 'QueueSubscribe' on Server")
}

func (s *Server) QueueSubscribeEx(subject string, queueGroup string) {
	log.Fatal("Unsupported 'QueueSubscribe' on Server")
}

func (s *Server) QueueSubscribeDup(subject string, queueGroup string) (subHandle interface{}) {
	log.Fatal("Unsupported 'QueueSubscribe' on Server")
	return false
}

func (s *Server) QueueSubscribeDupEx(subject string, queueGroup string) (subHandle interface{}) {
	log.Fatal("Unsupported 'QueueSubscribe' on Server")
	return false
}
