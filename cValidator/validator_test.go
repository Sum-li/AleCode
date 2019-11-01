package cValidator

import (
	"fmt"
	"testing"
)

func TestCValidator_Check(t *testing.T) {
	type ss struct {
		Data string `fun:"Str" data:"email"`
	}
	cv := NewCValidator()
	err := cv.Check(&ss{Data: "123@qq.com"})
	if err != nil {
		fmt.Println(err)
	}
}
