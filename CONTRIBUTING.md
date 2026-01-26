# Contributing to DataViz

Thank you for your interest in contributing to DataViz! This guide will help you get started.

## Development Setup

### Prerequisites

- Go 1.22 or later (1.25+ recommended)
- Git
- Make (optional, for convenience scripts)

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/dataviz.git
   cd dataviz
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Build the project:
   ```bash
   go build ./cmd/viz-cli
   go build ./cmd/dataviz-mcp
   ```

5. Run tests:
   ```bash
   go test ./...
   ```

## Project Structure

```
dataviz/
├── charts/          # Chart implementations (line, bar, pie, etc.)
├── mcp/             # MCP server implementation
│   ├── charts/      # MCP chart handlers
│   ├── types/       # MCP type definitions
│   └── mcp/         # MCP protocol implementation
├── cmd/
│   ├── viz-cli/     # CLI binary for terminal charts
│   └── dataviz-mcp/ # MCP server binary
├── internal/        # Internal packages (gallery, etc.)
├── examples/        # Example code and data files
├── examples-gallery/# Generated gallery SVGs
└── docs/            # Documentation
```

## Making Changes

### Code Style

- Follow standard Go conventions (gofmt, golint)
- Add comments for exported functions and types
- Keep functions focused and small
- Write tests for new functionality

### Testing

- Add tests for new features: `func TestNewFeature(t *testing.T)`
- Run tests: `go test ./...`
- Run tests with coverage: `go test -coverprofile=coverage.out ./...`
- Check coverage: `go tool cover -html=coverage.out`

### Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(charts): add candlestick chart support

Implements OHLC candlestick charts for financial data visualization.
Includes support for custom colors and volume bars.

Closes #42
```

```
fix(mcp): correct scatter plot marker rendering

Markers were not scaling correctly with chart dimensions.
Now properly sized relative to plot area.
```

## Pull Request Process

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Write code
   - Add tests
   - Update documentation

3. **Run checks locally**:
   ```bash
   go test ./...
   go vet ./...
   gofmt -w .
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat(charts): add your feature"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**:
   - Go to the original repository
   - Click "New Pull Request"
   - Select your fork and branch
   - Fill in the PR template
   - Wait for CI checks to pass
   - Address review feedback

### PR Guidelines

- **Title**: Clear, concise description of changes
- **Description**: Explain what, why, and how
- **Tests**: Include tests for new functionality
- **Documentation**: Update README or docs if needed
- **Breaking Changes**: Clearly mark and explain
- **Linked Issues**: Reference related issues with `Closes #N`

## Types of Contributions

### Bug Reports

Open an issue with:
- Clear title describing the bug
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS
- Relevant code snippets or error messages

### Feature Requests

Open an issue with:
- Clear description of the feature
- Use case and motivation
- Proposed API or interface (if applicable)
- Examples of similar features in other projects

### Documentation

- Fix typos or unclear explanations
- Add examples and tutorials
- Improve API documentation
- Update outdated information

### Code Contributions

#### Adding New Charts

1. Create chart implementation in `charts/`
2. Add MCP handler in `mcp/charts/`
3. Register tool in `mcp/tools.go`
4. Add example in `examples/`
5. Add gallery entry (if applicable)
6. Write tests
7. Update documentation

#### Improving Existing Charts

1. Identify the chart file in `charts/`
2. Make your improvements
3. Update tests
4. Add example demonstrating the improvement
5. Update documentation

## Development Workflow

### Local Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for specific package
go test ./charts/...

# Run specific test
go test -run TestLineChart ./charts/
```

### Building

```bash
# Build CLI
go build -o viz-cli ./cmd/viz-cli

# Build MCP server
go build -o dataviz-mcp ./cmd/dataviz-mcp

# Install locally
go install ./cmd/viz-cli
go install ./cmd/dataviz-mcp
```

### Generating Examples

```bash
# Generate all gallery examples
go run ./cmd/gallery-gen

# Test CLI with example data
./viz-cli examples/line-chart.json
```

## Code Review Process

1. Maintainers will review your PR
2. Address feedback by pushing new commits
3. Once approved, a maintainer will merge
4. Your contribution will be included in the next release

## Questions?

- Open an issue for questions
- Check existing issues and PRs
- Read the [documentation](docs/)

## License

By contributing, you agree that your contributions will be licensed under the BearWare 1.0 license (MIT-compatible).

Thank you for contributing! 🐻
