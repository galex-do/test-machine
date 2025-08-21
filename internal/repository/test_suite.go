package repository

import (
        "database/sql"

        "github.com/galex-do/test-machine/internal/models"
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
                        ORDER BY ts.created_at DESC
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
                        ORDER BY ts.created_at DESC
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

        return testSuites, nil
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