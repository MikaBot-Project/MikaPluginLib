package main

import (
	"MikaPluginLib/pluginIO"
	"log"
	"time"
)

func test(msg pluginIO.Message) {
	log.Println("post type:", msg.PostType)
}

func main() {
	log.Println("hello test log")
	go func() {
		var data pluginIO.Message
		for {
			data = <-pluginIO.MessageChan
			log.Println("type: ", data.PostType)
			log.Println("echo: ", data.RawMessage)
		}
	}()
	pluginIO.MessageRegister(test)
	pluginIO.NoticeRegister("poke", test)
	pluginIO.CommandRegister("!test", test)
	for {
		time.Sleep(10 * time.Second)
		pluginIO.SendData("send_msg", "test", "args")
	}
}
