# Contributing to go-prefix-trie

Thank you for considering contributing to go-prefix-trie! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming and inclusive community.

## How to Contribute

### Reporting Bugs

Before creating a bug report, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Description**: Clear description of the issue
- **Steps to Reproduce**: Detailed steps to reproduce the behavior
- **Expected Behavior**: What you expected to happen
- **Actual Behavior**: What actually happened
- **Environment**: Go version, OS, architecture
- **Code Sample**: Minimal code that reproduces the issue

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Use Case**: Why this enhancement would be useful
- **Proposed Solution**: How you envision the enhancement working
- **Alternatives**: Alternative solutions you've considered
- **Additional Context**: Any other context about the enhancement

### Pull Requests

1. **Fork the Repository**: Fork the repo and create your branch from `main`

2. **Set Up Development Environment**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-prefix-trie.git
   cd go-prefix-trie
   go mod download
   ```

3. **Make Your Changes**:
   - Follow the existing code style
   - Add tests for new functionality
   - Update documentation as needed
   - Ensure all tests pass

4. **Run Tests and Linters**:
   ```bash
   # Run all tests
   go test -v -race -cover ./...

   # Run linter
   golangci-lint run

   # Run benchmarks (if applicable)
   go test -bench=. -benchmem ./...
   ```

5. **Commit Your Changes**:
   - Use clear, descriptive commit messages
   - Follow conventional commits format if possible:
     - `feat: add new feature`
     - `fix: resolve bug`
     - `docs: update documentation`
     - `test: add tests`
     - `refactor: refactor code`
     - `perf: improve performance`

6. **Push and Create Pull Request**:
   ```bash
   git push origin your-feature-branch
   ```
   - Fill out the PR template completely
   - Link related issues
   - Provide context for reviewers

## Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Run `golangci-lint` before committing
- Write clear, self-documenting code
- Add comments for complex logic

### Testing

- Maintain test coverage above 95%
- Write table-driven tests when appropriate
- Test edge cases and error conditions
- Ensure tests pass with `-race` flag
- Add benchmarks for performance-critical code

### Documentation

- Add godoc comments for all exported types and functions
- Update README.md for user-facing changes
- Include examples for new features
- Keep documentation concise and accurate

### Thread Safety

- All public methods must be thread-safe
- Use fine-grained locking when possible
- Test concurrent operations thoroughly
- Document any thread-safety guarantees

## Project Structure

```
go-prefix-trie/
├── .github/
│   └── workflows/        # CI/CD workflows
├── prefix_trie.go        # Main implementation
├── prefix_trie_test.go   # Tests
├── prefix_trie_bench_test.go  # Benchmarks
├── README.md             # User documentation
├── CONTRIBUTING.md       # This file
├── LICENSE               # MIT License
└── go.mod                # Go module file
```

## Questions?

If you have questions about contributing, feel free to:
- Open an issue with the question label
- Reach out to the maintainers

## Recognition

Contributors will be recognized in the project. Thank you for helping make go-prefix-trie better!
