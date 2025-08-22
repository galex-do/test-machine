# Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migrations.

## Migration Files

Migrations are stored in the `migrations/` directory and follow the naming convention:
`YYYYMMDDHHMMSS_migration_name.sql`

### Current Migrations

1. **20250822090400_add_test_steps_table.sql**
   - Creates the `test_steps` table with proper foreign key constraints
   - Adds indexes for performance optimization

2. **20250822090500_add_sample_test_steps.sql**
   - Populates all test cases with comprehensive test steps
   - Includes detailed step descriptions and expected results

## How Migrations Work

- Migrations run automatically when the application starts
- Goose tracks applied migrations in the `goose_db_version` table
- Each migration has an `Up` section (applied) and `Down` section (rollback)
- Migrations are applied in chronological order based on timestamps

## Creating New Migrations

To create a new migration:

```bash
# Install goose CLI (if not already installed)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Create a new migration file
goose -dir migrations create migration_name sql
```

## Manual Migration Commands

```bash
# Apply all pending migrations
goose -dir migrations postgres "postgres://user:password@localhost/dbname" up

# Check migration status
goose -dir migrations postgres "postgres://user:password@localhost/dbname" status

# Rollback last migration
goose -dir migrations postgres "postgres://user:password@localhost/dbname" down
```

## Migration Format

Each migration file should follow this format:

```sql
-- +goose Up
-- Add your up migration here
CREATE TABLE example (id SERIAL PRIMARY KEY);

-- +goose Down  
-- Add your down migration here
DROP TABLE IF EXISTS example;
```