package main

import "fmt"

func main() {
	fmt.Printf("%v", lengthOfLongestSubstring("bbbbb"))
}

//自己的
func lengthOfLongestSubstring(s string) int {
	head, last, result := 0, 0, 0
	var m = make(map[uint8]int, 256)
	for last < len(s) {
		val, ok := m[s[last]]
		if !ok {
			m[s[last]] = last
			if last == len(s)-1 {
				if result < last-head+1 {
					return last - head + 1
				}
			}
		} else {
			if val >= head {
				if result < last-head {
					result = last - head
				}
				head = val + 1
				m[s[last]] = last
			} else {
				m[s[last]] = last
				if last == len(s)-1 {
					if result < last-head+1 {
						return last - head + 1
					}
				}
			}
		}
		last++
	}

	return result
}

//大佬的
func lengthOfLongestSubstringNice(s string) int {
	location := [256]int{}
	for i := range location {
		location[i] = -1
	}

	j, maxLen := 0, 0
	for i := 0; j < len(s); j++ {
		if location[s[j]] != -1 && location[s[j]] >= i {
			i = location[s[j]] + 1
		}
		location[s[j]] = j
		maxLen = max(j-i+1, maxLen)
	}

	return maxLen
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
