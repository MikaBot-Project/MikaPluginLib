package main

import (
	"MikaPluginLib/pluginIO"
	"log"
)

func test(msg pluginIO.Message) {
	log.Println("post type:", msg.PostType)
	switch msg.PostType {
	case "notice":
		log.Println("notice subtype", msg.SubType)
	case "message":
		if msg.GroupId == 0 {
			log.Println("messageId", pluginIO.SendMessage(msg.MessageArray, msg.UserId, 0))
		}
	case "command":
		log.Println("messageId", pluginIO.SendMessage("get cmd"+msg.MetaEventType, msg.UserId, msg.GroupId))
	}
}

func main() {
	log.Println("hello test log")
	pluginIO.MessageRegister(test)
	pluginIO.NoticeRegister("notify", test)
	pluginIO.CommandRegister("!test", test)
	var data pluginIO.Message
	for {
		data = <-pluginIO.MessageChan
		log.Println("type: ", data.PostType)
		log.Println("echo: ", data.RawMessage)
	}
}
