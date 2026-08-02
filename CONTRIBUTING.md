# Contributing to zte-cli

Thank you for your interest in contributing to zte-cli! 🎉

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git
- A ZTE F609 router for testing (optional)

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/CallMeJaja/zte-cli.git
cd zte-cli

# Install dependencies
go mod download

# Build
go build -o zte-cli .

# Run
./zte-cli --help
```

## How to Contribute

### Reporting Bugs

1. Check existing [issues](https://github.com/CallMeJaja/zte-cli/issues) first
2. Create a new issue using the [Bug Report template](https://github.com/CallMeJaja/zte-cli/issues/new?template=bug_report.md)
3. Include your router model, firmware version, and error output

### Suggesting Features

1. Check existing [issues](https://github.com/CallMeJaja/zte-cli/issues) first
2. Create a new issue using the [Feature Request template](https://github.com/CallMeJaja/zte-cli/issues/new?template=feature_request.md)
3. If possible, include the `.gch` page name for the feature

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Run linter: `golangci-lint run`
6. Commit with a descriptive message
7. Push to your fork
8. Create a Pull Request using the PR template

## Development Guidelines

### Code Style

- Follow standard Go conventions
- Use `gofmt` and `goimports`
- Run `golangci-lint` before committing

### Project Structure

```
zte-cli/
├── main.go              # Entry point, CLI routing
├── config/              # Config loading
├── router/              # Router communication
│   ├── client.go        # HTTP client, login
│   ├── pages.go         # Page constants
│   ├── parser.go        # HTML parsing
│   └── ...              # Feature-specific files
└── cli/                 # User interface
    ├── interactive.go   # TUI menu
    └── commands*.go     # CLI commands
```

### Adding a New Feature

1. **Add page constant** in `router/pages.go`:
   ```go
   const PageNewFeature = "new_feature.gch"
   ```

2. **Add parser** in `router/newfeature.go`:
   ```go
   type NewFeature struct { ... }
   func FetchNewFeature(client *Client) (*NewFeature, error) { ... }
   func FormatNewFeature(f *NewFeature) string { ... }
   ```

3. **Add CLI command** in `cli/commands_new.go`:
   ```go
   func RunNewFeatureCommand(client *router.Client) { ... }
   ```

4. **Add to main.go** routing:
   ```go
   case "new-feature":
       cli.RunNewFeatureCommand(client)
   ```

5. **Update help text** in `main.go`

6. **Update README.md** with new command documentation

### Finding Router Pages

To discover new `.gch` pages:

1. Open router web interface (`http://192.168.100.1`)
2. Open browser DevTools (F12)
3. Navigate to the feature you want to implement
4. Check Network tab for `.gch` page URLs
5. Check Elements tab for `Transfer_meaning()` calls

### Testing

```bash
# Run all tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Commit Messages

Use conventional commits:

```
feat: add new feature
fix: fix bug description
docs: update documentation
refactor: refactor code
test: add tests
chore: update dependencies
```

Examples:
```
feat: add WiFi channel selection
fix: fix login timeout handling
docs: update README with new commands
```

## Adding Support for Other Routers

The code is designed for ZTE F609, but the architecture supports other models:

1. Create a new router package (e.g., `router/f670/`)
2. Implement the same interfaces
3. Use router detection to select the correct implementation

## Questions?

Feel free to open an issue for any questions about contributing.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
