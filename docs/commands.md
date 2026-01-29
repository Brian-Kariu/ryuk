# Commands Reference

## `go-ryuk init`

Initialize Ryuk for the current project.

**What it does:**

1. Writes the current environment name and config path to `.ryuk.yaml`
2. Scans for existing `.env` files matching the environment name
3. If no valid Ryuk-managed files exist, creates `<environment>.env` with sentinel keys:
   ```
   TEST_ENV=test
   GORYUK=TRUE
   ```
4. Registers a config file watcher for `.ryuk.yaml`

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--environment` | `-e` | `test` | Environment to activate |

**Examples:**

```bash
go-ryuk init                    # Initialize with default "test" environment
go-ryuk init -e staging         # Initialize with "staging" environment
go-ryuk init -e prod            # Initialize with "prod" environment
```

**Files created:**

- `.ryuk.yaml` -- project configuration
- `<environment>.env` -- environment variables file (if none exists)

---

## `go-ryuk set <KEY> <VALUE>`

Add a key/value pair to the active environment's `.env` file.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `KEY` | Yes | Environment variable name |
| `VALUE` | Yes | Environment variable value |

Requires exactly 2 arguments. Duplicate keys are silently skipped.

**Examples:**

```bash
go-ryuk set DATABASE_URL postgres://localhost/mydb
go-ryuk set API_KEY sk-abc123
go-ryuk set -e staging DATABASE_URL postgres://staging/mydb
```

---

## `go-ryuk unset <KEY>...`

Remove one or more keys from the active environment's `.env` file.

**Arguments:**

| Argument | Required | Description |
|----------|----------|-------------|
| `KEY` | Yes (1+) | One or more environment variable names to remove |

Requires at least 1 argument. Missing keys are silently ignored.

**Examples:**

```bash
go-ryuk unset API_KEY
go-ryuk unset DATABASE_URL API_KEY SECRET_KEY
```

---

## `go-ryuk check`

Scan the directory tree for valid Ryuk-managed `.env` files.

Walks upward from the current directory looking for `.env` files, then validates each one for the `GORYUK=TRUE` sentinel key.

**Output:**

- `Found valid env file(s)` -- at least one valid file exists
- `No valid env file found configured for ryuk. Run go-ryuk init to configure ryuk for this project.` -- no valid files found

**Examples:**

```bash
go-ryuk check
```

---

## `go-ryuk detect`

Check whether the current project is configured for Ryuk. Functionally identical to `check` but with different messaging.

**Output:**

- `The current directory is ready to run ryuk commands!` -- project is configured
- `This directory does not have ryuk configured run 'go-ryuk init'` -- not configured

**Examples:**

```bash
go-ryuk detect
```

---

## Global flags

These flags apply to all commands:

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--environment` | `-e` | `test` | The active environment name |
| `--config` | `-c` | `true` | Enable/disable config file loading |
| `--help` | `-h` | | Show help for any command |
