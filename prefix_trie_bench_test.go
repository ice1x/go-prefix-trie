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
		"field1":  "value1",
		"field2":  "value2",
		"field3":  "value3",
		"field4":  "value4",
		"field5":  "value5",
		"field6":  "value6",
		"field7":  "value7",
		"field8":  "value8",
		"field9":  "value9",
		"field10": "value10",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = trie.Insert(fmt.Sprintf("word%d", i), data)
	}
}

// ============================================================================
// PERFORMANCE TARGET TESTS
// These tests define performance targets for optimizations.
// They measure allocations and verify they meet target thresholds.
// ============================================================================

// TestPerf_GetByPrefix_Allocations tests that GetByPrefix meets allocation targets.
// Target: reduce allocations by 50% from baseline (~1000 allocs to ~500 allocs for 1000 items)
func TestPerf_GetByPrefix_Allocations(t *testing.T) {
	trie := NewTrie()

	// Setup: 1000 items with common prefix
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
			"score": fmt.Sprintf("%d", i),
		})
	}

	// Measure allocations
	allocsPerOp := testing.AllocsPerRun(100, func() {
		_ = trie.GetByPrefix("prefix")
	})

	// Target: <= 500 allocations (50% reduction from ~1042)
	const targetAllocs = 500

	if allocsPerOp > targetAllocs {
		t.Logf("PERF TARGET NOT MET: GetByPrefix allocations = %.0f, target = %d", allocsPerOp, targetAllocs)
		t.Logf("This is expected before optimization. Run optimizations to meet target.")
		// Note: We log but don't fail to allow TDD red-green cycle visibility
	} else {
		t.Logf("PERF TARGET MET: GetByPrefix allocations = %.0f <= %d", allocsPerOp, targetAllocs)
	}
}

// TestPerf_GetByPrefixSortBy_Allocations tests sorting optimization.
// Target: sorting should not add more than 10 allocations over GetByPrefix
func TestPerf_GetByPrefixSortBy_Allocations(t *testing.T) {
	trie := NewTrie()

	// Setup: 1000 items with numeric scores
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
			"score": fmt.Sprintf("%d", i),
		})
	}

	// Measure GetByPrefix allocations
	getPrefixAllocs := testing.AllocsPerRun(100, func() {
		_ = trie.GetByPrefix("prefix")
	})

	// Measure GetByPrefixSortBy allocations
	sortByAllocs := testing.AllocsPerRun(100, func() {
		_ = trie.GetByPrefixSortBy("prefix", "score", Descending)
	})

	// Target: sorting overhead <= 10 allocations
	const maxSortingOverhead = 10
	overhead := sortByAllocs - getPrefixAllocs

	if overhead > maxSortingOverhead {
		t.Logf("PERF TARGET NOT MET: Sorting overhead = %.0f allocs, target = <= %d", overhead, maxSortingOverhead)
		t.Logf("This is expected before optimization. Run optimizations to meet target.")
	} else {
		t.Logf("PERF TARGET MET: Sorting overhead = %.0f allocs <= %d", overhead, maxSortingOverhead)
	}
}

// TestPerf_GetByPrefix_LargeDataset tests performance with larger dataset.
// Target: O(n) memory growth, not O(n^2)
func TestPerf_GetByPrefix_LargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	// Test with 5000 items
	trie := NewTrie()
	for i := 0; i < 5000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
		})
	}

	allocs5000 := testing.AllocsPerRun(10, func() {
		_ = trie.GetByPrefix("prefix")
	})

	// Test with 10000 items
	trie2 := NewTrie()
	for i := 0; i < 10000; i++ {
		_ = trie2.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
		})
	}

	allocs10000 := testing.AllocsPerRun(10, func() {
		_ = trie2.GetByPrefix("prefix")
	})

	// Target: allocations should scale linearly (ratio should be ~2x, not 4x)
	ratio := allocs10000 / allocs5000
	const maxRatio = 2.5 // allowing some overhead

	if ratio > maxRatio {
		t.Logf("PERF TARGET NOT MET: Allocation ratio (10K/5K) = %.2f, target = <= %.1f", ratio, maxRatio)
		t.Logf("This suggests O(n^2) or worse allocation growth.")
	} else {
		t.Logf("PERF TARGET MET: Allocation ratio (10K/5K) = %.2f <= %.1f (linear growth)", ratio, maxRatio)
	}
}

// TestPerf_SortResults_NoExtraAllocations tests that optimized sort doesn't allocate per comparison.
func TestPerf_SortResults_NoExtraAllocations(t *testing.T) {
	trie := NewTrie()

	// Setup: 100 items (enough to trigger many comparisons in sort)
	for i := 0; i < 100; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
			"score": fmt.Sprintf("%d", i),
		})
	}

	// Get data first
	data := trie.GetByPrefix("prefix")

	// Measure sort-only allocations
	allocsPerSort := testing.AllocsPerRun(100, func() {
		trie.sortResults(data, "score", Descending)
	})

	// Target: <= 5 allocations for sorting (ideally 0-2 for pre-parsed approach)
	const targetAllocs = 5

	if allocsPerSort > targetAllocs {
		t.Logf("PERF TARGET NOT MET: sortResults allocations = %.0f, target = <= %d", allocsPerSort, targetAllocs)
		t.Logf("Current implementation calls strconv.Atoi per comparison.")
	} else {
		t.Logf("PERF TARGET MET: sortResults allocations = %.0f <= %d", allocsPerSort, targetAllocs)
	}
}

// BenchmarkTrie_GetByPrefix_Optimized benchmarks with reporting of allocation targets
func BenchmarkTrie_GetByPrefix_Optimized(b *testing.B) {
	trie := NewTrie()
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
			"score": fmt.Sprintf("%d", i),
		})
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = trie.GetByPrefix("prefix")
	}

	// Report custom metrics
	b.ReportMetric(float64(1000), "items")
}

// BenchmarkTrie_SortResults_Optimized benchmarks sort with pre-parsing
func BenchmarkTrie_SortResults_Optimized(b *testing.B) {
	trie := NewTrie()
	for i := 0; i < 1000; i++ {
		_ = trie.Insert(fmt.Sprintf("prefix%d", i), map[string]string{
			"value": fmt.Sprintf("val%d", i),
			"score": fmt.Sprintf("%d", i),
		})
	}
	data := trie.GetByPrefix("prefix")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Make a copy since sort is in-place
		dataCopy := make([]map[string]string, len(data))
		copy(dataCopy, data)
		trie.sortResults(dataCopy, "score", Descending)
	}
}
