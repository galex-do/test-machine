# Database Migrations

This project uses [Goose](https://github.com/pressly/goose) for database migrations with Docker Compose support for easy management without local installation.

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

## Docker-based Migration Commands (Recommended)

The project includes a Docker-based migration system that requires no local Goose installation:

```bash
# Run all pending migrations
make migrate

# Check migration status
make migrate-status

# Rollback last migration  
make migrate-down

# Create a new migration file
make migrate-create NAME=add_new_feature

# Reset all migrations (WARNING: destructive)
make migrate-reset
```

## Docker Compose Migration Service

The migration service is defined in `docker-compose.yml` with the profile `migration`:

```bash
# Run migrations only
docker-compose --profile migration up migrate

# Run migrations and start application
docker-compose --profile migration up

# Check migration status
docker-compose run --rm migrate goose postgres "$DATABASE_URL" status
```

## Local Migration Commands (Optional)

If you have Goose installed locally:

```bash
# Install goose CLI
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations locally
make migrate-local-up

# Check status locally
make migrate-local-status

# Rollback locally
make migrate-local-down
```

## Creating New Migrations

### Using Make (Recommended)
```bash
make migrate-create NAME=add_user_roles
```

### Using Docker
```bash
docker-compose run --rm migrate goose create add_user_roles sql
```

### Manual Creation
Create a file with timestamp: `migrations/YYYYMMDDHHMMSS_migration_name.sql`

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

## Environment Variables

The migration system uses these environment variables:

- `DATABASE_USER` (default: postgres)
- `DATABASE_PASSWORD` (default: postgres)  
- `DATABASE_DB` (default: test)
- `DATABASE_HOST` (default: postgres for containers, localhost for local)

## Troubleshooting

### Container Issues
```bash
# Rebuild migration container
docker-compose build migrate

# View migration logs
docker-compose logs migrate

# Manual migration check
docker-compose run --rm migrate goose postgres "$DATABASE_URL" version
```

### Database Connection Issues
```bash
# Check database health
docker-compose ps postgres

# Connect to database directly
make shell-db
```