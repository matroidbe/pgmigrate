# pgmigrate

Declarative PostgreSQL schema management CLI - like Terraform for your database schema.

## Overview

pgmigrate is a command-line tool that wraps the [pg_migrate](https://github.com/matroidbe/pg_extensions/tree/main/extensions/pg_migrate) PostgreSQL extension to provide a Terraform-like workflow for managing database schema.

**Key Features:**
- **Declarative**: Define your schema in YAML, pgmigrate computes the diff
- **No state file**: Uses pg_catalog as the source of truth
- **Safety-aware**: Categorizes changes as Safe, Destructive, or Breaking
- **Audit trail**: All migrations recorded in database history table

## Installation

### From Source (requires Go 1.22+)

```bash
curl -fsSL https://raw.githubusercontent.com/matroidbe/pgmigrate/main/install.sh | bash
```

Or with go install:

```bash
go install github.com/matroidbe/pgmigrate/cmd/pgmigrate@latest
```

### Prerequisites

The `pg_migrate` PostgreSQL extension must be installed in your database:

```sql
CREATE EXTENSION pg_migrate;
```

## Quick Start

```bash
# Set your database connection
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb"

# Create a template schema.yaml
pgmigrate init

# Edit schema.yaml to define your schema
vim schema.yaml

# Preview what changes would be made
pgmigrate plan

# Apply the changes
pgmigrate apply
```

## Commands

### `pgmigrate init`

Creates a template `schema.yaml` file in the current directory.

```bash
pgmigrate init
pgmigrate init --force  # Overwrite existing file
```

### `pgmigrate plan [file]`

Shows what changes would be applied without making any changes.

```bash
pgmigrate plan                    # Use schema.yaml
pgmigrate plan myschema.yaml      # Use specific file
pgmigrate plan -o json            # Output as JSON
```

**Output format:**
```
pgmigrate will perform the following actions:

  + CREATE TABLE public.products
  + public.users.profile_url (add column)
  - public.legacy.old_data (drop column)
  ! public.users.email (alter type)

Plan: 2 to add, 1 to destroy, 1 breaking.
```

### `pgmigrate apply [file]`

Applies schema changes to the database.

```bash
pgmigrate apply                          # Apply safe changes only
pgmigrate apply --allow-destructive      # Include DROP operations
pgmigrate apply --auto-approve           # Skip confirmation prompt
```

**Safety levels:**
- `+` **Safe**: Additive changes (CREATE, ADD COLUMN) - applied automatically
- `-` **Destructive**: Data loss possible (DROP) - requires `--allow-destructive`
- `!` **Breaking**: May fail or corrupt (ALTER TYPE) - requires manual intervention

### `pgmigrate dump <schema> [schemas...]`

Exports current database schema as YAML.

```bash
pgmigrate dump public                    # Dump single schema
pgmigrate dump public api                # Dump multiple schemas
pgmigrate dump public -o schema.yaml     # Write to file
```

### `pgmigrate history`

Shows migration history from the database.

```bash
pgmigrate history           # Show last 20 migrations
pgmigrate history -n 50     # Show last 50 migrations
```

### `pgmigrate version`

Shows CLI and extension versions.

```bash
pgmigrate version
```

## Schema Format

```yaml
# Schemas managed by pgmigrate
managed_schemas:
  - public

# Table definitions
tables:
  public.users:
    columns:
      - name: id
        type: bigserial
        primary_key: true
      - name: email
        type: varchar(255)
        not_null: true
      - name: name
        type: text
      - name: active
        type: boolean
        default: "true"
      - name: created_at
        type: timestamptz
        not_null: true
        default: "now()"

    indexes:
      - name: users_email_idx
        columns: [email]
        unique: true
```

### Column Properties

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | Column name (required) |
| `type` | string | PostgreSQL data type (required) |
| `primary_key` | bool | Part of primary key |
| `not_null` | bool | NOT NULL constraint |
| `default` | string | Default value (SQL expression) |
| `references` | string | Foreign key: `"schema.table.column"` |

### Index Properties

| Property | Type | Description |
|----------|------|-------------|
| `name` | string | Index name (required) |
| `columns` | list | Column names (required) |
| `unique` | bool | Unique index |

## Connection

pgmigrate connects to PostgreSQL via the `DATABASE_URL` environment variable:

```bash
export DATABASE_URL="postgres://user:password@host:port/database"
```

Or use the `--database-url` flag:

```bash
pgmigrate plan --database-url "postgres://..."
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--database-url` | Override DATABASE_URL environment variable |
| `--verbose, -v` | Enable verbose output |
| `--no-color` | Disable colored output |

## Breaking Changes

Some schema changes cannot be applied automatically because they may fail or cause data loss:

- Changing column types (may fail if data can't be converted)
- Adding NOT NULL to existing column (may fail if NULLs exist)
- Reducing column size (may truncate data)

For these changes, use `pgmigrate.dba_migrate()` directly in psql:

```sql
SELECT pgmigrate.dba_migrate(
    'ALTER TABLE users ALTER COLUMN email TYPE varchar(500)',
    'Expand email column size'
);
```

## License

LicenseRef-Matroid-SAL-1.0
