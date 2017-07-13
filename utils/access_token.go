package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/hundredlee/wechat_pusher/config"
	"github.com/hundredlee/wechat_pusher/models"
	"github.com/hundredlee/wechat_pusher/redis"
	"github.com/hundredlee/wechat_pusher/statics"
	"io/ioutil"
	"net/http"
)

var (
	conf *config.Config = config.Instance()
)

func GetAccessToken() string {

	appId,ok := conf.ConMap["WeChat.APPID"].(string)

	if !ok {
		panic("appid not string")
	}

	cacheKey := "access_token." + appId

	accessTokenCache := redis.Get(cacheKey)

	if len(string(accessTokenCache)) > 0 && redis.TTL(cacheKey) > 0 {
		return string(accessTokenCache)
	}

	url := fmt.Sprintf(
		statics.WECHAT_GET_ACCESS_TOKEN,
		conf.ConMap["WeChat.APPID"], conf.ConMap["WeChat.SECRET"])

	response, _ := http.Get(url)

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println()

	if errCode, err := jsonparser.GetInt(body, "errcode"); err == nil && errCode != 0 {
		if errMsg, _, _, err := jsonparser.Get(body, "errmsg"); err != nil {
			panic(errors.New(string(errMsg)))
		}
	}

	var token models.Token

	if err := json.Unmarshal(body, &token); err != nil {
		panic(err)
	}

	redis.Set(cacheKey,token.AccessToken,false,token.ExpiresIn)

	return token.AccessToken
}
