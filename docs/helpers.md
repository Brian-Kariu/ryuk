# Helper Functions

The helper layer contains the core logic of go-ryuk. All helpers accept explicit parameters, return errors, and produce no user-facing output.

## `Scan(startDir, environment string) ([]string, error)`

**File:** `cmd/scan.go`

Walks the directory tree upward from `startDir` looking for `.env` files.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `startDir` | `string` | Directory to start searching from. Empty string defaults to `os.Getwd()`. |
| `environment` | `string` | Filter files by environment name. Empty string matches all `.env` files. |

**Returns:** A slice of absolute paths to matching `.env` files, and an error.

**Behavior:**

1. Start at `startDir` (or current working directory if empty)
2. List all files in the directory
3. Keep files whose name contains `.env`
4. If `environment` is non-empty, further filter to files containing the environment name
5. If matches found, return them
6. If no matches, move to the parent directory and repeat
7. Stop at the filesystem root

**Example:**

```
project/
  test.env        <-- found by Scan("", "test")
  staging.env     <-- found by Scan("", "staging")
  src/
    (cwd)         <-- Scan starts here, walks up to project/
```

---

## `Validate(envFiles []string) ([]validEnv, error)`

**File:** `cmd/validate.go`

Checks which `.env` files are Ryuk-managed by looking for the `GORYUK` sentinel key.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `envFiles` | `[]string` | Paths to `.env` files to validate |

**Returns:** A slice of `validEnv` structs (file path + variable count) for files containing `GORYUK=TRUE`, and an error.

**The `validEnv` struct:**

```go
type validEnv struct {
    file    string  // absolute path to .env file
    envVars int     // total number of variables in the file
}
```

---

## `Create(key, value, filename string) error`

**File:** `cmd/create.go`

Appends a `KEY=VALUE` line to a `.env` file. Skips the write if the key already exists.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `key` | `string` | Variable name |
| `value` | `string` | Variable value |
| `filename` | `string` | Path to the `.env` file (created if it doesn't exist) |

**Behavior:**

1. Open the file in append mode (create if missing)
2. Read existing variables via `godotenv.Read()`
3. If the key already exists, skip
4. Otherwise, append `KEY=VALUE\n`
5. Flush the buffer

---

## `RemoveEnv(envs []string, filename string) error`

**File:** `cmd/remove_env.go`

Removes one or more keys from a `.env` file.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `envs` | `[]string` | Keys to remove |
| `filename` | `string` | Path to the `.env` file |

**Behavior:**

1. Read the file into a map via `godotenv.Read()`
2. Delete each specified key from the map
3. Write the modified map back to the file via `godotenv.Write()`

---

## `CheckDir() (bool, error)`

**File:** `cmd/check_dir.go`

Convenience function that combines `Scan()` and `Validate()` to determine if the current directory tree contains any Ryuk-managed `.env` files.

**Returns:** `true` if at least one valid file exists, `false` otherwise.

**Used by:** `check` and `detect` commands.
