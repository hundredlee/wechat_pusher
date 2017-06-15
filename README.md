# wechat_pusher

## 功能列表
- 消息推送
    - 模板消息推送
    	- model -> message.go
    	- task -> template_task.go
    - 图片推送(TODO)
    - 文字推送(TODO)
    - 图文推送(TODO)
- 日志存储
- 计划任务

## 如何开始？
### 第一步:当然是go get
- `go get github.com/hundredlee/wechat_pusher.git`

- 项目结构如下：

```
├── README.md
├── config
│   └── config.go
├── config.conf
├── config.conf.example
├── enum
│   └── task_type.go
├── glide.lock
├── glide.yaml
├── hlog
│   ├── filelog.go
│   ├── filelog_test.go
│   └── hlog.go
├── main.go
├── main.go.example
├── models
│   ├── message.go
│   └── token.go
├── redis
│   ├── redis.go
│   └── redis_test.go
├── statics
│   └── global.go
├── task
│   ├── task.go
│   └── template_task.go
├── utils
│   ├── access_token.go
│   ├── crontab.go
│   └── push.go
└── vendor
    └── github.com

```

### 第二步:创建一个项目

#### 创建配置文件
- 项目根目录有一个config.conf.example，重命名为config.conf即可
- 内容如下：

```
[WeChat]
APPID=
SECRET=
TOKEN=

[Redis]
POOL_SIZE=
TIMEOUT=
HOST=
PASS=
DB=

[Log]
LOG_PATH=

```

- WeChat部分
	- APPID && SECRET && TOKEN  这些是微信开发者必须了解的东西。不细讲
- Redis部分
	- POOL_SIZE 连接池大小 ，整型 int
	- TIMEOUT 连接超时时间 ，整型 int
	- HOST  连接的IP 字符串 string
	- PASS   密码 字符串 string
	- DB    数据库选择 整型 int
- Log部分
	- LOG_PATH  日志存放文件夹，例如值为wechat_log，那么完整的目录应该是 GOPATH/wechat_log
	
- 调用的时候这么写:

```Go

conf := config.Instance()
//例如wechat 的 appid
appId := conf.ConMap["WeChat.APPID"]

```


#### 模板怎么配置
- 以模板消息作为例子说明：
- message.go 是模板消息的结构
- template_task.go 是将一个模板消息封装成任务（template_task.go 是实现了接口task.go的）
```Go
mess := models.Message{
		ToUser:     "openid",
		TemplateId: "templateid",
		Url:        "url",
		Data: models.Data{
			First:   models.Raw{"xxx", "#173177"},
			Subject: models.Raw{"xxx", "#173177"},
			Sender:  models.Raw{"xxx", "#173177"},
			Remark:  models.Raw{"xxx", "#173177"}}}

//封装成一个任务，TemplateTask表示模板消息任务
task := task.TemplateTask{}
task.SetTask(mess)

```
- 以上代码是模板消息的配置，这个微信开发者应该都能看懂。


#### 如何创建一个任务

- 例如我们要创建一个模板消息定时推送任务
	- 第一步，封装任务
	- 第二步，添加任务，并设置任务类型、并发执行的个数、失败尝试次数等。
	- 第三步，启动任务
-  我们用示例代码演示整个完整的过程

```Go
package main

import (
	"github.com/hundredlee/wechat_pusher/enum"
	"github.com/hundredlee/wechat_pusher/models"
	"github.com/hundredlee/wechat_pusher/task"
	"github.com/hundredlee/wechat_pusher/utils"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var tasks []task.Task
	tasks = make([]task.Task, 100)
	mess := models.Message{
		ToUser:     "oBv9cuLU5zyI27CtzI4VhV6Xabms",
		TemplateId: "UXb6s5dahNC5Zt-xQIxbLJG1BdP8mP73LGLhNXl68J8",
		Url:        "http://baidu.com",
		Data: models.Data{
			First:   models.Raw{"xxx", "#173177"},
			Subject: models.Raw{"xxx", "#173177"},
			Sender:  models.Raw{"xxx", "#173177"},
			Remark:  models.Raw{"xxx", "#173177"}}}
	task := task.TemplateTask{}
	task.SetTask(mess)

	for i := 0; i < 100; i++ {
		tasks[i] = &task
	}

    utils.NewPush(&utils.Push{
    	Tasks:tasks,
    	TaskType:enum.TASK_TYPE_TEMPLATE,
    	Retries:4,
    	BufferNum:10,
    }).Add("45 * * * * *")

    utils.StartCron()

}

```


- Add方法里面填写的是执行的时间
    - ("10 * * * * *") 表示每分钟的第十秒钟执行一次。
    - ("@hourly") 每小时执行一次
- 具体请参照 https://github.com/robfig/cron/blob/master/doc.go
- 本推送服务的计划任务是由 https://github.com/robfig/cron 实现的。

### Run
- 很简单，当你组装好所有的task以后，直接运行一句话就可以了。

```Go
    utils.NewPush(&utils.Push{
    	Tasks:tasks,
    	TaskType:enum.TASK_TYPE_TEMPLATE,
    	Retries:4,
    	BufferNum:10,
    }).Add("45 * * * * *")

    utils.StartCron()

```

- `utils.StartCron()`

## Contributor
- HundredLee https://github.com/hundredlee
- Cospotato  https://github.com/cospotato


