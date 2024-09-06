//package prefixtrie
//
//import (
//	"reflect"
//	"testing"
//)
//
//func TestTrie_Insert(t *testing.T) {
//	trie := NewTrie()
//	trie.Insert("hello", map[string]string{"value": "world"})
//	trie.Insert("help", map[string]string{"value": "me"})
//
//	if len(trie.Root.Children) != 1 {
//		t.Errorf("Expected 1 child in root, got %d", len(trie.Root.Children))
//	}
//
//	hNode := trie.Root.Children['h']
//	if hNode == nil {
//		t.Fatal("Expected 'h' node to exist")
//	}
//
//	if len(hNode.Children) != 1 {
//		t.Errorf("Expected 1 child in 'h' node, got %d", len(hNode.Children))
//	}
//}
//
//func TestTrie_GetByPrefixSortDescBy(t *testing.T) {
//	trie := NewTrie()
//	trie.Insert("hello", map[string]string{"value": "world", "score": "10"})
//	trie.Insert("help", map[string]string{"value": "me", "score": "5"})
//	trie.Insert("hell", map[string]string{"value": "yeah", "score": "7"})
//
//	tests := []struct {
//		name     string
//		prefix   string
//		sortKey  string
//		expected []map[string]string
//	}{
//		{
//			name:    "Sort by score",
//			prefix:  "hel",
//			sortKey: "score",
//			expected: []map[string]string{
//				{"value": "world", "score": "10"},
//				{"value": "yeah", "score": "7"},
//				{"value": "me", "score": "5"},
//			},
//		},
//		{
//			name:    "No results",
//			prefix:  "hi",
//			sortKey: "score",
//			expected: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			result := trie.GetByPrefixSortDescBy(tt.prefix, tt.sortKey)
//			if !reflect.DeepEqual(result, tt.expected) {
//				t.Errorf("GetByPrefixSortDescBy(%q, %q) = %v, want %v", tt.prefix, tt.sortKey, result, tt.expected)
//			}
//		})
//	}
//}
//
//func TestTrie_Concurrency(t *testing.T) {
//	trie := NewTrie()
//	done := make(chan bool)
//
//	go func() {
//		for i := 0; i < 1000; i++ {
//			trie.Insert("concurrent", map[string]string{"value": "test"})
//		}
//		done <- true
//	}()
//
//	go func() {
//		for i := 0; i < 1000; i++ {
//			trie.GetByPrefixSortDescBy("con", "value")
//		}
//		done <- true
//	}()
//
//	<-done
//	<-done
//
//	result := trie.GetByPrefixSortDescBy("concurrent", "value")
//	if len(result) != 1000 {
//		t.Errorf("Expected 1000 results, got %d", len(result))
//	}
//}

package prefixtrie

import (
	"reflect"
	"testing"
)

func TestTrie_Insert(t *testing.T) {
	trie := NewTrie()
	trie.Insert("hello", map[string]string{"value": "world"})
	trie.Insert("help", map[string]string{"value": "me"})

	if len(trie.Root.Children) != 1 {
		t.Errorf("After inserting 'hello' and 'help', expected root to have 1 child, got %d", len(trie.Root.Children))
	}

	hNode := trie.Root.Children['h']
	if hNode == nil {
		t.Fatal("Expected 'h' node to exist after inserting 'hello' and 'help'")
	}

	if len(hNode.Children) != 1 {
		t.Errorf("Expected 'h' node to have 1 child, got %d", len(hNode.Children))
	}
}

func TestTrie_GetByPrefixSortDescBy(t *testing.T) {
	trie := NewTrie()
	trie.Insert("hello", map[string]string{"value": "world", "score": "10"})
	trie.Insert("help", map[string]string{"value": "me", "score": "5"})
	trie.Insert("hell", map[string]string{"value": "yeah", "score": "7"})

	tests := []struct {
		name     string
		prefix   string
		sortKey  string
		expected []map[string]string
	}{
		{
			name:    "Sort by score",
			prefix:  "hel",
			sortKey: "score",
			expected: []map[string]string{
				{"value": "world", "score": "10"},
				{"value": "yeah", "score": "7"},
				{"value": "me", "score": "5"},
			},
		},
		{
			name:     "No results",
			prefix:   "hi",
			sortKey:  "score",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trie.GetByPrefixSortDescBy(tt.prefix, tt.sortKey)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetByPrefixSortDescBy(%q, %q) = %v, want %v", tt.prefix, tt.sortKey, result, tt.expected)
			}
		})
	}
}

func TestTrie_Concurrency(t *testing.T) {
	trie := NewTrie()
	done := make(chan bool)

	go func() {
		for i := 0; i < 1000; i++ {
			trie.Insert("concurrent", map[string]string{"value": "test"})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			trie.GetByPrefixSortDescBy("con", "value")
		}
		done <- true
	}()

	<-done
	<-done

	result := trie.GetByPrefixSortDescBy("concurrent", "value")
	if len(result) != 1000 {
		t.Errorf("Expected 1000 results, got %d", len(result))
	}
}
