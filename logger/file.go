package logger

import (
	"log"
	"os"
)

var (
	// 文件默认存放在项目的 runtime 目录下
	DefaultLogSavePath = "runtime/logs/"
	DefaultLogSaveName = "app"
)

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+DefaultLogSavePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
