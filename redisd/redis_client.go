// Package core deal with setting up redis client, parsing incoming T38 commands,
// executing commands and returning the resp and error and handling joining of new
// setter + getter peer and beanstalkd/gearman services
package redisd

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	TILE38_REDIS_IDLE_TIMEOUT = time.Duration(5) * time.Second
	TILE38_REDIS_NETWORK      = "tcp"
)

// GetRedisConn returns a raw redis connection
func GetRedisConn(address string) (redis.Conn, error) {
	dialOption := redis.DialConnectTimeout(TILE38_REDIS_IDLE_TIMEOUT)
	return redis.Dial(TILE38_REDIS_NETWORK, address, dialOption)
}
