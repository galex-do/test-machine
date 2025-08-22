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

// Update updates an existing test suite
func (s *TestSuiteService) Update(id int, req *models.UpdateTestSuiteRequest) (*models.TestSuite, error) {
        if req.Name == "" {
                return nil, errors.New("name is required")
        }
        return s.repo.Update(id, req)
}

// Delete deletes a test suite
func (s *TestSuiteService) Delete(id int) error {
        return s.repo.Delete(id)
}