package models

import "time"

// Project represents a test project
type Project struct {
        ID              int       `json:"id"`
        Name            string    `json:"name"`
        Description     string    `json:"description"`
        RepositoryID    *int      `json:"repository_id,omitempty"`
        CreatedAt       time.Time `json:"created_at"`
        UpdatedAt       time.Time `json:"updated_at"`
        TestSuitesCount int       `json:"test_suites_count,omitempty"`
        Repository      *Repository `json:"repository,omitempty"`
}

// TestSuite represents a collection of test cases
type TestSuite struct {
        ID            int       `json:"id"`
        Name          string    `json:"name"`
        Description   string    `json:"description"`
        ProjectID     int       `json:"project_id"`
        CreatedAt     time.Time `json:"created_at"`
        UpdatedAt     time.Time `json:"updated_at"`
        Project       *Project  `json:"project,omitempty"`
        TestCasesCount int      `json:"test_cases_count,omitempty"`
        TestCases     []TestCase `json:"test_cases,omitempty"`
}

// TestCase represents an individual test case
type TestCase struct {
        ID            int        `json:"id"`
        Title         string     `json:"title"`
        Description   string     `json:"description"`
        Priority      string     `json:"priority"`
        Status        string     `json:"status"`
        TestSuiteID   int        `json:"test_suite_id"`
        CreatedAt     time.Time  `json:"created_at"`
        UpdatedAt     time.Time  `json:"updated_at"`
        TestSuite     *TestSuite `json:"test_suite,omitempty"`
        TestSteps     []TestStep `json:"test_steps,omitempty"`
        TestStepsCount int       `json:"test_steps_count,omitempty"`
}

// TestStep represents an individual step in a test case
type TestStep struct {
        ID             int       `json:"id"`
        TestCaseID     int       `json:"test_case_id"`
        StepNumber     int       `json:"step_number"`
        Description    string    `json:"description"`
        ExpectedResult string    `json:"expected_result"`
        CreatedAt      time.Time `json:"created_at"`
        UpdatedAt      time.Time `json:"updated_at"`
}


// CreateProjectRequest represents the request to create a new project
type CreateProjectRequest struct {
        Name         string `json:"name"`
        Description  string `json:"description"`
        RepositoryID *int   `json:"repository_id"`
}

// UpdateProjectRequest represents the request to update a project
type UpdateProjectRequest struct {
        Name         string `json:"name"`
        Description  string `json:"description"`
        RepositoryID *int   `json:"repository_id"`
}

// CreateTestSuiteRequest represents the request to create a new test suite
type CreateTestSuiteRequest struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        ProjectID   int    `json:"project_id"`
}

// UpdateTestSuiteRequest represents the request to update a test suite
type UpdateTestSuiteRequest struct {
        Name        string `json:"name"`
        Description string `json:"description"`
}

// CreateTestCaseRequest represents the request to create a new test case
type CreateTestCaseRequest struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        Priority    string `json:"priority"`
        Status      string `json:"status"`
        TestSuiteID int    `json:"test_suite_id"`
}

// UpdateTestCaseRequest represents the request to update a test case
type UpdateTestCaseRequest struct {
        Title       string `json:"title"`
        Description string `json:"description"`
        Priority    string `json:"priority"`
        Status      string `json:"status"`
}

// CreateTestStepRequest represents the request to create a new test step
type CreateTestStepRequest struct {
        TestCaseID     int    `json:"test_case_id"`
        StepNumber     int    `json:"step_number"`
        Description    string `json:"description"`
        ExpectedResult string `json:"expected_result"`
}

// UpdateTestStepRequest represents the request to update a test step
type UpdateTestStepRequest struct {
        StepNumber     int    `json:"step_number"`
        Description    string `json:"description"`
        ExpectedResult string `json:"expected_result"`
}

// TestRun represents a collection of test cases to be executed
type TestRun struct {
        ID           int                `json:"id"`
        Name         string             `json:"name"`
        Description  string             `json:"description"`
        ProjectID    int                `json:"project_id"`
        RepositoryID *int               `json:"repository_id,omitempty"`
        BranchName   *string            `json:"branch_name,omitempty"`
        TagName      *string            `json:"tag_name,omitempty"`
        Status       string             `json:"status"`
        CreatedBy    *string            `json:"created_by,omitempty"`
        StartedAt    *time.Time         `json:"started_at,omitempty"`
        CompletedAt  *time.Time         `json:"completed_at,omitempty"`
        CreatedAt    time.Time          `json:"created_at"`
        UpdatedAt    time.Time          `json:"updated_at"`
        Project      *Project           `json:"project,omitempty"`
        Repository   *Repository        `json:"repository,omitempty"`
        TestCases    []TestRunCase      `json:"test_cases,omitempty"`
        Intervals    []TestRunInterval  `json:"intervals,omitempty"`
        TotalExecutionTime *int         `json:"total_execution_time,omitempty"` // in seconds
}

// TestRunInterval represents a time interval during test run execution
type TestRunInterval struct {
        ID        int        `json:"id"`
        TestRunID int        `json:"test_run_id"`
        StartTime time.Time  `json:"start_time"`
        EndTime   *time.Time `json:"end_time,omitempty"`
        CreatedAt time.Time  `json:"created_at"`
        UpdatedAt time.Time  `json:"updated_at"`
}

// TestRunCase represents a test case within a test run
type TestRunCase struct {
        ID           int       `json:"id"`
        TestRunID    int       `json:"test_run_id"`
        TestCaseID   int       `json:"test_case_id"`
        Status       string    `json:"status"`
        ResultNotes  *string   `json:"result_notes,omitempty"`
        ExecutedBy   *string   `json:"executed_by,omitempty"`
        StartedAt    *time.Time `json:"started_at,omitempty"`
        CompletedAt  *time.Time `json:"completed_at,omitempty"`
        CreatedAt    time.Time `json:"created_at"`
        UpdatedAt    time.Time `json:"updated_at"`
        TestCase     *TestCase `json:"test_case,omitempty"`
}

// TestExecution represents an individual test execution (renamed from TestRun)
type TestExecution struct {
        ID             int        `json:"id"`
        Name           string     `json:"name"`
        Description    string     `json:"description"`
        TestCaseID     int        `json:"test_case_id"`
        TestRunCaseID  *int       `json:"test_run_case_id,omitempty"`
        Status         string     `json:"status"`
        Result         *string    `json:"result,omitempty"`
        Notes          *string    `json:"notes,omitempty"`
        ExecutedBy     *string    `json:"executed_by,omitempty"`
        StartedAt      *time.Time `json:"started_at,omitempty"`
        CompletedAt    *time.Time `json:"completed_at,omitempty"`
        CreatedAt      time.Time  `json:"created_at"`
        UpdatedAt      time.Time  `json:"updated_at"`
        TestCase       *TestCase  `json:"test_case,omitempty"`
}

// CreateTestRunRequest represents the request to create a new test run
type CreateTestRunRequest struct {
        Name         string   `json:"name"`
        Description  string   `json:"description"`
        ProjectID    int      `json:"project_id"`
        RepositoryID *int     `json:"repository_id"`
        BranchName   *string  `json:"branch_name"`
        TagName      *string  `json:"tag_name"`
        TestCaseIDs  []int    `json:"test_case_ids"`
        CreatedBy    *string  `json:"created_by"`
}

// UpdateTestRunRequest represents the request to update a test run
type UpdateTestRunRequest struct {
        Name         *string    `json:"name,omitempty"`
        Description  *string    `json:"description,omitempty"`
        ProjectID    *int       `json:"project_id,omitempty"`
        RepositoryID *int       `json:"repository_id,omitempty"`
        BranchName   *string    `json:"branch_name,omitempty"`
        TagName      *string    `json:"tag_name,omitempty"`
        TestCaseIDs  []int      `json:"test_case_ids,omitempty"`
        CreatedBy    *string    `json:"created_by,omitempty"`
        Status       *string    `json:"status,omitempty"`
        StartedAt    *time.Time `json:"started_at,omitempty"`
        CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

// UpdateTestRunCaseRequest represents the request to update a test case within a run
type UpdateTestRunCaseRequest struct {
        Status      *string    `json:"status,omitempty"`
        ResultNotes *string    `json:"result_notes,omitempty"`
        ExecutedBy  *string    `json:"executed_by,omitempty"`
        StartedAt   *time.Time `json:"started_at,omitempty"`
        CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Key represents an authentication key for Git repositories
type Key struct {
        ID            int       `json:"id"`
        Name          string    `json:"name"`
        Description   string    `json:"description"`
        KeyType       string    `json:"key_type"`
        Username      *string   `json:"username,omitempty"`
        EncryptedData string    `json:"-"` // Never expose encrypted data in JSON
        CreatedAt     time.Time `json:"created_at"`
        UpdatedAt     time.Time `json:"updated_at"`
}

// CreateKeyRequest represents the request to create a new key
type CreateKeyRequest struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        KeyType     string  `json:"key_type"`
        Username    *string `json:"username"`
        SecretData  string  `json:"secret_data"` // This will be encrypted before storage
}

// UpdateKeyRequest represents the request to update a key
type UpdateKeyRequest struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        Username    *string `json:"username"`
        SecretData  *string `json:"secret_data,omitempty"` // Optional - only if changing
}

// KeyDecryptResponse represents the decrypted key data
type KeyDecryptResponse struct {
        Data string `json:"data"`
}

// Repository represents a Git repository that can be used by projects
type Repository struct {
        ID            int       `json:"id"`
        Name          string    `json:"name"`
        Description   string    `json:"description"`
        RemoteURL     string    `json:"remote_url"`
        KeyID         *int      `json:"key_id,omitempty"`
        DefaultBranch *string   `json:"default_branch,omitempty"`
        SyncedAt      *time.Time `json:"synced_at,omitempty"`
        CreatedAt     time.Time `json:"created_at"`
        UpdatedAt     time.Time `json:"updated_at"`
        Key           *Key      `json:"key,omitempty"`
        Branches      []Branch  `json:"branches,omitempty"`
        Tags          []Tag     `json:"tags,omitempty"`
}

// Branch represents a Git branch in a repository
type Branch struct {
        ID            int        `json:"id"`
        RepositoryID  int        `json:"repository_id"`
        Name          string     `json:"name"`
        CommitHash    *string    `json:"commit_hash,omitempty"`
        CommitDate    *time.Time `json:"commit_date,omitempty"`
        CommitMessage *string    `json:"commit_message,omitempty"`
        IsDefault     bool       `json:"is_default"`
        CreatedAt     time.Time  `json:"created_at"`
        UpdatedAt     time.Time  `json:"updated_at"`
}

// Tag represents a Git tag in a repository
type Tag struct {
        ID            int        `json:"id"`
        RepositoryID  int        `json:"repository_id"`
        Name          string     `json:"name"`
        CommitHash    *string    `json:"commit_hash,omitempty"`
        CommitDate    *time.Time `json:"commit_date,omitempty"`
        CommitMessage *string    `json:"commit_message,omitempty"`
        CreatedAt     time.Time  `json:"created_at"`
        UpdatedAt     time.Time  `json:"updated_at"`
}

// CreateRepositoryRequest represents the request to create a new repository
type CreateRepositoryRequest struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        RemoteURL   string `json:"remote_url"`
        KeyID       *int   `json:"key_id"`
}

// UpdateRepositoryRequest represents the request to update a repository
type UpdateRepositoryRequest struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        KeyID       *int   `json:"key_id"`
        // Note: RemoteURL is intentionally omitted - it's immutable after creation
}

// SyncRequest represents a request to sync a repository
type SyncRequest struct {
        RepositoryID int `json:"repository_id"`
}

// SyncResponse represents the response from a sync operation
type SyncResponse struct {
        Success      bool       `json:"success"`
        Message      string     `json:"message"`
        Repository   *Repository `json:"repository,omitempty"`
        BranchCount  int        `json:"branch_count"`
        TagCount     int        `json:"tag_count"`
}