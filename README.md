# Test Management Platform

A manual test management platform built with Go and PostgreSQL featuring test case creation, execution tracking, and comprehensive reporting capabilities.

## Features

- **Project Management**: Organize tests into projects
- **Test Suites**: Group related test cases together  
- **Test Case Management**: Create and manage individual test cases
- **Test Execution**: Track test runs and results with detailed metadata
- **Reporting**: Basic reporting and analytics
- **Web Interface**: Clean, responsive web UI
- **PostgreSQL Database**: Persistent data storage with relational integrity
- **Docker Ready**: Full containerization with database orchestration

## Running Locally

### Prerequisites

- Docker and Docker Compose installed on your machine

### Quick Start with Docker Compose

1. **Clone or download the project files**

2. **Start the full stack (PostgreSQL + App):**
   ```bash
   docker-compose up --build
   ```
   This will:
   - Start a PostgreSQL database container
   - Build and start the Go application
   - Automatically create database schema and sample data
   - Make the app available on port 5000

3. **Access the application:**
   - Web Interface: `http://localhost:5000`
   - PostgreSQL Database: `localhost:5432` (if needed)

4. **Stop the application:**
   ```bash
   docker-compose down
   ```

5. **Remove all data (including database):**
   ```bash
   docker-compose down -v
   ```

### Alternative Docker Commands

If you prefer to use Docker directly without Docker Compose:

1. **Build the image:**
   ```bash
   docker build -t test-management .
   ```

2. **Run the container:**
   ```bash
   docker run -p 5000:5000 test-management
   ```

3. **Stop the container:**
   ```bash
   docker stop <container_id>
   ```

## Development

The application includes sample data with:
- Two projects: "Web Application Testing" and "Mobile App Testing"
- Test suites for user authentication and e-commerce checkout
- Sample test cases with different priorities and statuses

### Project Structure

```
.
├── simple_main.go          # Main application file
├── templates/              # HTML templates
│   ├── index.html          # Dashboard
│   ├── project.html        # Project details
│   ├── test-suite.html     # Test suite management
│   ├── test-case.html      # Test case details
│   ├── test-run.html       # Test execution
│   └── reports.html        # Reports and analytics
├── static/                 # Static assets
│   ├── css/style.css       # Custom styles
│   └── js/app.js          # Client-side JavaScript
├── Dockerfile              # Docker build instructions
└── docker-compose.yml     # Docker Compose configuration
```

### API Endpoints

The application provides REST API endpoints:

- `GET /api/projects` - List all projects
- `POST /api/projects` - Create new project
- `GET /api/projects/{id}` - Get project details
- `GET /api/test-suites` - List test suites
- `POST /api/test-suites` - Create new test suite
- `GET /api/test-cases` - List test cases
- `POST /api/test-cases` - Create new test case

## Technology Stack

- **Backend**: Go with PostgreSQL driver (lib/pq)
- **Database**: PostgreSQL 15 with full ACID compliance
- **Frontend**: HTML templates with Bootstrap 5
- **Containerization**: Docker with multi-service orchestration
- **Development**: Hot-reload capable, persistent data storage

## Database Schema

The PostgreSQL database includes the following tables:
- **Projects**: Top-level test organization
- **Test Suites**: Grouped collections of test cases
- **Test Cases**: Individual test scenarios with priority and status
- **Test Runs**: Execution records with results and metadata

All tables include proper foreign key relationships, indexes for performance, and audit timestamps.

## Configuration

### Database Connection
The application uses the `DATABASE_URL` environment variable to connect to PostgreSQL:
```
DATABASE_URL=postgres://testmgmt:password@postgres:5432/testmanagement?sslmode=disable
```

### Environment Variables
- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Application port (defaults to 5000)

## Development

### Sample Data
The application automatically creates sample data on first run:
- Two projects: "Web Application Testing" and "Mobile App Testing"  
- Three test suites covering authentication, checkout, and navigation
- Eight test cases with various priorities and descriptions

### API Endpoints
- `GET/POST /api/projects` - Project management
- `GET/POST /api/test-suites` - Test suite management
- `GET/POST /api/test-cases` - Test case management
- All endpoints support full CRUD operations with PostgreSQL persistence