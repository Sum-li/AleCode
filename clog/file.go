package clog

import (
	"fmt"
	"os"
	"time"
)

type FileLogger struct {
	//日志的等级
	level         int
	//日志文件存放的路径
	logPath       string
	//日志文件的 基文件名
	logName       string
	//普通日志的文件句柄
	file          *os.File
	//错误日志的文件句柄
	warnFile      *os.File
	//日志写入文件的传输 管道
	LogDataChan   chan *LogData
	//日志文件的 切割方式 (hour, size)
	logSplitType  string
	//按大小切割的 文件大小
	logSplitSize  int64
	//按时间切割，最后一次切割的时间戳
	lastSplitHour int
}

func NewFileLogger(config Config) (log CLog, err error) {

	if config.Log_path == "" && config.Log_name == "" && config.Log_level == "" {
		err = fmt.Errorf("Incomplete configuration information\n")
		return
	}

	//通道的大小
	var chanSize int
	if config.Log_chan_size < 5000 {
		chanSize = 50000
	} else {
		chanSize = config.Log_chan_size
	}

	//切割方式(默认为 hour)
	var logSplitType = "hour"
	//切割的文件大小
	var logSplitSize int64

	//10485760为 1MB
	if config.Log_split_type == "size" {
		if config.Log_split_size < 10485760 {
			logSplitSize = 104857600
		} else {
			logSplitSize = config.Log_split_size
		}
		logSplitType = config.Log_split_type
	}

	level := getLogLevel(config.Log_level)
	log = &FileLogger{
		level:         level,
		logPath:       config.Log_path,
		logName:       config.Log_name,
		LogDataChan:   make(chan *LogData, chanSize),
		logSplitSize:  logSplitSize,
		logSplitType:  logSplitType,
		lastSplitHour: time.Now().Hour(),
	}
	return
}

func (f *FileLogger) Init() {

	filename := fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	//获得日志文件的文件句柄
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
	}
	f.file = file

	filename = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	//获得错误日志文件的文件句柄
	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
	}
	f.warnFile = file


	go f.writeLogBackground()
}

//通过hour切割日志文件
func (f *FileLogger) splitFileHour(warnFile bool) {
	//判断是否需要切割
	now := time.Now()
	if now.Hour() == f.lastSplitHour {
		return
	}
	//跟新切割时间
	f.lastSplitHour = now.Hour()
	//被切割的文件
	var backupFilename string
	//当前要写入日志的文件
	var filename string

	if warnFile {
		backupFilename = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)

		filename = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFilename = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), f.lastSplitHour)
		filename = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}

	file := f.file
	if warnFile {
		file = f.warnFile
	}

	_ = file.Close()
	_ = os.Rename(filename, backupFilename)

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
		return
	}

	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

//通过size切割日志文件
func (f *FileLogger) splitFileSize(warnFile bool) {

	file := f.file
	if warnFile {
		file = f.warnFile
	}

	//Stat()返回文件信息的结构体
	statInfo, err := file.Stat()
	if err != nil {
		panic(fmt.Sprintf("get file info failed, err:%v", err))
		return
	}

	//判断是否需要切割
	fileSize := statInfo.Size()
	if fileSize <= f.logSplitSize {
		return
	}

	var backupFilename string
	var filename string

	now := time.Now()
	if warnFile {
		backupFilename = fmt.Sprintf("%s/%s.log.wf_%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

		filename = fmt.Sprintf("%s/%s.log.wf", f.logPath, f.logName)
	} else {
		backupFilename = fmt.Sprintf("%s/%s.log_%04d%02d%02d%02d%02d%02d",
			f.logPath, f.logName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		filename = fmt.Sprintf("%s/%s.log", f.logPath, f.logName)
	}

	_ = file.Close()
	_ = os.Rename(filename, backupFilename)

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(fmt.Sprintf("open file %s failed, err:%v", filename, err))
		return
	}

	if warnFile {
		f.warnFile = file
	} else {
		f.file = file
	}
}

func (f *FileLogger) checkSplitFile(warnFile bool) {

	if f.logSplitType == "hour" {
		f.splitFileHour(warnFile)
		return
	}

	f.splitFileSize(warnFile)
}

func (f *FileLogger) writeLogBackground() {
	for logData := range f.LogDataChan {

		f.checkSplitFile(logData.WarnAndFatal)

		file  := f.file
		if logData.WarnAndFatal {
			file = f.warnFile
		}


		_,_ = fmt.Fprintf(file, "[%s %s] (%s:%s:%d) %s\r\n", logData.TimeStr,
			logData.LevelStr, logData.Filename, logData.FuncName, logData.LineNo, logData.Message)
	}
}

func (f *FileLogger) SetLevel(level int) {
	if level < LogLevelDebug || level > LogLevelFatal {
		level = LogLevelDebug
	}
	f.level = level
}

func (f *FileLogger) Debug(format string, args ...interface{}) {
	if f.level > LogLevelDebug {
		return
	}

	logData := writeLog(LogLevelDebug, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Trace(format string, args ...interface{}) {
	if f.level > LogLevelTrace {
		return
	}
	logData := writeLog(LogLevelTrace, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Info(format string, args ...interface{}) {
	if f.level > LogLevelInfo {
		return
	}
	logData := writeLog(LogLevelInfo, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Warn(format string, args ...interface{}) {
	if f.level > LogLevelWarn {
		return
	}

	logData := writeLog(LogLevelWarn, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Error(format string, args ...interface{}) {
	if f.level > LogLevelError {
		return
	}

	logData := writeLog(LogLevelError, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Fatal(format string, args ...interface{}) {
	if f.level > LogLevelFatal {
		return
	}

	logData := writeLog(LogLevelFatal, format, args...)
	select {
	case f.LogDataChan <- logData:
	default:
	}
}

func (f *FileLogger) Close() {
	_ = f.file.Close()
	_ = f.warnFile.Close()
}
