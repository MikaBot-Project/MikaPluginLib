package main

import (
	"log"
	"strconv"

	"github.com/MikaBot-Project/MikaPluginLib/pluginConfig"
	"github.com/MikaBot-Project/MikaPluginLib/pluginData"
	"github.com/MikaBot-Project/MikaPluginLib/pluginIO"
)

var config = struct {
	Text string `json:"text"`
}{
	Text: "test",
}

func test(msg pluginIO.Message) {
	log.Println("post type:", msg.PostType)
	switch msg.PostType {
	case "notice":
		log.Println("notice subtype", msg.SubType)
		if msg.SubType == "poke" && msg.TargetId == msg.SelfId {
			pluginIO.SendPoke(msg.UserId, msg.GroupId)
			var data []pluginIO.MessageItem
			pluginData.ReadJson("j/"+strconv.FormatInt(msg.UserId, 10), &data)
			pluginIO.SendMessage(data, msg.UserId, msg.GroupId)
		}
	case "message":
		pluginData.SaveBinary("b/"+strconv.FormatInt(msg.UserId, 10), msg.RawMessage)
		pluginData.SaveJson("j/"+strconv.FormatInt(msg.UserId, 10), msg.MessageArray)
		if msg.GroupId == 0 {
			log.Println("messageId", pluginIO.SendMessage(msg.MessageArray, msg.UserId, 0))
		}
		if msg.AtMe {
			log.Println("messageId", pluginIO.SendMessage(msg.MessageArray, msg.UserId, msg.GroupId))
		}
	case "command":
		if len(msg.CommandArgs) < 2 {
			log.Println("messageId", pluginIO.SendMessage(config.Text, msg.UserId, msg.GroupId))
			return
		}
		switch msg.CommandArgs[1] {
		case "set":
			if len(msg.CommandArgs) < 3 {
				return
			}
			config.Text = msg.CommandArgs[2]
			pluginConfig.SaveJson("test.json", config)
		case "msg":
			var data string
			pluginData.ReadBinary("b/"+strconv.FormatInt(msg.UserId, 10), &data)
			pluginIO.SendMessage(data, msg.UserId, msg.GroupId)
		}
	}
}

func main() {
	log.Println("hello test log")
	pluginIO.MessageRegister(test)
	pluginIO.NoticeRegister("notify", test)
	pluginIO.CommandRegister("!test", test)
	pluginConfig.ReadJson("test.json", &config)
	var data pluginIO.Message
	for {
		data = <-pluginIO.MessageChan
		log.Println("type: ", data.PostType)
		log.Println("echo: ", data.RawMessage)
	}
}
