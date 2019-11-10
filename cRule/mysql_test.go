package cRule

import (
	"AleCode/cDbops"
	"fmt"
	"testing"
)

var adapter_test *Adapter

//测试前的初始化操作
func TestMain(m *testing.M) {
	//测试前操作
	cDbops.Init()
	adapter_test = NewMysqlAdapter(cDbops.GetDb(), "cRules")
	m.Run()
	//测试后操作
}

//总测试
func TestRun(t *testing.T) {
	t.Run("AddPolicy", testAddPolicy)
	t.Run("RemovePolicy", testRemovePolicy)
	t.Run("SelectFields", testSelectFields)
	t.Run("Truncate", testTruncate)
}

func testAddPolicy(t *testing.T) {
	err := adapter_test.AddPolicy("", "p", []string{"ale", "user", "bbb"})
	if err != nil {
		fmt.Println(err)
	}
}

func testRemovePolicy(t *testing.T) {
	err := adapter_test.RemovePolicy("", "p", []string{"user", "test", "AAA"})
	if err != nil {
		fmt.Println(err)
	}
}

func testSelectFields(t *testing.T) {
	var s []*Rule
	err := adapter_test.selectFields(&s)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range s {
		fmt.Printf("%#v\n", *v)
	}
}

func testTruncate(t *testing.T) {
	_, err := adapter_test.db.Exec("truncate " + adapter_test.table)
	if err != nil {
		fmt.Println(err)
	}
}
