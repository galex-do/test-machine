package service

import (
        "fmt"
        "time"

        "github.com/galex-do/test-machine/internal/models"
        "github.com/galex-do/test-machine/internal/repository"
)

// TestRunService handles business logic for test runs
type TestRunService struct {
        repo        *repository.TestRunRepository
        projectRepo *repository.ProjectRepository
}

// NewTestRunService creates a new test run service
func NewTestRunService(repo *repository.TestRunRepository, projectRepo *repository.ProjectRepository) *TestRunService {
        return &TestRunService{
                repo:        repo,
                projectRepo: projectRepo,
        }
}

// GetAllTestRuns returns all test runs
func (s *TestRunService) GetAllTestRuns() ([]models.TestRun, error) {
        return s.repo.GetAll()
}

// GetTestRunByID returns a test run by ID
func (s *TestRunService) GetTestRunByID(id int) (*models.TestRun, error) {
        return s.repo.GetByID(id)
}

// CreateTestRun creates a new test run
func (s *TestRunService) CreateTestRun(req models.CreateTestRunRequest) (*models.TestRun, error) {
        // Validate that the project exists
        project, err := s.projectRepo.GetByID(req.ProjectID)
        if err != nil {
                return nil, err
        }
        if project == nil {
                return nil, fmt.Errorf("project not found")
        }

        // Validate test case IDs if needed
        if len(req.TestCaseIDs) == 0 {
                return nil, fmt.Errorf("at least one test case must be selected")
        }

        // Auto-generate name if empty
        if req.Name == "" {
                req.Name = s.generateTestRunName(project, req.BranchName, req.TagName)
        }

        return s.repo.Create(req)
}

// UpdateTestRun updates a test run
func (s *TestRunService) UpdateTestRun(id int, req models.UpdateTestRunRequest) (*models.TestRun, error) {
        return s.repo.Update(id, req)
}

// DeleteTestRun deletes a test run
func (s *TestRunService) DeleteTestRun(id int) error {
        return s.repo.Delete(id)
}

// UpdateTestRunCase updates a test case within a test run
func (s *TestRunService) UpdateTestRunCase(testRunID, testCaseID int, req models.UpdateTestRunCaseRequest) (*models.TestRunCase, error) {
        return s.repo.UpdateTestRunCase(testRunID, testCaseID, req)
}

// generateTestRunName creates an auto-generated name for test runs
func (s *TestRunService) generateTestRunName(project *models.Project, branchName, tagName *string) string {
        name := project.Name
        
        // Add branch or tag if specified
        if branchName != nil && *branchName != "" {
                name += "-" + *branchName
        } else if tagName != nil && *tagName != "" {
                name += "-" + *tagName
        }
        
        // Add datetime (format: YYYY-MM-DD-HHMMSS)
        now := time.Now()
        datetime := now.Format("2006-01-02-150405")
        name += "-" + datetime
        
        return name
}