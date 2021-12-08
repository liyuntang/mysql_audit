package common

import (
	"fmt"
	"io"
	"log"
	"os"
)

func WriteLog(logFile string) *log.Logger {
	// 判断logFile是否存在,根据logFile的状态（存在/不存在）定义flag，根据flag打开文件
	var flag int
	_, err := os.Stat(logFile)
	if os.IsNotExist(err) {
		// 说明logFile不存在
		fmt.Println("logFile", logFile, "is not exist,and now we create it")
		flag = os.O_CREATE|os.O_WRONLY
	} else {
		// 说明logFile存在
		flag = os.O_WRONLY|os.O_APPEND
	}
	file, err := os.OpenFile(logFile, flag, 0644)
	if err != nil {
		fmt.Println("open logFile", logFile, "is bad, err is", err)
		os.Exit(0)
	}
	return log.New(io.MultiWriter(file, os.Stdout), "[mysql_audit ]", log.Lshortfile|log.Ltime|log.Ldate)
}

