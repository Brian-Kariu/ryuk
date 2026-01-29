# Configuration

## `.ryuk.yaml`

The project configuration file is created in the current working directory when you run `go-ryuk init`. It stores the active environment and config file path.

**Example contents:**

```yaml
config: /home/user/myproject/.ryuk.yaml
environment: test
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `config` | string | Absolute path to this config file |
| `environment` | string | The active environment name |

### How it's loaded

1. On startup, `initConfig()` runs via `cobra.OnInitialize`
2. It calls `os.Getwd()` to find the current directory
3. It constructs the path `<cwd>/.ryuk.yaml`
4. If the file doesn't exist, it creates an empty one
5. Viper reads the file and merges values with CLI flags and environment variables

### Precedence

Viper resolves configuration values in this order (highest to lowest):

1. CLI flags (`--environment`, `--config`)
2. Environment variables (via `viper.AutomaticEnv()`)
3. `.ryuk.yaml` file values
4. Default values (`test` for environment, `true` for config)

## Environment `.env` files

Each environment has its own `.env` file named `<environment>.env` (e.g., `test.env`, `staging.env`, `prod.env`).

**Example `test.env`:**

```
TEST_ENV=test
GORYUK=TRUE
DATABASE_URL=postgres://localhost/mydb
API_KEY=sk-abc123
```

### Sentinel key

Every Ryuk-managed `.env` file must contain:

```
GORYUK=TRUE
```

This sentinel key is defined as the constant `SentinelKey` in `cmd/validate.go`. The `Validate()` function uses it to distinguish Ryuk-managed files from plain `.env` files that may exist in the same directory tree.

### File format

Files follow the standard dotenv format parsed by [godotenv](https://github.com/joho/godotenv):

- One `KEY=VALUE` pair per line
- Lines starting with `#` are comments
- Values can be quoted with single or double quotes
- Multiline values are supported with quotes
