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
	"time"
)

var accessToken = AccessTokenInstance(false)

type Push struct {
	bufferNum int
	retries   int
	tasks     []models.Task
}

func NewPush(tasks []models.Task) *Push {
	return &Push{tasks: tasks}
}

func (self *Push) SetRetries(retries int) *Push {
	self.retries = retries
	return self
}

func (self *Push) SetBufferNum(bufferNum int) *Push {
	self.bufferNum = bufferNum
	return self
}

func (self *Push) Add(schedule string) {

	getCronInstance().AddFunc(schedule, func() {
		if self.retries == 0 || self.bufferNum == 0 {
			panic("Please SetRetries or SetBufferNum")
		}

		var resourceChannel = make(chan bool, self.bufferNum)

		for _, task := range self.tasks {

			resourceChannel <- true

			go func(task models.Task) {

				retr := 0

				defer func() {
					if recover() != nil {
						log.Printf("error-log ++++ TaskInfo : %v ++++\n", task)
					}
				}()

				r, _ := json.Marshal(task.Message)
				url := fmt.Sprintf(statics.WECHAT_TEMPLATE_SEND, accessToken.GetToken())

			LABEL:
				resp, _ := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(r))

				body, _ := ioutil.ReadAll(resp.Body)
				errCode, _ := jsonparser.GetInt(body, "errcode")

				if errCode != 0 {
					if retr >= self.retries {
						log.Printf("error-log ++++ TaskInfo : %v -- ErrorCode : %d ++++\n", task, errCode)
					} else {

						if errCode == 40001 {
							accessToken.Refresh()
						}

						log.Printf("retry times : %d", retr)
						time.Sleep(3 * time.Second)
						retr++
						goto LABEL
					}
				}

				<-resourceChannel

			}(task)

		}
	})

}
