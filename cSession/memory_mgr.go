package cSession

import (
	"AleCode/clog"
	"errors"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

//内存session的session管理器
type memorySessionMgr struct {
	//用于储存session
	sessionMap map[string]Session
	rwlock sync.RWMutex
	//释放session垃圾的周期
	gcTime time.Duration
}

func newMemorySessionMgr(gct time.Duration) SessionMgr {
	return &memorySessionMgr{gcTime:gct}
}

//初始化session管理器
func (s *memorySessionMgr) Init(options ...interface{}) (err error) {
	s.sessionMap = make(map[string]Session,1024)
	go s.cgc()
	return
}

//获得session
func (s *memorySessionMgr) Get(sessionId string) (session Session,err error) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()

	session,ok := s.sessionMap[sessionId]
	if !ok {
		err = errors.New("session not exit")
		clog.Error("err:%v",err)
		return
	}
	return
}

func (s *memorySessionMgr) Del(sessionId string) {
	s.rwlock.Lock()
	delete(s.sessionMap,sessionId)
	s.rwlock.Unlock()
	return
}

func (s *memorySessionMgr) Update(sessionId string) (err error) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	session,ok := s.sessionMap[sessionId]
	if !ok {
		err = errors.New("session not exit")
		clog.Error("err:%v",err)
		return
	}
	session.Update()
	return
}

//创建session
func (s *memorySessionMgr) Create(mlt time.Duration) (session Session,err error) {
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	var(
		id uuid.UUID
		sessionId string
	)
	id,err =uuid.NewV4()
	if err != nil {
		clog.Error("get uuid failed,err:%v",err)
		return
	}
	sessionId = id.String()
	session = newMemorySession(sessionId,mlt)

	s.sessionMap[sessionId] = session
	return
}

func (s *memorySessionMgr) cgc() {
	s.rwlock.Lock()
	for id,session := range s.sessionMap {
		if session.IsOvertime() {
			delete(s.sessionMap,id)
		}
	}
	time.AfterFunc(s.gcTime,s.cgc)
	s.rwlock.Unlock()
}
