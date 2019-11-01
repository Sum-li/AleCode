package cSession

import (
	"AleCode/clog"
	"errors"
	"sync"
	"time"
)

type memorySession struct {
	//session中存的数据
	data map[string]interface{}
	//session的标识
	id string
	//操作标记（用于标记session被操作的状态）
	flag int
	//session的过期时间
	maxLifeTime time.Duration
	//最后访问时间
	lastAccessedTime time.Time
	rwlock sync.RWMutex
}

func newMemorySession(id string,mlt time.Duration) Session {
	return &memorySession{
		id: id,
		data: make(map[string]interface{},8),
		lastAccessedTime:time.Now(),
		maxLifeTime:mlt,
	}
}

//获取session的id
func (m *memorySession) Id() string {
	return m.id
}

//获取session是否被更改
func (m *memorySession) IsModify() bool {
	if m.flag == sessionFlagModify {
		return true
	}
	return false
}

func (m *memorySession) Set(key string,value interface{}) (err error) {
	m.rwlock.Lock()
	m.data[key] = value
	//标记session被更改
	m.flag = sessionFlagModify
	m.rwlock.Unlock()
	return
}

func (m *memorySession) Get(key string) (value interface{},err error) {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	value,ok := m.data[key]
	if !ok {
		err = errors.New("the data not exit")
		clog.Error("err:%v",err)
		return
	}
	return
}

func (m *memorySession) Del(key string) (err error) {
	m.rwlock.Lock()
	//标记session被更改
	m.flag = sessionFlagModify
	delete(m.data,key)
	m.rwlock.Unlock()
	return
}

func (m *memorySession) Save() (err error) {
	return
}

func (m *memorySession) IsOvertime() (overtime bool) {
	if time.Now().Unix()-m.lastAccessedTime.Unix() > int64(m.maxLifeTime) {
		overtime = true
	}
	return
}

func (m *memorySession) Update() {
	m.rwlock.Lock()
	m.lastAccessedTime = time.Now()
	m.rwlock.Unlock()
	return
}