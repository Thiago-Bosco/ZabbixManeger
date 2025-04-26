
package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Logger struct {
	output *log.Logger
}

type LogEntry struct {
	Timestamp time.Time     `json:"timestamp"`
	Level     string       `json:"level"`
	Message   string       `json:"message"`
	Data      interface{}  `json:"data,omitempty"`
}

func Novo() *Logger {
	return &Logger{
		output: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Info(msg string, data interface{}) {
	l.log("INFO", msg, data)
}

func (l *Logger) Error(msg string, data interface{}) {
	l.log("ERROR", msg, data)
}

func (l *Logger) log(level, msg string, data interface{}) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Data:      data,
	}
	
	json, _ := json.Marshal(entry)
	l.output.Println(string(json))
}
