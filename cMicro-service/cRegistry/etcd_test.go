package cRegistry

import (
	"fmt"
	"testing"
	"time"
)

func TestEtcd(t *testing.T) {
	err := Init("etcd", WithAddrs([]string{"http://127.0.0.1:2379"}), WithTimeout(time.Second),
		WithHeartBeat(5), WithRegistryPath("/ale/cool/"))
	if err != nil {
		fmt.Printf("init fail,err :%v\n", err)
		return
	}
	service := &Service{
		Name: "ale1",
		Nodes: []*Node{
			{Ip: "127.0.0.1", Port: "8001"},
			{Ip: "127.0.0.2", Port: "8002"},
			{Ip: "127.0.0.3", Port: "8003"},
			{Ip: "127.0.0.4", Port: "8004"},
		},
	}
	err = Register(service)
	if err != nil {
		fmt.Printf("register service %s failed,err:%v\n", service.Name, err)
	}
	for {
		time.Sleep(time.Second * 1)
		service, err := GetService("ale1")
		if err != nil {
			fmt.Printf("get service %s failed,err:%v", service.Name, err)
			return
		}
		fmt.Printf("name:%s\n", service.Name)
		for _, v := range service.Nodes {
			fmt.Printf("http://%s:%s\n", v.Ip, v.Port)
		}
	}
}
