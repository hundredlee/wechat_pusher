package task

type Task interface {
	GetTaskType() string
	SetTask(interface{})
	GetTask() interface{}
}
