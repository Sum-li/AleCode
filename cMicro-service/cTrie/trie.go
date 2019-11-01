package cTrie

import "strings"

type Node struct {
	//当前节点对应的字符
	char rune
	//对应key下的数据
	Data   interface{}
	parent *Node
	////当前节点的深度
	//Depth  int
	//当前节点的所有子节点
	childs map[rune]*Node
	//当前节点是否为有效的非中间节点
	term bool
}

type Trie struct {
	//树的头节点
	root *Node
	////树的深度
	//size int
}

func NewNode() *Node {
	return &Node{
		childs: make(map[rune]*Node, 32),
	}
}

func NewTrie() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

func (t *Trie) Add(key string, data interface{}) (err error) {
	key = strings.TrimSpace(key)
	var (
		runes = []rune(key)
		node  = t.root
	)
	for _, v := range runes {
		ret, ok := node.childs[v]
		if !ok {
			ret = NewNode()
			ret.char = v
			ret.parent = node
			//将节点添加到childs中，不然检查不到
			node.childs[v] = ret
			//ret.Depth = node.Depth + 1
		}
		node = ret
	}
	node.term = true
	node.Data = data
	return
}

func (t *Trie) findNode(key string) (result *Node, ok bool) {
	var (
		node  = t.root
		ret   *Node
		chars = []rune(key)
	)
	for _, v := range chars {
		ret, ok = node.childs[v]
		if !ok {
			return
		}
		node = ret
	}
	result = node
	ok = true
	return
}

//获取当前节点下的叶子节点
func (t *Trie) collectNode(node *Node) (result []*Node) {
	if node == nil {
		return
	}
	if node.term {
		result = append(result, node)
		return
	}
	var queue []*Node
	queue = append(queue, node)
	for i := 0; i < len(queue); i++ {
		//将所有的叶子节点取出
		if queue[i].term {
			result = append(result, queue[i])
			continue
		}
		//将所有的节点都存在queue中
		for _, v := range queue[i].childs {
			queue = append(queue, v)
		}
	}
	return
}

//搜索前缀相同的node
func (t *Trie) PrefixSearch(key string) (result []*Node) {
	node, ok := t.findNode(key)
	if !ok {
		return
	}
	result = t.collectNode(node)
	return
}

func (t *Trie) Check(text, replace string) (result string, hit bool) {
	if t.root == nil {
		return
	}
	var (
		chars = []rune(text)
		left  []rune
		node  = t.root
		start int
	)
	for i, v := range chars {
		ret, ok := node.childs[v]
		if !ok {
			left = append(left, chars[start:i+1]...)
			start = i + 1
			node = t.root
			continue
		}
		node = ret
		if ret.term {
			hit = true
			node = t.root
			left = append(left, []rune(replace)...)
			start = i + 1
			continue
		}
	}
	result = string(left)
	return
}
