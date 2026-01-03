package prefixtrie

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestTrie_Insert(t *testing.T) {
	trie := NewTrie()

	// Test normal insert
	err := trie.Insert("hello", map[string]string{"value": "world"})
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}

	err = trie.Insert("help", map[string]string{"value": "me"})
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}

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

	// Test error cases
	err = trie.Insert("", map[string]string{"value": "test"})
	if err != ErrEmptyKey {
		t.Errorf("Expected ErrEmptyKey for empty string, got %v", err)
	}

	err = trie.Insert("test", nil)
	if err != ErrNilData {
		t.Errorf("Expected ErrNilData for nil data, got %v", err)
	}
}

func TestTrie_Search(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world", "id": "1"})
	_ = trie.Insert("hello", map[string]string{"value": "earth", "id": "2"})
	_ = trie.Insert("help", map[string]string{"value": "me", "id": "3"})

	tests := []struct {
		name    string
		word    string
		wantLen int
		wantErr error
	}{
		{
			name:    "Find hello with 2 entries",
			word:    "hello",
			wantLen: 2,
			wantErr: nil,
		},
		{
			name:    "Find help with 1 entry",
			word:    "help",
			wantLen: 1,
			wantErr: nil,
		},
		{
			name:    "Not found",
			word:    "hell",
			wantLen: 0,
			wantErr: ErrNotFound,
		},
		{
			name:    "Empty key",
			word:    "",
			wantLen: 0,
			wantErr: ErrEmptyKey,
		},
		{
			name:    "Non-existent word",
			word:    "goodbye",
			wantLen: 0,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := trie.Search(tt.word)
			if err != tt.wantErr {
				t.Errorf("Search(%q) error = %v, want %v", tt.word, err, tt.wantErr)
			}
			if len(result) != tt.wantLen {
				t.Errorf("Search(%q) returned %d results, want %d", tt.word, len(result), tt.wantLen)
			}
		})
	}
}

func TestTrie_HasPrefix(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world"})
	_ = trie.Insert("help", map[string]string{"value": "me"})
	_ = trie.Insert("hero", map[string]string{"value": "brave"})

	tests := []struct {
		name   string
		prefix string
		want   bool
	}{
		{"Empty prefix", "", true},
		{"Prefix 'h'", "h", true},
		{"Prefix 'he'", "he", true},
		{"Prefix 'hel'", "hel", true},
		{"Prefix 'help'", "help", true},
		{"Prefix 'hello'", "hello", true},
		{"Prefix 'her'", "her", true},
		{"Non-existent prefix 'hi'", "hi", false},
		{"Non-existent prefix 'goodbye'", "goodbye", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trie.HasPrefix(tt.prefix)
			if got != tt.want {
				t.Errorf("HasPrefix(%q) = %v, want %v", tt.prefix, got, tt.want)
			}
		})
	}
}

func TestTrie_Delete(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world"})
	_ = trie.Insert("help", map[string]string{"value": "me"})

	// Delete existing word
	err := trie.Delete("hello")
	if err != nil {
		t.Errorf("Delete('hello') failed: %v", err)
	}

	// Verify deletion
	_, err = trie.Search("hello")
	if err != ErrNotFound {
		t.Errorf("After deletion, Search('hello') should return ErrNotFound, got %v", err)
	}

	// Verify other word still exists
	result, err := trie.Search("help")
	if err != nil {
		t.Errorf("Search('help') failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 'help' to still have 1 entry, got %d", len(result))
	}

	// Test error cases
	err = trie.Delete("")
	if err != ErrEmptyKey {
		t.Errorf("Delete('') should return ErrEmptyKey, got %v", err)
	}

	err = trie.Delete("nonexistent")
	if err != ErrNotFound {
		t.Errorf("Delete('nonexistent') should return ErrNotFound, got %v", err)
	}

	err = trie.Delete("hel") // prefix exists but not as complete word
	if err != ErrNotFound {
		t.Errorf("Delete('hel') should return ErrNotFound, got %v", err)
	}
}

func TestTrie_Update(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world", "version": "1"})
	_ = trie.Insert("hello", map[string]string{"value": "earth", "version": "2"})

	// Verify initial state
	result, _ := trie.Search("hello")
	if len(result) != 2 {
		t.Errorf("Expected 2 entries before update, got %d", len(result))
	}

	// Update
	err := trie.Update("hello", map[string]string{"value": "universe", "version": "3"})
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Verify update replaced all entries
	result, err = trie.Search("hello")
	if err != nil {
		t.Errorf("Search after update failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 entry after update, got %d", len(result))
	}
	if result[0]["value"] != "universe" {
		t.Errorf("Expected value 'universe', got '%s'", result[0]["value"])
	}

	// Test error cases
	err = trie.Update("", map[string]string{"value": "test"})
	if err != ErrEmptyKey {
		t.Errorf("Update('') should return ErrEmptyKey, got %v", err)
	}

	err = trie.Update("hello", nil)
	if err != ErrNilData {
		t.Errorf("Update with nil data should return ErrNilData, got %v", err)
	}

	err = trie.Update("nonexistent", map[string]string{"value": "test"})
	if err != ErrNotFound {
		t.Errorf("Update('nonexistent') should return ErrNotFound, got %v", err)
	}
}

func TestTrie_GetByPrefix(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world"})
	_ = trie.Insert("help", map[string]string{"value": "me"})
	_ = trie.Insert("hell", map[string]string{"value": "yeah"})
	_ = trie.Insert("hero", map[string]string{"value": "brave"})

	tests := []struct {
		name     string
		prefix   string
		wantLen  int
		wantNil  bool
	}{
		{"Prefix 'hel'", "hel", 3, false},
		{"Prefix 'he'", "he", 4, false},
		{"Prefix 'hello'", "hello", 1, false},
		{"Prefix 'her'", "her", 1, false},
		{"Non-existent prefix", "hi", 0, true},
		{"Empty prefix", "", 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trie.GetByPrefix(tt.prefix)
			if tt.wantNil {
				if result != nil {
					t.Errorf("GetByPrefix(%q) should return nil, got %d results", tt.prefix, len(result))
				}
			} else {
				if len(result) != tt.wantLen {
					t.Errorf("GetByPrefix(%q) returned %d results, want %d", tt.prefix, len(result), tt.wantLen)
				}
			}
		})
	}
}

func TestTrie_GetByPrefixSortBy(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world", "score": "10"})
	_ = trie.Insert("help", map[string]string{"value": "me", "score": "5"})
	_ = trie.Insert("hell", map[string]string{"value": "yeah", "score": "7"})

	tests := []struct {
		name     string
		prefix   string
		sortKey  string
		order    SortOrder
		expected []map[string]string
	}{
		{
			name:    "Sort by score descending",
			prefix:  "hel",
			sortKey: "score",
			order:   Descending,
			expected: []map[string]string{
				{"value": "world", "score": "10"},
				{"value": "yeah", "score": "7"},
				{"value": "me", "score": "5"},
			},
		},
		{
			name:    "Sort by score ascending",
			prefix:  "hel",
			sortKey: "score",
			order:   Ascending,
			expected: []map[string]string{
				{"value": "me", "score": "5"},
				{"value": "yeah", "score": "7"},
				{"value": "world", "score": "10"},
			},
		},
		{
			name:     "No results",
			prefix:   "hi",
			sortKey:  "score",
			order:    Descending,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trie.GetByPrefixSortBy(tt.prefix, tt.sortKey, tt.order)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetByPrefixSortBy(%q, %q, %v) = %v, want %v", tt.prefix, tt.sortKey, tt.order, result, tt.expected)
			}
		})
	}
}

func TestTrie_GetByPrefixSortDescBy(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world", "score": "10"})
	_ = trie.Insert("help", map[string]string{"value": "me", "score": "5"})
	_ = trie.Insert("hell", map[string]string{"value": "yeah", "score": "7"})

	result := trie.GetByPrefixSortDescBy("hel", "score")
	expected := []map[string]string{
		{"value": "world", "score": "10"},
		{"value": "yeah", "score": "7"},
		{"value": "me", "score": "5"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("GetByPrefixSortDescBy('hel', 'score') = %v, want %v", result, expected)
	}
}

func TestTrie_StringSorting(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("apple", map[string]string{"name": "apple", "category": "fruit"})
	_ = trie.Insert("apricot", map[string]string{"name": "apricot", "category": "fruit"})
	_ = trie.Insert("avocado", map[string]string{"name": "avocado", "category": "fruit"})

	// Test ascending string sort
	result := trie.GetByPrefixSortBy("a", "name", Ascending)
	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}
	if result[0]["name"] != "apple" || result[1]["name"] != "apricot" || result[2]["name"] != "avocado" {
		t.Errorf("String ascending sort failed: got %v", result)
	}

	// Test descending string sort
	result = trie.GetByPrefixSortBy("a", "name", Descending)
	if result[0]["name"] != "avocado" || result[1]["name"] != "apricot" || result[2]["name"] != "apple" {
		t.Errorf("String descending sort failed: got %v", result)
	}
}

func TestTrie_Count(t *testing.T) {
	trie := NewTrie()

	// Empty trie
	if count := trie.Count(); count != 0 {
		t.Errorf("Empty trie should have count 0, got %d", count)
	}

	// Add entries
	_ = trie.Insert("hello", map[string]string{"value": "world"})
	_ = trie.Insert("hello", map[string]string{"value": "earth"})
	_ = trie.Insert("help", map[string]string{"value": "me"})

	if count := trie.Count(); count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}

	// After delete
	_ = trie.Delete("hello")
	if count := trie.Count(); count != 1 {
		t.Errorf("After delete, expected count 1, got %d", count)
	}
}

func TestTrie_Clear(t *testing.T) {
	trie := NewTrie()
	_ = trie.Insert("hello", map[string]string{"value": "world"})
	_ = trie.Insert("help", map[string]string{"value": "me"})

	// Verify data exists
	if count := trie.Count(); count != 2 {
		t.Errorf("Before clear, expected count 2, got %d", count)
	}

	// Clear
	trie.Clear()

	// Verify empty
	if count := trie.Count(); count != 0 {
		t.Errorf("After clear, expected count 0, got %d", count)
	}

	result := trie.GetByPrefix("")
	if len(result) != 0 {
		t.Errorf("After clear, expected 0 results, got %d", len(result))
	}
}

func TestTrie_Concurrency(t *testing.T) {
	trie := NewTrie()
	var wg sync.WaitGroup
	iterations := 100

	// Concurrent inserts
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_ = trie.Insert(fmt.Sprintf("word%d", id), map[string]string{"value": fmt.Sprintf("val%d-%d", id, j)})
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				trie.GetByPrefixSortBy("word", "value", Descending)
				trie.HasPrefix("w")
			}
		}()
	}

	// Concurrent searches
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, _ = trie.Search(fmt.Sprintf("word%d", id))
			}
		}(i)
	}

	wg.Wait()

	// Verify no data corruption
	count := trie.Count()
	if count == 0 {
		t.Error("After concurrent operations, trie should not be empty")
	}
}

func TestTrie_ConcurrentUpdatesAndDeletes(t *testing.T) {
	trie := NewTrie()
	var wg sync.WaitGroup

	// Pre-populate
	for i := 0; i < 10; i++ {
		_ = trie.Insert(fmt.Sprintf("key%d", i), map[string]string{"value": "initial"})
	}

	// Concurrent updates and deletes
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("key%d", id)
				_ = trie.Update(key, map[string]string{"value": fmt.Sprintf("updated-%d", j)})
			}
		}(i)
	}

	for i := 5; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				key := fmt.Sprintf("key%d", id)
				_ = trie.Delete(key)
				_ = trie.Insert(key, map[string]string{"value": fmt.Sprintf("reinserted-%d", j)})
			}
		}(i)
	}

	wg.Wait()

	// Verify all keys still exist
	for i := 0; i < 10; i++ {
		result, err := trie.Search(fmt.Sprintf("key%d", i))
		if err != nil {
			t.Errorf("Key%d should exist after concurrent operations", i)
		}
		if len(result) == 0 {
			t.Errorf("Key%d should have data", i)
		}
	}
}

func TestTrie_EdgeCases(t *testing.T) {
	trie := NewTrie()

	// Unicode characters
	err := trie.Insert("你好", map[string]string{"value": "hello"})
	if err != nil {
		t.Errorf("Failed to insert Unicode: %v", err)
	}

	result, err := trie.Search("你好")
	if err != nil || len(result) == 0 {
		t.Error("Failed to search Unicode")
	}

	// Special characters
	err = trie.Insert("test@#$%", map[string]string{"value": "special"})
	if err != nil {
		t.Errorf("Failed to insert special chars: %v", err)
	}

	// Very long string
	longString := ""
	for i := 0; i < 1000; i++ {
		longString += "a"
	}
	err = trie.Insert(longString, map[string]string{"value": "long"})
	if err != nil {
		t.Errorf("Failed to insert long string: %v", err)
	}

	result, err = trie.Search(longString)
	if err != nil || len(result) == 0 {
		t.Error("Failed to search long string")
	}
}
