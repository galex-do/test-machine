package repository

import (
        "database/sql"
        "fmt"
        "strings"
        "time"

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

// GetAll returns all test runs with their basic information
func (r *TestRunRepository) GetAll() ([]models.TestRun, error) {
        query := `
                SELECT tr.id, tr.name, tr.description, tr.project_id, tr.repository_id, 
                       tr.branch_name, tr.tag_name, tr.status, tr.created_by, 
                       tr.started_at, tr.completed_at, tr.created_at, tr.updated_at,
                       p.id, p.name, p.description, p.created_at, p.updated_at,
                       r.id, r.name, r.description, r.remote_url, r.default_branch, 
                       r.synced_at, r.created_at, r.updated_at,
                       COALESCE(COUNT(trc.id), 0) as test_cases_count
                FROM test_runs tr
                JOIN projects p ON tr.project_id = p.id
                LEFT JOIN repositories r ON tr.repository_id = r.id
                LEFT JOIN test_run_cases trc ON tr.id = trc.test_run_id
                GROUP BY tr.id, tr.name, tr.description, tr.project_id, tr.repository_id, 
                         tr.branch_name, tr.tag_name, tr.status, tr.created_by, 
                         tr.started_at, tr.completed_at, tr.created_at, tr.updated_at,
                         p.id, p.name, p.description, p.created_at, p.updated_at,
                         r.id, r.name, r.description, r.remote_url, r.default_branch, 
                         r.synced_at, r.created_at, r.updated_at
                ORDER BY tr.created_at DESC
        `

        rows, err := r.db.Query(query)
        if err != nil {
                return nil, fmt.Errorf("failed to get test runs: %w", err)
        }
        defer rows.Close()

        var testRuns []models.TestRun
        for rows.Next() {
                var tr models.TestRun
                var project models.Project
                var repository models.Repository
                var testCasesCount int
                // Repository fields can be NULL, so use pointers
                var repoID sql.NullInt32
                var repoCreatedAt, repoUpdatedAt, repoSyncedAt sql.NullTime
                var repoName, repoDescription, repoRemoteURL, repoDefaultBranch sql.NullString
                
                err = rows.Scan(
                        &tr.ID, &tr.Name, &tr.Description, &tr.ProjectID, &tr.RepositoryID,
                        &tr.BranchName, &tr.TagName, &tr.Status, &tr.CreatedBy,
                        &tr.StartedAt, &tr.CompletedAt, &tr.CreatedAt, &tr.UpdatedAt,
                        &project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt,
                        &repoID, &repoName, &repoDescription, &repoRemoteURL, &repoDefaultBranch,
                        &repoSyncedAt, &repoCreatedAt, &repoUpdatedAt,
                        &testCasesCount,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan test run: %w", err)
                }
                
                tr.Project = &project
                tr.TestCasesCount = &testCasesCount
                
                // Only populate repository if it exists
                if repoID.Valid {
                        repository.ID = int(repoID.Int32)
                        repository.Name = repoName.String
                        repository.Description = repoDescription.String
                        repository.RemoteURL = repoRemoteURL.String
                        if repoDefaultBranch.Valid {
                                repository.DefaultBranch = &repoDefaultBranch.String
                        }
                        if repoSyncedAt.Valid {
                                repository.SyncedAt = &repoSyncedAt.Time
                        }
                        if repoCreatedAt.Valid {
                                repository.CreatedAt = repoCreatedAt.Time
                        }
                        if repoUpdatedAt.Valid {
                                repository.UpdatedAt = repoUpdatedAt.Time
                        }
                        tr.Repository = &repository
                }
                
                testRuns = append(testRuns, tr)
        }

        return testRuns, nil
}

// GetByID returns a test run by ID with all related data
func (r *TestRunRepository) GetByID(id int) (*models.TestRun, error) {
        query := `
                SELECT tr.id, tr.name, tr.description, tr.project_id, tr.repository_id, 
                       tr.branch_name, tr.tag_name, tr.status, tr.created_by, 
                       tr.started_at, tr.completed_at, tr.created_at, tr.updated_at,
                       p.id, p.name, p.description, p.created_at, p.updated_at
                FROM test_runs tr
                JOIN projects p ON tr.project_id = p.id
                WHERE tr.id = $1
        `

        var tr models.TestRun
        var project models.Project
        err := r.db.QueryRow(query, id).Scan(
                &tr.ID, &tr.Name, &tr.Description, &tr.ProjectID, &tr.RepositoryID,
                &tr.BranchName, &tr.TagName, &tr.Status, &tr.CreatedBy,
                &tr.StartedAt, &tr.CompletedAt, &tr.CreatedAt, &tr.UpdatedAt,
                &project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt,
        )
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, nil
                }
                return nil, fmt.Errorf("failed to get test run: %w", err)
        }

        tr.Project = &project

        // Load test cases for this run
        testCases, err := r.getTestRunCases(tr.ID)
        if err != nil {
                return nil, fmt.Errorf("failed to get test run cases: %w", err)
        }
        tr.TestCases = testCases

        return &tr, nil
}

// getTestRunCases loads test cases for a test run
func (r *TestRunRepository) getTestRunCases(testRunID int) ([]models.TestRunCase, error) {
        query := `
                SELECT trc.id, trc.test_run_id, trc.test_case_id, trc.status, trc.result_notes, 
                       trc.executed_by, trc.started_at, trc.completed_at, trc.created_at, trc.updated_at,
                       tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at
                FROM test_run_cases trc
                JOIN test_cases tc ON trc.test_case_id = tc.id
                WHERE trc.test_run_id = $1
                ORDER BY tc.title
        `

        rows, err := r.db.Query(query, testRunID)
        if err != nil {
                return nil, fmt.Errorf("failed to get test run cases: %w", err)
        }
        defer rows.Close()

        var testRunCases []models.TestRunCase
        for rows.Next() {
                var trc models.TestRunCase
                var testCase models.TestCase
                err = rows.Scan(
                        &trc.ID, &trc.TestRunID, &trc.TestCaseID, &trc.Status, &trc.ResultNotes,
                        &trc.ExecutedBy, &trc.StartedAt, &trc.CompletedAt, &trc.CreatedAt, &trc.UpdatedAt,
                        &testCase.ID, &testCase.Title, &testCase.Description, &testCase.Priority, &testCase.Status, &testCase.TestSuiteID, &testCase.CreatedAt, &testCase.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan test run case: %w", err)
                }
                trc.TestCase = &testCase
                testRunCases = append(testRunCases, trc)
        }

        return testRunCases, nil
}

// Create creates a new test run
func (r *TestRunRepository) Create(req models.CreateTestRunRequest) (*models.TestRun, error) {
        tx, err := r.db.Begin()
        if err != nil {
                return nil, fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()

        // Create the test run
        var testRun models.TestRun
        err = tx.QueryRow(`
                INSERT INTO test_runs (name, description, project_id, repository_id, branch_name, tag_name, created_by)
                VALUES ($1, $2, $3, $4, $5, $6, $7)
                RETURNING id, name, description, project_id, repository_id, branch_name, tag_name, status, 
                          created_by, started_at, completed_at, created_at, updated_at
        `, req.Name, req.Description, req.ProjectID, req.RepositoryID, req.BranchName, req.TagName, req.CreatedBy).Scan(
                &testRun.ID, &testRun.Name, &testRun.Description, &testRun.ProjectID, &testRun.RepositoryID,
                &testRun.BranchName, &testRun.TagName, &testRun.Status, &testRun.CreatedBy,
                &testRun.StartedAt, &testRun.CompletedAt, &testRun.CreatedAt, &testRun.UpdatedAt,
        )
        if err != nil {
                return nil, fmt.Errorf("failed to create test run: %w", err)
        }

        // Add test cases to the run
        for _, testCaseID := range req.TestCaseIDs {
                _, err = tx.Exec(`
                        INSERT INTO test_run_cases (test_run_id, test_case_id)
                        VALUES ($1, $2)
                `, testRun.ID, testCaseID)
                if err != nil {
                        return nil, fmt.Errorf("failed to add test case %d to run: %w", testCaseID, err)
                }
        }

        err = tx.Commit()
        if err != nil {
                return nil, fmt.Errorf("failed to commit transaction: %w", err)
        }

        // Return the created test run with full details
        return r.GetByID(testRun.ID)
}

// Update updates a test run
func (r *TestRunRepository) Update(id int, req models.UpdateTestRunRequest) (*models.TestRun, error) {
        tx, err := r.db.Begin()
        if err != nil {
                return nil, fmt.Errorf("failed to begin transaction: %w", err)
        }
        defer tx.Rollback()

        setParts := []string{}
        args := []interface{}{}
        argIndex := 1

        if req.Name != nil {
                setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
                args = append(args, *req.Name)
                argIndex++
        }
        if req.Description != nil {
                setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
                args = append(args, *req.Description)
                argIndex++
        }
        if req.ProjectID != nil {
                setParts = append(setParts, fmt.Sprintf("project_id = $%d", argIndex))
                args = append(args, *req.ProjectID)
                argIndex++
        }
        if req.RepositoryID != nil {
                setParts = append(setParts, fmt.Sprintf("repository_id = $%d", argIndex))
                args = append(args, *req.RepositoryID)
                argIndex++
        }
        if req.BranchName != nil {
                setParts = append(setParts, fmt.Sprintf("branch_name = $%d", argIndex))
                args = append(args, *req.BranchName)
                argIndex++
        }
        if req.TagName != nil {
                setParts = append(setParts, fmt.Sprintf("tag_name = $%d", argIndex))
                args = append(args, *req.TagName)
                argIndex++
        }
        if req.CreatedBy != nil {
                setParts = append(setParts, fmt.Sprintf("created_by = $%d", argIndex))
                args = append(args, *req.CreatedBy)
                argIndex++
        }
        if req.Status != nil {
                setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
                args = append(args, *req.Status)
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

        // Always update the updated_at timestamp
        setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
        args = append(args, time.Now())
        argIndex++

        args = append(args, id)

        // Update test run if there are changes
        if len(setParts) > 1 { // More than just updated_at
                query := fmt.Sprintf(`
                        UPDATE test_runs 
                        SET %s
                        WHERE id = $%d
                `, strings.Join(setParts, ", "), argIndex)

                _, err = tx.Exec(query, args...)
                if err != nil {
                        return nil, fmt.Errorf("failed to update test run: %w", err)
                }
        }

        // Update test cases if provided
        if len(req.TestCaseIDs) > 0 {
                // First, remove all existing test cases for this run
                _, err = tx.Exec("DELETE FROM test_run_cases WHERE test_run_id = $1", id)
                if err != nil {
                        return nil, fmt.Errorf("failed to remove existing test cases: %w", err)
                }

                // Then add the new test cases
                for _, testCaseID := range req.TestCaseIDs {
                        _, err = tx.Exec(`
                                INSERT INTO test_run_cases (test_run_id, test_case_id)
                                VALUES ($1, $2)
                        `, id, testCaseID)
                        if err != nil {
                                return nil, fmt.Errorf("failed to add test case %d to run: %w", testCaseID, err)
                        }
                }
        }

        err = tx.Commit()
        if err != nil {
                return nil, fmt.Errorf("failed to commit transaction: %w", err)
        }

        return r.GetByID(id)
}

// Delete deletes a test run
func (r *TestRunRepository) Delete(id int) error {
        _, err := r.db.Exec("DELETE FROM test_runs WHERE id = $1", id)
        if err != nil {
                return fmt.Errorf("failed to delete test run: %w", err)
        }
        return nil
}

// UpdateTestRunCase updates a test case within a test run
func (r *TestRunRepository) UpdateTestRunCase(testRunID, testCaseID int, req models.UpdateTestRunCaseRequest) (*models.TestRunCase, error) {
        setParts := []string{}
        args := []interface{}{}
        argIndex := 1

        if req.Status != nil {
                setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
                args = append(args, *req.Status)
                argIndex++
        }
        if req.ResultNotes != nil {
                setParts = append(setParts, fmt.Sprintf("result_notes = $%d", argIndex))
                args = append(args, *req.ResultNotes)
                argIndex++
        }
        if req.ExecutedBy != nil {
                setParts = append(setParts, fmt.Sprintf("executed_by = $%d", argIndex))
                args = append(args, *req.ExecutedBy)
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

        if len(setParts) == 0 {
                return nil, fmt.Errorf("no fields to update")
        }

        setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
        args = append(args, time.Now())
        argIndex++

        args = append(args, testRunID, testCaseID)

        query := fmt.Sprintf(`
                UPDATE test_run_cases 
                SET %s
                WHERE test_run_id = $%d AND test_case_id = $%d
        `, fmt.Sprintf("%s", setParts[0]), argIndex-1, argIndex)

        for i := 1; i < len(setParts); i++ {
                query = fmt.Sprintf(`
                        UPDATE test_run_cases 
                        SET %s
                        WHERE test_run_id = $%d AND test_case_id = $%d
                `, fmt.Sprintf("%s, %s", setParts[0], setParts[i]), argIndex-1, argIndex)
        }

        _, err := r.db.Exec(query, args...)
        if err != nil {
                return nil, fmt.Errorf("failed to update test run case: %w", err)
        }

        // Return the updated test run case
        var trc models.TestRunCase
        var testCase models.TestCase
        err = r.db.QueryRow(`
                SELECT trc.id, trc.test_run_id, trc.test_case_id, trc.status, trc.result_notes, 
                       trc.executed_by, trc.started_at, trc.completed_at, trc.created_at, trc.updated_at,
                       tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at
                FROM test_run_cases trc
                JOIN test_cases tc ON trc.test_case_id = tc.id
                WHERE trc.test_run_id = $1 AND trc.test_case_id = $2
        `, testRunID, testCaseID).Scan(
                &trc.ID, &trc.TestRunID, &trc.TestCaseID, &trc.Status, &trc.ResultNotes,
                &trc.ExecutedBy, &trc.StartedAt, &trc.CompletedAt, &trc.CreatedAt, &trc.UpdatedAt,
                &testCase.ID, &testCase.Title, &testCase.Description, &testCase.Priority, &testCase.Status, &testCase.TestSuiteID, &testCase.CreatedAt, &testCase.UpdatedAt,
        )
        if err != nil {
                return nil, fmt.Errorf("failed to get updated test run case: %w", err)
        }

        trc.TestCase = &testCase
        return &trc, nil
}