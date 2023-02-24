package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/warmplanet/proto/go/common"
	"github.com/warmplanet/proto/go/depth"
	"github.com/warmplanet/proto/go/sdk"
	"github.com/warmplanet/proto/go/sdk/broker"
)

func testPubSub(cc *sdk.TtNatsConfig, workerSize int) {
	var s broker.Sub
	var p broker.Pub
	handler := func(subject string, data []byte) []byte {
		fmt.Println("Recv", subject, string(data))
		time.Sleep(10 * time.Second)
		return nil
	}
	dur := sdk.Duration{}
	dur.UnmarshalText([]byte("10s"))
	subConfig := sdk.SubConfig{WorkerPoolSize: workerSize}
	subjects := []string{"dd.1.2.3.*", "dd.1.2.4.*"}
	s.SubConfig = subConfig

	s.Init(cc, broker.SubHandlers{Data: handler})
	p.Init(cc)

	for idx, subject := range subjects {
		if idx%2 == 0 {
			s.Subscribe(subject)
		} else {
			s.QueueSubscribe(subject, "aaa")
		}
	}

	i := 0
	for {
		//fmt.Println("sleep 1s...")
		time.Sleep(time.Second)
		if i%2 == 0 {
			_ = p.Publish("dd.1.2.3.ETH/USDT", []byte(strconv.Itoa(i)))
			_ = p.Publish("dd.1.2.4.ETH/USDT", []byte(strconv.Itoa(i)))
			//s.Unsubscribe("dd.1.2.3.*")
			//s.Unsubscribe("dd.1.2.4.*")
		} else {
			//s.Subscribe("dd.1.2.3.ETH/USDT")
			//s.Subscribe("dd.1.2.4.*")
		}
		i++
	}
}

func testClientServer(cc *sdk.TtNatsConfig, workerSize int) {
	var s broker.Server
	var c broker.Client
	handler := func(subject string, data []byte) []byte {
		//fmt.Println(subject, string(data))
		r, _ := strconv.Atoi(string(data))
		return []byte(strconv.Itoa(r + 1))
	}

	subjects := []string{"dd.1.2.3.*", "dd.1.2.4.*"}
	s.SubConfig = sdk.SubConfig{WorkerPoolSize: workerSize}
	s.Init(cc, []broker.SrvHandlers{{Data: handler}, {Data: handler}}, subjects)
	c.Init(cc)

	i := 0
	for {
		//fmt.Println("sleep 1s...")
		time.Sleep(time.Second)

		r1, err1 := c.Request("dd.1.2.3.ETH/USDT", []byte(strconv.Itoa(i)), time.Second)
		r2, err2 := c.Request("dd.1.2.4.ETH/USDT", []byte(strconv.Itoa(i)), time.Second)
		fmt.Println(i, err1, err2)
		if r1 != nil {
			fmt.Println("recv from dd.1.2.3.ETH/USDT " + string(r1))
		}
		if r2 != nil {
			fmt.Println("recv from dd.1.2.4.ETH/USDT " + string(r2))
		}
		i++
	}
}

func TestCs(t *testing.T) {

	broker.BrokerLogger = sdk.NewLogger("logs", "depth", 8*time.Hour)

	hdr := common.MsgHeader{Timestamp: uint64(time.Now().UnixMilli())}
	dd := depth.Depth{Hdr: &hdr}
	fmt.Println(dd.Hdr.Timestamp)

	cc := sdk.TtNatsConfig{}
	sdk.LoadConfigFile("../conf/nats.toml", &cc)

	//testPubSub(&cc, 10)
	testClientServer(&cc, 10)
}
