package prefixtrie_test

import (
	"errors"
	"fmt"

	prefixtrie "github.com/ice1x/go-prefix-trie"
)

// Example demonstrates the basic workflow: build a trie, then look a word up
// with an exact Search. Data returned by Search preserves insertion order, so
// the output is deterministic.
func Example() {
	trie := prefixtrie.NewTrie()

	_ = trie.Insert("hello", map[string]string{"lang": "en"})
	_ = trie.Insert("hola", map[string]string{"lang": "es"})

	results, err := trie.Search("hello")
	if err != nil {
		fmt.Println("search error:", err)
		return
	}
	fmt.Println(results[0]["lang"])
	// Output: en
}

// ExampleTrie_Insert shows that multiple data entries can be attached to the
// same word; Count reports the total number of stored entries.
func ExampleTrie_Insert() {
	trie := prefixtrie.NewTrie()

	_ = trie.Insert("go", map[string]string{"kind": "verb"})
	_ = trie.Insert("go", map[string]string{"kind": "noun"})
	_ = trie.Insert("golang", map[string]string{"kind": "language"})

	fmt.Println("entries:", trie.Count())
	// Output: entries: 3
}

// ExampleTrie_Search retrieves every data entry stored under an exact word.
// The slice keeps the order in which the entries were inserted.
func ExampleTrie_Search() {
	trie := prefixtrie.NewTrie()

	_ = trie.Insert("port", map[string]string{"n": "80"})
	_ = trie.Insert("port", map[string]string{"n": "443"})

	results, _ := trie.Search("port")
	for _, r := range results {
		fmt.Println(r["n"])
	}
	// Output:
	// 80
	// 443
}

// ExampleTrie_Search_notFound shows the sentinel error returned when a word is
// absent from the trie.
func ExampleTrie_Search_notFound() {
	trie := prefixtrie.NewTrie()

	_, err := trie.Search("missing")
	fmt.Println(errors.Is(err, prefixtrie.ErrNotFound))
	// Output: true
}

// ExampleTrie_HasPrefix reports whether any stored word starts with the given
// prefix. The empty prefix always matches.
func ExampleTrie_HasPrefix() {
	trie := prefixtrie.NewTrie()
	_ = trie.Insert("carpet", map[string]string{})

	fmt.Println(trie.HasPrefix("car"))
	fmt.Println(trie.HasPrefix("cat"))
	fmt.Println(trie.HasPrefix(""))
	// Output:
	// true
	// false
	// true
}

// ExampleTrie_GetByPrefixSortBy collects every entry under a prefix and sorts
// it by a numeric field. Because the sort key is numeric and distinct, the
// output order is deterministic (unlike the unordered GetByPrefix).
func ExampleTrie_GetByPrefixSortBy() {
	trie := prefixtrie.NewTrie()
	_ = trie.Insert("task-a", map[string]string{"name": "a", "priority": "3"})
	_ = trie.Insert("task-b", map[string]string{"name": "b", "priority": "1"})
	_ = trie.Insert("task-c", map[string]string{"name": "c", "priority": "2"})

	results := trie.GetByPrefixSortBy("task", "priority", prefixtrie.Descending)
	for _, r := range results {
		fmt.Printf("%s=%s\n", r["name"], r["priority"])
	}
	// Output:
	// a=3
	// c=2
	// b=1
}

// ExampleTrie_Update replaces all data stored under a word with a single new
// entry.
func ExampleTrie_Update() {
	trie := prefixtrie.NewTrie()
	_ = trie.Insert("key", map[string]string{"v": "old"})

	_ = trie.Update("key", map[string]string{"v": "new"})

	results, _ := trie.Search("key")
	fmt.Println(len(results), results[0]["v"])
	// Output: 1 new
}

// ExampleTrie_Delete removes all data associated with a word. A subsequent
// Search reports ErrNotFound.
func ExampleTrie_Delete() {
	trie := prefixtrie.NewTrie()
	_ = trie.Insert("temp", map[string]string{"v": "x"})

	_ = trie.Delete("temp")

	_, err := trie.Search("temp")
	fmt.Println(errors.Is(err, prefixtrie.ErrNotFound))
	// Output: true
}

// ExampleTrie_Count returns the total number of data entries across the whole
// trie.
func ExampleTrie_Count() {
	trie := prefixtrie.NewTrie()
	_ = trie.Insert("a", map[string]string{})
	_ = trie.Insert("ab", map[string]string{})
	_ = trie.Insert("abc", map[string]string{})

	fmt.Println(trie.Count())
	// Output: 3
}
