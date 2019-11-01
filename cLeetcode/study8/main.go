package main

import (
	"fmt"
	"math"
)

//str :int32  空格32
//			  - 45
//			  + 43
//			  0  1  2  3  4  5  6  7  8  9
//			  48 49 50 51 52 53 54 55 56 57

func main() {
	fmt.Println(myAtoi("9223372036854775808"))
}

func myAtoi(str string) int {
	var (
		resList = make([]int, len(str))
		flag    bool
		flag2   bool
		res     float64
		num     int
	)
	for i, v := range str {
		if v == 32 {
			continue
		}
		str = str[i:]
		if str[0] != 45 && str[0] != 43 && str[0] < 47 && str[0] > 58 {
			return 0
		}
		if str[0] == 45 {
			flag = true
			str = str[1:]
			break
		}
		if str[0] == 43 {
			str = str[1:]
			break
		}
		break
	}
	for i, v := range str {
		num = i
		switch v {
		case 48:
			resList[i] = 0
		case 49:
			resList[i] = 1
		case 50:
			resList[i] = 2
		case 51:
			resList[i] = 3
		case 52:
			resList[i] = 4
		case 53:
			resList[i] = 5
		case 54:
			resList[i] = 6
		case 55:
			resList[i] = 7
		case 56:
			resList[i] = 8
		case 57:
			resList[i] = 9
		default:
			num--
			flag2 = true
		}
		if flag2 {
			break
		}
	}
	if len(resList) == 0 {
		return 0
	}
	resList = resList[:num+1]
	for _, v := range resList {
		res = res*10 + float64(v)
	}
	if flag {
		res = -res
		if res < -math.Exp2(31) {
			return int(-math.Exp2(31))
		}
	}
	if res > math.Exp2(31)-1 {
		return int(math.Exp2(31) - 1)
	}
	return int(res)
}
