package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	baseLogger *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func init() {
	filePath := getLogFileFullPath()
	F = openLogFile(filePath)

	baseLogger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	baseLogger.Println(v)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	baseLogger.Println(v)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	baseLogger.Println(v)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	baseLogger.Println(v)
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	baseLogger.Fatalln(v)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	baseLogger.SetPrefix(logPrefix)
}
