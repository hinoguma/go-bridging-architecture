package log

import (
	"app/internal/crosscutting/errors"
	"app/internal/crosscutting/timer"
	"encoding/json"
	"fmt"
	"time"
)

var logger Logger = NewStdLogger()

func GetGlobalLogger() Logger {
	return logger
}

type Logger interface {
	Info(message string, requestFuncs ...LogRequestFunc)
	Error(message string, requestFuncs ...LogRequestFunc)
}

type LogRequest struct {
	Level     LogLevel
	Message   string
	Tags      map[string]any
	Err       error
	Time      *time.Time
	RequestID string
}

func (req *LogRequest) JsonString() string {
	if req.Time == nil {
		t := timer.Now()
		req.Time = &t
	}

	tagsStr := "{}"
	if req.Tags != nil {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			tagsStr = "{}"
		} else {
			tagsStr = string(tagsBytes)
		}
	}
	errStr := "{}"
	if req.Err != nil {
		errStr = errors.ToJsonString(req.Err)
	}

	return fmt.Sprintf(
		`{"level": "%s", "message": "%s", "error": %s, "tags": %s, "time": "%s"}`,
		req.Level, req.Message, errStr, tagsStr, req.Time,
	)
}

type LogRequestFunc func(req *LogRequest)

func WithErr(err error) LogRequestFunc {
	return func(req *LogRequest) {
		req.Err = err
	}
}

func WithTags(tags map[string]any) LogRequestFunc {
	return func(req *LogRequest) {
		req.Tags = tags
	}
}

func WithRequestID(requestID string) LogRequestFunc {
	return func(req *LogRequest) {
		req.RequestID = requestID
	}
}

type LogLevel string

const (
	InfoLevel  LogLevel = "INFO"
	ErrorLevel LogLevel = "ERROR"
)

type StdLogger struct{}

func NewStdLogger() Logger {
	return StdLogger{}
}

func (logger StdLogger) Info(message string, requestFuncs ...LogRequestFunc) {
	req := LogRequest{
		Message: message,
	}
	for _, f := range requestFuncs {
		f(&req)
	}
	req.Level = InfoLevel
	if req.Time == nil {
		t := timer.Now()
		req.Time = &t
	}
	fmt.Println(req.JsonString())
}

func (logger StdLogger) Error(message string, requestFuncs ...LogRequestFunc) {
	req := LogRequest{
		Message: message,
	}
	for _, f := range requestFuncs {
		f(&req)
	}
	req.Level = ErrorLevel
	if req.Time == nil {
		t := timer.Now()
		req.Time = &t
	}
	fmt.Println(req.JsonString())
}
