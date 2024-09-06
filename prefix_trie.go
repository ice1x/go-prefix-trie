package prefixtrie

import (
	"sort"
	"strconv"
	"sync"
)

type Node struct {
	Char     rune
	Data     []map[string]string
	Children map[rune]*Node
	mu       sync.RWMutex
}

func NewNode(char rune) *Node {
	return &Node{
		Char:     char,
		Data:     []map[string]string{},
		Children: make(map[rune]*Node),
	}
}

type Trie struct {
	Root *Node
}

func NewTrie() *Trie {
	return &Trie{
		Root: NewNode(0),
	}
}

func (t *Trie) Insert(word string, data map[string]string) {
	node := t.Root
	for _, char := range word {
		node.mu.Lock()
		if _, exists := node.Children[char]; !exists {
			node.Children[char] = NewNode(char)
		}
		node.mu.Unlock()
		node = node.Children[char]
	}
	node.mu.Lock()
	node.Data = append(node.Data, data)
	node.mu.Unlock()
}

//func (t *Trie) GetByPrefixSortDescBy(prefix, key string) []map[string]string {
//	node := t.Root
//	for _, char := range prefix {
//		node.mu.RLock()
//		child, exists := node.Children[char]
//		node.mu.RUnlock()
//		if !exists {
//			return nil
//		}
//		node = child
//	}
//
//	result := t.collectDataByChild(node)
//	sort.Slice(result, func(i, j int) bool {
//		return result[i][key] > result[j][key]
//	})
//	return result
//}

func (t *Trie) GetByPrefixSortDescBy(prefix, key string) []map[string]string {
	node := t.Root
	for _, char := range prefix {
		node.mu.RLock()
		child, exists := node.Children[char]
		node.mu.RUnlock()
		if !exists {
			return nil
		}
		node = child
	}

	result := t.collectDataByChild(node)
	sort.Slice(result, func(i, j int) bool {
		// Convert the values to integers for sorting
		iVal, err1 := strconv.Atoi(result[i][key])
		jVal, err2 := strconv.Atoi(result[j][key])
		if err1 != nil || err2 != nil {
			// Fall back to string comparison if conversion fails
			return result[i][key] > result[j][key]
		}
		return iVal > jVal
	})
	return result
}

func (t *Trie) collectDataByChild(node *Node) []map[string]string {
	node.mu.RLock()
	defer node.mu.RUnlock()

	result := append([]map[string]string(nil), node.Data...)
	for _, child := range node.Children {
		result = append(result, t.collectDataByChild(child)...)
	}
	return result
}
