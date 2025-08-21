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