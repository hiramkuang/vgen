# Contributing to VGen

Thank you for your interest in contributing to VGen! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it before contributing.

## How Can I Contribute?

### Reporting Bugs

Before submitting a bug report, please check if the issue has already been reported. If not, create a new issue with the following information:

- A clear and descriptive title
- A detailed description of the problem
- Steps to reproduce the issue
- Expected behavior vs. actual behavior
- Information about your environment (Go version, OS, etc.)

### Suggesting Enhancements

If you have an idea for a new feature or improvement, please create an issue with:

- A clear and descriptive title
- A detailed explanation of the proposed feature
- Examples of how the feature would be used
- Any potential drawbacks or considerations

### Code Contributions

1. Fork the repository
2. Create a new branch for your feature or bug fix
3. Write your code following the project's coding style
4. Add or update tests as necessary
5. Ensure all tests pass
6. Commit your changes with a clear and descriptive commit message
7. Push your branch to your fork
8. Open a pull request with a clear title and description

## Development Setup

1. Install Go (version 1.16 or later)
2. Fork and clone the repository
3. Run `go mod tidy` to ensure dependencies are up to date

## Testing

Before submitting your changes, make sure all tests pass:

```bash
go test ./...
```

Please add tests for any new functionality you implement.

## Style Guide

This project follows the standard Go formatting conventions:

- Use `go fmt` to format your code
- Follow the guidelines in [Effective Go](https://golang.org/doc/effective_go.html)
- Write clear and concise comments for exported functions and types
- Use meaningful variable and function names

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build
2. Update the README.md with details of changes to the interface, including new environment variables, exposed ports, useful file locations and container parameters
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent
4. Your Pull Request will be reviewed by maintainers, who may request changes or ask questions
5. Once approved, your Pull Request will be merged

## Adding New Validation Rules

To add a new validation rule:

1. Add parsing logic in [internal/parser/tag.go](internal/parser/tag.go)
2. Add code generation logic in [internal/generator/generate.go](internal/generator/generate.go)
3. Add tests in the examples directory
4. Update the README.md with documentation for the new rule

## Questions?

If you have any questions about contributing, feel free to create an issue asking for clarification. We're here to help!