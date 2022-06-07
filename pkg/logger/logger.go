package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"
)

var (
	finishing = ""
)

type Logger struct {
	log      *log.Logger
	Severity int
	mu       *sync.Mutex
	File     io.Writer
}

func New(prefix string, severity int) *Logger {
	if severity >= 1 {
		finishing = "\r"
	}
	return &Logger{
		log:      log.New(log.Writer(), prefix, 0),
		Severity: severity,
		mu:       &sync.Mutex{},
		File:     ioutil.Discard,
	}
}

func (l Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.log.SetPrefix(prefix)
}

func (l *Logger) _log(message string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	message = fmt.Sprintf(message, args...)
	l.log.Print(message + finishing)
	l.File.Write([]byte(message + "\n"))
}

func (l Logger) Debug(message string, arg ...interface{}) {
	if l.Severity >= 3 {
		l._log(message, arg...)
	}
}

func (l Logger) Informational(message string, arg ...interface{}) {
	if l.Severity >= 2 {
		l._log(message, arg...)
	}
}

func (l Logger) Default(message string, arg ...interface{}) {
	if l.Severity >= 1 {
		l._log(message, arg...)
	}
}

func (l Logger) Critical(message string, arg ...interface{}) {
	if l.Severity >= 0 {
		l._log(message, arg...)
	}
}
