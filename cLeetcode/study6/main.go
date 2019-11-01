package main

import "fmt"

func main() {
	fmt.Println(convert("AB", 1))
}

func convert(s string, numRows int) string {
	var (
		zmap  = make(map[int]string, numRows)
		index int
		flag  bool
		res   string
	)
	for _, v := range s {
		zmap[index] += string(v)
		if index == numRows-1 {
			flag = true
		}
		if index == 0 {
			flag = false
		}
		if flag {
			index--
		} else {
			index++
		}
	}
	for i := 0; i < len(zmap); i++ {
		res += zmap[i]
	}
	return res
}
