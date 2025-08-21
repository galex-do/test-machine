package service

import (
	"errors"

	"github.com/galex-do/test-machine/internal/models"
	"github.com/galex-do/test-machine/internal/repository"
)

// TestRunService handles business logic for test runs
type TestRunService struct {
	repo *repository.TestRunRepository
}

// NewTestRunService creates a new test run service
func NewTestRunService(repo *repository.TestRunRepository) *TestRunService {
	return &TestRunService{repo: repo}
}

// GetAll returns all test runs, optionally filtered by test case ID
func (s *TestRunService) GetAll(testCaseID *int) ([]models.TestRun, error) {
	return s.repo.GetAll(testCaseID)
}

// GetByID returns a test run by ID
func (s *TestRunService) GetByID(id int) (*models.TestRun, error) {
	return s.repo.GetByID(id)
}

// Create creates a new test run
func (s *TestRunService) Create(req *models.CreateTestRunRequest) (*models.TestRun, error) {
	if req.Name == "" || req.TestCaseID == 0 {
		return nil, errors.New("name and test_case_id are required")
	}
	return s.repo.Create(req)
}

// Update updates an existing test run
func (s *TestRunService) Update(id int, req *models.UpdateTestRunRequest) (*models.TestRun, error) {
	return s.repo.Update(id, req)
}