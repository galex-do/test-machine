package repository

import (
        "database/sql"
        "fmt"

        "github.com/galex-do/test-machine/internal/models"
        "github.com/galex-do/test-machine/internal/utils"
)

// TestSuiteRepository handles database operations for test suites
type TestSuiteRepository struct {
        db *sql.DB
}

// NewTestSuiteRepository creates a new test suite repository
func NewTestSuiteRepository(db *sql.DB) *TestSuiteRepository {
        return &TestSuiteRepository{db: db}
}

// GetAll returns all test suites with test case counts, optionally filtered by project ID
func (r *TestSuiteRepository) GetAll(projectID *int) ([]models.TestSuite, error) {
        var query string
        var args []interface{}

        if projectID != nil {
                query = `
                        SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(COUNT(tc.id), 0) as test_cases_count
                        FROM test_suites ts
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN test_cases tc ON ts.id = tc.test_suite_id
                        WHERE ts.project_id = $1
                        GROUP BY ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                 p.id, p.name, p.description, p.created_at, p.updated_at
                        ORDER BY ts.created_at ASC
                `
                args = []interface{}{*projectID}
        } else {
                query = `
                        SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(COUNT(tc.id), 0) as test_cases_count
                        FROM test_suites ts
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN test_cases tc ON ts.id = tc.test_suite_id
                        GROUP BY ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                 p.id, p.name, p.description, p.created_at, p.updated_at
                        ORDER BY ts.created_at ASC
                `
        }

        rows, err := r.db.Query(query, args...)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var testSuites []models.TestSuite
        for rows.Next() {
                var ts models.TestSuite
                var p models.Project
                err := rows.Scan(
                        &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                        &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                        &ts.TestCasesCount,
                )
                if err != nil {
                        return nil, err
                }
                ts.Project = &p
                testSuites = append(testSuites, ts)
        }

        // Load test cases for each test suite when filtering by project (for test run creation)
        if projectID != nil {
                for i := range testSuites {
                        testCases, err := r.getTestCasesByTestSuite(testSuites[i].ID)
                        if err != nil {
                                return nil, err
                        }
                        testSuites[i].TestCases = testCases
                }
        }

        return testSuites, nil
}

// GetAllPaginated returns test suites with pagination
func (r *TestSuiteRepository) GetAllPaginated(pagination models.PaginationRequest, projectID *int) (*models.PaginatedResult, error) {
        // Build count query
        var countQuery string
        var countArgs []interface{}
        
        if projectID != nil {
                countQuery = `SELECT COUNT(*) FROM test_suites WHERE project_id = $1`
                countArgs = []interface{}{*projectID}
        } else {
                countQuery = `SELECT COUNT(*) FROM test_suites`
        }
        
        // Get total count
        var total int
        err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
        if err != nil {
                return nil, fmt.Errorf("failed to count test suites: %w", err)
        }

        // Calculate pagination
        offset, limit := utils.GetOffsetAndLimit(pagination.Page, pagination.PageSize)
        paginationResp := utils.CalculatePagination(pagination.Page, pagination.PageSize, total)

        // Build data query
        var query string
        var args []interface{}

        if projectID != nil {
                query = `
                        SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(COUNT(tc.id), 0) as test_cases_count
                        FROM test_suites ts
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN test_cases tc ON ts.id = tc.test_suite_id
                        WHERE ts.project_id = $1
                        GROUP BY ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                 p.id, p.name, p.description, p.created_at, p.updated_at
                        ORDER BY ts.created_at ASC
                        LIMIT $2 OFFSET $3
                `
                args = []interface{}{*projectID, limit, offset}
        } else {
                query = `
                        SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at,
                               COALESCE(COUNT(tc.id), 0) as test_cases_count
                        FROM test_suites ts
                        JOIN projects p ON ts.project_id = p.id
                        LEFT JOIN test_cases tc ON ts.id = tc.test_suite_id
                        GROUP BY ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                 p.id, p.name, p.description, p.created_at, p.updated_at
                        ORDER BY ts.created_at ASC
                        LIMIT $1 OFFSET $2
                `
                args = []interface{}{limit, offset}
        }

        rows, err := r.db.Query(query, args...)
        if err != nil {
                return nil, fmt.Errorf("failed to get test suites: %w", err)
        }
        defer rows.Close()

        var testSuites []models.TestSuite
        for rows.Next() {
                var ts models.TestSuite
                var p models.Project
                err := rows.Scan(
                        &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                        &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                        &ts.TestCasesCount,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan test suite: %w", err)
                }
                ts.Project = &p
                testSuites = append(testSuites, ts)
        }

        return &models.PaginatedResult{
                Data:       testSuites,
                Pagination: paginationResp,
        }, nil
}

// GetByID returns a test suite by ID with test case count
func (r *TestSuiteRepository) GetByID(id int) (*models.TestSuite, error) {
        var ts models.TestSuite
        var p models.Project
        err := r.db.QueryRow(`
                SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                       p.id, p.name, p.description, p.created_at, p.updated_at,
                       COALESCE(COUNT(tc.id), 0) as test_cases_count
                FROM test_suites ts
                JOIN projects p ON ts.project_id = p.id
                LEFT JOIN test_cases tc ON ts.id = tc.test_suite_id
                WHERE ts.id = $1
                GROUP BY ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                         p.id, p.name, p.description, p.created_at, p.updated_at
        `, id).Scan(
                &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                &ts.TestCasesCount,
        )

        if err == sql.ErrNoRows {
                return nil, nil
        }
        if err != nil {
                return nil, err
        }

        ts.Project = &p
        return &ts, nil
}

// Create creates a new test suite
func (r *TestSuiteRepository) Create(req *models.CreateTestSuiteRequest) (*models.TestSuite, error) {
        var testSuite models.TestSuite
        err := r.db.QueryRow(
                "INSERT INTO test_suites (name, description, project_id) VALUES ($1, $2, $3) RETURNING id, name, description, project_id, created_at, updated_at",
                req.Name, req.Description, req.ProjectID,
        ).Scan(&testSuite.ID, &testSuite.Name, &testSuite.Description, &testSuite.ProjectID, &testSuite.CreatedAt, &testSuite.UpdatedAt)

        if err != nil {
                return nil, err
        }

        return &testSuite, nil
}

// Update updates an existing test suite
func (r *TestSuiteRepository) Update(id int, req *models.UpdateTestSuiteRequest) (*models.TestSuite, error) {
        var testSuite models.TestSuite
        err := r.db.QueryRow(
                "UPDATE test_suites SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, name, description, project_id, created_at, updated_at",
                req.Name, req.Description, id,
        ).Scan(&testSuite.ID, &testSuite.Name, &testSuite.Description, &testSuite.ProjectID, &testSuite.CreatedAt, &testSuite.UpdatedAt)

        if err == sql.ErrNoRows {
                return nil, nil
        }
        if err != nil {
                return nil, err
        }

        return &testSuite, nil
}

// Delete deletes a test suite
func (r *TestSuiteRepository) Delete(id int) error {
        result, err := r.db.Exec("DELETE FROM test_suites WHERE id = $1", id)
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

// getTestCasesByTestSuite loads test cases for a specific test suite
func (r *TestSuiteRepository) getTestCasesByTestSuite(testSuiteID int) ([]models.TestCase, error) {
        query := `
                SELECT id, title, description, priority, status, test_suite_id, created_at, updated_at
                FROM test_cases
                WHERE test_suite_id = $1
                ORDER BY title
        `

        rows, err := r.db.Query(query, testSuiteID)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var testCases []models.TestCase
        for rows.Next() {
                var tc models.TestCase
                err = rows.Scan(&tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt)
                if err != nil {
                        return nil, err
                }
                testCases = append(testCases, tc)
        }

        return testCases, nil
}
