package main

import (
	"fmt"
	"path"
	"runtime"
)

func GetStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func Errorf(format string, err error) error {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf(fmt.Sprintf("[runtime.Caller failed]%s", format), err)
	} else {
		funcName := runtime.FuncForPC(pc).Name()
		fileName := path.Base(file)
		return fmt.Errorf(fmt.Sprintf("[%s.%s:%d]%s", fileName, funcName, line, format), err)
	}

}
