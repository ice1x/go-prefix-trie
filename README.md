# go-prefix-trie

A high-performance, thread-safe prefix trie (prefix tree) implementation in Go for efficient string prefix matching and retrieval.

[![Tests](https://github.com/ice1x/go-prefix-trie/actions/workflows/test.yml/badge.svg)](https://github.com/ice1x/go-prefix-trie/actions/workflows/test.yml)
[![CodeQL](https://github.com/ice1x/go-prefix-trie/actions/workflows/codeql.yml/badge.svg)](https://github.com/ice1x/go-prefix-trie/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ice1x/go-prefix-trie?v=1.1.0)](https://goreportcard.com/report/github.com/ice1x/go-prefix-trie)
[![codecov](https://codecov.io/gh/ice1x/go-prefix-trie/branch/main/graph/badge.svg)](https://codecov.io/gh/ice1x/go-prefix-trie)
[![Go Version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![GoDoc](https://pkg.go.dev/badge/github.com/ice1x/go-prefix-trie)](https://pkg.go.dev/github.com/ice1x/go-prefix-trie)

## Features

- **Thread-Safe**: All operations are protected with fine-grained locking for concurrent access
- **Flexible Sorting**: Sort results by any field in ascending or descending order
- **Rich API**: Insert, Search, Delete, Update, HasPrefix, and more
- **Error Handling**: Comprehensive error handling with custom error types
- **High Performance**: Optimized memory allocation and recursive algorithms
- **Unicode Support**: Full support for Unicode characters
- **Zero Dependencies**: Uses only Go standard library

## Installation

```bash
go get github.com/ice1x/go-prefix-trie
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/ice1x/go-prefix-trie"
)

func main() {
    // Create a new trie
    trie := prefixtrie.NewTrie()

    // Insert data
    trie.Insert("hello", map[string]string{"value": "world", "score": "10"})
    trie.Insert("help", map[string]string{"value": "me", "score": "5"})
    trie.Insert("hero", map[string]string{"value": "brave", "score": "8"})

    // Search by prefix and sort
    results := trie.GetByPrefixSortBy("hel", "score", prefixtrie.Descending)
    for _, result := range results {
        fmt.Printf("%s: %s (score: %s)\n", result["value"], result["value"], result["score"])
    }
}
```

## API Documentation

### Creating a Trie

```go
trie := prefixtrie.NewTrie()
```

### Insert

Add a word with associated data to the trie.

```go
err := trie.Insert("hello", map[string]string{"value": "world", "id": "1"})
if err != nil {
    // Handle error (ErrEmptyKey or ErrNilData)
}
```

**Errors:**
- `ErrEmptyKey`: when the key is empty
- `ErrNilData`: when data is nil

### Search

Find exact matches for a given word.

```go
results, err := trie.Search("hello")
if err != nil {
    // Handle error (ErrEmptyKey or ErrNotFound)
}
// results contains all data associated with "hello"
```

**Errors:**
- `ErrEmptyKey`: when the key is empty
- `ErrNotFound`: when the word doesn't exist

### HasPrefix

Check if any words in the trie start with the given prefix.

```go
exists := trie.HasPrefix("hel") // returns true if any words start with "hel"
```

### GetByPrefix

Retrieve all data for words starting with the given prefix (unsorted).

```go
results := trie.GetByPrefix("hel")
// returns all entries that start with "hel"
```

### GetByPrefixSortBy

Retrieve all data for words starting with the given prefix, sorted by a specific field.

```go
// Sort by score in descending order
results := trie.GetByPrefixSortBy("hel", "score", prefixtrie.Descending)

// Sort by name in ascending order
results = trie.GetByPrefixSortBy("hel", "name", prefixtrie.Ascending)
```

**Sort Orders:**
- `prefixtrie.Ascending`: Sort from lowest to highest
- `prefixtrie.Descending`: Sort from highest to lowest

**Note:** The function automatically detects numeric values and sorts them numerically. String values are sorted lexicographically.

### Delete

Remove all data associated with a specific word.

```go
err := trie.Delete("hello")
if err != nil {
    // Handle error (ErrEmptyKey or ErrNotFound)
}
```

**Errors:**
- `ErrEmptyKey`: when the key is empty
- `ErrNotFound`: when the word doesn't exist

### Update

Replace all data for a given word with new data.

```go
err := trie.Update("hello", map[string]string{"value": "universe", "score": "99"})
if err != nil {
    // Handle error (ErrEmptyKey, ErrNilData, or ErrNotFound)
}
```

**Errors:**
- `ErrEmptyKey`: when the key is empty
- `ErrNilData`: when data is nil
- `ErrNotFound`: when the word doesn't exist

### Count

Get the total number of data entries in the trie.

```go
count := trie.Count()
```

### Clear

Remove all data from the trie.

```go
trie.Clear()
```

## Examples

### Autocomplete System

```go
trie := prefixtrie.NewTrie()

// Add words with popularity scores
trie.Insert("apple", map[string]string{"name": "apple", "popularity": "1000"})
trie.Insert("application", map[string]string{"name": "application", "popularity": "500"})
trie.Insert("apply", map[string]string{"name": "apply", "popularity": "800"})

// Get autocomplete suggestions sorted by popularity
suggestions := trie.GetByPrefixSortBy("app", "popularity", prefixtrie.Descending)
for _, suggestion := range suggestions {
    fmt.Println(suggestion["name"])
}
// Output:
// apple
// apply
// application
```

### Word Dictionary

```go
trie := prefixtrie.NewTrie()

// Add words with definitions
trie.Insert("hello", map[string]string{
    "word": "hello",
    "definition": "a greeting",
    "partOfSpeech": "interjection",
})

// Search for exact word
results, err := trie.Search("hello")
if err == nil {
    fmt.Println(results[0]["definition"]) // Output: a greeting
}

// Update definition
trie.Update("hello", map[string]string{
    "word": "hello",
    "definition": "used as a greeting or to begin a phone conversation",
    "partOfSpeech": "interjection",
})

// Delete word
trie.Delete("hello")
```

### Concurrent Access

```go
trie := prefixtrie.NewTrie()

// Safe concurrent writes
go func() {
    for i := 0; i < 1000; i++ {
        trie.Insert(fmt.Sprintf("word%d", i), map[string]string{"id": fmt.Sprintf("%d", i)})
    }
}()

// Safe concurrent reads
go func() {
    for i := 0; i < 1000; i++ {
        trie.GetByPrefix("word")
    }
}()
```

## Performance

The trie is optimized for:
- Fast prefix lookups: O(m) where m is the length of the prefix
- Memory-efficient storage with shared prefixes
- Minimal allocations during data collection

Run benchmarks:

```bash
go test -bench=. -benchmem
```

Example results:
```
BenchmarkTrie_Insert-8                    500000    2847 ns/op    1024 B/op    10 allocs/op
BenchmarkTrie_Search-8                   1000000    1234 ns/op     512 B/op     5 allocs/op
BenchmarkTrie_GetByPrefix-8               100000   15678 ns/op    8192 B/op    50 allocs/op
BenchmarkTrie_ConcurrentInsert-8         2000000     678 ns/op     512 B/op     8 allocs/op
```

## Testing

Run all tests:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

Run race detector:

```bash
go test -race
```

## Thread Safety

All methods are thread-safe and can be called concurrently from multiple goroutines. The implementation uses fine-grained locking at the node level to maximize concurrency while preventing data races.

## Error Handling

The package defines three error types:

- `ErrEmptyKey`: Returned when an empty key is provided
- `ErrNilData`: Returned when nil data is provided
- `ErrNotFound`: Returned when a key is not found

Example:

```go
_, err := trie.Search("nonexistent")
if err == prefixtrie.ErrNotFound {
    fmt.Println("Word not found")
}
```

## Use Cases

- Autocomplete systems
- Spell checkers
- IP routing tables
- Dictionary implementations
- Search engines
- Command-line completion
- Contact lists with prefix search
- URL shorteners with prefix matching

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**ice1x** - [GitHub](https://github.com/ice1x)

## Acknowledgments

- Inspired by classic trie data structure implementations
- Built with performance and concurrency in mind
