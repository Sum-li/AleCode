package main

import "testing"

func BenchmarkLongestPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		longestPalindrome("babad")
	}
}
