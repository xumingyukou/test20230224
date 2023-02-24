package main

import (
	"clients/conn"
	"fmt"
	"time"
)

type Request struct {
	Jsonrpc string                 `json:"jsonrpc,omitempty"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
	ID      int                    `json:"id,omitempty"`
}

func main() {
	var (
		readTimeout = 100 * time.Second
		url         = "wss://ws.lightstream.bitflyer.com/json-rpc"
		//proxyUrl    = "http://127.0.0.1:7890"
		reqx = Request{
			//Jsonrpc: "2.0",
			Method: "subscribe",
			Params: map[string]interface{}{
				"channel": "lightning_executions_BTC_JPY",
			},
			//ID: 1,
		}
	)
	wsBuilder := conn.NewWsBuilderWithReadTimeout(readTimeout).WsUrl(url).ProtoHandleFunc(handler)
	wsClient := wsBuilder.Build()
	if wsClient == nil {
		fmt.Println("create client err")
		return
	}
	err := Subscribe(wsClient, reqx)
	if err != nil {
		fmt.Println(err)
		return
	}
	select {}
}

func Subscribe(client *conn.WsConn, reqx interface{}) error {
	err := client.Subscribe(reqx)
	if err != nil {
		fmt.Println("subscribe err:", err)
	}
	return err
}

func handler(data []byte) error {
	fmt.Println(string(data))
	return nil
}
