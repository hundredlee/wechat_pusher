package redis_test

import (
	"testing"
	"github.com/hundredlee/wechat_pusher/redis"
	"fmt"
)

func TestNewRedis(t *testing.T) {

}

func TestTTL(t *testing.T) {
	fmt.Println(redis.TTL("xxx"))
}
