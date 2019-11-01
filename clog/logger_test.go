package clog

import (
	"testing"
)

func TestFileLogger(t *testing.T) {
	logger,_ := NewFileLogger(Config{Log_level:"debug",Log_name:"testLog",Log_path:"C:\\Users\\Administrator\\Desktop"})
	logger.Debug("user id[%d] is come from china", 324234)
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Close()
}

func TestConsoleLogger(t *testing.T) {
	logger,_ := NewConsoleLogger(Config{Log_level:"debug"})
	logger.Debug("user id[%d] is come from china", 324234)
	logger.Warn("test warn log")
	logger.Fatal("test fatal log")
	logger.Close()
}
