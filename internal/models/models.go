package models

import "time"

// Project represents a test project
type Project struct {
        ID              int       `json:"id"`
        Name            string    `json:"name"`
        Description     string    `json:"description"`
        CreatedAt       time.Time `json:"created_at"`
        UpdatedAt       time.Time `json:"updated_at"`
        TestSuitesCount int       `json:"test_suites_count,omitempty"`
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
        ID          int        `json:"id"`
        Title       string     `json:"title"`
        Description string     `json:"description"`
        Priority    string     `json:"priority"`
        Status      string     `json:"status"`
        TestSuiteID int        `json:"test_suite_id"`
        CreatedAt   time.Time  `json:"created_at"`
        UpdatedAt   time.Time  `json:"updated_at"`
        TestSuite   *TestSuite `json:"test_suite,omitempty"`
        TestSteps   []TestStep `json:"test_steps,omitempty"`
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
        Name        string `json:"name"`
        Description string `json:"description"`
}

// UpdateProjectRequest represents the request to update a project
type UpdateProjectRequest struct {
        Name        string `json:"name"`
        Description string `json:"description"`
}

// CreateTestSuiteRequest represents the request to create a new test suite
type CreateTestSuiteRequest struct {
        Name        string `json:"name"`
        Description string `json:"description"`
        ProjectID   int    `json:"project_id"`
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