# wechat_pusher

## 开源
> 项目已经在Github开源，还没有完全地完善，但是用是没问题的。
> 求各位大神的star啊，这是我的第一个比较完整的Go程序。(*^__^*)

- https://github.com/hundredlee/wechat_pusher

## 怎么用？
### 第一步当然是go get
- `go get github.com/hundredlee/wechat_pusher.git`
- 当然你也可以直接clone整个项目，直接导入IDE中跑一下试试

### 项目结构

```

├── README.md
├── config
│   └── config.go
├── config.conf.example
├── glide.lock
├── glide.yaml
├── main.go.example
├── models
│   ├── message.go
│   ├── task.go
│   └── token.go
├── statics
│   └── global.go
├── utils
│   ├── access_token.go
│   └── push.go
└── vendor

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

- 具体怎么填，我就不说了。这是接触过微信开发的童鞋都知道的东西。


### 模板配置怎么配置
- 我们看看models文件夹里面有message.go文件，里面其实就是模板的格式。
- 具体怎么用，看看main.go.example文件里面的示例。

```
package main

import (
	"fmt"
	"github.com/hundredlee/wechat_pusher/models"
	"github.com/hundredlee/wechat_pusher/utils"
	"runtime"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var tasks []models.Task
	tasks = make([]models.Task, 100)
	mess := models.Message{
		ToUser: "openid",
		TemplateId: "templateId",
		Url: "http://baidu.com",
		Data: models.Data{
			First: models.Raw{"xxx", "#173177"},
			Subject: models.Raw{"xxx", "#173177"},
			Sender: models.Raw{"xxx", "#173177"},
			Remark: models.Raw{"xxx", "#173177"}}}
	task := models.Task{Message: mess}
	for i := 0; i < 100; i++ {
		task.Message.Data.First.Value = fmt.Sprintf("%d", i)
		tasks[i] = task
	}

	utils.NewPush(tasks).SetRetries(4).SetBufferNum(10).Run()
}

```

### Run
- 很简单，当你组装好所有的task以后，直接运行一句话就可以了。
- `utils.NewPush(tasks).SetRetries(4).SetBufferNum(10).Run()`

### 打算？
- 目前还是比较简单的推送，然后日志相对来说比较完整。但是缺少了计划任务功能。大家可以star一下，等我更新计划任务的功能。

