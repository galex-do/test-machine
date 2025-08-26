package models

import "time"

// Project represents a test project
type Project struct {
        ID              int       `json:"id"`
        Name            string    `json:"name"`
        Description     string    `json:"description"`
        GitProject      *string   `json:"git_project,omitempty"`
        KeyID           *int      `json:"key_id,omitempty"`
        CreatedAt       time.Time `json:"created_at"`
        UpdatedAt       time.Time `json:"updated_at"`
        TestSuitesCount int       `json:"test_suites_count,omitempty"`
        Key             *Key      `json:"key,omitempty"`
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

// TestRun represents an execution of a test case
type TestRun struct {
        ID          int        `json:"id"`
        Name        string     `json:"name"`
        Description string     `json:"description"`
        TestCaseID  int        `json:"test_case_id"`
        Status      string     `json:"status"`
        Result      *string    `json:"result"`
        Notes       *string    `json:"notes"`
        ExecutedBy  *string    `json:"executed_by"`
        StartedAt   *time.Time `json:"started_at"`
        CompletedAt *time.Time `json:"completed_at"`
        CreatedAt   time.Time  `json:"created_at"`
        UpdatedAt   time.Time  `json:"updated_at"`
        TestCase    *TestCase  `json:"test_case,omitempty"`
}

// CreateProjectRequest represents the request to create a new project
type CreateProjectRequest struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        GitProject  *string `json:"git_project"`
        KeyID       *int    `json:"key_id"`
}

// UpdateProjectRequest represents the request to update a project
type UpdateProjectRequest struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        GitProject  *string `json:"git_project"`
        KeyID       *int    `json:"key_id"`
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

// CreateTestRunRequest represents the request to create a new test run
type CreateTestRunRequest struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        TestCaseID  int     `json:"test_case_id"`
        ExecutedBy  *string `json:"executed_by"`
}

// UpdateTestRunRequest represents the request to update a test run
type UpdateTestRunRequest struct {
        Status      *string    `json:"status"`
        Result      *string    `json:"result"`
        Notes       *string    `json:"notes"`
        StartedAt   *time.Time `json:"started_at"`
        CompletedAt *time.Time `json:"completed_at"`
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

// Repository represents a Git repository synced with a project
type Repository struct {
        ID            int       `json:"id"`
        ProjectID     int       `json:"project_id"`
        RemoteURL     string    `json:"remote_url"`
        DefaultBranch *string   `json:"default_branch,omitempty"`
        SyncedAt      *time.Time `json:"synced_at,omitempty"`
        CreatedAt     time.Time `json:"created_at"`
        UpdatedAt     time.Time `json:"updated_at"`
        Branches      []Branch  `json:"branches,omitempty"`
        Tags          []Tag     `json:"tags,omitempty"`
}

// Branch represents a Git branch in a repository
type Branch struct {
        ID           int       `json:"id"`
        RepositoryID int       `json:"repository_id"`
        Name         string    `json:"name"`
        CommitHash   *string   `json:"commit_hash,omitempty"`
        IsDefault    bool      `json:"is_default"`
        CreatedAt    time.Time `json:"created_at"`
        UpdatedAt    time.Time `json:"updated_at"`
}

// Tag represents a Git tag in a repository
type Tag struct {
        ID           int       `json:"id"`
        RepositoryID int       `json:"repository_id"`
        Name         string    `json:"name"`
        CommitHash   *string   `json:"commit_hash,omitempty"`
        CreatedAt    time.Time `json:"created_at"`
        UpdatedAt    time.Time `json:"updated_at"`
}

// SyncRequest represents a request to sync a project with its Git repository
type SyncRequest struct {
        ProjectID int `json:"project_id"`
}

// SyncResponse represents the response from a sync operation
type SyncResponse struct {
        Success      bool       `json:"success"`
        Message      string     `json:"message"`
        Repository   *Repository `json:"repository,omitempty"`
        BranchCount  int        `json:"branch_count"`
        TagCount     int        `json:"tag_count"`
}