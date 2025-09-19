package main

import (
	"MikaPluginLib/pluginIO"
	"log"
	"time"
)

func main() {
	log.Println("hello test log")
	pluginIO.SendData("init", "test", "args")
	for {
		time.Sleep(1000)
		pluginIO.SendData("send_msg", "test", "args")
	}
}
