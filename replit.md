# Overview

This is a Test Management Platform built as a web application for organizing and tracking software testing activities. The platform provides functionality for managing test projects, test suites, test cases, and test runs with a comprehensive dashboard and reporting system. It appears to be designed for QA teams and software testers to organize their testing workflows and track test execution results.

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

**Deployment Considerations:**
- Static file serving for CSS, JavaScript, and other assets
- Template compilation and caching
- Database connection management
- Environment-based configuration