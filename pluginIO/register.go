package pluginIO

func MessageRegister(callback func(Message)) {
	SendData("register", "message")
	messageMap = append(messageMap, callback)
}

func CommandRegister(cmdName string, handler func(Message)) {
	SendData("register", "cmd", cmdName)
	commandMap[cmdName] = handler
}

func NoticeRegister(noticeType string, handler func(Message)) {
	SendData("register", "notice", noticeType)
	noticeMap[noticeType] = append(noticeMap[noticeType], handler)
}
