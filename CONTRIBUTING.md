# Contributing to gogame

Thank you for your interest in contributing to gogame! This document provides guidelines for contributing to the project.

## ⚠️ Alpha Software Notice

gogame is in **alpha development stage**. APIs are unstable and subject to change without notice. Breaking changes are frequent as the engine design evolves.

## Before Contributing

1. **Check existing issues** - Look for existing issues or create a new one to discuss your proposed changes
2. **Read the documentation** - Familiarize yourself with the codebase structure and patterns
3. **Understand the project goals** - This is a macOS-focused 2D game engine prioritizing simplicity and performance

## Development Setup

### Prerequisites

- macOS 12.0 (Monterey) or newer
- Go 1.25.3 or newer
- Homebrew package manager
- SDL2 and SDL2_image libraries

### Installation

```bash
# Clone the repository
git clone https://github.com/dshills/gogame.git
cd gogame

# Install SDL2 dependencies
brew install sdl2 sdl2_image pkg-config

# Download Go dependencies
go mod download

# Verify installation
go build ./...
go test ./...
```

## How to Contribute

### Reporting Bugs

Create an issue with:
- Clear, descriptive title
- Steps to reproduce the bug
- Expected vs actual behavior
- macOS version, Go version, SDL2 version
- Code sample if applicable

### Suggesting Features

Create an issue with:
- Clear description of the feature
- Use cases and benefits
- Potential implementation approach
- Alignment with project goals

### Pull Requests

1. **Fork and create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Follow coding standards**
   - Write idiomatic Go code following [Effective Go](https://golang.org/doc/effective_go)
   - Add comments for exported functions, types, and packages
   - Handle errors explicitly (no silent failures)
   - Use meaningful variable and function names

3. **Write tests**
   - Add unit tests for new functionality
   - Ensure existing tests pass: `go test ./...`
   - Table-driven tests are preferred
   - Aim for meaningful test coverage

4. **Run linters**
   ```bash
   # Format code
   gofmt -w .
   
   # Run golangci-lint (if installed)
   golangci-lint run
   ```

5. **Update documentation**
   - Update README.md if adding major features
   - Add code examples for new APIs
   - Update CAMERA_GUIDE.md or create new guides as needed

6. **Commit your changes**
   - Write clear, descriptive commit messages
   - Reference issue numbers where applicable
   - Keep commits focused and atomic

7. **Submit pull request**
   - Provide clear description of changes
   - Reference related issues
   - Explain testing performed
   - Note any breaking changes

## Code Style Guidelines

### General Principles

- **Simplicity first** - Prefer simple, readable code over clever solutions
- **YAGNI** - Don't add features or abstractions before they're needed
- **Composition over inheritance** - Use interfaces and composition
- **Explicit over implicit** - Make intentions clear in code

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `core`, `graphics`, `math`)
- **Exported types**: PascalCase (e.g., `Engine`, `Sprite`, `Vector2`)
- **Unexported types**: camelCase (e.g., `internalState`)
- **Interfaces**: noun or adjective (e.g., `Behavior`, `Renderer`)

### Error Handling

```go
// Good - explicit error handling
if err := file.Close(); err != nil {
    return fmt.Errorf("failed to close file: %w", err)
}

// Good - intentional error ignore with comment
_ = texture.Destroy() // Safe to ignore in cleanup path

// Bad - silent error ignore
file.Close()
```

### Testing

- Place tests in `_test.go` files
- Use table-driven tests for multiple test cases
- Mock external dependencies (e.g., SDL2 calls)
- Test both success and error paths
- Use descriptive test names: `TestEntityUpdate_WithNilBehavior`

## Project Structure

```
gogame/
├── engine/
│   ├── core/        # Engine, Scene, Entity, game loop
│   ├── graphics/    # Renderer, Sprite, Texture, Camera
│   ├── input/       # InputManager, actions, keycodes
│   ├── physics/     # Collision detection, Collider
│   └── math/        # Vector2, Rectangle, Transform, Color
├── examples/        # Example games and demos
├── specs/           # Technical specifications (Specify framework)
└── tests/           # Unit, integration, and benchmark tests
```

## Development Workflow

1. Create an issue describing the change
2. Fork the repository and create a feature branch
3. Make your changes following the guidelines above
4. Write tests for your changes
5. Run tests and linters
6. Update documentation
7. Submit a pull request with clear description
8. Address review feedback

## Questions?

- Open an issue for questions about contributing
- Check existing documentation in `specs/` directory
- Review examples in `examples/` directory

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming community for all contributors.

## License

By contributing to gogame, you agree that your contributions will be licensed under the MIT License.

## Thank You!

Your contributions help make gogame better for everyone. Whether it's a bug report, feature suggestion, or code contribution, we appreciate your involvement!
