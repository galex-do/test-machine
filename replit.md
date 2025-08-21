# Overview

This is a Test Management Platform built as a web application for organizing and tracking software testing activities. The platform provides functionality for managing test projects, test suites, test cases, and test runs with a comprehensive dashboard and reporting system. It is designed for QA teams and software testers to organize their testing workflows and track test execution results.

## Recent Changes (August 2025)
- Built working Go implementation using standard library (simple_main.go)
- **MAJOR REFACTOR**: Migrated from in-memory storage to PostgreSQL database
- Added PostgreSQL driver (lib/pq) and full database integration
- Created comprehensive database schema with proper relationships and indexes
- Updated Docker containerization with PostgreSQL support and orchestration
- Modified Docker Compose to include PostgreSQL service with persistent volumes
- Added database initialization scripts and sample data seeding
- **COMPLETE RESTRUCTURE**: Refactored from single file to standard Go project structure
- Implemented proper Go module structure with github.com/galex-do/test-machine
- Created layered architecture: cmd/, internal/ packages following Go conventions
- Added repository pattern with separate packages: models, services, handlers, database
- Separated concerns with proper service layer and database abstraction
- Updated Dockerfile to build from new cmd/server/main.go entry point
- Cleaned up legacy files (simple_main.go, main.go, old packages)
- Successfully deployed on port 5000 with persistent PostgreSQL backend
- All CRUD operations now use authentic PostgreSQL data with ACID compliance

# User Preferences

Preferred communication style: Simple, everyday language.

# System Architecture

## Frontend Architecture

The frontend follows a traditional server-side rendered architecture using HTML templates with Go's template engine. The application uses a multi-page application (MPA) approach rather than a single-page application (SPA).

**Template Structure:**
- Consistent layout across all pages with shared navigation and sidebar
- Bootstrap 5 for responsive UI components and styling
- Font Awesome for iconography
- Custom CSS for platform-specific styling and theme variables

**Key Pages:**
- Dashboard (index.html) - Main overview page
- Project management (project.html) - Individual project views
- Test suite management (test-suite.html) - Test suite organization
- Test case management (test-case.html) - Individual test case details
- Test run tracking (test-run.html) - Test execution tracking
- Reports (reports.html) - Analytics and reporting dashboard

**Client-Side Features:**
- Search functionality for test cases
- Interactive status tracking with color-coded badges
- Date formatting utilities
- Alert system for user feedback
- Loading states for async operations

## Backend Architecture

The backend is built with Go following standard Go project conventions and clean architecture principles. The system uses PostgreSQL for persistent data storage with a layered architecture pattern.

**Project Structure:**
- `cmd/server/main.go` - Application entry point and dependency injection
- `internal/config/` - Configuration management with environment variables
- `internal/database/` - Database connection and connection pooling
- `internal/models/` - Data models and request/response structures
- `internal/repository/` - Data access layer with PostgreSQL integration
- `internal/service/` - Business logic layer with validation
- `internal/handlers/` - HTTP handlers for web and API routes

**Architecture Pattern:**
- Repository pattern for data access abstraction
- Service layer for business logic separation
- Dependency injection through constructor functions
- Clean separation of concerns across layers

**Routing Structure:**
- `/` - Dashboard/home page
- `/reports` - Reporting and analytics
- `/projects/{id}` - Individual project pages
- `/test-suites/{id}` - Test suite management
- `/test-cases/{id}` - Test case details
- `/test-runs/{id}` - Test execution tracking

**Data Models:**
The application manages several core entities:
- Projects - Top-level organization units
- Test Suites - Collections of related test cases
- Test Cases - Individual test scenarios
- Test Runs - Execution instances of test cases

**Status Management:**
The system tracks various status types:
- Test execution statuses: Pass, Fail, Blocked, Skip, Not Executed
- Progress statuses: In Progress, Completed
- Activity statuses: Active, Inactive

## Data Storage

The application requires a database for persistent storage of test management data. The structure suggests relational data with hierarchical relationships between projects, test suites, test cases, and test runs.

## Authentication and Authorization

The current templates don't show explicit authentication mechanisms, suggesting either:
- Authentication is handled at a higher level (reverse proxy, middleware)
- Authentication features are not yet implemented
- The application is designed for internal use without complex auth requirements

# External Dependencies

## Frontend Dependencies

**CSS Frameworks:**
- Bootstrap 5.1.3 - UI component library and responsive grid system
- Font Awesome 6.0.0 - Icon library for consistent iconography

**CDN Resources:**
- All external resources are loaded from CDNs rather than bundled locally
- No build process or package manager appears to be in use

## Backend Dependencies

**Go Standard Library:**
- html/template - Server-side HTML templating
- net/http - HTTP server functionality (implied)

**Potential Database:**
- The application structure suggests a SQL database is needed
- Likely PostgreSQL, MySQL, or SQLite based on common Go patterns

## Development Dependencies

**Static Assets:**
- Custom CSS file (style.css) for platform-specific styling
- Custom JavaScript file (app.js) for client-side functionality
- No build tools or transpilation appears to be configured

## Containerization and Deployment

**Docker Support:**
- Multi-stage Dockerfile with Go 1.21 Alpine build environment
- Optimized final image using minimal Alpine Linux base
- Proper static asset copying (templates and static files)
- Exposed on port 5000 for web access

**Local Deployment Options:**
1. **Docker Compose** (Recommended):
   ```bash
   docker-compose up --build
   ```
   - Includes networking configuration
   - Volume mounting for persistent data
   - Automatic restart policies

2. **Direct Docker Build**:
   ```bash
   docker build -t test-management .
   docker run -p 5000:5000 test-management
   ```

**Configuration Files:**
- `Dockerfile` - Multi-stage build with Go compilation and Alpine runtime
- `docker-compose.yml` - Complete orchestration configuration
- `.dockerignore` - Optimized build context exclusions
- `build.sh` - Local build script for development
- `README.md` - Complete deployment instructions

**Current Implementation:**
- Uses `simple_main.go` with Go standard library only
- In-memory data storage with pre-initialized sample data
- No external database dependencies for containerized deployment
- All static assets and templates included in container image