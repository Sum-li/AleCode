package cSession

import (
	"AleCode/clog"
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

type redisSession struct {
	id string
	pool *redis.Pool
	data map[string]interface{}
	rwlock sync.RWMutex
	flag int
	//session的过期时间
	maxLifeTime time.Duration
	//最后访问时间
	lastAccessedTime time.Time
}

func newRedisSession(id string,mlt time.Duration,pool *redis.Pool) Session {
	return &redisSession{
		id:id,
		pool:pool,
		data:make(map[string]interface{}),
		flag:sessionFlagNone,
		maxLifeTime:mlt,
		lastAccessedTime:time.Now(),
	}
}

func (r *redisSession) Set(key string,value interface{}) (err error) {
	r.rwlock.Lock()
	r.data[key] = value
	r.flag = sessionFlagModify
	r.rwlock.Unlock()
	return
}

func (r *redisSession) loadFromRedis() (err error) {
	conn := r.pool.Get()
	res,err := conn.Do("GET",r.id)
	if err !=nil {
		clog.Error("get from redis failed,err:%v",err)
		return
	}
	data,err := redis.String(res,err)
	if err !=nil {
		clog.Error("get from redis failed,err:%v",err)
		return
	}
	err = json.Unmarshal([]byte(data),&r.data)
	if err !=nil {
		clog.Error("data unmarshal failed,err:%v",err)
		return
	}
	return
}

func (r *redisSession) Get(key string) (res interface{},err error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	if r.flag == sessionFlagNone {
		err = r.loadFromRedis()
		if err != nil {
			clog.Error("load from redis failed,err:%v",err)
			return
		}
	}
	res,ok := r.data[key]
	if !ok {
		err = errors.New("the data not exit")
		clog.Error("err:%v",err)
		return
	}
	return
}

func (r *redisSession) Id() string {
	return r.id
}

func (r *redisSession) Del(key string) (err error) {
	r.rwlock.Lock()
	r.flag = sessionFlagModify
	delete(r.data,key)
	r.rwlock.Unlock()
	return
}

func (r *redisSession) Save() (err error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	if r.flag != sessionFlagModify {
		return
	}

	data,err := json.Marshal(r.data)
	if err != nil {
		clog.Error("data marshal failed,err:%v",err)
		return
	}
	conn := r.pool.Get()
	_,err = conn.Do("SET",r.id,string(data))
	if err != nil {
		clog.Error("insert into redis failed,err:%v",err)
		return
	}
	return
}

func (r *redisSession) IsModify() bool {
	if r.flag == sessionFlagModify {
		return true
	}
	return false
}

func (r *redisSession) IsOvertime() (overtime bool) {
	if time.Now().Unix()-r.lastAccessedTime.Unix() > int64(r.maxLifeTime) {
		overtime = true
	}
	return
}

func (r *redisSession) Update() {
	r.rwlock.Lock()
	r.lastAccessedTime = time.Now()
	r.rwlock.Unlock()
	return
}