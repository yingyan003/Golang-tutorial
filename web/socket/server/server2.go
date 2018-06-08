package main

import (
	"golang.org/x/net/websocket"
	"fmt"
	"encoding/json"
	"time"
	"net/http"
)

//前言：
//此程序用于模拟：前端动态列表展示
//因项目中列表实际是动态变化的，但如不用socket传输时，用户需要刷新页面，重新获取列表才可获取最新展示。
//而采用socket后，用户无需刷新页面即可看到动态的最新列表

//Key可以是用户搜索列表的关键词，Key仅用于测试json解析的效果
//结果显示，Key中声明Key2的指针，那么json.Unmarsh解析Key前，需实例化Key，否则json解析报空指针异常。
//但实例化后也无法成功解析uid到Key2.Uid。若希望一次性解析成功，Key中可直接内嵌Key2,而不是声明指针，即Key struct{ Key2 }

type Key struct {
	Key string
	*Key2
}

type Key2 struct {
	Uid int32
}

func Handler(conn *websocket.Conn) {
	k := new(Key)
	k.Key2 = new(Key2)
	//此参数不做实际用处，仅用于测试：socket断开连接有可能是因为程序出错导致，并不一定由超时引起。
	kch := make(chan *Key)

	buff := make([]byte, 8)
	_, err := conn.Read(buff)
	if err != nil {
		fmt.Println("read socket error: ", err)
	}
	fmt.Println("test buff", string(buff[:]))
	//time.Sleep(3 * time.Second)

	//todo 有坑! 超时时间慎重设置，如果太小，如几秒，则管道连接断开，获取值时会报"close connection"之类的错误，cleint端会报"broken pipe"的错误
	timer := time.NewTimer(5 * time.Minute)
	timeout := timer.C

	//在此实时获取用户在页面输入的"搜索关键词"，获取列表的具体逻辑可以在下一个for循环的default中操作，这里用简单的打印代替了
	go func(kch chan *Key) {
		for {
			fmt.Println("----enter first default")

			//todo 有坑! 长度设置太小，不足以容下一个完整的json内容，则json解析会报错。
			//client端发送： {"key":"1"}，长度是11，刚开始设置的8，最多容下{"key":"，
			//所有第一次读取到{"key":"，第二次读取到1"}，导致json解析失败。
			buff := make([]byte, 32)

			//todo 如果client write([]byte("")),写空串，那么socket中是没有数据传输的，server端会阻塞在此
			n, err := conn.Read(buff)
			if err != nil {
				fmt.Println("read socket error: ", err)
				//continue
			}
			fmt.Println("default 1: n=:", n, " buff", string(buff[:n]))
			time.Sleep(2 * time.Second)

			if n != 0 {
				param := buff[:n]
				err := json.Unmarshal(param, k)
				if err != nil {
					fmt.Println("json unmarshal error:", err)
				}
				fmt.Println("default 1 ；n: ", n, " k: ", k, " k.key: ", k.Key, " uid:", k.Uid)
				//kch <- k
				timer.Reset(5 * time.Minute)
			}
		}

	}(kch)

	for {
		//select一般与管道配合使用：
		//1.从上到下，找到第一个可执行的case（即该case对应的管道（<-）可以取出数据），执行它。后面就算有可执行的也忽略
		//2.如果没有可执行的case，就执行default.否则不会执行defalt
		//3.继续下一次循环
		//注意：
		// 1.select的循环，就算break也退出不了外层的for循环
		// 2.select的case不能是bool类型，管道一定支持，其他自测吧
		select {
		//当timer.Rest设置的超时时间已到，此期间内server没接受到任何来自client的数据，则断开socket连接
		case <-timeout:
			{
				conn.Write([]byte("timeout!"))
				fmt.Println("timeout")
				time.Sleep(3 * time.Second)
				conn.Close()
				fmt.Println("socket close")
			}
			//纯粹用于测试，可忽略
			/*case k = <-kch:
				{
					//todo 坑! socket断开连接有可能是因为程序出错导致，并不一定由超时引起。
					// 如以下例子所示，当设置了*Key2时，执行到下面的语句时，k.Uid会报空指针异常（因为Key2需要make），这时程序panic出错了！但是程序不会退出,可管道自动断开了！
					//type Key struct {
					//Key string
					//*Key2 todo
					//}
					//type Key2 struct {
					//	Uid int32
					//}

					fmt.Println("default 2; case k= <- kch: k=", k, " uid:", k.Uid)
				}*/
			//大致逻辑：可以在此实现获取列表的操作。如果前面的goroutine没有获取到前端返回的"搜索关键词"就返回全部的列表，否则返回匹配搜索关键词的列表。
		default:
			{
				fmt.Println("----enter second default")
				fmt.Println("default 2 k=", k, " k.key: ", k.Key, " uid:", k.Uid)
				time.Sleep(1 * time.Second)
			}
		}
	}

}

func StartServer2() {
	http.Handle("/server2", websocket.Handler(Handler))
	fmt.Println("server3 begin to listen")
	//127.0.0.1与localhost等同（前提：hosts文件已绑定）
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func main() {
	StartServer2()
}
