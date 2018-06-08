package main

import (
	"time"
	"strconv"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"reflect"
	"Golang-tutorial/redigo/redis"
)

const CH_RESOURCE = "K8SRESOURCE"
const CH_CHECK = "K8SCHECK"

func main() {
	myredis.NewRedis()
	//pub()需要go routine 起，否则循环回不来，sub没法执行
	go pub()
	go sub(CH_CHECK)
	go sub(CH_CHECK)
	sub(CH_CHECK)
}

//每次publish都会申请连接并释放连接。当连接关闭后，发布的消息是否还能被订阅？
func pub() {
	i := 0
	//ticker := time.NewTicker(5 * time.Second)
	//for _ = range ticker.C {
		i++
		//get k8s resource
		data := []byte("k8s resource pod" + strconv.Itoa(i))
		redis := myredis.Redis
		ok, err := redis.Publish(CH_RESOURCE, data)
		if !ok {
			fmt.Println("publish failed. err=", err)
			//continue
		}
		fmt.Println("publish success. channel=", CH_RESOURCE, " data=", string(data))

		//------
		data = []byte("k8s resource svc" + strconv.Itoa(i))
		ok, err = redis.Publish(CH_RESOURCE, data)
		if !ok {
			fmt.Println("publish failed. err=", err)
			//continue
		}
		fmt.Println("publish success. channel=", CH_RESOURCE, " data=", string(data))

		//------
		data = []byte("k8s resource node" + strconv.Itoa(i))
		ok, err = redis.Publish(CH_RESOURCE, data)
		if !ok {
			fmt.Println("publish failed. err=", err)
			//continue
		}
		fmt.Println("publish success. channel=", CH_RESOURCE, " data=", string(data))
	//}
}

func sub(channels ...string) {
	//这个获取连接不能放在for循环里，不然会瞬间试着从连接池中获取连接，
	//会超过连接池中的最大空闲连接数，请求瞬间多了易出问题
	ok, subConn := myredis.Redis.GetSubConn()
	if !ok {
		fmt.Println("error：GetSubConn failed")
	}
	//todo subConn仍未关闭

	//订阅也不要放在for循环中，订阅只订阅一次即可，receive在for循环中就行
	//至于订阅失败的处理情况再说
	ch := redis.Args{}.AddFlat(channels)
	if err := subConn.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		fmt.Println("error：订阅频道失败。channel=", ch, channels)
		return
	}
	fmt.Println("订阅频道 success. channel=", channels)

	//go routine不要放在for循环中，receive在一个go routine中循环执行即可。
	//go routine放for中，会瞬间起n多个go routine，每个都Receive，会导致redis运行时错误bufio数组越界，
	//且起多个go routine影响系统性能，系统很快会崩的。
	go func() {
		for {
			//订阅频道时，最开始会接收到redis.Subscription类型，订阅多少个就会收到多少个，
			//如果频道中无数据传输，且无错误，会阻塞。直到可以receive。
			data := subConn.Receive()
			fmt.Println("receive type:", reflect.TypeOf(data))
			switch t := data.(type) {
			case error:
				fmt.Println("type: error. err=", t)
			case redis.Message:
				fmt.Println("type: message.")
				fmt.Println("receive: ", string(t.Data))

			case redis.Subscription:
				fmt.Println("type: subcription. len(channel)=", t.Count)
			}
		}
	}()

	for {
		time.Sleep(time.Second * 2)
	}

}
