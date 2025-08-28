package repository

import (
        "database/sql"
        "fmt"

        "github.com/galex-do/test-machine/internal/models"
        "github.com/galex-do/test-machine/internal/utils"
)

// TestCaseRepository handles database operations for test cases
type TestCaseRepository struct {
        db *sql.DB
}

// NewTestCaseRepository creates a new test case repository
func NewTestCaseRepository(db *sql.DB) *TestCaseRepository {
        return &TestCaseRepository{db: db}
}

// GetAll returns all test cases, optionally filtered by test suite ID
func (r *TestCaseRepository) GetAll(testSuiteID *int) ([]models.TestCase, error) {
        var query string
        var args []interface{}

        if testSuiteID != nil {
                query = `
                        SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                               ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(step_counts.step_count, 0) as test_steps_count
                        FROM test_cases tc
                        JOIN test_suites ts ON tc.test_suite_id = ts.id
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN (
                                SELECT test_case_id, COUNT(*) as step_count
                                FROM test_steps
                                GROUP BY test_case_id
                        ) step_counts ON tc.id = step_counts.test_case_id
                        WHERE tc.test_suite_id = $1
                        ORDER BY tc.created_at DESC
                `
                args = []interface{}{*testSuiteID}
        } else {
                query = `
                        SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                               ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(step_counts.step_count, 0) as test_steps_count
                        FROM test_cases tc
                        JOIN test_suites ts ON tc.test_suite_id = ts.id
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN (
                                SELECT test_case_id, COUNT(*) as step_count
                                FROM test_steps
                                GROUP BY test_case_id
                        ) step_counts ON tc.id = step_counts.test_case_id
                        ORDER BY tc.created_at DESC
                `
        }

        rows, err := r.db.Query(query, args...)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var testCases []models.TestCase
        for rows.Next() {
                var tc models.TestCase
                var ts models.TestSuite
                var p models.Project
                err := rows.Scan(
                        &tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt,
                        &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                        &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                        &tc.TestStepsCount,
                )
                if err != nil {
                        return nil, err
                }
                ts.Project = &p
                tc.TestSuite = &ts
                testCases = append(testCases, tc)
        }

        return testCases, nil
}

// GetAllPaginated returns test cases with pagination
func (r *TestCaseRepository) GetAllPaginated(pagination models.PaginationRequest, testSuiteID *int) (*models.PaginatedResult, error) {
        // Build count query
        var countQuery string
        var countArgs []interface{}
        
        if testSuiteID != nil {
                countQuery = `SELECT COUNT(*) FROM test_cases WHERE test_suite_id = $1`
                countArgs = []interface{}{*testSuiteID}
        } else {
                countQuery = `SELECT COUNT(*) FROM test_cases`
        }
        
        // Get total count
        var total int
        err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
        if err != nil {
                return nil, fmt.Errorf("failed to count test cases: %w", err)
        }

        // Calculate pagination
        offset, limit := utils.GetOffsetAndLimit(pagination.Page, pagination.PageSize)
        paginationResp := utils.CalculatePagination(pagination.Page, pagination.PageSize, total)

        // Build data query
        var query string
        var args []interface{}

        if testSuiteID != nil {
                query = `
                        SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                               ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(step_counts.step_count, 0) as test_steps_count
                        FROM test_cases tc
                        JOIN test_suites ts ON tc.test_suite_id = ts.id
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN (
                                SELECT test_case_id, COUNT(*) as step_count
                                FROM test_steps
                                GROUP BY test_case_id
                        ) step_counts ON tc.id = step_counts.test_case_id
                        WHERE tc.test_suite_id = $1
                        ORDER BY tc.created_at DESC
                        LIMIT $2 OFFSET $3
                `
                args = []interface{}{*testSuiteID, limit, offset}
        } else {
                query = `
                        SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                               ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(step_counts.step_count, 0) as test_steps_count
                        FROM test_cases tc
                        JOIN test_suites ts ON tc.test_suite_id = ts.id
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN (
                                SELECT test_case_id, COUNT(*) as step_count
                                FROM test_steps
                                GROUP BY test_case_id
                        ) step_counts ON tc.id = step_counts.test_case_id
                        ORDER BY tc.created_at DESC
                        LIMIT $1 OFFSET $2
                `
                args = []interface{}{limit, offset}
        }

        rows, err := r.db.Query(query, args...)
        if err != nil {
                return nil, fmt.Errorf("failed to get test cases: %w", err)
        }
        defer rows.Close()

        var testCases []models.TestCase
        for rows.Next() {
                var tc models.TestCase
                var ts models.TestSuite
                var p models.Project
                err := rows.Scan(
                        &tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt,
                        &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                        &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                        &tc.TestStepsCount,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan test case: %w", err)
                }
                ts.Project = &p
                tc.TestSuite = &ts
                testCases = append(testCases, tc)
        }

        return &models.PaginatedResult{
                Data:       testCases,
                Pagination: paginationResp,
        }, nil
}

// GetByID returns a test case by ID
func (r *TestCaseRepository) GetByID(id int) (*models.TestCase, error) {
        var tc models.TestCase
        var ts models.TestSuite
        var p models.Project
        err := r.db.QueryRow(`
                SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                       ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                       p.id, p.name, p.description, p.created_at, p.updated_at
                FROM test_cases tc
                JOIN test_suites ts ON tc.test_suite_id = ts.id
                JOIN projects p ON ts.project_id = p.id
                WHERE tc.id = $1
        `, id).Scan(
                &tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt,
                &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
        )

        if err == sql.ErrNoRows {
                return nil, nil
        }
        if err != nil {
                return nil, err
        }

        ts.Project = &p
        tc.TestSuite = &ts
        return &tc, nil
}

// Create creates a new test case
func (r *TestCaseRepository) Create(req *models.CreateTestCaseRequest) (*models.TestCase, error) {
        // Set default priority if not provided
        priority := req.Priority
        if priority == "" {
                priority = "Medium"
        }

        var testCase models.TestCase
        err := r.db.QueryRow(
                "INSERT INTO test_cases (title, description, priority, test_suite_id) VALUES ($1, $2, $3, $4) RETURNING id, title, description, priority, status, test_suite_id, created_at, updated_at",
                req.Title, req.Description, priority, req.TestSuiteID,
        ).Scan(&testCase.ID, &testCase.Title, &testCase.Description, &testCase.Priority, &testCase.Status, &testCase.TestSuiteID, &testCase.CreatedAt, &testCase.UpdatedAt)

        if err != nil {
                return nil, err
        }

        return &testCase, nil
}

// GetTestSteps returns all test steps for a test case
func (r *TestCaseRepository) GetTestSteps(testCaseID int) ([]models.TestStep, error) {
        query := `
                SELECT id, test_case_id, step_number, description, expected_result, created_at, updated_at
                FROM test_steps
                WHERE test_case_id = $1
                ORDER BY step_number ASC
        `
        
        rows, err := r.db.Query(query, testCaseID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()
        
        var testSteps []models.TestStep
        for rows.Next() {
                var step models.TestStep
                err := rows.Scan(
                        &step.ID, &step.TestCaseID, &step.StepNumber, &step.Description, 
                        &step.ExpectedResult, &step.CreatedAt, &step.UpdatedAt,
                )
                if err != nil {
                        return nil, err
                }
                testSteps = append(testSteps, step)
        }
        
        if err = rows.Err(); err != nil {
                return nil, err
        }
        
        return testSteps, nil
}

// CreateTestStep creates a new test step
func (r *TestCaseRepository) CreateTestStep(req *models.CreateTestStepRequest) (*models.TestStep, error) {
        var testStep models.TestStep
        err := r.db.QueryRow(
                "INSERT INTO test_steps (test_case_id, step_number, description, expected_result) VALUES ($1, $2, $3, $4) RETURNING id, test_case_id, step_number, description, expected_result, created_at, updated_at",
                req.TestCaseID, req.StepNumber, req.Description, req.ExpectedResult,
        ).Scan(&testStep.ID, &testStep.TestCaseID, &testStep.StepNumber, &testStep.Description, &testStep.ExpectedResult, &testStep.CreatedAt, &testStep.UpdatedAt)

        if err != nil {
                return nil, err
        }

        return &testStep, nil
}

// UpdateTestStep updates an existing test step
func (r *TestCaseRepository) UpdateTestStep(id int, req *models.UpdateTestStepRequest) (*models.TestStep, error) {
        var testStep models.TestStep
        err := r.db.QueryRow(
                "UPDATE test_steps SET step_number = $1, description = $2, expected_result = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 RETURNING id, test_case_id, step_number, description, expected_result, created_at, updated_at",
                req.StepNumber, req.Description, req.ExpectedResult, id,
        ).Scan(&testStep.ID, &testStep.TestCaseID, &testStep.StepNumber, &testStep.Description, &testStep.ExpectedResult, &testStep.CreatedAt, &testStep.UpdatedAt)

        if err != nil {
                return nil, err
        }

        return &testStep, nil
}

// Update updates an existing test case
func (r *TestCaseRepository) Update(id int, req *models.UpdateTestCaseRequest) (*models.TestCase, error) {
        var testCase models.TestCase
        err := r.db.QueryRow(
                "UPDATE test_cases SET title = $1, description = $2, priority = $3, status = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5 RETURNING id, title, description, priority, status, test_suite_id, created_at, updated_at",
                req.Title, req.Description, req.Priority, req.Status, id,
        ).Scan(&testCase.ID, &testCase.Title, &testCase.Description, &testCase.Priority, &testCase.Status, &testCase.TestSuiteID, &testCase.CreatedAt, &testCase.UpdatedAt)

        if err == sql.ErrNoRows {
                return nil, nil
        }
        if err != nil {
                return nil, err
        }

        return &testCase, nil
}

// DeleteTestStep deletes a test step
func (r *TestCaseRepository) DeleteTestStep(id int) error {
        result, err := r.db.Exec("DELETE FROM test_steps WHERE id = $1", id)
        if err != nil {
                return err
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return err
        }

        if rowsAffected == 0 {
                return sql.ErrNoRows
        }

        return nil
}

// Delete deletes a test case and all its related test steps
func (r *TestCaseRepository) Delete(id int) error {
        tx, err := r.db.Begin()
        if err != nil {
                return err
        }
        defer tx.Rollback()

        // First delete all test steps related to this test case
        _, err = tx.Exec("DELETE FROM test_steps WHERE test_case_id = $1", id)
        if err != nil {
                return err
        }

        // Then delete the test case itself
        result, err := tx.Exec("DELETE FROM test_cases WHERE id = $1", id)
        if err != nil {
                return err
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return err
        }

        if rowsAffected == 0 {
                return sql.ErrNoRows
        }

        return tx.Commit()
}