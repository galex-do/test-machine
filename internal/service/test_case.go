package service

import (
        "errors"

        "github.com/galex-do/test-machine/internal/models"
        "github.com/galex-do/test-machine/internal/repository"
)

// TestCaseService handles business logic for test cases
type TestCaseService struct {
        repo *repository.TestCaseRepository
}

// NewTestCaseService creates a new test case service
func NewTestCaseService(repo *repository.TestCaseRepository) *TestCaseService {
        return &TestCaseService{repo: repo}
}

// GetAll returns all test cases, optionally filtered by test suite ID
func (s *TestCaseService) GetAll(testSuiteID *int) ([]models.TestCase, error) {
        return s.repo.GetAll(testSuiteID)
}

// GetByID returns a test case by ID
func (s *TestCaseService) GetByID(id int) (*models.TestCase, error) {
        return s.repo.GetByID(id)
}

// Create creates a new test case
func (s *TestCaseService) Create(req *models.CreateTestCaseRequest) (*models.TestCase, error) {
        if req.Title == "" || req.TestSuiteID == 0 {
                return nil, errors.New("title and test_suite_id are required")
        }
        return s.repo.Create(req)
}

// GetTestSteps returns all test steps for a test case
func (s *TestCaseService) GetTestSteps(testCaseID int) ([]models.TestStep, error) {
        return s.repo.GetTestSteps(testCaseID)
}

// CreateTestStep creates a new test step
func (s *TestCaseService) CreateTestStep(req *models.CreateTestStepRequest) (*models.TestStep, error) {
        if req.TestCaseID == 0 || req.StepNumber <= 0 || req.Description == "" || req.ExpectedResult == "" {
                return nil, errors.New("test_case_id, step_number, description, and expected_result are required")
        }
        return s.repo.CreateTestStep(req)
}

// UpdateTestStep updates an existing test step
func (s *TestCaseService) UpdateTestStep(id int, req *models.UpdateTestStepRequest) (*models.TestStep, error) {
        if req.StepNumber <= 0 || req.Description == "" || req.ExpectedResult == "" {
                return nil, errors.New("step_number, description, and expected_result are required")
        }
        return s.repo.UpdateTestStep(id, req)
}

// DeleteTestStep deletes a test step
func (s *TestCaseService) DeleteTestStep(id int) error {
        return s.repo.DeleteTestStep(id)
}