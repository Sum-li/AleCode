package clog

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogData struct {
	//打印错误信息的具体内容
	Message      string
	//日志时间
	TimeStr      string
	//日志等级
	LevelStr     string
	//日志输出的文件名
	Filename     string
	//出错位置所处的函数
	FuncName     string
	//出错行数
	LineNo       int
	//是否输出到错误日志文件中去
	WarnAndFatal bool
}

//获取出错位置相关信息
func GetLineInfo() (fileName string, funcName string, lineNo int) {
	//skip表示提升的栈数，及相对于当前出错位置函数的第几次调用位置（0 表示当前函数本身）
	// pc:被调用函数的指针，file: 函数所处的文件路径，line: 当前所处的行数，ok: 是否能获取到pc,file,line等值
	pc, file, line, ok := runtime.Caller(4)
	if ok {
		fileName = path.Base(file)
		funcName = path.Base(runtime.FuncForPC(pc).Name())
		lineNo = line
	}
	return
}

/*
1. 当业务调用打日志的方法时，我们把日志相关的数据写入到chan（队列）
2. 然后我们有一个后台的线程不断的从chan里面获取这些日志，最终写入到文件。
*/
func writeLog(level int, format string, args ...interface{}) *LogData {
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05")
	levelStr := getLevelText(level)

	fileName, funcName, lineNo := GetLineInfo()
	msg := fmt.Sprintf(format, args...)

	logData := &LogData{
		Message:      msg,
		TimeStr:      nowStr,
		LevelStr:     levelStr,
		Filename:     fileName,
		FuncName:     funcName,
		LineNo:       lineNo,
		WarnAndFatal: false,
	}

	if level == LogLevelError || level == LogLevelWarn || level == LogLevelFatal {
		logData.WarnAndFatal = true
	}

	return logData
}
