package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Printf("%v", longestPalindromeNice("babad"))
}

func longestPalindrome(s string) string {
	str := make([]int32, 2*len(s)+1)
	var max int = 0
	var result []int32
	for i, v := range s {
		str[2*i+1] = v
	}
	for i := range str {
		j := 1
		for j <= i && j < len(str)-i && str[i+j] == str[i-j] {
			j++
		}
		if j > max {
			max = j
			result = str[i-j+1 : i+j]
		}
	}
	return strings.Replace(string(result), "\x00", "", -1)
}

func longestPalindromeNice(s string) string {
	if len(s) <= 1 {
		return s
	}

	ms := []byte{}   // 用于存储转换后的 Manacher 字符串的 slice
	radii := []int{} // 记录 Manacher 字符串每个索引的回文半径

	// 转换成 Manacher 字符串并由ms存储
	for i := 0; i < 2*len(s)+1; i++ {
		if i%2 == 0 {
			ms = append(ms, '#')
		} else {
			ms = append(ms, s[i/2])
		}
	}

	// mid: 回文中心 maxRight: 回文中心对应的最右回文索引 maxMid: 记录最大回文中心的下标
	mid, maxRight, maxMid := 0, 0, 0

	for i := 0; i < len(ms); i++ {
		r := 1 // 最小回文半径

		if i < maxRight { // 从已知条件中获取最大回文半径
			if radii[2*mid-i] < maxRight-i+1 { // 该种情况理论上是不用扩散的,因为i关于pos对称的左侧另一部分回文并没有超出回文mid的回文半径
				r = radii[2*mid-i]
			} else {
				r = maxRight - i + 1 // i关于pos对称的左侧另一部分回文超出了mid的回文半径,因此i的回文半径还需要进一步通过扩散来判断
			}
		}

		for i-r >= 0 && i+r < len(ms) && ms[i-r] == ms[i+r] { // 通过扩散来进一步判断回文半径
			r++
		}
		radii = append(radii, r)

		if radii[i] > radii[maxMid] { // 记录最大回文半径中心
			maxMid = i
		}
		if i+radii[i]-1 > maxRight { // 更新回文中心及其回文最右索引
			maxRight = i + radii[i] - 1
			mid = i
		}
	}

	if radii[maxMid]%2 == 0 { // 根据回文中心和回文半径确定回文字符串
		return s[maxMid/2-(radii[maxMid]-1)/2 : maxMid/2+(radii[maxMid]-1)/2+1]
	}
	return s[maxMid/2-(radii[maxMid]-1)/2 : maxMid/2+(radii[maxMid]-1)/2]
}
