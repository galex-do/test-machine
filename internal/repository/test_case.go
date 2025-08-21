package repository

import (
	"database/sql"

	"github.com/galex-do/test-machine/internal/models"
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
			       p.id, p.name, p.description, p.created_at, p.updated_at
			FROM test_cases tc
			JOIN test_suites ts ON tc.test_suite_id = ts.id
			JOIN projects p ON ts.project_id = p.id
			WHERE tc.test_suite_id = $1
			ORDER BY tc.created_at DESC
		`
		args = []interface{}{*testSuiteID}
	} else {
		query = `
			SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
			       ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
			       p.id, p.name, p.description, p.created_at, p.updated_at
			FROM test_cases tc
			JOIN test_suites ts ON tc.test_suite_id = ts.id
			JOIN projects p ON ts.project_id = p.id
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