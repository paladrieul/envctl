# envctl

A CLI tool for managing environment variables across multiple deployment targets with encryption support.

## Installation

```bash
go install ./...
```

## Commands

### `env set`
Set one or more environment variables in a target.
```bash
envctl env set --target production DATABASE_URL=postgres://...
```

### `env get`
Get the value of an environment variable.
```bash
envctl env get --target production DATABASE_URL
```

### `env list`
List all environment variables in a target.
```bash
envctl env list --target production
```

### `env del`
Delete an environment variable from a target.
```bash
envctl env del --target production DATABASE_URL
```

### `rename`
Rename an environment variable key within a target, preserving its value.
```bash
envctl rename --target production OLD_KEY NEW_KEY

# Overwrite the destination key if it already exists
envctl rename --target production --overwrite OLD_KEY EXISTING_KEY
```

### `copy`
Copy environment variables from one target to another.
```bash
envctl copy --from staging --to production
```

### `diff`
Show differences between two targets.
```bash
envctl diff --a staging --b production
```

### `import`
Import environment variables from a file.
```bash
envctl import --target production --format dotenv .env
envctl import --target production --format json env.json
```

### `export` (via `env list`)
Export environment variables to a file format.
```bash
envctl env list --target production --format dotenv
envctl env list --target production --format json
envctl env list --target production --format shell
```

### `encrypt` / `decrypt`
Encrypt or decrypt a value using a passphrase.
```bash
envctl encrypt --passphrase mysecret myvalue
envctl decrypt --passphrase mysecret <ciphertext>
```

## Configuration

Passphrases can be provided via:
- `--passphrase` flag
- `ENVCTL_PASSPHRASE` environment variable
- A passphrase file at `~/.config/envctl/passphrase`
