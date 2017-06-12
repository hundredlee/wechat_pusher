/**
 * @Description redis pool
 * @Author HundredLee
 * @Email hundred9411@gmail.com
 */
package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/hundredlee/wechat_pusher/config"
	"time"
)

var (
	redisPool *redis.Pool
	conf      *config.Config = config.Instance()
)

func instance() *redis.Pool {

	if redisPool != nil {
		return redisPool
	}

	host := conf.ConMap["Redis.HOST"]
	pass := conf.ConMap["Redis.PASS"]
	db := conf.ConMap["Redis.DB"]
	if host == nil {
		panic("Redis Config not complete")
	}

	redisPool = &redis.Pool{
		MaxActive:   int(conf.ConMap["Redis.POOL_SIZE"]),
		IdleTimeout: time.Duration(int(conf.ConMap["Redis.TIMEOUT"])) * time.Second,
		Dial: func() (conn redis.Conn, err error) {
			c, err := redis.Dial("tpc", string(host), redis.DialPassword(string(pass)), redis.DialDatabase(int(db)))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	conn := redisPool.Get()
	defer conn.Close()

	return redisPool

}

func Exists(key string) bool {
	conn := instance().Get()
	defer conn.Close()
	if conn == nil {
		panic("redis init faild")
	}

	n, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return n > 0
}
