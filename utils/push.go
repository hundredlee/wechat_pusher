package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/hundredlee/wechat_pusher/hlog"
	"github.com/hundredlee/wechat_pusher/statics"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"github.com/hundredlee/wechat_pusher/task"
	"github.com/hundredlee/wechat_pusher/enum"
)

var accessToken = AccessTokenInstance(false)
var fileLog = hlog.LogInstance()

type Push struct {
	bufferNum int
	retries   int
	tasks     []task.Task
	taskType  string
}

func NewPush(tasks []task.Task) *Push {
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

func (self *Push) SetTaskType (taskType string) *Push{
	self.taskType = taskType
	return self
}

func (self *Push) Add(schedule string) {


	//if tasks len equal 0 || the first object is not right taskType panic
	if len(self.tasks) <= 0{
		panic("task is not allow empty")
	}

	firstTask := self.tasks[0]
	switch self.taskType {
	case enum.TASK_TYPE_IMAGE:
		if _,ok := firstTask.(*task.TemplateTask); !ok {
			panic("not allow other TaskType struct in this TaskType")
		}
	//TODO other taskType
	}

	getCronInstance().AddFunc(schedule, func() {
		if self.retries == 0 || self.bufferNum == 0 {
			panic("Please SetRetries or SetBufferNum")
		}

		if self.taskType == ""{
			panic("Please Set TaskType")
		}

		fileLog.LogInfo("Start schedule " + schedule + " TaskNumber:" + strconv.Itoa(len(self.tasks)) + " TaskType:" + self.taskType)

		var resourceChannel = make(chan bool, self.bufferNum)

		for _, task := range self.tasks {

			resourceChannel <- true

			go run(task,self.retries,resourceChannel)

		}
	})

}

func run(task task.Task,retries int,resourceChannel chan bool) {
	retr := 0

	defer func() {
		if recover() != nil {
			fileLog.LogError(fmt.Sprintf("task : %v", task))
		}
	}()

	r, _ := json.Marshal(task.GetTask())
	url := fmt.Sprintf(statics.WECHAT_TEMPLATE_SEND, accessToken.GetToken())

LABEL:
	resp, _ := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(r))

	body, _ := ioutil.ReadAll(resp.Body)
	errCode, _ := jsonparser.GetInt(body, "errcode")

	if errCode != 0 {
		if retr >= retries {
			fileLog.LogError(fmt.Sprintf("TaskInfo : %v -- ErrorCode : %d -- TryTimeOut : %d", task, errCode, retr))
		} else {

			if errCode == 40001 {
				fileLog.LogError("AccessToken expired and refresh")
				accessToken.Refresh()
			}

			time.Sleep(3 * time.Second)
			retr++
			goto LABEL
		}
	}else{
		fileLog.LogInfo(fmt.Sprintf("%v -- push success",task))
	}

	<-resourceChannel
}
