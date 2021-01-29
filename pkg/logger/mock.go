package logger

import (
	"bytes"
	"log"
)

// Mock to logger.
type Mock struct {
	*log.Logger
}

func (m Mock) Error(args ...interface{}) {
	m.Println(args...)
}

func (m Mock) Errorf(template string, args ...interface{}) {
	m.Printf(template, args...)
}

func (m Mock) Debug(args ...interface{}) {
	m.Println(args...)
}

func (m Mock) Debugf(template string, args ...interface{}) {
	m.Printf(template, args...)
}

func (m Mock) Info(args ...interface{}) {
	m.Println(args...)
}

func (m Mock) Infof(template string, args ...interface{}) {
	m.Printf(template, args...)
}

func (m Mock) Warn(args ...interface{}) {
	m.Println(args...)
}

func (m Mock) Warnf(template string, args ...interface{}) {
	m.Printf(template, args...)
}

func (m Mock) SetLevel(level logLevel) {}

// NewMock returns a new mock logger.
func NewMock() Logger {
	return &Mock{
		Logger: log.New(bytes.NewBuffer([]byte{}), "", 0),
	}
}
