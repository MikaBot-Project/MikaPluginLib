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
		SendData("send_msg",
			strconv.FormatInt(userId, 10),
			strconv.FormatInt(groupId, 10),
			message.(string), echo)
	case []MessageItem:
		data, _ := json.Marshal(message.([]MessageItem))
		SendData("send_msg",
			strconv.FormatInt(userId, 10),
			strconv.FormatInt(groupId, 10),
			string(data), echo)
	}
	defer delete(sendRecvMap, echo)
	var exists = false
	for !exists {
		_, exists = sendRecvMap[echo]
	}
	_ = json.Unmarshal([]byte(sendRecvMap[echo].RawMessage), &messageId)
	return messageId
}

func SendApi(apiName string, data []byte) []byte {
	echo := fmt.Sprintf("send_api_%s", RandomString(64))
	SendData("send_api", apiName, string(data), echo)
	defer delete(sendRecvMap, echo)
	var exists = false
	for !exists {
		_, exists = sendRecvMap[echo]
	}
	return []byte(sendRecvMap[echo].RawMessage)
}
