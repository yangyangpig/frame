package main

import (
	// "errors"
	// "fmt"
	// "math/rand"
	// "time"
	// "github.com/natefinch/lumberjack"
	// "log"
	// "os"
	"putil/log"
)

func main() {
	plog.SetOutPutFileRotate("/Users/dQingQuan/log/bird.log", 500, 3, 28)
	plog.Debug("bird test")

	// l := &lumberjack.Logger{
	// 	Filename:   "/Users/QingQuan/log/fooa.log",
	// 	MaxSize:    500, // megabytes
	// 	MaxBackups: 3,
	// 	MaxAge:     28,    //days
	// 	Compress:   false, // disabled by default
	// }

}
