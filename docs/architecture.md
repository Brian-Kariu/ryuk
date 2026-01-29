# Architecture

## Overview

go-ryuk is a single-binary CLI application built with the Cobra + Viper pattern, the standard approach for Go CLI tools. All source code lives in a single `cmd` package.

## Dependency graph

```
main.go
  └── cmd.Execute()
        └── rootCmd (cobra.Command)
              ├── PersistentPreRun: preExecute()
              │     └── viper.Unmarshal(&configInstance)
              ├── cobra.OnInitialize: initConfig()
              │     ├── createFileIfNotExists(".ryuk.yaml")
              │     └── viper.ReadInConfig()
              │
              ├── initCmd
              │     ├── Scan()
              │     ├── Validate()
              │     └── Create()
              │
              ├── checkCmd
              │     └── CheckDir()
              │           ├── Scan()
              │           └── Validate()
              │
              ├── detectCmd
              │     └── CheckDir()
              │
              ├── setCmd
              │     └── Create()
              │
              └── unsetCmd
                    └── RemoveEnv()
```

## Layers

The codebase has two conceptual layers:

### Command layer (user-facing)

Files: `root.go`, `init.go`, `check.go`, `detect.go`, `set.go`, `unset.go`

Each file defines a Cobra command with `Use`, `Short`, `Long`, `Args`, and `RunE` fields. Commands are responsible for:

- Parsing and validating arguments
- Calling helper functions
- Printing user-facing output
- Returning errors to Cobra for display

### Helper layer (logic)

Files: `scan.go`, `validate.go`, `create.go`, `remove_env.go`, `check_dir.go`

Helper functions contain the core logic. They:

- Accept explicit parameters (no global state)
- Return errors instead of printing
- Perform file I/O through `os` and `godotenv`

## Configuration flow

```
                ┌─────────────┐
                │  .ryuk.yaml │
                └──────┬──────┘
                       │ viper.ReadInConfig()
                       v
              ┌─────────────────┐
              │  Viper (global) │
              └────────┬────────┘
                       │ viper.Unmarshal()
                       v
              ┌─────────────────┐
              │ configInstance   │
              │   .CfgFile      │
              │   .Environment  │
              └─────────────────┘
```

1. `cobra.OnInitialize` triggers `initConfig()` before any command runs
2. `initConfig()` locates/creates `.ryuk.yaml` and loads it into Viper
3. `PersistentPreRun` triggers `preExecute()` which unmarshals Viper state into the `configInstance` struct
4. Commands read the active environment via `viper.GetString("environment")`

## File discovery flow

```
Scan(startDir, environment)
  │
  ├── Start at startDir (or cwd)
  ├── List directory entries
  ├── Filter by ".env" in name
  ├── Filter by environment name (if provided)
  ├── If matches found → return
  └── Otherwise → move to parent directory, repeat
```

## Sentinel validation

Every Ryuk-managed `.env` file must contain `GORYUK=TRUE`. The `Validate()` function reads each `.env` file via `godotenv.Read()` and checks for this key. Files without it are ignored.

## Dependencies

| Package | Purpose |
|---------|---------|
| `spf13/cobra` | CLI framework (commands, flags, help text) |
| `spf13/viper` | Configuration management (YAML, env vars, flags) |
| `joho/godotenv` | Reading/writing `.env` files |
| `fsnotify/fsnotify` | Config file change watching (used in init) |
