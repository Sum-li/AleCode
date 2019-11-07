package cRegistry

import (
	"context"
	"fmt"
)

////单例模式
//var registrymgr *registryMgr
//
//func init() {
//	registrymgr = newRegistryMgr()
//}
//
////添加注册中兴
//func AddRegistry(registry Registry) (err error) {
//	return registrymgr.addRegistry(registry)
//}
//
////初始化注册中心
//func InitRegistry(ctx context.Context, name string, opts ...Option) (Registry, error) {
//	return registrymgr.initRegistry(ctx, name, opts...)
//}

//单例
var registry Registry

func Init(name string, opts ...Option) error {
	switch name {
	case "etcd":
		err := initEtcd(opts...)
		return err
	default:
		return fmt.Errorf("registry named %s is not exist", name)
	}
}

func initEtcd(opts ...Option) error {
	registry = NewEtcdRegistry()
	err := registry.Init(context.TODO(), opts...)
	return err
}

func Name() string {
	return registry.Name()
}

func Register(service *Service) error {
	return registry.Register(context.TODO(), service)
}

func UnRegister(service *Service) error {
	return registry.UnRegister(context.TODO(), service)
}

func GetService(name string) (*Service, error) {
	return registry.GetService(context.TODO(), name)
}
