// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package log

import (
	"fmt"
	"github.com/leovs/yc-go-sdk/errors"
)

var DebugMode = false

const (
	infoTag  = "[INFO]"
	errorTag = "[ERROR]"
	debugTag = "[DEBUG]"
)

func Println(a ...interface{}) {
	fmt.Println(a...)
}

func Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func Error(format string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s\n", errorTag, format), a...)
}

func ErrorMessage(err *errors.Message) {
	fmt.Printf(fmt.Sprintf("%s %s#%d\n", errorTag, err.Msg, err.Code))
}

func Info(format string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf("%s %s\n", infoTag, format), a...)
}

func Debug(format string, a ...interface{}) {
	if DebugMode {
		fmt.Printf(fmt.Sprintf("%s %s\n", debugTag, format), a...)
	}
}
func Panic(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
