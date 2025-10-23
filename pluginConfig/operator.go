package pluginConfig

import "github.com/MikaBot-Project/MikaPluginLib/pluginIO"

var onReloadHandler []func()

func operatorHandler(msg pluginIO.Message) {
	if msg.SubType == "reload" {
		for fileName, config := range readJsonMap {
			ReadJson(fileName, config)
		}
		for path, config := range readAllJsonMap {
			ReadAllJson(path, config.(*map[string]any))
		}
	}
	for _, handler := range onReloadHandler {
		handler()
	}
}

func AddOnReloadHandler(handler func()) {
	onReloadHandler = append(onReloadHandler, handler)
}
