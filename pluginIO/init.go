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

func (item *MessageItem) Get(name string) string {
	if val, ok := item.Data[name]; ok {
		return val.(string)
	}
	return ""
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
	RawMessage    string        `json:"raw_message"`
	NoticeType    string        `json:"notice_type"`
	TargetId      int64         `json:"target_id"`
	MetaEventType string        `json:"meta_event_type"`
	AtMe          bool          `json:"at_me"`
	CommandArgs   []string      `json:"command_args"`
}

var inReader *bufio.Reader
var outWriter *bufio.Writer
var sendRecvMap map[string]Message
var MessageChan chan Message
var commandMap map[string]func(Message)
var noticeMap map[string][]func(Message)
var messageMap []func(Message)

func init() {
	inReader = bufio.NewReader(os.Stdin)
	outWriter = bufio.NewWriter(os.Stdout)
	MessageChan = make(chan Message)
	sendRecvMap = make(map[string]Message)
	commandMap = make(map[string]func(Message))
	noticeMap = make(map[string][]func(Message))
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
			err = json.Unmarshal(data, &msg)
			if err != nil {
				log.Println("unmarshal error:", err)
				return
			}
			switch msg.PostType {
			case "return":
				sendRecvMap[msg.MessageType] = msg
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
	SendData("init", "v1")
}

func SendData(cmd string, args ...string) {
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
