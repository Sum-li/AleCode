package cRule

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/casbin/casbin"
)

var (
	adapter  *Adapter
	enforcer *casbin.Enforcer
	err      error
)

func Init(db *sql.DB) {
	adapter = NewMysqlAdapter(db, table)
	enforcer, err = casbin.NewEnforcer(model_path, adapter)
	if err != nil {
		fmt.Printf("init enforcer failed,err:%v\n", err)
	}
}

//验证是否含有某权限
func HasRule(vals ...interface{}) bool {
	ok, err := enforcer.Enforce(vals...)
	if err != nil {
		fmt.Printf("casbin has error,err:%v\n", err)
		return false
	}
	return ok
}

//加载所有权限
func LoadPolicy() error {
	return enforcer.LoadPolicy()
}

//保存所有权限
func SavePolicy() error {
	return enforcer.SavePolicy()
}

//添加权限
func AddPolicy(ptype string, rule []string) error {
	switch ptype {
	case "p":
		ok, err := enforcer.AddPolicy(rule)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("the policy is already exists")
		}
	case "g":
		ok, err := enforcer.AddGroupingPolicy(rule)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("the policy is already exists")
		}
	default:
		return errors.New("the ptype is not exists")
	}
	return nil
}

//删除权限
func RemovePolicy(ptype string, rule []string) error {
	switch ptype {
	case "p":
		ok, err := enforcer.RemovePolicy(rule)
		if err != nil {
			return fmt.Errorf("casbin remove success but db remove failed,err:%v", err)
		}
		if !ok {
			return errors.New("remove policy failed")
		}
	case "g":
		ok, err := enforcer.RemoveGroupingPolicy(rule)
		if err != nil {
			return fmt.Errorf("casbin remove success but db remove failed,err:%v", err)
		}
		if !ok {
			return errors.New("remove policy failed")
		}
	default:
		return errors.New("the ptype is not exists")
	}
	return nil
}
