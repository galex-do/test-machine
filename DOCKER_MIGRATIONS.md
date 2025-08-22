# Docker-based Database Migrations

This guide shows how to use the containerized migration system without installing Goose locally.

## Quick Start

```bash
# Run all pending migrations
make migrate

# Check what migrations have been applied
make migrate-status

# Create a new migration
make migrate-create NAME=add_user_roles

# Rollback the last migration (if needed)
make migrate-down
```

## Architecture

The migration system uses:
- **Dockerfile.migrate**: Custom container with Goose installed
- **docker-compose.yml**: Migration service with `migration` profile
- **Makefile**: Convenient commands for common operations

## Migration Service Features

- ✅ No local Goose installation required
- ✅ Consistent environment across all systems
- ✅ Health checks ensure database is ready
- ✅ Automatic dependency management
- ✅ Profile-based execution (migrations don't run by default)

## Advanced Usage

### Run migrations manually with Docker Compose
```bash
# Run migration service only
docker-compose --profile migration up migrate --build

# Check migration status
docker-compose run --rm migrate goose postgres "$DATABASE_URL" status

# Create new migration using Docker
docker-compose run --rm migrate goose create add_feature sql
```

### Environment Variables
The system respects these environment variables:
- `DATABASE_USER` (default: postgres)
- `DATABASE_PASSWORD` (default: postgres)
- `DATABASE_DB` (default: test)

### Production Usage
For production environments, set appropriate environment variables:
```bash
DATABASE_USER=prod_user DATABASE_PASSWORD=secure_pass make migrate
```

## Troubleshooting

### Migration Container Won't Start
```bash
# Rebuild the migration container
docker-compose build migrate

# Check container logs
docker-compose logs migrate
```

### Database Connection Issues
```bash
# Verify database is running
docker-compose ps postgres

# Check database health
docker-compose logs postgres

# Connect to database manually
make shell-db
```

### Migration Conflicts
```bash
# Check current migration status
make migrate-status

# If needed, rollback problematic migration
make migrate-down

# Then fix the migration file and reapply
make migrate
```