package cValidator

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

const (
	url = `^((https|http|ftp|rtsp|mms)?:\/\/)[^\s]+`
	ip  = `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|
			[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`
	email = `\w[-\w.+]*@([A-Za-z0-9][-A-Za-z0-9]+\.)+[A-Za-z]{2,14}`
	tel   = `0?(13|14|15|17|18|19)[0-9]{9}`
)

/**
利用反射和正则表达式对 结构体 进行参数校验

@author：Ale
@date：19.10.31

思路：通过tag获得要调用的检测函数，并通过反射获取字段并检验
检验：通过正则表达式对字符串进行检验，数值型的自己通过函数检验

约定：tag： 函数为 fun 参数为 data
*/

type cValidator struct {
}

func NewCValidator() *cValidator {
	return &cValidator{}
}

func (v *cValidator) Check(data interface{}) error {
	var (
		validator, value, field, method reflect.Value
		ctype                           reflect.Type
		sField                          reflect.StructField
		funName, param                  string
		ok                              bool
	)
	value = reflect.ValueOf(data).Elem()
	ctype = value.Type()

	//如果不是结构体，则不进行校验
	if ctype.Kind() != reflect.Struct {
		return errors.New("请传入一个结构体的指针")
	}
	//获取结构体中字段的tag信息，并进行验证
	for i := 0; i < value.NumField(); i++ {
		field = value.Field(i)
		sField = ctype.Field(i)
		funName, ok = sField.Tag.Lookup("fun")
		//如果存在，说明需要校验
		if ok {
			validator = reflect.ValueOf(v)
			var params []reflect.Value
			dataVal := field.Interface()
			params = append(params, reflect.ValueOf(dataVal))
			param, ok = sField.Tag.Lookup("data")
			if ok {
				paramlist := strings.Split(param, ",")
				for _, v := range paramlist {
					params = append(params, reflect.ValueOf(v))
				}
				method = validator.MethodByName(funName)
				//执行校验函数，并获取执行结果
				res := method.Call(params)
				pass := res[0].Bool()
				if !pass {
					msg := res[1].String()
					return errors.New(msg)
				}
			}
		}
	}
	return nil
}

func (v *cValidator) Str(data, pattern string) (bool, string) {
	var rex string
	switch pattern {
	case "url":
		rex = url
	case "ip":
		rex = ip
	case "tel":
		rex = tel
	case "email":
		rex = email
	default:
		matched, err := regexp.MatchString(pattern, data)
		if err != nil {
			return false, err.Error()
		}
		if matched {
			return true, ""
		} else {
			return false, "参数格式不正确"
		}
	}
	matched, err := regexp.MatchString(rex, data)
	if err != nil {
		return false, err.Error()
	}
	if matched {
		return true, ""
	} else {
		return false, pattern + "的格式错误"
	}
}

//func (v *cValidator) Int(data,x,y interface{}) (bool,string) {
//	intx,err := strconv.Atoi(x.(string))
//	if err != nil {
//		return false,"tag参数错误"
//	}
//	inty,err := strconv.Atoi(y.(string))
//	if err != nil {
//		return false,"tag参数错误"
//	}
//	if data.(int)>intx || data.(int)<inty {
//		return false,"值超出范围"
//	}
//	return true,""
//}
