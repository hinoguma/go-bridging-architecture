package log

import (
	"app/pkg/crosscutting/errors"
	"app/pkg/crosscutting/timer"
	"fmt"
	"time"
)

var logger Logger = NewStdLogger()

func GetGlobalLogger() Logger {
	return logger
}

type Logger interface {
	Info(req LogRequest)
	Error(req ErrorLogRequest)
}

type LogRequest struct {
	Message string
	Tags    map[string]any
	Time    *time.Time
}

type ErrorLogRequest struct {
	Err  error
	Time *time.Time
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

func (logger StdLogger) Info(req LogRequest) {
	if req.Time == nil {
		t := timer.Now()
		req.Time = &t
	}
	fmt.Println(
		fmt.Sprintf(
			`{"level":"%s","message":"%s","error": {},"tags":{},"Time":"%s"}`,
			ErrorLevel, req.Message, req.Time,
		),
	)
}

func (logger StdLogger) Error(req ErrorLogRequest) {
	if req.Time == nil {
		t := timer.Now()
		req.Time = &t
	}
	fmt.Println(
		fmt.Sprintf(
			`{"level":"%s","message":"error","error": %s,"tags":{}","Time":"%s"}`,
			ErrorLevel, errors.ToJsonString(req.Err), req.Time,
		),
	)
}

type StdLogFormat struct {
	Level   LogLevel       `json:"level"`
	Message string         `json:"message"`
	Tags    map[string]any `json:"tags,omitempty"`
	Time    timer.Ymd_Hms  `json:"Time"`
}
