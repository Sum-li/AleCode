package cValidator

import (
	"fmt"
	"testing"
)

func TestCValidator_Check(t *testing.T) {
	type ss struct {
		Data string `fun:"Str" data:"url"`
	}
	cv := NewCValidator()
	err := cv.Check(&ss{Data: "htt"})
	if err != nil {
		fmt.Println(err)
	}
}
