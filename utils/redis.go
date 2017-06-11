/**
 * @Description redis pool
 * @Author HundredLee
 * @Email hundred9411@gmail.com
 */
package utils

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

var redisPool *redis.Pool

func NewRedis() *redis.Pool{
	host := conf.ConMap["Redis.HOST"]
	if host == nil {
		log.Println("Redis Not Config")
	}
}
