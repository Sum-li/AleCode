package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	fmt.Println(reverse(1534236469))
}

func reverse(x int) int {
	if float64(x) > math.Exp2(31)-1 || float64(x) < -math.Exp2(31) {
		return 0
	}
	var (
		xs   string
		flag bool
		res  int
	)
	if x < 0 {
		flag = true
		x = -x
	}
	for x/10 != 0 {
		xs += strconv.Itoa(x % 10)
		x = x / 10
	}
	xs += strconv.Itoa(x % 10)
	res, _ = strconv.Atoi(xs)
	if flag {
		res = -res
	}
	if float64(res) > math.Exp2(31)-1 || float64(res) < -math.Exp2(31) {
		return 0
	}
	return res
}
