package pluginIO

func MessageRegister(callback func(Message)) {
	sendData("register", "message")
	messageMap = append(messageMap, callback)
}

func CommandRegister(cmdName string, handler func(Message)) {
	sendData("register", "cmd", cmdName)
	commandMap[cmdName] = handler
}

func NoticeRegister(noticeType string, handler func(Message)) {
	sendData("register", "notice", noticeType)
	noticeMap[noticeType] = append(noticeMap[noticeType], handler)
}
