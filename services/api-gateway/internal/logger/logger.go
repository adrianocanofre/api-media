package logger

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type entry struct {
	Time    string         `json:"time"`
	Level   Level          `json:"level"`
	Service string         `json:"service"`
	Msg     string         `json:"msg"`
	Fields  map[string]any `json:"fields,omitempty"`
}

type Logger struct {
	service string
	out     io.Writer
}

func New(service string) *Logger {
	return &Logger{
		service: service,
		out:     os.Stdout,
	}
}

func (l *Logger) log(level Level, msg string, fields map[string]any) {
	e := entry{
		Time:    time.Now().UTC().Format(time.RFC3339),
		Level:   level,
		Service: l.service,
		Msg:     msg,
		Fields:  fields,
	}

	b, err := json.Marshal(e)
	if err != nil {
		return
	}

	l.out.Write(append(b, '\n'))
}

func (l *Logger) Info(msg string, fields map[string]any) {
	l.log(LevelInfo, msg, fields)
}

func (l *Logger) Warn(msg string, fields map[string]any) {
	l.log(LevelWarn, msg, fields)
}

func (l *Logger) Error(msg string, fields map[string]any) {
	l.log(LevelError, msg, fields)
}
