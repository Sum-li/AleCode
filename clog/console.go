package clog

import (
	"fmt"
	"os"
)

type consoleLogger struct {
	//日志等级
	level int
}

func NewConsoleLogger(config Config) (log CLog,err error) {
	var level int
	if config.Log_level == "" {
		level = LogLevelDebug
	} else {
		level = getLogLevel(config.Log_level)
	}
	log = &consoleLogger{
		level: level,
	}
	return
}

func (c *consoleLogger) Init() {}

func (c *consoleLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	c.level = level
}

func (c *consoleLogger) Debug(format string, args ...interface{}) {
	if c.level > LogLevelDebug {
		return
	}

	logData := writeLog(LogLevelDebug, format, args...)
	_,_ = fmt.Fprintf(os.Stdout, "%s [%s %s] (%s:%s:%d) %s\n",Gray, logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *consoleLogger) Trace(format string, args ...interface{}) {
	if c.level > LogLevelTrace {
		return
	}

	logData := writeLog(LogLevelTrace, format, args...)
	_,_ = fmt.Fprintf(os.Stdout, "%s [%s %s] (%s:%s:%d) %s\n",Silvery ,logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}
func (c *consoleLogger) Info(format string, args ...interface{}) {
	if c.level > LogLevelInfo {
		return
	}

	logData := writeLog(LogLevelInfo, format, args...)
	_,_ = fmt.Fprintf(os.Stdout, "%s [%s %s] (%s:%s:%d) %s\n",Green ,logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *consoleLogger) Warn(format string, args ...interface{}) {
	if c.level > LogLevelWarn {
		return
	}

	logData := writeLog(LogLevelWarn, format, args...)
	_,_ = fmt.Fprintf(os.Stdout, "%s [%s %s] (%s:%s:%d) %s\n",Blue ,logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *consoleLogger) Error(format string, args ...interface{}) {
	if c.level > LogLevelError {
		return
	}

	logData := writeLog(LogLevelError, format, args...)
	_,_ = fmt.Fprintf(os.Stdout, "%s [%s %s] (%s:%s:%d) %s\n",Pink ,logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}
func (c *consoleLogger) Fatal(format string, args ...interface{}) {
	if c.level > LogLevelFatal {
		return
	}

	logData := writeLog(LogLevelFatal, format, args...)
	_,_ = fmt.Fprintf(os.Stdout, "%s [%s %s] (%s:%s:%d) %s\n",Red ,logData.TimeStr,
		logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
}

func (c *consoleLogger) Close() {

}
