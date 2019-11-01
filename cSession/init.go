package cSession

import (
	"AleCode/clog"
	"errors"
	"fmt"
	"time"
)

var sessionMgr SessionMgr

//初始化
func InitSessionMgr(provider string,gct time.Duration,options ...interface{}) (err error) {
	switch provider {
	case "memory":
		sessionMgr = newMemorySessionMgr(gct)
	case "redis":
		sessionMgr = newRedisSessionMgr(gct)
	default:
		err = errors.New(fmt.Sprintf("%s not support",provider))
		clog.Error("err:%v",err)
		return
	}
	err = sessionMgr.Init(options...)
	if err != nil {
		clog.Error("init sessionMgr failed,err:%v",err)
		return
	}
	return
}

func Create(mlt time.Duration) (session Session,err error) {
	return sessionMgr.Create(mlt)
}

func Get(sessionId string) (session Session,err error) {
	return sessionMgr.Get(sessionId)
}

func Del(sessionId string) {
	sessionMgr.Del(sessionId)
}

func Update(sessionId string) error {
	return sessionMgr.Update(sessionId)
}
