package pluginIO

func MessageRegister(callback func(Message)) {
	sendData(struct {
		Action       string `json:"action"`
		RegisterType string `json:"register_type"`
	}{
		Action:       "register",
		RegisterType: "message",
	})
	messageMap = append(messageMap, callback)
}

func CommandRegister(cmdName string, handler func(Message)) {
	sendData(struct {
		Action       string `json:"action"`
		RegisterType string `json:"register_type"`
		SubType      string `json:"sub_type"`
	}{
		Action:       "register",
		RegisterType: "command",
		SubType:      cmdName,
	})
	commandMap[cmdName] = handler
}

func NoticeRegister(noticeType string, handler func(Message)) {
	sendData(struct {
		Action       string `json:"action"`
		RegisterType string `json:"register_type"`
		SubType      string `json:"sub_type"`
	}{
		Action:       "register",
		RegisterType: "notice",
		SubType:      noticeType,
	})
	noticeMap[noticeType] = append(noticeMap[noticeType], handler)
}
