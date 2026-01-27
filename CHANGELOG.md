# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-01-27

### Performance Improvements

#### `collectDataByChild` Optimization - 90% Allocation Reduction

**Problem:**
The original `collectDataByChild` function created a new slice at each recursion level and used slice spread operator (`...`) to merge results. For a trie with 1000 entries, this resulted in ~1042 allocations per `GetByPrefix` call.

```go
// BEFORE: Inefficient - creates new slice at each level
func (t *Trie) collectDataByChild(node *Node) []map[string]string {
    node.mu.RLock()
    defer node.mu.RUnlock()

    result := make([]map[string]string, 0, len(node.Data)+len(node.Children))
    result = append(result, node.Data...)

    for _, child := range node.Children {
        // This creates a new slice and copies all elements at EVERY recursion level
        result = append(result, t.collectDataByChild(child)...)
    }
    return result
}
```

**Solution:**
Introduced a two-function pattern where the public function allocates a single slice with estimated capacity, and an internal recursive function appends directly to that slice via pointer.

```go
// AFTER: Optimized - single slice allocation, append via pointer
func (t *Trie) collectDataByChild(node *Node) []map[string]string {
    result := make([]map[string]string, 0, 64)
    t.collectDataRecursive(node, &result)
    return result
}

func (t *Trie) collectDataRecursive(node *Node, result *[]map[string]string) {
    node.mu.RLock()
    *result = append(*result, node.Data...)

    children := make([]*Node, 0, len(node.Children))
    for _, child := range node.Children {
        children = append(children, child)
    }
    node.mu.RUnlock()

    for _, child := range children {
        t.collectDataRecursive(child, result)
    }
}
```

**Key Changes:**
1. **Single allocation**: Only one result slice is created at the top level
2. **Pointer passing**: Recursive function receives `*[]map[string]string` instead of returning a new slice
3. **Lock optimization**: Lock is released before recursing into children to reduce lock contention
4. **Children snapshot**: Children references are captured while holding lock, then processed after release

**Benchmark Results:**

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| `GetByPrefix` allocations | 1,042 | 105 | **-90%** |
| `GetByPrefix` latency | 94,040 ns/op | 55,733 ns/op | **-41%** |
| `GetByPrefix` memory | 66,208 B/op | 25,024 B/op | **-62%** |
| `ConcurrentGetByPrefix` latency | 79,699 ns/op | 32,774 ns/op | **-59%** |

### Added

#### Performance Target Tests

New test functions that define and verify performance targets:

- `TestPerf_GetByPrefix_Allocations` - Verifies allocation count meets target (≤500)
- `TestPerf_GetByPrefixSortBy_Allocations` - Verifies sorting overhead is minimal (≤10 extra allocs)
- `TestPerf_GetByPrefix_LargeDataset` - Verifies O(n) memory growth (not O(n²))
- `TestPerf_SortResults_NoExtraAllocations` - Verifies sort function allocations (≤5)

These tests use `testing.AllocsPerRun()` to measure actual allocations and compare against targets.

#### Unit Tests for `collectDataByChild`

Comprehensive test coverage for the optimized data collection:

- `TestCollectDataByChild_BasicFunctionality` - Basic prefix collection
- `TestCollectDataByChild_EmptyNode` - Nodes without data (prefix-only)
- `TestCollectDataByChild_DeepTree` - Deeply nested structures (10 levels)
- `TestCollectDataByChild_WideTree` - Wide structures (100 siblings)
- `TestCollectDataByChild_MultipleDataPerNode` - Multiple entries per key
- `TestCollectDataByChild_PreserveDataIntegrity` - Data not corrupted during collection
- `TestCollectDataByChild_ConcurrentCollection` - Thread safety with parallel reads
- `TestCollectDataByChild_ConcurrentModification` - Safety during concurrent writes

#### New Benchmarks

- `BenchmarkTrie_GetByPrefix_Optimized` - Reports allocation targets with custom metrics
- `BenchmarkTrie_SortResults_Optimized` - Measures sort-only performance

### Technical Details

#### Why This Matters

In Go, every `append()` that exceeds slice capacity triggers:
1. New memory allocation (typically 2x previous capacity)
2. Copy of all existing elements to new memory
3. Garbage collection of old memory

The original implementation did this at every recursion level:
- Level 1: allocate + copy
- Level 2: allocate + copy (includes Level 1 data)
- Level 3: allocate + copy (includes Level 1+2 data)
- ...and so on

This resulted in O(n²) copy operations in the worst case.

The optimized version:
- Single allocation at top level
- Direct append to same slice at all levels
- O(n) total operations

#### Thread Safety Considerations

The optimization also improves the locking strategy:

**Before:** Lock held during entire recursive descent
```go
node.mu.RLock()
defer node.mu.RUnlock()
// ... lock held while recursing into ALL children
```

**After:** Lock released before recursion
```go
node.mu.RLock()
// capture data and children references
node.mu.RUnlock()
// recurse without holding parent lock
```

This reduces lock contention in concurrent scenarios, explaining the 59% improvement in `ConcurrentGetByPrefix`.

## [1.0.0] - 2026-01-27

### Added

- Initial stable release
- Thread-safe prefix trie implementation with fine-grained locking
- Core operations: `Insert`, `Search`, `Delete`, `Update`
- Prefix queries: `HasPrefix`, `GetByPrefix`, `GetByPrefixSortBy`
- Flexible sorting with `Ascending` and `Descending` order
- Automatic numeric vs string comparison in sorting
- `Count` and `Clear` utility methods
- Custom error types: `ErrEmptyKey`, `ErrNilData`, `ErrNotFound`
- Full Unicode support
- 99.2% test coverage
- Comprehensive benchmark suite
- MIT License
