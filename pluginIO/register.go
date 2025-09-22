package pluginIO

func MessageRegister() {
	SendData("register", "message")
}

func CommandRegister(cmdName string) {
	SendData("register", "cmd", cmdName)
}

func NoticeRegister(noticeType string) {
	SendData("register", "notice", noticeType)
}
