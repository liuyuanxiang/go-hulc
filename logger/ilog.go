package logger

import (
	"encoding/json"
	"fmt"
	"log"
)

type ILog struct {
	savePath string
	saveName string
	isDev    bool

	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger

	reqLog *log.Logger
	dbLog  *log.Logger
}

// NewILog 返回一个基于 iLog 格式的日志记录器实现
func NewILog(opts ...ILogOption) *ILog {
	logger := &ILog{
		savePath: DefaultLogSavePath,
		saveName: DefaultLogSaveName,
	}

	for _, o := range opts {
		o(logger)
	}

	// iLog 相关的 5 个日志文件路径准备
	infoFile := fmt.Sprintf("%s%s_info.log", logger.savePath, logger.saveName)
	warnFile := fmt.Sprintf("%s%s_warning.log", logger.savePath, logger.saveName)
	errorFile := fmt.Sprintf("%s%s_error.log", logger.savePath, logger.saveName)
	reqFile := fmt.Sprintf("%s%s_request.log", logger.savePath, logger.saveName)
	dbFile := fmt.Sprintf("%s%s_database.log", logger.savePath, logger.saveName)

	// 创建 5 个对应文件的日志处理器，分别用于不同类型的日志记录
	logger.infoLog = log.New(openLogFile(infoFile), "", log.LstdFlags)
	logger.warnLog = log.New(openLogFile(warnFile), "", log.LstdFlags)
	logger.errorLog = log.New(openLogFile(errorFile), "", log.LstdFlags)
	logger.reqLog = log.New(openLogFile(reqFile), "", log.LstdFlags)
	logger.dbLog = log.New(openLogFile(dbFile), "", log.LstdFlags)
	return logger
}

type ILogOption func(*ILog)

func (l *ILog) Debug(v ...interface{}) {
	if !l.isDev {
		l.Info(v...)
	}
}

func (l *ILog) Info(v ...interface{}) {
	l.infoLog.Println(json.Marshal(newInfoLog(fmt.Sprint(v...))))
}

func (l *ILog) Warn(v ...interface{}) {
	l.infoLog.Println(json.Marshal(newWarningLog(fmt.Sprint(v...))))
}

func (l *ILog) Error(v ...interface{}) {
	l.infoLog.Println(json.Marshal(newErrorLog(fmt.Sprint(v...))))
}

func (l *ILog) DB(duration float64, v ...interface{}) {
	l.infoLog.Println(json.Marshal(newDatabaseLog(fmt.Sprint(v...), duration)))
}

func IsDev(isDev bool) ILogOption {
	return func(lg *ILog) {
		lg.isDev = isDev
	}
}

func SetSavePath(path string) ILogOption {
	return func(lg *ILog) {
		lg.savePath = path
	}
}

type requestLog struct {
	TraceID         string
	SpanID          string
	Method          string
	URL             string
	Duration        float64
	RequestTime     string
	QueryString     string
	GETParams       map[string]interface{}
	POSTParams      map[string]interface{}
	Cookie          map[string]interface{}
	UserAgent       string
	ReferURL        string
	ClientIP        string
	CPU             int64
	Memory          int64
	RequestHeaders  map[string]string
	ResponseHeaders map[string]string
}

type infoLog struct {
	TraceID string
	SpanID  string
	Message string
	Time    string
}

func newInfoLog(msg string) infoLog {
	return infoLog{Message: msg}
}

type databaseLog struct {
	TraceID  string
	SpanID   string
	DB       string
	Message  string
	Duration float64
	Time     string
}

func newDatabaseLog(msg string, duration float64) databaseLog {
	return databaseLog{Message: msg, Duration: duration}
}

type warningLog struct {
	TraceID string
	SpanID  string
	Message string
	Time    string
}

func newWarningLog(msg string) warningLog {
	return warningLog{Message: msg}
}

type errorLog struct {
	TraceID string
	SpanID  string
	Message string
	Time    string
}

func newErrorLog(msg string) errorLog {
	return errorLog{Message: msg}
}
