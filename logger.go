package main

import (
	"fmt"
	"log"
)

type Logger interface {
	Error(error)
	ApiError(interface{})
	Log(string)
}

type McTwistLogger struct{}

func (l *McTwistLogger) Error(err error) {
	panic(err)
}

func (l *McTwistLogger) ApiError(v interface{}) {
	panic(fmt.Errorf("mcflyapi responded with an error: %+v\n", v))
}

func (l *McTwistLogger) Log(msg string) {
	log.Println(msg)
}
