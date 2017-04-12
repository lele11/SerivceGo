package utils

import (
	"fmt"
	"os"
	"runtime/debug"

	"time"
)

const (
	_ = iota
	LoggerDebug
	LoggerInfo
	LoggerError
)

func NewLogger(prefix, path string, level int) *Logger {
	l := &Logger{
		info:       make(chan string, 1000),
		filePrefix: prefix,
		path:       path,
	}
	go l.out()
	return l
}

type Logger struct {
	path       string
	filePrefix string
	level      int
	info       chan string
	open       *os.File
	day        int64
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level >= LoggerDebug {
		return
	}
	l.string("DEBUG", format, args...)
}
func (l *Logger) Info(format string, args ...interface{}) {
	if l.level >= LoggerInfo {
		return
	}
	l.string("INFO", format, args...)
}
func (l *Logger) Error(format string, args ...interface{}) {
	if l.level >= LoggerError {
		return
	}
	l.string("ERROR", format, args...)
	l.string("ERROR", "%s", debug.Stack())
}
func (l *Logger) string(kind string, format string, args ...interface{}) {
	s := fmt.Sprintf("%s [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), kind, fmt.Sprintf(format, args...))
	l.info <- s
}
func (l *Logger) out() {
	for {
		select {
		case s := <-l.info:
			l.write(s)
		}
	}
}

func (l *Logger) write(line string) {
	if l.day <= time.Now().Unix() {
		l.open = nil
	}
	if l.open == nil {
		var e error
		fileName := fmt.Sprintf("%s_%s", l.filePrefix, time.Now().Format("20060102"))
		l.open, e = os.OpenFile(fileName, os.O_APPEND, 0666)
		if e != nil {
			l.open, _ = os.Create(fileName)
		} else {
			fmt.Println(e)
		}
		l.day = GetDayBegin().Unix()
	}
	l.open.WriteString(line)
}
