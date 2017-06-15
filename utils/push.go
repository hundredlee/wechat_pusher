package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/hundredlee/wechat_pusher/hlog"
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
	BufferNum int
	Retries   int
	Tasks     []task.Task
	TaskType  string
}

func NewPush(push *Push) *Push {
	return push
}

func (self *Push) SetRetries(retries int) *Push {
	self.Retries = retries
	return self
}

func (self *Push) SetBufferNum(bufferNum int) *Push {
	self.BufferNum = bufferNum
	return self
}

func (self *Push) SetTaskType (taskType string) *Push{
	self.TaskType = taskType
	return self
}

func (self *Push) Add(schedule string) {


	//if tasks len equal 0 || the first object is not right taskType panic
	if len(self.Tasks) <= 0{
		panic("task is not allow empty")
	}

	if self.Retries == 0 || self.BufferNum == 0 {
		panic("Please SetRetries or SetBufferNum")
	}

	if self.TaskType == ""{
		panic("Please Set TaskType")
	}

	firstTask := self.Tasks[0]
	switch self.TaskType {
	case enum.TASK_TYPE_TEMPLATE:
		if _,ok := firstTask.(*task.TemplateTask); !ok {
			panic("not allow other TaskType struct in this TaskType")
		}
	//TODO other taskType
	}

	getCronInstance().AddFunc(schedule, func() {

		fileLog.LogInfo("Start schedule " + schedule + " TaskNumber:" + strconv.Itoa(len(self.Tasks)) + " TaskType:" + self.TaskType)

		var resourceChannel = make(chan bool, self.BufferNum)

		for _, task := range self.Tasks {

			resourceChannel <- true

			go run(task,self.Retries,resourceChannel,self.TaskType)

		}
	})

}

func run(task task.Task,retries int,resourceChannel chan bool,taskType string) {
	retr := 0

	defer func() {
		if recover() != nil {
			fileLog.LogError(fmt.Sprintf("task : %v", task))
		}
	}()

	r, _ := json.Marshal(task.GetTask())
	url := fmt.Sprintf(enum.URL_MAP[taskType], accessToken.GetToken())

LABEL:
	resp, _ := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(r))

	body, _ := ioutil.ReadAll(resp.Body)
	errCode, _ := jsonparser.GetInt(body, "errcode")

	if errCode != 0 {
		if retr >= retries {
			fileLog.LogError(fmt.Sprintf("TaskInfo : %v -- ErrorCode : %d -- TryTimeOut : %d", task, errCode, retr))
		} else {

			if errCode == 40001 || accessToken.Expired(){
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
