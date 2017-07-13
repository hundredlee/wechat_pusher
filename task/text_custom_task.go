package task

import (
	"github.com/hundredlee/wechat_pusher/models"
	"github.com/hundredlee/wechat_pusher/enum"
)

type TextCustomTask struct {
	taskValue models.TextCustom
}

func (self *TextCustomTask) GetTaskType() string {
	return enum.TASK_TYPE_TEXT_CUSTOM
}

func (self *TextCustomTask) SetTask(taskValue interface{}) {

	v,ok := taskValue.(models.TextCustom)
	if ok {
		self.taskValue = v
	}
}

func (self *TextCustomTask) GetTask() interface {}{
	return self.taskValue
}
