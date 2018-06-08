package main

import (
	"context"
	"fmt"
	"time"
	"github.com/gomodule/redigo/redis"
	"reflect"
)


//官方例子

//官方文档：https://godoc.org/github.com/garyburd/redigo/redis#example-PubSubConn




// listenPubSubChannels listens for messages on Redis pubsub channels. The
// onStart function is called after the channels are subscribed. The onMessage
// function is called for each message.
func listenPubSubChannels(ctx context.Context, redisServerAddr string,
	onStart func() error,
	onMessage func(channel string, data []byte) error,
	channels ...string) error {

	fmt.Println("----------enter listenPubSubChannels")

	// A ping is set to the server with this period to test for the health of
	// the connection and server.
	const healthCheckPeriod = time.Minute

	c, err := redis.Dial("tcp", redisServerAddr,
		// Read timeout on server should be greater than ping period.
		redis.DialReadTimeout(healthCheckPeriod+10*time.Second),
		redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		return err
	}
	defer c.Close()

	psc := redis.PubSubConn{Conn: c}

	ch:=redis.Args{}.AddFlat(channels)
	fmt.Println("1. 订阅频道---------Subscribe channels: ",ch)
	if err := psc.Subscribe(redis.Args{}.AddFlat(channels)...); err != nil {
		return err
	}
	fmt.Println("完成订阅---------finish Subscribe channels")

	done := make(chan error, 1)

	// Start a goroutine to receive notifications from the server.
	go func() {
		fmt.Println("---------enter gorutine : 循环receive()")
		for {
			r:=psc.Receive()
			fmt.Println("--------receive type :",reflect.TypeOf(r))
			switch n := r.(type) {
			case error:
				fmt.Println("type: error. done <- n",n)
				done <- n
				return
			case redis.Message:
				fmt.Println("type: message. call onMessage()")
				if err := onMessage(n.Channel, n.Data); err != nil {
					fmt.Println("type: message. call onMessage() faled. done <- n",n)
					done <- err
					return
				}
			//	psc.Receive()返回的数据类型，如果不是错误error，一开始是redis.Subscription类型，
			//  当有消息发布到频道时，才是redis.Message类型。无消息发布时回到redis.Subscription类型。
			//  当退订频道时，redis.Subscription.Count记录的当前频道数会随着订阅的取消依次减相应的数量。
			case redis.Subscription:
				fmt.Println("type: Subscription. count：",n.Count)
				switch n.Count {
				case len(channels):
					fmt.Println("---当当前的频道数等于订阅的数量，开始发布. type: message. call onStart()")
					// Notify application when all channels are subscribed.
					if err := onStart(); err != nil {
						fmt.Println("type: message. call onStart() faled. done <- err",err)
						done <- err
						return
					}
				case 0:
					fmt.Println("当当前频道数=0时，退出receive循环。case 0：done <- nil")
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()
loop:
	for err == nil {
		select {
		case <-ticker.C:
			fmt.Println("test connect to server")
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				fmt.Println("test connect to server: ping failed")
				break loop
			}
		case d:=<-ctx.Done():
			fmt.Println("资源被释放后，退出循环。loop: ctx.Done has value. break loop. value=",d)
			break loop
		case err := <-done:
			fmt.Println("loop: ctx.Done has err. return err")
			// Return error from the receive goroutine.
			return err
		}
	}

	fmt.Println("3. 退订频道-----psc.Unsubscribe()")
	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	psc.Unsubscribe()

	d:=<-done
	fmt.Println("----- <- done: done=",d)
	// Wait for goroutine to complete.
	return d
}

func publish() {
	fmt.Println("2. 发布消息-------------enter publish")
	c, err := redis.Dial("tcp","10.151.30.50:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	c.Do("PUBLISH", "c1", "hello")
	c.Do("PUBLISH", "c2", "world")
	c.Do("PUBLISH", "c1", "goodbye")
	fmt.Println("----------------finish publish")
}

// This example shows how receive pubsub notifications with cancelation and
// health checks.
func main() {
	redisServerAddr:="10.151.30.50:6379"
	//该函数返回一个带有新Done通道的context，这个context是其父本的复制。
	//当返回的cancel函数被调用，或者其父本context的Done通道被关闭时，该context的Done通道关闭，无论前边提到的哪个先执行。
	//取消该context上下文就会释放与其关联的资源，所以当该context中运行的操作完成时，就应该调用cancel取消函数。
	ctx, cancel := context.WithCancel(context.Background())

	err := listenPubSubChannels(ctx,
		redisServerAddr,
		func() error {
			// The start callback is a good place to backfill missed
			// notifications. For the purpose of this example, a goroutine is
			// started to send notifications.
			go publish()
			return nil
		},
		func(channel string, message []byte) error {
			fmt.Printf("channel: %s, message: %s\n", channel, message)

			// For the purpose of this example, cancel the listener's context
			// after receiving last message sent by publish().
			if string(message) == "goodbye" {
				fmt.Println("当发布的消息为指定的结束消息时，调用ctx.cancel.释放相关资源。message == goodbye。 call cancel()")
				cancel()
			}
			return nil
		},
		"c1", "c2")

	if err != nil {
		fmt.Println(err)
		return
	}

}