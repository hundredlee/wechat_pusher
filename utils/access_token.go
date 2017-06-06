package utils

import (
	"YibaiPusher/config"
	"YibaiPusher/models"
	"YibaiPusher/statics"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	conf     *config.Config = config.Instance()
	instance *AccessToken
)

type IAccessToken interface {
	GetToken() string
	Refresh()
	Expired() bool
	TTL() int
}

type AccessToken struct {
	token      models.Token
	createTime time.Time
	lifeCycle  int16
}

func AccessTokenInstance(force bool) *AccessToken {

	if force || instance == nil {
		accessToken := requestAccessToken()
		instance = accessToken
	}

	return instance
}

func (self *AccessToken) GetToken() string {
	return self.token.AccessToken
}

func (self *AccessToken) Refresh() {
	AccessTokenInstance(true)
}

func (self *AccessToken) Expired() bool {

	dur := int64(self.lifeCycle) - time.Now().Sub(self.createTime).Nanoseconds()
	if dur > 0 {
		return false
	}
	return true

}

func (self *AccessToken) TTL() int64 {
	dur := int64(self.lifeCycle) - time.Now().Sub(self.createTime).Nanoseconds()
	if dur < 0 {
		return -1
	}
	return dur
}

func requestAccessToken() *AccessToken {

	url := fmt.Sprintf(
		statics.WECHAT_GET_ACCESS_TOKEN,
		conf.ConMap["WeChat.APPID"], conf.ConMap["WeChat.SECRET"])

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	if errCode, err := jsonparser.GetInt(body, "errcode"); err == nil && errCode != 0 {
		if errMsg, _, _, err := jsonparser.Get(body, "errmsg"); err == nil {
			panic(errors.New(string(errMsg)))
		}
	}

	var token models.Token

	if err := json.Unmarshal(body, &token); err != nil {
		panic(err)
	}

	return &AccessToken{token: token, createTime: time.Now(), lifeCycle: token.ExpiresIn}
}
