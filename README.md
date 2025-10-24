# MikaPluginLib

为MikaPanel框架编写插件时，建议使用这个库和框架主体进行通信。

使用这个框架的bot的主要功能都以插件的形式实现。

框架仅是对前端框架（NapCatQQ）的简单封装。

## 使用

搭建好go环境并创建好项目后，

在项目目录下运行以下指令，即可导入仓库

~~~bash
go mod init
go get github.com/MikaBot-Project/MikaPluginLib@latest
~~~

## pluginIO

将以下代码输入文件中即可引入该模块

```go
import "github.com/MikaBot-Project/MikaPluginLib/pluginIO"
```

> [!WARNING]
>
> 建议在main函数结尾添加以下代码
>
> ```go
> var data pluginIO.Message
> for {
> 	data = <-pluginIO.MessageChan
> 	log.Println("type: ", data.PostType)
> }
> ```
>
> 以保证插件主线程持续运行

### 注册

pluginIO模块使用回调注册模式

在程序开始运行时，需要注册相应的回调处理函数

```go
//使用示例
pluginIO.MessageRegister(func(msg pluginIO.Message){})
pluginIO.NoticeRegister("notify", func(msg pluginIO.Message){})
pluginIO.CommandRegister("!test", func(msg pluginIO.Message){})
```

函数`pluginIO.MessageRegister(callback func(Message))`将插件注册到message事件

当框架接收到消息时，会将数据传递到插件，插件接收到数据后，会调用注册时传递的函数

可以多次调用注册多个函数



函数`func NoticeRegister(noticeType string, handler func(Message))`将插件注册到Notice事件

当框架收到通知时，会将数据传递到插件，插件收到数据后会根据notice type调用相应的函数

同一notice type可以注册多个回调函数

> [!NOTE]
>
> notice type取值可以参考：[通知事件 (Notice Event) | NapCatQQ](https://napneko.github.io/onebot/event#通知事件-notice-event)
>
> 特别的：当账号是群聊管理员时，可以接收群聊加入请求，这个事件的notice type为`"group_add"`



函数`CommandRegister(cmdName string, handler func(Message))`将插件注册到指令消息

当框架接收到消息并检测到已注册的指令时，会将数据传递到插件

> [!CAUTION]
>
> 注意：当接收到指令后，数据不会传递到消息插件，且一个指令只能够存在一个插件，插件同一时刻只能注册一个回调函数。

### 回调函数

所有注册函数所接收的回调函数仅接收一个参数`msg pluginIO.Message`

其定义如下：

```go
type Message struct {
	Time          int64         `json:"time"`
	SelfId        int64         `json:"self_id"`
	PostType      string        `json:"post_type"`
	UserId        int64         `json:"user_id"`
	GroupId       int64         `json:"group_id"`
	GroupName     string        `json:"group_name"`
	Font          int           `json:"font"`
	RealSeq       int           `json:"real_seq"`
	MessageSeq    int           `json:"message_seq"`
	MessageType   string        `json:"message_type"`
	SubType       string        `json:"sub_type"`
	MessageId     int64         `json:"message_id"`
	MessageArray  []MessageItem `json:"message"`
	MessageFormat string        `json:"message_format"`
	RawMessage    string        `json:"raw_message"`
	Sender        struct {
		UserId   int64  `json:"user_id"`
		NickName string `json:"nickname"`
		Sex      string `json:"sex"`
		GroupId  int64  `json:"group_id"`
		Card     string `json:"card"`
		Role     string `json:"role"`
	} `json:"sender"`
	NoticeType    string   `json:"notice_type"`
	TargetId      int64    `json:"target_id"`
	TempSource    string   `json:"temp_source"`
	MetaEventType string   `json:"meta_event_type"`
	RequestType   string   `json:"request_type"`
	Comment       string   `json:"comment"`
	Flag          string   `json:"flag"`
	AtMe          bool     `json:"at_me"`
	CommandArgs   []string `json:"command_args"`
}
```

> [!NOTE]
>
> 其值具体可以参考[事件基础结构 | NapCatQQ](https://napneko.github.io/onebot/basic_event)

下面是回调函数的一个实现实例

```go
//具体可以参考 /test.go 文件
func test(msg pluginIO.Message) {
	log.Println("post type:", msg.PostType)
	switch msg.PostType {
	case "notice":
		log.Println("notice subtype", msg.SubType)
	case "message":
		log.Println("get message", msg.UserId)
	case "command":
		log.Println("get command", msg.CommandArgs)
	}
}
```

`pluginIO.Message`结构体中的`MessageArray`为`[]MessageItem`类型

`MessageItem`结构具体定义如下：

```go
type MessageItem struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}
//将data中对应字段数据转换为string类型
func (item *MessageItem) GetString(name string) string
//将data中对应字段数据转换为int类型
func (item *MessageItem) GetNumber(name string) int
//将数据存储到对应字段中
func (item *MessageItem) Set(name string, value any)
```

> [!NOTE]
>
> 具体定义和取值可以参考[消息段类型详解文档 | NapCatQQ](https://napneko.github.io/onebot/sement)

### 发送数据

发送数据函数定义如下

```go
SendMessage[T string | []MessageItem](msg T, userId int64, groupId int64)[]int
SendApi(apiName string, data []byte) []byte
SendPoke(userId int64, groupId int64) //发送戳一戳
```

`SendMessage`函数可以同时接收CQcode格式字符串和`[]MessageItem`数据

当`groupId`等于`0`时，将发送私聊消息到`userId`对应用户

当`groupId`不等于`0`时，将发送群聊消息到`groupId`对应群组



SendApi函数可以调用框架前端（NapCatQQ）的api

> [!NOTE]
>
> api的名称、参数和返回数据可以参考[NapCat 接口文档 - NapCat](https://napcat.apifox.cn/5430207m0)

## pluginConfig

此模块是插件的配置管理模块

将以下代码输入文件中即可引入该模块

```go
import "github.com/MikaBot-Project/MikaPluginLib/pluginConfig"
```


读取配置文件函数如下：

```
func ReadJson(fileName string, config any)
func ReadAllJson[T any](path string, config *map[string]T)
```

保存配置的函数如下

```
func SaveJson(fileName string, config any)
```

通过这个模块读取的配置，可以随时进行更新而不需要重启

## pluginData

此模块专用于保存插件运行时的数据，以实现持久化

将以下代码输入文件中即可引入该模块

```go
import "github.com/MikaBot-Project/MikaPluginLib/pluginData"
```


读取保存的数据

```
func ReadBinary[T any](name string, data *T)
func ReadJson[T any](name string, data *T)
```

保存数据

```
func SaveBinary(name string, data any)
func SaveJson(name string, data any)
```

`name`是文件路径，支持多级文件夹
