package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

func main() {
	serverURL := "ws://localhost:8080/handler2"
	//todo origin必须是"http://xxx"的模式,origin用于判断某个网页是否具有访问websocket的权限
	origin := "http://localhost"
	ws, err := websocket.Dial(serverURL, "", origin)
	if err != nil {
		fmt.Println(err)
	}

	//定时向server端写数据
	go Write(ws)

	//循环读取server写回的数据
	for {
		var msg []byte
		msg = make([]byte, 32)
		//若socket中无数据，程序阻塞在此，等待有数据可读
		_, err := ws.Read(msg[:])
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("Read: %s\n", string(msg[:]))
	}
}

func Write(conn *websocket.Conn) {
	for {
		n, err := conn.Write([]byte("hello world"))//一个字母和一个空格都占一个byte的长度
		if err != nil {
			fmt.Println("client send failed; err:", err)
		}
		fmt.Println("client send success; n=", n)
		//延时2秒
		time.Sleep(1 * time.Second)
	}
}
