package clog

//日志的等级，初始化时，通过不同的等级打印不同的日志
const (
	LogLevelDebug = iota
	LogLevelTrace
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)



//创建Log对象的时候的参数
type Config struct {
	//日志等级 (debug,trace,info,warn,error,fatal)
	Log_level string
	//日志输出到文件时，日志文件所在的位置
	Log_path string
	//日志输出到文件时，日志文件的 基文件名
	Log_name string
	//日志输出到文件时，将日志通过管道传输，管道的大小，默认50000
	Log_chan_size int
	//日志输出到文件时，日志文件的切分方式 (hour,size)
	Log_split_type string
	//日志输出到文件时，切割的文件大小（默认 104857600）
	Log_split_size int64
}

func getLogLevel(level string) int {
	switch level {
	case "debug":
		return LogLevelDebug
	case "trace":
		return LogLevelTrace
	case "info":
		return LogLevelInfo
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "fatal":
		return LogLevelFatal
	}
	return LogLevelDebug
}

func getLevelText(level int) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelTrace:
		return "TRACE"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	}
	return "UNKNOWN"
}

////日志切分的方式{1. 小时切分，2. 文件大小切分}
//const (
//	LogSplitTypeHour = iota
//	LogSplitTypeSize
//)

//打印日志的颜色
//右边为 ANSI转义序列 详见 https://zh.wikipedia.org/wiki/ANSI%E8%BD%AC%E4%B9%89%E5%BA%8F%E5%88%97
const (
	//Fatal
	Red    = "\033[;31m"
	//Info
	Green  = "\033[;32m"
	//Warn
	Blue   = "\033[;94m"
	//Error
	Pink   = "\033[;95m"
	//Debug
	Gray = "\033[;90m"
	//Trace
	Silvery   = "\033[;37m"
)
