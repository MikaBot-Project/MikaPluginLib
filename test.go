package main

import (
	"MikaPluginLib/pluginConfig"
	"MikaPluginLib/pluginIO"
	"log"
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
	case "message":
		if msg.GroupId == 0 {
			log.Println("messageId", pluginIO.SendMessage(msg.MessageArray, msg.UserId, 0))
		}
		if msg.AtMe {
			log.Println("messageId", pluginIO.SendMessage(msg.MessageArray, msg.UserId, msg.GroupId))
		}
	case "command":
		log.Println("messageId", pluginIO.SendMessage("get cmd "+msg.CommandArgs[0], msg.UserId, msg.GroupId))
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
