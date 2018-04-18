package main

import (
	"golang.org/x/net/websocket"
	"fmt"
	"time"
	"net/http"
)

//socket handler
//todo socket读写有2种方式
//todo 方式1. websocket.Message.Receive(conn, 数据接收体)/websocket.Message.Send(conn, 数据接收体)
//todo 方式2. conn.Read(数据接收体)/conn.Write(数据接收体)

//方式1
func Handler1(conn *websocket.Conn) {
	//打印客户端和服务端地址
	fmt.Printf("a new socket connection from %s -> %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())

	//socket读写频率由代码自行控制。这里用了死循环，不断等待接受来自client端的数据。如果不设置循环，只数据传输一次后，程序就在不断地监听，等待下一socket传输数据
	//一次send的数据发往socket,另一端一次receive最多只能接受一次send的数据（方式2同理）。如:
	//send（data:10bytes）
	//rec:=make([]byte,5)
	//receive(conn,&rec)
	//因为一次receive的长度只有5bytes,而一次send的数据有10bytes,此时，receive一次只能接收5bytes的数据，仍有5bytes在socket中，等待下一次被receive。
	//如果receive的数据接受体长度超过send一次的长度，如15bytes,那receive一次最多也只能接收一次send的数据。
	//也就是说就算此时socket中已经有2次send的数据，一次receive也最多只能读取一次send的数据，剩余数据等待下次receive。
	for {
		time.Sleep(3*time.Second)
		var data string
		//当client向socket写数据时，server在此接收数据，从conn中接收数据到reply中。
		//如果socket中没有数据，或数据receive不出来，程序会阻塞在此。
		err := websocket.Message.Receive(conn, &data) //注意：string类型的接收体需要引用的形式："&"取地址
		if err != nil {
			fmt.Println("receive err:", err.Error())
			break
		}
		fmt.Println("Received from client: " + data)
		//向client端写数据
		if err = websocket.Message.Send(conn, data); err != nil {
			fmt.Println("send err:", err.Error())
			break
		}
	}
}

//方式2
func Handler2(conn *websocket.Conn) {
	//只读写一次
	var data []byte
	//todo 有坑！
	//当声明"var data []byte"时，必须为data分配内存（make）！不分配无法接受数据，会报错。
	//注意分配的长度！最好不要小于一次write()的数据长度，否则一次read的数据少于一次write的数据，当数据用于json解析时，容易发生错误。
	data = make([]byte, 32)
	//读取client写入的数据。如果socket中没有数据，或server在此read不出数据，程序会阻塞在此。
	//data如果没有分配内存（make操作），那么执行Read(data)时，是无法从socket中读出数据的，因为只声明的data相当于一个空指针，无法存储数据。
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println("read err:", err.Error())
	}
	fmt.Println("read from client: length=", n, " data=", data)
	//向client端写数据
	if _, err = conn.Write(data); err != nil {
		fmt.Println("write err:", err.Error())
	}

	//设置循环读写
	for {
		//等待2秒钟
		time.Sleep(2 * time.Second)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println("read err:", err.Error())
		}
		fmt.Println("read from client: length=", n, " data=", string(data[:]))
		if _, err = conn.Write(data); err != nil {
			fmt.Println("write err:", err.Error())
		}
	}
}

func StartServer1(){
	//设置访问api
	http.Handle("/handler1", websocket.Handler(Handler1))
	http.Handle("/handler2", websocket.Handler(Handler2))
	fmt.Println("begin to listen")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe error:", err)
		return
	}
}

func main(){
	//注意：
	//启动时，先启动server端，后启client端，否则client端会报connect连接错误
	//当server和client端任意一端关闭连接时，此时无论socket中是否有数据，未关闭的一端会返回error:EOF
	StartServer1()
}

