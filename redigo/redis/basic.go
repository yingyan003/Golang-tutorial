package myredis

import "os"

//redis
const(
	ENV_REDIS_HOST = "REDIS_HOST"
	REDIS_HOST = "10.151.30.50:6379"
	MAX_IDLE = 1
	MAX_ACTIVE = 100
	IDLE_TIMEOUT = 0//180
)

func LoadEnvVar(key, value string)string{
	var v string
	if v=os.Getenv(key);v==""{
		v=value
	}
	return v
}
