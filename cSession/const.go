package cSession

import "time"

const (
	//session未被加载
	sessionFlagNone = iota
	//session被更改
	sessionFlagModify
	//session默认过期时间 两小时
	sessionLiftTime = time.Hour*2
)