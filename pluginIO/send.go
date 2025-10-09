package pluginIO

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
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

func SendMessage[T string | []MessageItem](msg T, userId int64, groupId int64) (messageId []int) {
	echo := fmt.Sprintf("send_msg_%s", RandomString(64))
	message := any(msg)
	switch message.(type) {
	case string:
		sendData("send_msg",
			strconv.FormatInt(userId, 10),
			strconv.FormatInt(groupId, 10),
			message.(string), echo)
	case []MessageItem:
		data, _ := json.Marshal(message.([]MessageItem))
		sendData("send_msg",
			strconv.FormatInt(userId, 10),
			strconv.FormatInt(groupId, 10),
			string(data), echo)
	}
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

func SendApi(apiName string, data []byte) []byte {
	echo := fmt.Sprintf("send_api_%s", RandomString(64))
	sendData("send_api", apiName, string(data), echo)
	defer sendRecvMap.Delete(echo)
	var exists = false
	var msg interface{}
	for !exists {
		time.Sleep(100 * time.Millisecond)
		msg, exists = sendRecvMap.Get(echo)
	}
	return []byte(msg.(Message).RawMessage)
}

func SendPoke(userId int64, groupId int64) {
	sendData("send_poke", strconv.FormatInt(userId, 10), strconv.FormatInt(groupId, 10))
}

func SendData(cmd string, args []string, hasReturn bool) string {
	echo := ""
	if hasReturn {
		echo = RandomString(64)
		args = append(args, echo)
	}
	sendData(cmd, args...)
	if hasReturn {
		defer sendRecvMap.Delete(echo)
		var exists = false
		var msg interface{}
		for !exists {
			time.Sleep(100 * time.Millisecond)
			msg, exists = sendRecvMap.Get(echo)
		}
		return msg.(Message).RawMessage
	} else {
		return ""
	}
}

func SendOperator(target, operator string, args []string) {
	args = append([]string{target, operator}, args...)
	sendData("operator", args...)
}
