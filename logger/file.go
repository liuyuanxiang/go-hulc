package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	LogSavePath = "/webser/www/logs/"
	LogSaveName = "ilog"
	LogSaveExt  = "log"
	TimeFormat  = "20060102"
)

func getLogFileSavePath() string {
	return LogSavePath
}

func getLogFileFullPath(category string) string {
	var fileName string
	if category == "" {
		fileName = fmt.Sprintf("%s.%s", LogSaveName, LogSaveExt)
	} else {
		fileName = fmt.Sprintf("%s_%s.%s", LogSaveName, category, LogSaveExt)
	}
	return fmt.Sprintf("%s%s", getLogFileSavePath(), fileName)
}

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
	err := os.MkdirAll(dir+"/"+getLogFileSavePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
