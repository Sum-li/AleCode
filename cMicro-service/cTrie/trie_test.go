package cTrie

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestTrie(t *testing.T) {
	var (
		trie = NewTrie()
		file *os.File
		err  error
		rd   *bufio.Reader
	)
	file, err = os.Open("C:/Sumli/Project/Go/src/AleCode/cMicro-service/cTrie/testdata.txt")
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	defer file.Close()
	rd = bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("err:%v", err)
		}
		line = strings.Replace(line, "\n", "", -1)
		err = trie.Add(line, line)
		if err != nil {
			fmt.Printf("err:%v", err)
		}
	}
	res, _ := trie.Check("我喜欢排卵", "***")
	fmt.Printf("res:%v", res)
}
