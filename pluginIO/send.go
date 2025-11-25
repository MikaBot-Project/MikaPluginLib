package pluginIO

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var randomTime = int64(0)

func RandomString(length int) string {
	for randomTime == time.Now().UnixNano() {
	}
	randomTime = time.Now().UnixNano()
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

func SendMessage[T string | []MessageItem](msg T, userId int64, groupId int64, selfId int64) (messageId []int) {
	echo := fmt.Sprintf("send_msg_%s", RandomString(64))
	message := any(msg)
	send := struct {
		Action  string `json:"action"`
		UserId  int64  `json:"user_id"`
		GroupId int64  `json:"group_id"`
		SelfId  int64  `json:"self_id"`
		SubType string `json:"sub_type"`
		Data    []byte `json:"data"`
		Echo    []byte `json:"echo"`
	}{
		Action:  "send_msg",
		Echo:    []byte(echo),
		UserId:  userId,
		GroupId: groupId,
		SelfId:  selfId,
	}
	switch message.(type) {
	case string:
		send.SubType = "string"
		send.Data = []byte(message.(string))
	case []MessageItem:
		send.SubType = "array"
		send.Data, _ = json.Marshal(msg)
	}
	sendData(send)
	defer sendRecvMap.Delete(echo)
	var exists = false
	var recv interface{}
	for !exists {
		time.Sleep(100 * time.Millisecond)
		recv, exists = sendRecvMap.Get(echo)
	}
	_ = json.Unmarshal([]byte(recv.(Message).RawMessage), &messageId)
	return messageId
}

func SendApi(apiName string, data []byte, selfId int64) []byte {
	return SendApiEcho(apiName, data, fmt.Sprintf("send_api_%s", RandomString(64)), selfId)
}

func SendApiEcho(apiName string, data []byte, echo any, selfId int64) []byte {
	var echoData []byte
	switch echo.(type) {
	case []byte:
		echoData = echo.([]byte)
	default:
		echoData, _ = json.Marshal(echo)
	}
	sendData(struct {
		Action  string `json:"action"`
		ApiName string `json:"api_name"`
		SelfId  int64  `json:"self_id"`
		Data    []byte `json:"data"`
		Echo    []byte `json:"echo"`
	}{
		Action:  "send_api",
		ApiName: apiName,
		Data:    data,
		Echo:    echoData,
		SelfId:  selfId,
	})
	defer sendRecvMap.Delete(string(echoData))
	var exists = false
	var msg interface{}
	for !exists {
		time.Sleep(100 * time.Millisecond)
		msg, exists = sendRecvMap.Get(string(echoData))
	}
	return []byte(msg.(Message).RawMessage)

}

func SendPoke(userId int64, groupId int64, selfId int64) {
	sendData(struct {
		Action  string `json:"action"`
		UserId  int64  `json:"user_id"`
		SelfId  int64  `json:"self_id"`
		GroupId int64  `json:"group_id"`
	}{
		Action:  "send_poke",
		UserId:  userId,
		SelfId:  selfId,
		GroupId: groupId,
	})
}

func SendOperator(target, operator string, args []string) {
	sendData(struct {
		Action    string   `json:"action"`
		ApiName   string   `json:"api_name"`
		SubType   string   `json:"sub_type"`
		Arguments []string `json:"arguments"`
	}{
		Action:    "operator",
		ApiName:   target,
		SubType:   operator,
		Arguments: args,
	})
}
