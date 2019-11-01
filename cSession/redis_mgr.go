package cSession

import (
	"AleCode/clog"
	"errors"
	"github.com/garyburd/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"sync"
	"time"
)

type redisSessionMgr struct {
	//redis 的连接池
	pool *redis.Pool
	//存session
	sessionMap map[string]Session
	rwlock sync.RWMutex
	//释放session垃圾的周期
	gcTime time.Duration
}

func newRedisSessionMgr(gct time.Duration) SessionMgr {
	return &redisSessionMgr{gcTime:gct}
}

func (r *redisSessionMgr) Init(options ...interface{}) (err error) {
	r.sessionMap = make(map[string]Session,1024)
	if len(options) > 0 {
		r.pool = options[0].(*redis.Pool)
	}
	go r.cgc()
	return
}

func (r *redisSessionMgr) Create(mlt time.Duration) (session Session,err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	var(
		id uuid.UUID
		sessionId string
	)
	id,err = uuid.NewV4()
	if err != nil {
		clog.Error("get uuid failed,err:%v",err)
		return
	}
	sessionId = id.String()
	session = newRedisSession(sessionId,mlt,r.pool)

	r.sessionMap[sessionId] = session
	return
}

func (r *redisSessionMgr) Get(sessionId string) (session Session,err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	session,ok := r.sessionMap[sessionId]
	if !ok {
		session = newRedisSession(sessionId,sessionLiftTime,r.pool)
		r.sessionMap[sessionId] = session
		return
	}
	return
}

func (r *redisSessionMgr) cgc() {
	r.rwlock.Lock()
	for id,session := range r.sessionMap{
		if session.IsOvertime() {
			delete(r.sessionMap,id)
		}
	}
	time.AfterFunc(r.gcTime,r.cgc)
	r.rwlock.Unlock()
}

func (r *redisSessionMgr) Del(sessionId string) {
	r.rwlock.Lock()
	delete(r.sessionMap,sessionId)
	r.rwlock.Unlock()
	return
}

func (r *redisSessionMgr) Update(sessionId string) (err error){
	r.rwlock.Lock()
	defer r.rwlock.Unlock()
	session,ok := r.sessionMap[sessionId]
	if !ok {
		err = errors.New("session not exit")
		clog.Error("err:%v",err)
		return
	}
	session.Update()
	return
}