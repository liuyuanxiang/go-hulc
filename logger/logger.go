package logger

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

type Level int

var (
	// F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	dbLogger    *log.Logger

	logPrefix  = ""
	levelFlags = []string{"debug", "info", "warning", "error", "database"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
	SQL
)

func init() {
	infoLogger = initInfoLogger()
	warnLogger = initWarnLogger()
	errorLogger = initErrorLogger()
	debugLogger = initDebugLogger()
	dbLogger = initDbLogger()
}

func initInfoLogger() *log.Logger {
	infoFilePath := getLogFileFullPath("INFO")
	infoFile := openLogFile(infoFilePath)

	return log.New(infoFile, DefaultPrefix, log.LstdFlags)
}

func initWarnLogger() *log.Logger {
	infoFilePath := getLogFileFullPath("INFO")
	infoFile := openLogFile(infoFilePath)

	return log.New(infoFile, DefaultPrefix, log.LstdFlags)
}

func initErrorLogger() *log.Logger {
	infoFilePath := getLogFileFullPath("INFO")
	infoFile := openLogFile(infoFilePath)

	return log.New(infoFile, DefaultPrefix, log.LstdFlags)
}

func initDebugLogger() *log.Logger {
	infoFilePath := getLogFileFullPath("INFO")
	infoFile := openLogFile(infoFilePath)

	return log.New(infoFile, DefaultPrefix, log.LstdFlags)
}

func initDbLogger() *log.Logger {
	infoFilePath := getLogFileFullPath("INFO")
	infoFile := openLogFile(infoFilePath)

	return log.New(infoFile, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	debugLogger.Println(v...)
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	infoLogger.Println(v...)
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	warnLogger.Println(v...)
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	errorLogger.Println(v...)
}

func Fatal(v ...interface{}) {
	setPrefix(ERROR)
	errorLogger.Fatalln(v...)
}

func DB(v ...interface{}) {
	setPrefix(SQL)
	dbLogger.Println(v...)
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	infoLogger.SetPrefix(logPrefix)
}
