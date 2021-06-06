package logger

import (
	"log"
	"os"
)

var (
	DefaultLogSavePath = "/webser/www/logs/application/"
	DefaultLogSaveName = "ilog"
	DefaultLogSaveExt  = "log"
)

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		log.Fatalf("logs path not exist")
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}
