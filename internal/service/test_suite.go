package service

import (
	"errors"

	"github.com/galex-do/test-machine/internal/models"
	"github.com/galex-do/test-machine/internal/repository"
)

// TestSuiteService handles business logic for test suites
type TestSuiteService struct {
	repo *repository.TestSuiteRepository
}

// NewTestSuiteService creates a new test suite service
func NewTestSuiteService(repo *repository.TestSuiteRepository) *TestSuiteService {
	return &TestSuiteService{repo: repo}
}

// GetAll returns all test suites, optionally filtered by project ID
func (s *TestSuiteService) GetAll(projectID *int) ([]models.TestSuite, error) {
	return s.repo.GetAll(projectID)
}

// GetByID returns a test suite by ID
func (s *TestSuiteService) GetByID(id int) (*models.TestSuite, error) {
	return s.repo.GetByID(id)
}

// Create creates a new test suite
func (s *TestSuiteService) Create(req *models.CreateTestSuiteRequest) (*models.TestSuite, error) {
	if req.Name == "" || req.ProjectID == 0 {
		return nil, errors.New("name and project_id are required")
	}
	return s.repo.Create(req)
}