package main

import (
	"MikaPluginLib/pluginIO"
	"log"
	"time"
)

func main() {
	log.Println("hello test log")
	pluginIO.SendData("init", "test", "args")
	go func() {
		var data pluginIO.Message
		for {
			data = <-pluginIO.MessageChan
			log.Println("type: ", data.PostType)
			log.Println("echo: ", data.RawMessage)
		}
	}()
	for {
		time.Sleep(10 * time.Second)
		pluginIO.SendData("send_msg", "test", "args")
	}
}
