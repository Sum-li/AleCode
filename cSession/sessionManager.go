package cSession

import "time"

type SessionMgr interface {
	Init(options ...interface{}) (err error)
	Create(mlt time.Duration) (session Session, err error)
	Get(sessionId string) (session Session, err error)
	cgc()
	Del(sessinId string)
	Update(sessionId string) error
}
