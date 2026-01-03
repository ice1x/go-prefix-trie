# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it by emailing the maintainer or creating a private security advisory on GitHub.

**Please do not create public issues for security vulnerabilities.**

### How to Report

1. **Via GitHub Security Advisories** (Recommended)
   - Go to the [Security tab](https://github.com/ice1x/go-prefix-trie/security)
   - Click "Report a vulnerability"
   - Fill in the details

2. **Via Email**
   - Send details to the repository maintainer
   - Include a description of the vulnerability
   - Include steps to reproduce if possible

### What to Include

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if any)

## Security Measures

This project implements several security measures:

### Automated Security Scanning

- **CodeQL Analysis**: Advanced semantic code analysis runs on every push
- **gosec**: Go security checker scans for common security issues
- **govulncheck**: Checks for known vulnerabilities in dependencies
- **Dependency Review**: Reviews pull requests for dependency vulnerabilities

### Development Practices

- **Thread Safety**: All public methods are thread-safe
- **Input Validation**: All user inputs are validated
- **Error Handling**: Comprehensive error handling throughout
- **Code Review**: All changes require review before merging

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Security Updates

Security updates are released as soon as possible after a vulnerability is confirmed. Check the [Releases](https://github.com/ice1x/go-prefix-trie/releases) page for security patches.

## Enabling Code Scanning (for Repository Owners)

To enable GitHub Code Scanning for this repository:

1. Go to **Settings** → **Code security and analysis**
2. Find **Code scanning**
3. Click **Set up** → **Advanced**
4. The workflow is already configured in `.github/workflows/codeql.yml`
5. Code scanning will run automatically on every push

### For Public Repositories
Code scanning is free and available immediately.

### For Private Repositories
Requires GitHub Advanced Security (part of GitHub Enterprise).

## Acknowledgments

We appreciate the security community's efforts in responsibly disclosing vulnerabilities.
