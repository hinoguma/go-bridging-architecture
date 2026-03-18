package errors

import "sync"

type errorConfig struct {
	sync.Mutex
	contextRequestIDKey string
}

func (conf *errorConfig) SetContextRequestIDKey(key string) {
	conf.Lock()
	defer conf.Unlock()
	conf.contextRequestIDKey = key
}

var config = errorConfig{
	contextRequestIDKey: "requestId",
}

func GetContextRequestIDKey() string {
	return config.contextRequestIDKey
}

func SetContextRequestIDKey(key string) {
	config.SetContextRequestIDKey(key)
}
