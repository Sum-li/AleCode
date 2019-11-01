package cRedis

import (
	"AleCode/clog"
	"github.com/garyburd/redigo/redis"
	"time"
)

var pool *redis.Pool

func newPool(server,password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 64,
		MaxActive: 1000,
		IdleTimeout: 240*time.Second,
		Dial: func() (conn redis.Conn, err error) {
			conn,err = redis.Dial("tcp",server)
			if err !=nil {
				clog.Error("get redis conn failed,err:%v",err)
				return
			}
			if password != "" {
				if _,err = conn.Do("AUTH",password); err != nil {
					clog.Error("redis set password failed,err:%v",err)
					_ = conn.Close()
					return
				}
			}
			return
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_,err := conn.Do("PING")
			clog.Error("redis ping timeout,err:%v",err)
			return err
		},
	}
}

func init() {
	_ = clog.InitLogger("console",clog.Config{Log_level:"debug"})
	pool = newPool("10.14.4.240:6379","")
	clog.Debug("redis pool stats:%#v",pool.Stats())
}
func GetPool() *redis.Pool {
	return pool
}