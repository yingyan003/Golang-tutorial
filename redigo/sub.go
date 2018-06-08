package main

import (
	"time"
	"strconv"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"Golang-tutorial/redigo/redis"
)

const CH_RESOURCE2 = "K8SRESOURCE"
const CH_CHECK2 = "K8SCHECK"

func main() {
	myredis.NewRedis()
	sub2(CH_RESOURCE2, "AA")
}

func sub2(channels ...string) {
	ok, subConn := myredis.Redis.GetSubConn()
	if !ok {
		fmt.Println("error：GetSubConn failed")
	}
	//todo subConn仍未关闭

	//这里的订阅与频道是否存在不关，只要提供频道名称，即可订阅。
	//当订阅的频道无数据传输时，会阻塞在下面的Receive()中
	ch := redis.Args{}.AddFlat(channels)
	if err := subConn.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		fmt.Println("error：订阅频道失败。channel=", ch, channels)
		return
	}
	fmt.Println("订阅频道 success. channel=", channels)

	go func() {
		i := 0
		for {
			//time.Sleep(time.Second*5)
			data := subConn.Receive()
			fmt.Println("receive type:", reflect.TypeOf(data))
			switch t := data.(type) {
			case error:
				fmt.Println("type: error. err=", t)
			//1. 假设频道中有数据发布，且发布后关闭连接。而连接关闭后才开始订阅，那么频道中之前发布的消息不会再被订阅。
			//2. 如果先订阅频道，频道中才发布消息，就算发布结束就关闭连接，消息也会被订阅到（连接关闭后消息一样会被订阅）
			//3. 基于1的情况下，但发布后连接不关闭，消息也一样订阅不到
			case redis.Message:
				fmt.Println("type: message.")
				fmt.Println("receive: ", string(t.Data))
				i++
				pub2(i,t.Data)
			case redis.Subscription:
				fmt.Println("type: subcription. len(channel)=", t.Count)//订阅的频道数，与频道是否有消息发布无关
			}
		}
	}()
	//time.Sleep(time.Second*1)
	for {
		time.Sleep(time.Second * 5)
	}

}

func pub2(i int, receive []byte) {
	//get k8s resource
	data := []byte("k8s check " + strconv.Itoa(i)+": "+string(receive))
	redis := myredis.Redis
	ok, err := redis.Publish(CH_CHECK2, data)
	if !ok {
		fmt.Println("publish failed. err=", err)
		return
	}
	fmt.Println("publish success. channel=", CH_CHECK2, " data=", string(data))
}
