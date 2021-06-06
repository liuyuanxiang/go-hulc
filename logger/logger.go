package logger

type LogInterface interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	DB(float64, ...interface{})
}

type Level int

var (
	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	lg LogInterface
)

func init() {
	lg = NewILog()
}

// Logger 返回当前 logger 包下正在使用的 log 实例
// 默认使用 ilog 的实现
func Logger() LogInterface {
	return lg
}

// Debug 记录调试类型日志信息
func Debug(v ...interface{}) { lg.Debug(v...) }

// Info 记录常规日志信息
func Info(v ...interface{}) { lg.Info(v...) }

// Warn 记录警告类型日志信息
func Warn(v ...interface{}) { lg.Warn(v...) }

// Error 记录错误类型日志信息
func Error(v ...interface{}) { lg.Error(v...) }

// DB 记录数据库执行记录相关信息
func DB(duration float64, v ...interface{}) { lg.DB(duration, v...) }
