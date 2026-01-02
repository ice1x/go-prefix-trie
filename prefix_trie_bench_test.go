package prefixtrie

import (
	"fmt"
	"testing"
)

func BenchmarkTrie_Insert(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test", "score": "100"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}
}

func BenchmarkTrie_Search(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test", "score": "100"}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = trie.Search(fmt.Sprintf("word%d", i%1000))
	}
}

func BenchmarkTrie_HasPrefix(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test"}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.HasPrefix("word")
	}
}

func BenchmarkTrie_GetByPrefix(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test", "score": "100"}

	// Pre-populate with various prefixes
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), data)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.GetByPrefix("prefix")
	}
}

func BenchmarkTrie_GetByPrefixSortBy(b *testing.B) {
	trie := NewTrie()

	// Pre-populate with various scores
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
			"score": fmt.Sprintf("%d", i),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.GetByPrefixSortBy("prefix", "score", Descending)
	}
}

func BenchmarkTrie_Delete(b *testing.B) {
	// Pre-populate
	trie := NewTrie()
	data := map[string]string{"value": "test"}
	for i := 0; i < b.N; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Delete(fmt.Sprintf("word%d", i))
	}
}

func BenchmarkTrie_Update(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test"}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}

	newData := map[string]string{"value": "updated"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Update(fmt.Sprintf("word%d", i%1000), newData)
	}
}

func BenchmarkTrie_ConcurrentInsert(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test"}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = trie.Insert(fmt.Sprintf("word%d", i), data)
			i++
		}
	})
}

func BenchmarkTrie_ConcurrentSearch(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test"}

	// Pre-populate
	for i := 0; i < 10000; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = trie.Search(fmt.Sprintf("word%d", i%10000))
			i++
		}
	})
}

func BenchmarkTrie_ConcurrentGetByPrefix(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test"}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), data)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = trie.GetByPrefix("prefix")
		}
	})
}

func BenchmarkTrie_Count(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"value": "test"}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Count()
	}
}

// Benchmarks for different data sizes
func BenchmarkTrie_Insert_SmallData(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{"k": "v"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Insert(fmt.Sprintf("w%d", i), data)
	}
}

func BenchmarkTrie_Insert_LargeData(b *testing.B) {
	trie := NewTrie()
	data := map[string]string{
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
		"field4": "value4",
		"field5": "value5",
		"field6": "value6",
		"field7": "value7",
		"field8": "value8",
		"field9": "value9",
		"field10": "value10",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}
}
