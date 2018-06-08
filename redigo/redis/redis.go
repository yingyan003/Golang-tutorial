package myredis

import (
"github.com/garyburd/redigo/redis"
"sync"
"fmt"
	"log"
)

type redisClient struct {
	pool    *redis.Pool
	subConn *redis.PubSubConn
	//pubConn *redis.Conn
}

var Redis *redisClient
var once sync.Once

func NewRedis() {
	once.Do(func() {
		fmt.Println("-------enter once.do newRedis")
		Redis=new(redisClient)
		Redis.pool = newPool()
	})
}

func newPool() *redis.Pool {
	host := LoadEnvVar(ENV_REDIS_HOST, REDIS_HOST)
	pool := &redis.Pool{
		//todo 这里的连接数貌似并不等于poo.Get()
		MaxIdle:     MAX_IDLE,
		MaxActive:   MAX_ACTIVE,
		IdleTimeout: IDLE_TIMEOUT,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				log.Fatalf("redisClient Dial host failed: host=%s, err=%v", host, err)
				return nil, err
			}
			//log.Infof("NewRedis newPool Success: host=%s", host)
			return c, nil
		},
	}
	return pool
}

func (r *redisClient) Get(key string) string {
	//从连接池pool中获取一个可用的空闲连接，实际执行的是redis.Pool{Dial: func}中的func函数
	conn := r.pool.Get()
	if conn == nil {
		log.Fatalf("redisClient Get failed: the connection get from pool is nil.")
		return ""
	}
	defer conn.Close()

	ok, err := r.Exists(key)
	if err != nil {
		log.Fatalf("redisClient Get failed: key=%s Exists error. err=%s", key, err)
		return ""
	}

	if !ok {
		log.Printf("redisClient Get: key=%s Not exist", key)
		return ""
	}

	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Fatalf("redisClient Get failed: Do error. err=%s", err)
		return ""
	}

	return result
}

func (r *redisClient) Set(key, value string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Fatalf("redisClient Set failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		log.Fatalf("redisClient Set failed: Do error. err=%s", err)
		return false, err
	}

	return true, nil
}

func (r *redisClient) SetWithExpire(key, value, expire string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Fatalf("redisClient SetWithExpire failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		log.Fatalf("redisClient SetWithExpire failed: Do Set error. error=%s", err)
		return false, err
	}

	_, err = conn.Do("EXPIRE", key, expire)
	if err != nil {
		log.Fatalf("redisClient SetWithExpire failed: Do EXPIRE error. error=%s", err)
		return false, err
	}

	return true, nil
}

func (r *redisClient) Delete(key string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Fatalf("redisClient Delete failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		log.Fatalf("redisClient Delete failed: Do error. err=%s", err)
		return false, err
	}

	return true, nil
}

func (r *redisClient) Exists(key string) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Fatalf("redisClient Exist failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	result, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		log.Fatalf("redisClient Exist failed: Do error. err=%s", err)
		return false, err
	}

	if result == 0 {
		return false, nil
	}

	return true, nil
}

func (r *redisClient) Publish(channel string, message []byte) (bool, error) {
	conn := r.pool.Get()
	if conn == nil {
		log.Fatalf("redisClient Publish failed: the connection get from pool is nil.")
		return false, nil
	}
	defer conn.Close()

	_, err := conn.Do("PUBLISH", channel, message)
	if err != nil {
		log.Fatalf("redisClient Publish failed: Do error. err=%s", err)
		return false, err
	}

	//fmt.Println("Publish:  defer close connection" )
	return true, nil
}

func (r *redisClient) GetSubConn() (bool, *redis.PubSubConn) {
	//conn := r.pool.Get()
	//if conn == nil {
	//	log.Fatalf("redisClient Subcribe failed: the connection get from pool is nil.")
	//	return false, nil
	//}

	//defer conn.Close()

	Redis.subConn = &redis.PubSubConn{Conn: r.pool.Get()}

	return true, Redis.subConn
}
