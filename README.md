# go-ryuk

A CLI tool for managing per-project environment variables stored in `.env` files. Ryuk handles creating, storing, switching, and sharing environment configurations across deployment targets (test, staging, prod).

## Features

- **Per-environment `.env` files** -- separate configs for test, staging, production, etc.
- **Automatic discovery** -- walks the directory tree upward to find `.env` files
- **Sentinel-based validation** -- uses a `GORYUK=TRUE` key to identify Ryuk-managed files
- **Simple key/value management** -- set and unset environment variables from the command line
- **YAML configuration** -- project settings stored in `.ryuk.yaml`

## Installation

### From source

Requires Go 1.24+.

```bash
git clone https://github.com/Brian-Kariu/go-ryuk.git
cd go-ryuk
go build -o go-ryuk .
```

Move the binary to a directory in your `$PATH`:

```bash
mv go-ryuk /usr/local/bin/
```

## Quick start

```bash
# Initialize ryuk in your project directory
go-ryuk init

# Set environment variables
go-ryuk set DATABASE_URL postgres://localhost/mydb
go-ryuk set API_KEY sk-abc123

# Remove a variable
go-ryuk unset API_KEY

# Verify the project is configured
go-ryuk check

# Switch environments
go-ryuk init -e staging
go-ryuk set DATABASE_URL postgres://staging-host/mydb
```

## Commands

| Command | Description |
|---------|-------------|
| `init` | Initialize Ryuk for the current project. Creates `.ryuk.yaml` and a `<environment>.env` file. |
| `set <KEY> <VALUE>` | Add a key/value pair to the active environment's `.env` file. |
| `unset <KEY>...` | Remove one or more keys from the active environment's `.env` file. |
| `check` | Scan the directory tree for valid Ryuk-managed `.env` files. |
| `detect` | Check whether the current project is configured for Ryuk. |

## Global flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--environment` | `-e` | `test` | The active environment name (test, staging, prod, etc.) |
| `--config` | `-c` | `true` | Enable/disable config file loading |

## How it works

1. `go-ryuk init` creates a `.ryuk.yaml` config file and an `<environment>.env` file (e.g. `test.env`) with two sentinel entries:

   ```
   TEST_ENV=test
   GORYUK=TRUE
   ```

2. The `GORYUK=TRUE` line marks the file as Ryuk-managed. Commands like `check` and `detect` look for this sentinel to validate `.env` files.

3. `set` and `unset` read and write the environment-specific `.env` file using the [godotenv](https://github.com/joho/godotenv) library.

4. The `Scan()` function walks up the directory tree from the current working directory until it finds `.env` files, so you can run Ryuk commands from subdirectories.

## Project structure

```
go-ryuk/
  main.go              Entry point
  cmd/
    root.go            Root command, Viper config, persistent flags
    init.go            init subcommand
    check.go           check subcommand
    check_dir.go       CheckDir() helper
    detect.go          detect subcommand
    set.go             set subcommand
    unset.go           unset subcommand
    scan.go            Scan() -- directory walker for .env files
    validate.go        Validate() -- sentinel key checker
    create.go          Create() -- appends key=value to .env
    remove_env.go      RemoveEnv() -- deletes keys from .env
```

## Development

```bash
go build -o go-ryuk .       # Build
go run main.go <command>     # Run without building
go test ./...                # Run all tests
go vet ./...                 # Static analysis
go mod tidy                  # Clean up dependencies
```

## License

MIT License. See [LICENSE](main.go) header for details.
