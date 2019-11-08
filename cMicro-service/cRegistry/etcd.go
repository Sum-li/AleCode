package cRegistry

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

//key = RegistryPath + / + name

//etcd注册中心
type etcdRegistry struct {
	//注册中心的配置
	config      *config
	client      *clientv3.Client
	serviceChan chan *Service
	serviceMap  map[string]*etcdService
	//原子性数据结构，用于服务发现时在本地存储服务信息
	value atomic.Value
	lock  sync.Mutex
}

func NewEtcdRegistry() *etcdRegistry {
	return &etcdRegistry{
		serviceChan: make(chan *Service, MaxServiceNum),
		serviceMap:  make(map[string]*etcdService, MaxServiceNum),
	}
}

type etcdService struct {
	registered bool
	//用于保活
	leaseId clientv3.LeaseID
	Service *Service
}

type allService map[string]*Service

//注册中心的名字
func (e *etcdRegistry) Name() string {
	return "etcd"
}

//初始化注册中心
func (e *etcdRegistry) Init(ctx context.Context, opts ...Option) (err error) {
	var all_service allService = make(map[string]*Service)
	e.value.Store(&all_service)
	e.config = &config{}
	for _, opt := range opts {
		opt(e.config)
	}
	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.config.Addrs,
		DialTimeout: e.config.Timeout,
	})
	if err != nil {
		err = fmt.Errorf("init etcd failed,err:%v", err)
	}
	go e.run()
	return
}

//服务注册
func (e *etcdRegistry) Register(ctx context.Context, service *Service) (err error) {
	select {
	case e.serviceChan <- service:
	default:
		err = fmt.Errorf("the service quantity reaches the maximum")
		return
	}
	return
}

//服务注销
func (e *etcdRegistry) UnRegister(ctx context.Context, service *Service) (err error) {
	e.lock.Lock()
	delete(e.serviceMap, service.Name)
	e.lock.Unlock()
	return
}

//获取服务
func (e *etcdRegistry) GetService(ctx context.Context, name string) (service *Service, err error) {
	//现在缓存中找有没有
	service, ok := e.getServiceLocal(ctx, name)
	if ok {
		return
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	service, ok = e.getServiceLocal(ctx, name)
	if ok {
		return
	}
	key := path.Join(e.config.RegistryPath, name)
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		fmt.Printf("get service from etcd failed,err:%v", err)
		return
	}
	if len(resp.Kvs) != 1 {
		err = fmt.Errorf("service named:%s is not exist", name)
		return
	}
	var serviceimpl Service
	err = json.Unmarshal(resp.Kvs[0].Value, &serviceimpl)
	if err != nil {
		fmt.Printf("json unmaeshl failed,data:%v,err:%v", string(resp.Kvs[0].Value), err)
		return
	}
	service = &serviceimpl
	all_service := *e.value.Load().(*allService)
	all_service[name] = service
	e.value.Store(&all_service)
	return
}

//将服务注册到etcd中去
func (e *etcdRegistry) register(service *etcdService) (err error) {
	resp, err := e.client.Grant(context.TODO(), e.config.HeartBeat)
	if err != nil {
		return fmt.Errorf("client etcd failed,err:%v", err)
	}
	data, err := json.Marshal(service.Service)
	if err != nil {
		return fmt.Errorf("marshal to json failed,err:%v", err)
	}
	key := path.Join(e.config.RegistryPath, service.Service.Name)
	_, err = e.client.Put(context.TODO(), key, string(data), clientv3.WithLease(resp.ID))
	if err != nil {
		return fmt.Errorf("load into etcd failed,err:%v", err)
	}
	_, err = e.client.KeepAliveOnce(context.TODO(), resp.ID)
	if err != nil {
		return fmt.Errorf("etcd keepalive failed,err:%v", err)
	}
	service.registered = true
	service.leaseId = resp.ID
	//加入本地缓存
	all_service := *e.value.Load().(*allService)
	all_service[service.Service.Name] = service.Service
	e.value.Store(&all_service)
	return
}

//etcd服务保活
func (e *etcdRegistry) keepAlive(service *etcdService) {
	_, err := e.client.KeepAliveOnce(context.TODO(), service.leaseId)
	if err != nil {
		service.registered = false
	}
	return
}

func (e *etcdRegistry) registerOrKeepAlive() {
	if len(e.serviceMap) == 0 {
		time.Sleep(time.Microsecond * 3000)
		return
	}
	for _, service := range e.serviceMap {
		if service.registered {
			e.keepAlive(service)
			continue
		}
		//将服务注册到etcd中
		err := e.register(service)
		if err != nil {
			fmt.Printf("register into etcd failed,err:%v", err)
			service.registered = false
			continue
		}
	}
}

//后台进程将服务注册到etcd中
func (e *etcdRegistry) run() {
	ticker := time.NewTicker(MaxSyncServiceInterval)
	for {
		select {
		case service := <-e.serviceChan:
			etcd_service, ok := e.serviceMap[service.Name]
			if ok {
				etcd_service.Service = service
				etcd_service.registered = false
				break
			}
			etcd_service = &etcdService{
				registered: false,
				Service:    service,
			}
			e.serviceMap[service.Name] = etcd_service
		case <-ticker.C:
			e.loadFromEtcd()
		default:
			e.registerOrKeepAlive()
		}
	}
}

//从etcd中拉取服务
func (e *etcdRegistry) loadFromEtcd() {
	serviceMap := *e.value.Load().(*allService)
	for _, service := range serviceMap {
		var newService Service
		key := path.Join(e.config.RegistryPath, service.Name)
		resp, err := e.client.Get(context.TODO(), key)
		if err != nil {
			fmt.Printf("load from etcd failed,err:%v", err)
			continue
		}
		if len(resp.Kvs) != 1 {
			err := e.UnRegister(context.TODO(), service)
			fmt.Printf("service named:%s is not exist,err:%v\n", service.Name, err)
			continue
		}
		err = json.Unmarshal(resp.Kvs[0].Value, &newService)
		if err != nil {
			fmt.Printf("json unmaeshl failed,data:%v,err:%v", string(resp.Kvs[0].Value), err)
			return
		}
		serviceMap[service.Name] = &newService
	}
	e.value.Store(&serviceMap)
}

//从本地获得服务
func (e *etcdRegistry) getServiceLocal(ctx context.Context, name string) (service *Service, ok bool) {
	serviceMap := *e.value.Load().(*allService)
	service, ok = serviceMap[name]
	return
}
