# wechat_pusher

## 功能列表
- 微信模板消息推送
- 可添加计划任务

## 怎么用？
### 第一步当然是go get
- `go get github.com/hundredlee/wechat_pusher.git`
- 当然你也可以直接clone整个项目，直接导入IDE中跑一下试试

### 项目结构

```

├── README.md
├── config
│   └── config.go
├── config.conf
├── config.conf.example
├── glide.lock
├── glide.yaml
├── main.go
├── models
│   ├── message.go
│   ├── task.go
│   └── token.go
├── statics
│   └── global.go
├── utils
│   ├── access_token.go
│   ├── crontab.go
│   └── push.go
└── vendor
    └──

```

### 配置文件
-  我们可以看到根目录有一个config.conf.example，重命名为config.conf即可
- 内容如下：

```
[WeChat]
APPID=
SECRET=
TOKEN=
TEMPLATE=

```

- 具体怎么填，不细说。这是接触过微信开发的童鞋都知道的东西。

- 调用的时候这么写:

```Go

conf := config.Instance()
//例如wechat 的 appid
appId := conf.ConMap["WeChat.APPID"]

```


### 模板配置怎么配置
- 我们看看models文件夹里面有message.go文件，里面其实就是模板的格式。
- 具体怎么用，看看main.go.example文件里面的示例。

```Go
package main

import (
	"fmt"
	"runtime"
	"github.com/hundredlee/wechat_pusher/models"
	"github.com/hundredlee/wechat_pusher/utils"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var tasks []models.Task
	tasks = make([]models.Task, 100)
	mess := models.Message{
		ToUser:     "oBv9cuLU5zyI27CtzI4VhV6Xabms",
		TemplateId: "UXb6s5dahNC5Zt-xQIxbLJG1BdP8mP73LGLhNXl68J8",
		Url:        "http://baidu.com",
		Data: models.Data{
			First:   models.Raw{"xxx", "#173177"},
			Subject: models.Raw{"xxx", "#173177"},
			Sender:  models.Raw{"xxx", "#173177"},
			Remark:  models.Raw{"xxx", "#173177"}}}
	task := models.Task{Message: mess}
	for i := 0; i < 100; i++ {
		task.Message.Data.First.Value = fmt.Sprintf("%d", i)
		tasks[i] = task
	}

	utils.NewPush(tasks).SetRetries(4).SetBufferNum(10).Add("50 * * * * *")
	utils.NewPush(tasks).SetRetries(4).SetBufferNum(10).Add("10 * * * * *")
	utils.NewPush(tasks).SetRetries(4).SetBufferNum(10).Add("20 * * * * *")
	utils.StartCron()
}


```

### 定时任务

- `utils.NewPush(tasks).SetRetries(4).SetBufferNum(10).Add("10 * * * * *")`
- Add方法里面填写的是执行的时间
    - ("10 * * * * *") 表示每分钟的第十秒钟执行一次。
    - ("@hourly") 每小时执行一次
- 具体请参照 https://github.com/robfig/cron/blob/master/doc.go
- 本推送服务的计划任务是由 https://github.com/robfig/cron 实现的。

### Run
- 很简单，当你组装好所有的task以后，直接运行一句话就可以了。
- `utils.NewPush(tasks).SetRetries(4).SetBufferNum(10).Add("10 * * * * *")`
- `utils.StartCron()`
