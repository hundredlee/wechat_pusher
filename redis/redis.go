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
	"os"
	"github.com/hundredlee/wechat_pusher/hlog"
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
	poolsize := conf.ConMap["Redis.POOL_SIZE"]
	timeout := conf.ConMap["Redis.TIMEOUT"]
	if host == nil {
		panic("Redis Config not complete")
	}

	redisPool = &redis.Pool{
		MaxActive:   int(poolsize.(int)),
		IdleTimeout: time.Duration(timeout.(int)) * time.Second,
		Dial: func() (redis.Conn, error) {

			var c redis.Conn
			var err error
			if pass != nil && db != nil {
				c, err = redis.Dial("tpc", host.(string), redis.DialPassword(pass.(string)), redis.DialDatabase(db.(int)))
			} else if db != nil {
				c, err = redis.Dial("tpc", host.(string), redis.DialDatabase(db.(int)))
			} else if pass != nil {
				c, err = redis.Dial("tpc", host.(string), redis.DialPassword(pass.(string)))
			} else {
				c, err = redis.Dial("tcp", host.(string))
			}

			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	conn := redisPool.Get()
	defer conn.Close()

	r,err := redis.String(conn.Do("PING","test"))
	if err != nil {
		panic(err)
	}

	if r !=  "test"{
		panic("redis connect failed")
		os.Exit(-1)
	}else{
		hlog.LogInstance().LogInfo("redis connect success")
	}

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
