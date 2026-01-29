# Development Guide

## Prerequisites

- Go 1.24 or later

## Build and run

```bash
# Build the binary
go build -o go-ryuk .

# Run without building
go run main.go <command>

# Run a specific command
go run main.go init
go run main.go set FOO bar
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./cmd/ -v

# Run a specific test
go test ./cmd/ -run TestScan_FindsEnvFile

# Run tests with race detection
go test -race ./...
```

### Existing test coverage

| File | Tests |
|------|-------|
| `cmd/scan_test.go` | `TestScan_FindsEnvFile`, `TestScan_WalksUpward`, `TestScan_FiltersByEnvironment` |
| `cmd/validate_test.go` | `TestValidate_AcceptsGoryukFile`, `TestValidate_RejectsPlainEnvFile` |

Helper functions accept a `startDir` parameter so tests can use `t.TempDir()` instead of the real filesystem.

## Linting

```bash
go vet ./...
```

## Dependency management

```bash
# Add a new dependency
go get github.com/example/package

# Update all dependencies
go get -u ./...

# Clean up go.mod and go.sum
go mod tidy
```

## Project conventions

### Error handling

- Helper functions return `error` as the last return value
- Commands use `RunE` (not `Run`) to propagate errors to Cobra
- Errors are wrapped with `fmt.Errorf("context: %w", err)` for stack context
- User-facing output is only produced in command handlers, not in helpers

### Naming

- Exported helper functions use PascalCase: `Scan()`, `Validate()`, `Create()`, `RemoveEnv()`, `CheckDir()`
- Command variables use camelCase: `initCmd`, `checkCmd`, `setCmd`
- The sentinel constant is `SentinelKey` in `cmd/validate.go`

### Adding a new command

1. Create `cmd/<name>.go`
2. Define a `var <name>Cmd = &cobra.Command{...}` with `RunE`
3. Add `Args` validation if the command takes arguments
4. Register it in an `init()` function: `rootCmd.AddCommand(<name>Cmd)`
5. Call helper functions for logic; keep the command handler thin

### Adding a new helper

1. Create or edit a file in `cmd/`
2. Accept all inputs as explicit parameters (no package-level state)
3. Return `error` as the last return value
4. Do not print to stdout -- let callers handle output
5. Add tests in `cmd/<name>_test.go` using `t.TempDir()` for file operations
