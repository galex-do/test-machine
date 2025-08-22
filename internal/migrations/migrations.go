package migrations

import (
        "database/sql"
        "fmt"
        "log"
)

// Migration represents a database migration
type Migration struct {
        Version     int
        Description string
        Query       string
}

// Migrator handles database migrations
type Migrator struct {
        db *sql.DB
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
        return &Migrator{db: db}
}

// migrations contains all database migrations
var migrations = []Migration{
        {
                Version:     1,
                Description: "Add test_steps table",
                Query: `
-- Create test_steps table
CREATE TABLE IF NOT EXISTS test_steps (
    id SERIAL PRIMARY KEY,
    test_case_id INTEGER NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    expected_result TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(test_case_id, step_number)
);

-- Create index for test_steps
CREATE INDEX IF NOT EXISTS idx_test_steps_test_case_id ON test_steps(test_case_id);

-- Insert sample test steps data if test_steps table is empty
INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT * FROM (VALUES
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
    (8, 2, 'Return to home screen using back button', 'Transition back is smooth and maintains app state')
) AS v(test_case_id, step_number, description, expected_result)
WHERE NOT EXISTS (SELECT 1 FROM test_steps)
AND EXISTS (SELECT 1 FROM test_cases WHERE id = v.test_case_id);
                `,
        },
        {
                Version:     2,
                Description: "Add missing test steps for all test cases",
                Query: `
-- Add test steps for test case 2 (Invalid Password Login) if missing
INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 2, 1, 'Navigate to the login page', 'Login page loads successfully'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 2 AND step_number = 1);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 2, 2, 'Enter valid username but invalid password', 'Password field shows validation error'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 2 AND step_number = 2);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 2, 3, 'Click the login button', 'Error message displays: "Invalid password"'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 2 AND step_number = 3);

-- Add test steps for test case 3 (Password Reset Flow)
INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 3, 1, 'Navigate to login page and click "Forgot Password"', 'Password reset form is displayed'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 3 AND step_number = 1);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 3, 2, 'Enter valid email address', 'Email field accepts input and shows no validation errors'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 3 AND step_number = 2);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 3, 3, 'Click submit button', 'Success message shows: "Password reset link sent to your email"'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 3 AND step_number = 3);

-- Add test steps for test case 4 (Add Item to Cart)
INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 4, 1, 'Browse to product catalog page', 'Product catalog loads with available items'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 4 AND step_number = 1);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 4, 2, 'Select a product and click "Add to Cart"', 'Product is added to cart and cart counter increases'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 4 AND step_number = 2);

-- Add test steps for test case 5 (Checkout Process)
INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 5, 1, 'Add items to cart and navigate to checkout', 'Checkout page displays with cart items and total'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 5 AND step_number = 1);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 5, 2, 'Fill in shipping and billing information', 'Forms accept input and validate required fields'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 5 AND step_number = 2);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 5, 3, 'Review order and click "Place Order"', 'Order confirmation page displays with order number'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 5 AND step_number = 3);

-- Add test steps for test case 6 (Payment Processing)
INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 6, 1, 'Proceed to payment section during checkout', 'Payment form loads with supported payment methods'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 6 AND step_number = 1);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 6, 2, 'Enter valid credit card information', 'Payment form validates and accepts card details'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 6 AND step_number = 2);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 6, 3, 'Submit payment', 'Payment processes successfully and confirmation is displayed'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 6 AND step_number = 3);
                `,
        },
        {
                Version:     3,
                Description: "Fix test case 1 steps and clean up inconsistent data",
                Query: `
-- Remove any manual test steps that don't follow the standard format
DELETE FROM test_steps WHERE test_case_id = 1 AND description = 'wwgw';

-- Add proper test steps for test case 1 (Valid User Login)
INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 1, 1, 'Navigate to the login page', 'Login page loads successfully with username and password fields'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 1 AND step_number = 1);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 1, 2, 'Enter valid username and password', 'Credentials are accepted and fields show no validation errors'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 1 AND step_number = 2);

INSERT INTO test_steps (test_case_id, step_number, description, expected_result)
SELECT 1, 3, 'Click the login button', 'User is redirected to dashboard and sees welcome message'
WHERE NOT EXISTS (SELECT 1 FROM test_steps WHERE test_case_id = 1 AND step_number = 3);
                `,
        },
}

// Run executes all pending migrations
func (m *Migrator) Run() error {
        // Create migrations table if it doesn't exist
        if err := m.createMigrationsTable(); err != nil {
                return fmt.Errorf("failed to create migrations table: %v", err)
        }

        // Get current migration version
        currentVersion, err := m.getCurrentVersion()
        if err != nil {
                return fmt.Errorf("failed to get current version: %v", err)
        }

        log.Printf("Current migration version: %d", currentVersion)

        // Run pending migrations
        for _, migration := range migrations {
                if migration.Version > currentVersion {
                        log.Printf("Running migration %d: %s", migration.Version, migration.Description)
                        
                        if err := m.runMigration(migration); err != nil {
                                return fmt.Errorf("failed to run migration %d: %v", migration.Version, err)
                        }
                        
                        log.Printf("Migration %d completed successfully", migration.Version)
                }
        }

        log.Println("All migrations completed successfully")
        return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
        query := `
                CREATE TABLE IF NOT EXISTS migrations (
                        version INTEGER PRIMARY KEY,
                        description TEXT NOT NULL,
                        applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                );
        `
        _, err := m.db.Exec(query)
        return err
}

// getCurrentVersion returns the current migration version
func (m *Migrator) getCurrentVersion() (int, error) {
        var version int
        err := m.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM migrations").Scan(&version)
        if err != nil {
                return 0, err
        }
        return version, nil
}

// runMigration executes a single migration
func (m *Migrator) runMigration(migration Migration) error {
        tx, err := m.db.Begin()
        if err != nil {
                return err
        }
        defer tx.Rollback()

        // Execute migration query
        _, err = tx.Exec(migration.Query)
        if err != nil {
                return err
        }

        // Record migration in migrations table
        _, err = tx.Exec(
                "INSERT INTO migrations (version, description) VALUES ($1, $2)",
                migration.Version, migration.Description,
        )
        if err != nil {
                return err
        }

        return tx.Commit()
}