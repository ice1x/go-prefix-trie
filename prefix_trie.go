// Package prefixtrie provides a thread-safe prefix trie (prefix tree) implementation
// for efficient string prefix matching and retrieval with flexible sorting options.
package prefixtrie

import (
	"errors"
	"sort"
	"strconv"
	"sync"
)

var (
	// ErrEmptyKey is returned when an empty key is provided to Insert, Delete, Search, or Update
	ErrEmptyKey = errors.New("key cannot be empty")
	// ErrNilData is returned when nil data is provided to Insert or Update
	ErrNilData = errors.New("data cannot be nil")
	// ErrNotFound is returned when a key is not found in Delete, Search, or Update operations
	ErrNotFound = errors.New("key not found")
)

// SortOrder defines the sorting direction
type SortOrder int

const (
	// Ascending sort order
	Ascending SortOrder = iota
	// Descending sort order
	Descending
)

// Node represents a single node in the prefix trie
type Node struct {
	Char     rune
	Data     []map[string]string
	Children map[rune]*Node
	mu       sync.RWMutex
	isEnd    bool // marks if this node represents the end of a word
}

// NewNode creates a new node with the given character
func NewNode(char rune) *Node {
	return &Node{
		Char:     char,
		Data:     []map[string]string{},
		Children: make(map[rune]*Node),
		isEnd:    false,
	}
}

// Trie represents the prefix trie data structure
type Trie struct {
	Root *Node
	mu   sync.RWMutex
}

// NewTrie creates and initializes a new Trie
func NewTrie() *Trie {
	return &Trie{
		Root: NewNode(0),
	}
}

// Insert adds a word with associated data to the trie
// Returns an error if the word is empty or data is nil
func (t *Trie) Insert(word string, data map[string]string) error {
	if word == "" {
		return ErrEmptyKey
	}
	if data == nil {
		return ErrNilData
	}

	node := t.Root
	for _, char := range word {
		node.mu.Lock()
		if _, exists := node.Children[char]; !exists {
			node.Children[char] = NewNode(char)
		}
		nextNode := node.Children[char]
		node.mu.Unlock()
		node = nextNode
	}

	node.mu.Lock()
	node.Data = append(node.Data, data)
	node.isEnd = true
	node.mu.Unlock()

	return nil
}

// Search finds exact matches for a given word
// Returns all data associated with the word, or ErrNotFound if the word doesn't exist
func (t *Trie) Search(word string) ([]map[string]string, error) {
	if word == "" {
		return nil, ErrEmptyKey
	}

	node := t.Root
	for _, char := range word {
		node.mu.RLock()
		child, exists := node.Children[char]
		node.mu.RUnlock()
		if !exists {
			return nil, ErrNotFound
		}
		node = child
	}

	node.mu.RLock()
	defer node.mu.RUnlock()

	if !node.isEnd || len(node.Data) == 0 {
		return nil, ErrNotFound
	}

	// Return a copy to prevent external modifications
	result := make([]map[string]string, len(node.Data))
	copy(result, node.Data)
	return result, nil
}

// HasPrefix checks if any words in the trie start with the given prefix
func (t *Trie) HasPrefix(prefix string) bool {
	if prefix == "" {
		return true // empty prefix matches everything
	}

	node := t.Root
	for _, char := range prefix {
		node.mu.RLock()
		child, exists := node.Children[char]
		node.mu.RUnlock()
		if !exists {
			return false
		}
		node = child
	}
	return true
}

// Delete removes all data associated with a specific word
// Returns ErrNotFound if the word doesn't exist
func (t *Trie) Delete(word string) error {
	if word == "" {
		return ErrEmptyKey
	}

	node := t.Root
	for _, char := range word {
		node.mu.RLock()
		child, exists := node.Children[char]
		node.mu.RUnlock()
		if !exists {
			return ErrNotFound
		}
		node = child
	}

	node.mu.Lock()
	defer node.mu.Unlock()

	if !node.isEnd {
		return ErrNotFound
	}

	node.Data = []map[string]string{}
	node.isEnd = false
	return nil
}

// Update replaces all data for a given word with new data
// Returns ErrNotFound if the word doesn't exist
func (t *Trie) Update(word string, data map[string]string) error {
	if word == "" {
		return ErrEmptyKey
	}
	if data == nil {
		return ErrNilData
	}

	node := t.Root
	for _, char := range word {
		node.mu.RLock()
		child, exists := node.Children[char]
		node.mu.RUnlock()
		if !exists {
			return ErrNotFound
		}
		node = child
	}

	node.mu.Lock()
	defer node.mu.Unlock()

	if !node.isEnd {
		return ErrNotFound
	}

	node.Data = []map[string]string{data}
	return nil
}

// GetByPrefix retrieves all data for words starting with the given prefix
// Returns nil if no matches are found
func (t *Trie) GetByPrefix(prefix string) []map[string]string {
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

	return t.collectDataByChild(node)
}

// GetByPrefixSortBy retrieves all data for words starting with the given prefix
// and sorts them by the specified key in the given order
// Returns nil if no matches are found
func (t *Trie) GetByPrefixSortBy(prefix, key string, order SortOrder) []map[string]string {
	result := t.GetByPrefix(prefix)
	if result == nil {
		return nil
	}

	t.sortResults(result, key, order)
	return result
}

// GetByPrefixSortDescBy retrieves all data for words starting with the given prefix
// and sorts them by the specified key in descending order
// Deprecated: Use GetByPrefixSortBy with Descending order instead
// Returns nil if no matches are found
func (t *Trie) GetByPrefixSortDescBy(prefix, key string) []map[string]string {
	return t.GetByPrefixSortBy(prefix, key, Descending)
}

// sortResults sorts the results by the given key and order
func (t *Trie) sortResults(result []map[string]string, key string, order SortOrder) {
	sort.Slice(result, func(i, j int) bool {
		// Try to convert values to integers for numeric sorting
		iVal, err1 := strconv.Atoi(result[i][key])
		jVal, err2 := strconv.Atoi(result[j][key])

		var comparison bool
		if err1 != nil || err2 != nil {
			// Fall back to string comparison if conversion fails
			if order == Ascending {
				comparison = result[i][key] < result[j][key]
			} else {
				comparison = result[i][key] > result[j][key]
			}
		} else {
			// Numeric comparison
			if order == Ascending {
				comparison = iVal < jVal
			} else {
				comparison = iVal > jVal
			}
		}
		return comparison
	})
}

// collectDataByChild recursively collects all data from a node and its descendants
// Optimized to reduce allocations by passing slice pointer
func (t *Trie) collectDataByChild(node *Node) []map[string]string {
	// Estimate initial capacity based on typical trie usage
	result := make([]map[string]string, 0, 64)
	t.collectDataRecursive(node, &result)
	return result
}

// collectDataRecursive is the internal recursive function that appends to the provided slice
// This avoids creating new slices at each recursion level
func (t *Trie) collectDataRecursive(node *Node, result *[]map[string]string) {
	node.mu.RLock()

	// Add current node's data
	*result = append(*result, node.Data...)

	// Get children references while holding lock
	children := make([]*Node, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, child)
	}

	node.mu.RUnlock()

	// Recursively collect from children (lock released to avoid holding during recursion)
	for _, child := range children {
		t.collectDataRecursive(child, result)
	}
}

// Count returns the total number of data entries in the trie
func (t *Trie) Count() int {
	return t.countNode(t.Root)
}

// countNode recursively counts all data entries from a node and its descendants
func (t *Trie) countNode(node *Node) int {
	node.mu.RLock()
	defer node.mu.RUnlock()

	count := len(node.Data)
	for _, child := range node.Children {
		count += t.countNode(child)
	}
	return count
}

// Clear removes all data from the trie
func (t *Trie) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Root = NewNode(0)
}
