package clog

type CLog interface {
	//初始化日志库
	Init()
	//设置日志的打印等级
	SetLevel(level int)
	//打印不同等级的日志
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	//释放日志所占用的资源
	Close()
}