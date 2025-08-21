# Overview

This is a Test Management Platform built as a web application for organizing and tracking software testing activities. The platform provides functionality for managing test projects, test suites, test cases, and test runs with a comprehensive dashboard and reporting system. It is designed for QA teams and software testers to organize their testing workflows and track test execution results.

## Recent Changes (August 2025)
- Built working Go implementation using standard library (simple_main.go)
- Added complete Docker containerization support for local deployment
- Created multi-stage Dockerfile with optimized Alpine Linux build
- Added Docker Compose configuration for easy local development
- Implemented in-memory data storage with sample test data
- Successfully deployed on port 5000 with web interface

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

Based on the template structure and JavaScript, the backend appears to be built with Go, using Go's html/template package for server-side rendering. The architecture follows a traditional MVC pattern.

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