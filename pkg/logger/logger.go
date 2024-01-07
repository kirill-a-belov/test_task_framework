package logger

import (
	"fmt"
	stdLog "log"
	"os"
	"strings"
)

type Logger interface {
	Error(err error, args ...interface{})
	Info(args ...interface{})
}

type log struct {
	l *stdLog.Logger
}

func (l *log) Error(err error, args ...interface{}) {
	var result strings.Builder

	for _, arg := range append(args, err) {
		if _, err := result.WriteString(fmt.Sprintf(" %v ", arg)); err != nil {
			panic(err)
		}
	}

	l.l.Println(result.String())
}

func (l *log) Info(args ...interface{}) {
	l.Error(nil, args...)
}

func New(prefix string) Logger {
	return &log{
		l: stdLog.New(os.Stdout, fmt.Sprintf("%s: ", prefix), stdLog.LstdFlags),
	}
}
