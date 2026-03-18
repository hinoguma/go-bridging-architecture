package log

func Info(message string, requestFuncs ...LogRequestFunc) {
	GetGlobalLogger().Info(message, requestFuncs...)
}

func Error(message string, requestFuncs ...LogRequestFunc) {
	GetGlobalLogger().Error(message, requestFuncs...)
}
