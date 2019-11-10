package cRule

import (
	"AleCode/cDbops"
	"fmt"
	"strconv"
	"testing"
)

func TestCasbin(t *testing.T) {
	cDbops.Init()
	Init(cDbops.GetDb())
	err := LoadPolicy()
	if err != nil {
		fmt.Printf("failed,err1:%v", err)
		return
	}
	err = AddPolicy("g", []string{"ale", "user"})
	for i := 1; i < 10; i++ {
		err = AddPolicy("p", []string{"user", "test", strconv.Itoa(i)})
		if err != nil {
			fmt.Printf("failed,err2:%v", err)
			return
		}
	}
	for i := 1; i < 10; i++ {
		fmt.Printf("result:%v\n", HasRule("ale", "test", strconv.Itoa(i)))
	}
	for i := 1; i < 10; i++ {
		if i%2 == 0 {
			err = RemovePolicy("p", []string{"user", "test", strconv.Itoa(i)})
			if err != nil {
				fmt.Printf("failed,err3:%v", err)
				return
			}
		}
	}
	for i := 1; i < 10; i++ {
		fmt.Printf("result:%v\n", HasRule("ale", "test", strconv.Itoa(i)))
	}
}
