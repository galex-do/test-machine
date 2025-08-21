package models

import (
	"time"
	"gorm.io/gorm"
)

// Project represents a test project
type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	TestSuites  []TestSuite    `json:"test_suites,omitempty" gorm:"foreignKey:ProjectID"`
}

// TestSuite represents a collection of test cases
type TestSuite struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	ProjectID   uint           `json:"project_id" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Project     Project        `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	TestCases   []TestCase     `json:"test_cases,omitempty" gorm:"foreignKey:TestSuiteID"`
}

// TestCase represents a test case
type TestCase struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Title          string         `json:"title" gorm:"not null"`
	Description    string         `json:"description"`
	Priority       string         `json:"priority" gorm:"default:'Medium'"` // High, Medium, Low
	Status         string         `json:"status" gorm:"default:'Active'"`   // Active, Inactive, Deprecated
	TestSuiteID    uint           `json:"test_suite_id" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	TestSuite      TestSuite      `json:"test_suite,omitempty" gorm:"foreignKey:TestSuiteID"`
	TestSteps      []TestStep     `json:"test_steps,omitempty" gorm:"foreignKey:TestCaseID"`
	TestExecutions []TestExecution `json:"test_executions,omitempty" gorm:"foreignKey:TestCaseID"`
}

// TestStep represents a step in a test case
type TestStep struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	StepNumber     int            `json:"step_number" gorm:"not null"`
	Description    string         `json:"description" gorm:"not null"`
	ExpectedResult string         `json:"expected_result" gorm:"not null"`
	TestCaseID     uint           `json:"test_case_id" gorm:"not null"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	TestCase       TestCase       `json:"test_case,omitempty" gorm:"foreignKey:TestCaseID"`
}

// TestRun represents a test execution session
type TestRun struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	Name           string          `json:"name" gorm:"not null"`
	Description    string          `json:"description"`
	TestSuiteID    uint            `json:"test_suite_id" gorm:"not null"`
	Status         string          `json:"status" gorm:"default:'In Progress'"` // In Progress, Completed, Cancelled
	StartedAt      *time.Time      `json:"started_at"`
	CompletedAt    *time.Time      `json:"completed_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `json:"-" gorm:"index"`
	TestSuite      TestSuite       `json:"test_suite,omitempty" gorm:"foreignKey:TestSuiteID"`
	TestExecutions []TestExecution `json:"test_executions,omitempty" gorm:"foreignKey:TestRunID"`
}

// TestExecution represents the execution of a test case within a test run
type TestExecution struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	TestRunID    uint           `json:"test_run_id" gorm:"not null"`
	TestCaseID   uint           `json:"test_case_id" gorm:"not null"`
	Status       string         `json:"status" gorm:"default:'Not Executed'"` // Not Executed, Pass, Fail, Blocked, Skip
	Notes        string         `json:"notes"`
	ExecutedAt   *time.Time     `json:"executed_at"`
	ExecutedBy   string         `json:"executed_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	TestRun      TestRun        `json:"test_run,omitempty" gorm:"foreignKey:TestRunID"`
	TestCase     TestCase       `json:"test_case,omitempty" gorm:"foreignKey:TestCaseID"`
}

// CreateProjectRequest represents request data for creating a project
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateTestSuiteRequest represents request data for creating a test suite
type CreateTestSuiteRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ProjectID   uint   `json:"project_id" binding:"required"`
}

// CreateTestCaseRequest represents request data for creating a test case
type CreateTestCaseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	TestSuiteID uint   `json:"test_suite_id" binding:"required"`
}

// CreateTestStepRequest represents request data for creating a test step
type CreateTestStepRequest struct {
	StepNumber     int    `json:"step_number" binding:"required"`
	Description    string `json:"description" binding:"required"`
	ExpectedResult string `json:"expected_result" binding:"required"`
}

// CreateTestRunRequest represents request data for creating a test run
type CreateTestRunRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	TestSuiteID uint   `json:"test_suite_id" binding:"required"`
}

// UpdateTestExecutionRequest represents request data for updating test execution
type UpdateTestExecutionRequest struct {
	Status     string `json:"status" binding:"required"`
	Notes      string `json:"notes"`
	ExecutedBy string `json:"executed_by"`
}

// ReportSummary represents test execution summary for reports
type ReportSummary struct {
	TotalProjects    int64                    `json:"total_projects"`
	TotalTestSuites  int64                    `json:"total_test_suites"`
	TotalTestCases   int64                    `json:"total_test_cases"`
	TotalTestRuns    int64                    `json:"total_test_runs"`
	ExecutionSummary TestExecutionSummary     `json:"execution_summary"`
	RecentRuns       []TestRun                `json:"recent_runs"`
	ProjectSummaries []ProjectExecutionSummary `json:"project_summaries"`
}

// TestExecutionSummary represents overall test execution statistics
type TestExecutionSummary struct {
	TotalExecutions int64 `json:"total_executions"`
	Passed          int64 `json:"passed"`
	Failed          int64 `json:"failed"`
	Blocked         int64 `json:"blocked"`
	Skipped         int64 `json:"skipped"`
	NotExecuted     int64 `json:"not_executed"`
	PassRate        float64 `json:"pass_rate"`
}

// ProjectExecutionSummary represents test execution summary per project
type ProjectExecutionSummary struct {
	ProjectID   uint                 `json:"project_id"`
	ProjectName string               `json:"project_name"`
	Summary     TestExecutionSummary `json:"summary"`
}
