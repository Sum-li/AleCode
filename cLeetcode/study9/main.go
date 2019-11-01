package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println(isPalindrome(125521))
}

func isPalindrome(x int) bool {
	var res = true
	var xs = strconv.Itoa(x)
	for i := 0; i <= len(xs)/2; i++ {
		if xs[i] != xs[len(xs)-1-i] {
			res = false
			break
		}
	}
	return res
}
