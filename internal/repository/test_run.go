package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/galex-do/test-machine/internal/models"
)

// TestRunRepository handles database operations for test runs
type TestRunRepository struct {
	db *sql.DB
}

// NewTestRunRepository creates a new test run repository
func NewTestRunRepository(db *sql.DB) *TestRunRepository {
	return &TestRunRepository{db: db}
}

// parseNullString converts sql.NullString to *string
func parseNullString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

// parseNullTime converts sql.NullTime to *time.Time
func parseNullTime(t sql.NullTime) *models.TestRun {
	// This function signature seems wrong, let me fix it
	return nil
}

// GetAll returns all test runs, optionally filtered by test case ID
func (r *TestRunRepository) GetAll(testCaseID *int) ([]models.TestRun, error) {
	var query string
	var args []interface{}

	if testCaseID != nil {
		query = `
			SELECT tr.id, tr.name, tr.description, tr.test_case_id, tr.status, tr.result, tr.notes, tr.executed_by, tr.started_at, tr.completed_at, tr.created_at, tr.updated_at,
			       tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at
			FROM test_runs tr
			JOIN test_cases tc ON tr.test_case_id = tc.id
			WHERE tr.test_case_id = $1
			ORDER BY tr.created_at DESC
		`
		args = []interface{}{*testCaseID}
	} else {
		query = `
			SELECT tr.id, tr.name, tr.description, tr.test_case_id, tr.status, tr.result, tr.notes, tr.executed_by, tr.started_at, tr.completed_at, tr.created_at, tr.updated_at,
			       tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at
			FROM test_runs tr
			JOIN test_cases tc ON tr.test_case_id = tc.id
			ORDER BY tr.created_at DESC
		`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testRuns []models.TestRun
	for rows.Next() {
		var tr models.TestRun
		var tc models.TestCase
		var result, notes, executedBy sql.NullString
		var startedAt, completedAt sql.NullTime

		err := rows.Scan(
			&tr.ID, &tr.Name, &tr.Description, &tr.TestCaseID, &tr.Status, &result, &notes, &executedBy, &startedAt, &completedAt, &tr.CreatedAt, &tr.UpdatedAt,
			&tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tr.Result = parseNullString(result)
		tr.Notes = parseNullString(notes)
		tr.ExecutedBy = parseNullString(executedBy)
		if startedAt.Valid {
			tr.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			tr.CompletedAt = &completedAt.Time
		}
		tr.TestCase = &tc
		testRuns = append(testRuns, tr)
	}

	return testRuns, nil
}

// GetByID returns a test run by ID
func (r *TestRunRepository) GetByID(id int) (*models.TestRun, error) {
	var tr models.TestRun
	var tc models.TestCase
	var result, notes, executedBy sql.NullString
	var startedAt, completedAt sql.NullTime

	err := r.db.QueryRow(`
		SELECT tr.id, tr.name, tr.description, tr.test_case_id, tr.status, tr.result, tr.notes, tr.executed_by, tr.started_at, tr.completed_at, tr.created_at, tr.updated_at,
		       tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at
		FROM test_runs tr
		JOIN test_cases tc ON tr.test_case_id = tc.id
		WHERE tr.id = $1
	`, id).Scan(
		&tr.ID, &tr.Name, &tr.Description, &tr.TestCaseID, &tr.Status, &result, &notes, &executedBy, &startedAt, &completedAt, &tr.CreatedAt, &tr.UpdatedAt,
		&tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tr.Result = parseNullString(result)
	tr.Notes = parseNullString(notes)
	tr.ExecutedBy = parseNullString(executedBy)
	if startedAt.Valid {
		tr.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		tr.CompletedAt = &completedAt.Time
	}
	tr.TestCase = &tc
	return &tr, nil
}

// Create creates a new test run
func (r *TestRunRepository) Create(req *models.CreateTestRunRequest) (*models.TestRun, error) {
	var tr models.TestRun
	var result, notes, executedBy sql.NullString
	var startedAt, completedAt sql.NullTime

	err := r.db.QueryRow(
		"INSERT INTO test_runs (name, description, test_case_id, executed_by) VALUES ($1, $2, $3, $4) RETURNING id, name, description, test_case_id, status, result, notes, executed_by, started_at, completed_at, created_at, updated_at",
		req.Name, req.Description, req.TestCaseID, req.ExecutedBy,
	).Scan(&tr.ID, &tr.Name, &tr.Description, &tr.TestCaseID, &tr.Status, &result, &notes, &executedBy, &startedAt, &completedAt, &tr.CreatedAt, &tr.UpdatedAt)

	if err != nil {
		return nil, err
	}

	tr.Result = parseNullString(result)
	tr.Notes = parseNullString(notes)
	tr.ExecutedBy = parseNullString(executedBy)
	if startedAt.Valid {
		tr.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		tr.CompletedAt = &completedAt.Time
	}

	return &tr, nil
}

// Update updates an existing test run
func (r *TestRunRepository) Update(id int, req *models.UpdateTestRunRequest) (*models.TestRun, error) {
	// Build dynamic update query
	setParts := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{}
	argIndex := 1

	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}
	if req.Result != nil {
		setParts = append(setParts, fmt.Sprintf("result = $%d", argIndex))
		args = append(args, *req.Result)
		argIndex++
	}
	if req.Notes != nil {
		setParts = append(setParts, fmt.Sprintf("notes = $%d", argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}
	if req.StartedAt != nil {
		setParts = append(setParts, fmt.Sprintf("started_at = $%d", argIndex))
		args = append(args, *req.StartedAt)
		argIndex++
	}
	if req.CompletedAt != nil {
		setParts = append(setParts, fmt.Sprintf("completed_at = $%d", argIndex))
		args = append(args, *req.CompletedAt)
		argIndex++
	}

	query := fmt.Sprintf(
		"UPDATE test_runs SET %s WHERE id = $%d RETURNING id, name, description, test_case_id, status, result, notes, executed_by, started_at, completed_at, created_at, updated_at",
		strings.Join(setParts, ", "), argIndex,
	)
	args = append(args, id)

	var tr models.TestRun
	var result, notes, executedBy sql.NullString
	var startedAt, completedAt sql.NullTime

	err := r.db.QueryRow(query, args...).Scan(
		&tr.ID, &tr.Name, &tr.Description, &tr.TestCaseID, &tr.Status, &result, &notes, &executedBy, &startedAt, &completedAt, &tr.CreatedAt, &tr.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tr.Result = parseNullString(result)
	tr.Notes = parseNullString(notes)
	tr.ExecutedBy = parseNullString(executedBy)
	if startedAt.Valid {
		tr.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		tr.CompletedAt = &completedAt.Time
	}

	return &tr, nil
}