package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const tradeStream string = "wss://stream.binance.com:9443/ws/btcusdt@trade"

type Client struct {
	conn *websocket.Conn
}

func GetBinanceClient() Client {
	conn, _, err := websocket.DefaultDialer.Dial(tradeStream, nil)
	fmt.Println(tradeStream, 1111)
	if err != nil {
		log.Fatal("Не удалось подключиться к сокету ", err)
	}
	return Client{conn: conn}
}

func (client Client) ReadTrade() map[string]interface{} {
	trade := make(map[string]interface{})
	err := client.conn.ReadJSON(trade)
	if err != nil {
		fmt.Println(err)
	}
	return trade
}

func main() {
	c := GetBinanceClient()
	err := c.conn.WriteMessage(1, []byte("{\"method\": \"SUBSCRIBE\",\"params\":[\"btcusdt@trade\"],\"id\": 1}"))
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
		time.Sleep(time.Second)
	}
}
