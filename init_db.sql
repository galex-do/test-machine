-- Test Management Platform Database Schema
-- Updated to reflect current database structure

-- Keys table for authentication credentials
CREATE TABLE IF NOT EXISTS keys (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    key_type VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    encrypted_data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Repositories table for Git repository management
CREATE TABLE IF NOT EXISTS repositories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT 'Untitled Repository',
    description TEXT DEFAULT '',
    remote_url VARCHAR(255) NOT NULL,
    key_id INTEGER REFERENCES keys(id) ON DELETE SET NULL,
    default_branch VARCHAR(255),
    synced_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_repository_url UNIQUE (remote_url)
);

-- Projects table (updated with repository reference)
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    repository_id INTEGER REFERENCES repositories(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Branches table for Git branch tracking
CREATE TABLE IF NOT EXISTS branches (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    commit_hash VARCHAR(255),
    commit_date TIMESTAMP,
    commit_message TEXT,
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT branches_repository_id_name_key UNIQUE (repository_id, name)
);

-- Tags table for Git tag tracking
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    repository_id INTEGER NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    commit_hash VARCHAR(255),
    commit_date TIMESTAMP,
    commit_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT tags_repository_id_name_key UNIQUE (repository_id, name)
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

-- Test Steps table
CREATE TABLE IF NOT EXISTS test_steps (
    id SERIAL PRIMARY KEY,
    test_case_id INTEGER NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    expected_result TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT test_steps_test_case_id_step_number_key UNIQUE(test_case_id, step_number)
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
CREATE INDEX IF NOT EXISTS idx_repositories_key_id ON repositories(key_id);
CREATE INDEX IF NOT EXISTS idx_projects_repository_id ON projects(repository_id);
CREATE INDEX IF NOT EXISTS idx_branches_repository_id ON branches(repository_id);
CREATE INDEX IF NOT EXISTS idx_tags_repository_id ON tags(repository_id);
CREATE INDEX IF NOT EXISTS idx_test_suites_project_id ON test_suites(project_id);
CREATE INDEX IF NOT EXISTS idx_test_cases_test_suite_id ON test_cases(test_suite_id);
CREATE INDEX IF NOT EXISTS idx_test_steps_test_case_id ON test_steps(test_case_id);
CREATE INDEX IF NOT EXISTS idx_test_runs_test_case_id ON test_runs(test_case_id);
CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_repositories_name ON repositories(name);
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

    IF NOT EXISTS (SELECT 1 FROM test_steps) THEN
        INSERT INTO test_steps (test_case_id, step_number, description, expected_result) VALUES
            (1, 1, 'Navigate to the login page', 'Login page loads successfully with username and password fields'),
            (1, 2, 'Enter valid username and password', 'Credentials are accepted and fields show no validation errors'),
            (1, 3, 'Click the login button', 'User is redirected to dashboard and sees welcome message'),
            (2, 1, 'Navigate to the login page', 'Login page loads successfully'),
            (2, 2, 'Enter valid username but invalid password', 'Password field shows validation error'),
            (2, 3, 'Click the login button', 'Error message displays: "Invalid password"'),
            (7, 1, 'Open the mobile application', 'Application launches successfully and shows the main screen'),
            (7, 2, 'Tap on the navigation menu button', 'Navigation menu slides out from the left side'),
            (7, 3, 'Verify all menu items are visible', 'All navigation options are displayed: Home, Profile, Settings, Logout'),
            (8, 1, 'Navigate from home screen to profile', 'Screen transitions smoothly with appropriate animation'),
            (8, 2, 'Return to home screen using back button', 'Transition back is smooth and maintains app state');
    END IF;
END $$;