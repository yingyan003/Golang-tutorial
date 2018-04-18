package main

import (
	"time"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
)

type Key struct {
	Key string
	*Key2
}

type Key2 struct {
	Uid int32
}

func main() {
	origin := "http://localhost"
	ws, err := websocket.Dial("ws://127.0.0.1:8080/server3", "", origin)
	if err != nil {
		fmt.Println("dail socket server err: ", err)
	}

	_, err = ws.Write([]byte("--test--"))
	if err != nil {
		fmt.Println("client write 1 error: ", err)
	}

	k := new(Key)
	for {
		k.Key = "1"
		k.Key2 = new(Key2)
		k.Uid = 3
		time.Sleep(2 * time.Second)
		jk, err := json.Marshal(k)
		if err != nil {
			fmt.Println("json marshal err: ", err)
		}
		fmt.Println("json :", k, "k.jey:", k.Key, " k.uid：", k.Uid)
		_, err = ws.Write(jk)
		if err != nil {
			fmt.Println("client write 1 error: ", err)
		}

		time.Sleep(2 * time.Second)
		k.Key = "2"
		jk, _ = json.Marshal(k)
		_, err = ws.Write(jk)
		if err != nil {
			fmt.Println("client write 2 error: ", err)
		}
		fmt.Println("json :", k, "k.jey:", k.Key, " k.uid：", k.Uid)
		fmt.Println("client sleep 5")
		k.Key = ""
		//jk,_=json.Marshal(k)
		_, err = ws.Write([]byte("   "))
		if err != nil {
			fmt.Println("client write 1 error: ", err)
		}
		time.Sleep(5 * time.Second)
	}

}
