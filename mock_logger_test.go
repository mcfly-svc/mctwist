package main_test

import "fmt"

type MockLogger struct {
	Output      string
	ErrorOutput string
}

func (l *MockLogger) Error(err error) {
	l.ErrorOutput = err.Error()
}

func (l *MockLogger) ApiError(v interface{}) {
	l.ErrorOutput = fmt.Sprintf("mcflyapi responded with an error: %+v\n", v)
}

func (l *MockLogger) Log(msg string) {
	l.Output = msg
}
