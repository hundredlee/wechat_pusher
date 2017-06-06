package utils

import (
	"github.com/hundredlee/wechat_pusher/models"
	"fmt"
	"encoding/json"
	"github.com/hundredlee/wechat_pusher/statics"
	"net/http"
	"bytes"
)

var accessToken = AccessTokenInstance(true)

type Push struct {
	tasks   []models.Task
}

func NewPush(tasks []models.Task) *Push {
	return &Push{tasks: tasks}
}

func (self *Push) Run () {

	var resourceChannel = make(chan bool,len(self.tasks) /2)

	for _,task:=range self.tasks{

		resourceChannel <- true

		go func(task models.Task) {

			fmt.Println(task.Message.Data.First.Value)
			r, _ := json.Marshal(task.Message)
			url := fmt.Sprintf(statics.WECHAT_TEMPLATE_SEND, accessToken.GetToken())
			http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(r))
			<- resourceChannel

		}(task)

	}

}
