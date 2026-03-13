package log

func Info(req LogRequest) {
	GetGlobalLogger().Info(req)
}

func Error(req ErrorLogRequest) {
	GetGlobalLogger().Error(req)
}
