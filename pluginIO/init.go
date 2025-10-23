package pluginIO

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type MessageItem struct {
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

func (item *MessageItem) GetString(name string) string {
	if val, ok := item.Data[name]; ok {
		return val.(string)
	}
	return ""
}

func (item *MessageItem) GetNumber(name string) int {
	if val, ok := item.Data[name]; ok {
		return val.(int)
	}
	return 0
}

func (item *MessageItem) Set(name string, value any) {
	if item.Data == nil {
		item.Data = make(map[string]any)
	}
	item.Data[name] = value
}

type Message struct {
	Time          int64         `json:"time"`
	SelfId        int64         `json:"self_id"`
	PostType      string        `json:"post_type"`
	UserId        int64         `json:"user_id"`
	GroupId       int64         `json:"group_id"`
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

var inReader *bufio.Reader
var outWriter *bufio.Writer
var sendRecvMap *SafeMap
var MessageChan chan Message
var commandMap map[string]func(Message)
var noticeMap map[string][]func(Message)
var messageMap []func(Message)
var OperatorMap map[string]func(Message)

func init() {
	inReader = bufio.NewReader(os.Stdin)
	outWriter = bufio.NewWriter(os.Stdout)
	MessageChan = make(chan Message)
	sendRecvMap = NewSafeMap()
	commandMap = make(map[string]func(Message))
	noticeMap = make(map[string][]func(Message))
	OperatorMap = make(map[string]func(Message))
	go func() {
		var msg Message
		var data []byte
		var err error
		for {
			data, err = inReader.ReadSlice('\n')
			if len(data) == 0 {
				continue
			}
			if err != nil {
				log.Println("read error:", err)
				return
			}
			msg = Message{}
			err = json.Unmarshal(data, &msg)
			if err != nil {
				log.Println("unmarshal error:", err)
				return
			}
			switch msg.PostType {
			case "operator":
				handler, ok := OperatorMap[msg.MessageType]
				if ok {
					handler(msg)
				}
			case "return":
				sendRecvMap.Set(msg.MessageType, msg)
			case "message":
				for _, callback := range messageMap {
					go callback(msg)
				}
			case "command":
				go commandMap[msg.CommandArgs[0]](msg)
			case "notice":
				for _, callback := range noticeMap[msg.NoticeType] {
					go callback(msg)
				}
			default:
				MessageChan <- msg
			}
		}
	}()
}

func sendData(cmd string, args ...string) {
	_, err := outWriter.Write([]byte(cmd))
	if err != nil {
		log.Println("write error:", err)
		return
	}
	for _, arg := range args {
		_, err = outWriter.Write([]byte(":##:"))
		if err != nil {
			log.Println("write error:", err)
			return
		}
		_, err = outWriter.Write([]byte(arg))
		if err != nil {
			log.Println("write error:", err)
			return
		}
	}
	_, err = outWriter.Write([]byte("\n"))
	if err != nil {
		log.Println("write error:", err)
		return
	}
	err = outWriter.Flush()
	if err != nil {
		log.Println("flush error:", err)
		return
	}
}
