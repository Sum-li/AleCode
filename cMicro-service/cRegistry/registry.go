package cRegistry

import "context"

type Registry interface {
	//插件的名字
	Name() string
	//初始化插件
	Init(context.Context, ...Option) error
	//服务注册
	Register(context.Context, *Service) error
	//服务注销
	UnRegister(context.Context, *Service) error
	//获取服务
	GetService(context.Context, string) (*Service, error)
}
