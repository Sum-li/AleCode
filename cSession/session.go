package cSession

type Session interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Save() error
	IsModify() bool
	Id() string
	IsOvertime() bool
	Update()
	//todo:redis 的同步操作
	//todo:redis 过期处理，即redis的清除操作
	/*
	todo: 最好对redis的部分进行重构
	 */
}
