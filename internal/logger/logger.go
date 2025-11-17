package logger

import (
	"log"
	"os"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

func New(filename string) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	l := log.New(f, "wrenchi: ", log.LstdFlags|log.Lmicroseconds)

	return &Logger{
		file:   f,
		logger: l,
	}, nil
}

func (l *Logger) Printf(format string, v ...any) {
	l.logger.Printf(format, v...)
}
