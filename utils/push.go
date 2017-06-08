package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/hundredlee/wechat_pusher/models"
	"github.com/hundredlee/wechat_pusher/statics"
	"io/ioutil"
	"log"
	"net/http"
)

var accessToken = AccessTokenInstance(true)

type Push struct {
	tasks []models.Task
}

func NewPush(tasks []models.Task) *Push {
	return &Push{tasks: tasks}
}

func (self *Push) Run(bufferNum int) {

	var resourceChannel = make(chan bool, bufferNum)

	for _, task := range self.tasks {

		resourceChannel <- true

		go func(task models.Task) {

			defer func() {
				if recover() != nil {
					log.Printf("error-log ++++ TaskInfo : %v ++++\n", task)
				}
			}()

			r, _ := json.Marshal(task.Message)
			url := fmt.Sprintf(statics.WECHAT_TEMPLATE_SEND, accessToken.GetToken()+"xxxx")
			resp, _ := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(r))

			body, _ := ioutil.ReadAll(resp.Body)

			errCode, _ := jsonparser.GetInt(body, "errcode")

			if errCode != 0 {
				log.Printf("error-log ++++ TaskInfo : %v -- ErrorCode : %d ++++\n", task, errCode)
			}

			<-resourceChannel

		}(task)

	}

}
