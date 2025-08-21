# Test Management Platform

A manual test management platform built with Go featuring test case creation, execution tracking, and basic reporting capabilities.

## Features

- **Project Management**: Organize tests into projects
- **Test Suites**: Group related test cases together  
- **Test Case Management**: Create and manage individual test cases
- **Test Execution**: Track test runs and results
- **Reporting**: Basic reporting and analytics
- **Web Interface**: Clean, responsive web UI

## Running Locally

### Prerequisites

- Docker and Docker Compose installed on your machine

### Quick Start with Docker

1. **Clone or download the project files**

2. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

3. **Access the application:**
   - Open your browser and go to: `http://localhost:5000`

4. **Stop the application:**
   ```bash
   docker-compose down
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

- **Backend**: Go (standard library)
- **Frontend**: HTML templates with Bootstrap 5
- **Storage**: In-memory (for demo purposes)
- **Containerization**: Docker

## Notes

- This is a demo version using in-memory storage
- Data will be lost when the container is stopped
- For production use, consider adding a database backend