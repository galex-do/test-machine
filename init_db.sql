-- Test Management Platform Database Schema

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Test Suites table
CREATE TABLE IF NOT EXISTS test_suites (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Test Cases table
CREATE TABLE IF NOT EXISTS test_cases (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    priority VARCHAR(50) DEFAULT 'Medium' CHECK (priority IN ('Low', 'Medium', 'High', 'Critical')),
    status VARCHAR(50) DEFAULT 'Active' CHECK (status IN ('Active', 'Inactive', 'Archived')),
    test_suite_id INTEGER NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Test Runs table
CREATE TABLE IF NOT EXISTS test_runs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    test_case_id INTEGER NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'Not Started' CHECK (status IN ('Not Started', 'In Progress', 'Completed')),
    result VARCHAR(50) CHECK (result IN ('Pass', 'Fail', 'Blocked', 'Skip', 'Not Executed')),
    notes TEXT,
    executed_by VARCHAR(255),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_test_suites_project_id ON test_suites(project_id);
CREATE INDEX IF NOT EXISTS idx_test_cases_test_suite_id ON test_cases(test_suite_id);
CREATE INDEX IF NOT EXISTS idx_test_runs_test_case_id ON test_runs(test_case_id);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_test_cases_status ON test_cases(status);
CREATE INDEX IF NOT EXISTS idx_test_runs_status ON test_runs(status);

-- Sample data insertion (only if tables are empty)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM projects) THEN
        INSERT INTO projects (name, description) VALUES
            ('Web Application Testing', 'Testing suite for the main web application'),
            ('Mobile App Testing', 'Testing suite for the mobile application');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM test_suites) THEN
        INSERT INTO test_suites (name, description, project_id) VALUES
            ('User Authentication', 'Test cases for user login, registration, and password reset', 1),
            ('E-commerce Checkout', 'Test cases for shopping cart and checkout process', 1),
            ('Mobile Navigation', 'Test cases for mobile app navigation and UI', 2);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM test_cases) THEN
        INSERT INTO test_cases (title, description, priority, test_suite_id) VALUES
            ('Valid User Login', 'Test successful login with valid credentials', 'High', 1),
            ('Invalid Password Login', 'Test login failure with invalid password', 'High', 1),
            ('Password Reset Flow', 'Test password reset functionality', 'Medium', 1),
            ('Add Item to Cart', 'Test adding products to shopping cart', 'Medium', 2),
            ('Checkout Process', 'Test complete checkout flow', 'High', 2),
            ('Payment Processing', 'Test payment gateway integration', 'Critical', 2),
            ('Navigation Menu', 'Test mobile app navigation menu', 'High', 3),
            ('Screen Transitions', 'Test smooth transitions between screens', 'Medium', 3);
    END IF;
END $$;