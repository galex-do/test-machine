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
        intervalRepo *repository.TestRunIntervalRepository
}

// NewTestRunService creates a new test run service
func NewTestRunService(repo *repository.TestRunRepository, projectRepo *repository.ProjectRepository, intervalRepo *repository.TestRunIntervalRepository) *TestRunService {
        return &TestRunService{
                repo:        repo,
                projectRepo: projectRepo,
                intervalRepo: intervalRepo,
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
        // Check if test run exists and validate status
        testRun, err := s.repo.GetByID(id)
        if err != nil {
                return nil, err
        }
        if testRun == nil {
                return nil, fmt.Errorf("test run not found")
        }

        // Prevent editing completed test runs
        if testRun.Status == "Completed" {
                return nil, fmt.Errorf("cannot edit completed test runs")
        }

        return s.repo.Update(id, req)
}

// DeleteTestRun deletes a test run
func (s *TestRunService) DeleteTestRun(id int) error {
        // Check if test run exists and validate status
        testRun, err := s.repo.GetByID(id)
        if err != nil {
                return err
        }
        if testRun == nil {
                return fmt.Errorf("test run not found")
        }

        // Only allow deleting test runs that haven't started
        if testRun.Status != "Not Started" {
                return fmt.Errorf("can only delete test runs that haven't started")
        }

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

// StartTestRun starts a test run execution and creates a new time interval
func (s *TestRunService) StartTestRun(id int) (*models.TestRun, error) {
        // Check if test run exists and is in valid state
        testRun, err := s.repo.GetByID(id)
        if err != nil {
                return nil, err
        }
        if testRun == nil {
                return nil, fmt.Errorf("test run not found")
        }

        // Check if test run is already running
        hasActive, err := s.intervalRepo.HasActiveInterval(id)
        if err != nil {
                return nil, err
        }
        if hasActive {
                return nil, fmt.Errorf("test run is already running")
        }

        // Update test run status to "In Progress" and set started_at if not already set
        status := "In Progress"
        req := models.UpdateTestRunRequest{
                Status: &status,
        }
        if testRun.StartedAt == nil {
                now := time.Now()
                req.StartedAt = &now
        }

        testRun, err = s.repo.Update(id, req)
        if err != nil {
                return nil, err
        }

        // Create new execution interval
        _, err = s.intervalRepo.Create(id)
        if err != nil {
                return nil, fmt.Errorf("failed to create execution interval: %w", err)
        }

        return testRun, nil
}

// PauseTestRun pauses a test run execution and closes the current interval
func (s *TestRunService) PauseTestRun(id int) (*models.TestRun, error) {
        // Check if test run exists and is running
        testRun, err := s.repo.GetByID(id)
        if err != nil {
                return nil, err
        }
        if testRun == nil {
                return nil, fmt.Errorf("test run not found")
        }

        // Check if test run has active interval
        hasActive, err := s.intervalRepo.HasActiveInterval(id)
        if err != nil {
                return nil, err
        }
        if !hasActive {
                return nil, fmt.Errorf("test run is not currently running")
        }

        // Close the active interval
        err = s.intervalRepo.CloseActiveInterval(id)
        if err != nil {
                return nil, err
        }

        // Update test run status to "Not Started" (paused state)
        status := "Not Started"
        req := models.UpdateTestRunRequest{
                Status: &status,
        }

        return s.repo.Update(id, req)
}

// FinishTestRun finishes a test run execution and closes any active intervals
func (s *TestRunService) FinishTestRun(id int) (*models.TestRun, error) {
        // Check if test run exists
        testRun, err := s.repo.GetByID(id)
        if err != nil {
                return nil, err
        }
        if testRun == nil {
                return nil, fmt.Errorf("test run not found")
        }

        // Only allow finishing test runs that are "In Progress"
        if testRun.Status != "In Progress" {
                return nil, fmt.Errorf("can only finish test runs that are in progress")
        }

        // Close any active interval
        hasActive, err := s.intervalRepo.HasActiveInterval(id)
        if err != nil {
                return nil, err
        }
        if hasActive {
                err = s.intervalRepo.CloseActiveInterval(id)
                if err != nil {
                        return nil, err
                }
        }

        // Update test run status to "Completed" and set completed_at
        now := time.Now()
        status := "Completed"
        req := models.UpdateTestRunRequest{
                Status:      &status,
                CompletedAt: &now,
        }

        return s.repo.Update(id, req)
}

// GetTestRunWithTimeTracking returns a test run with execution intervals and total time
func (s *TestRunService) GetTestRunWithTimeTracking(id int) (*models.TestRun, error) {
        testRun, err := s.repo.GetByID(id)
        if err != nil {
                return nil, err
        }
        if testRun == nil {
                return nil, fmt.Errorf("test run not found")
        }

        // Load execution intervals
        intervals, err := s.intervalRepo.GetByTestRunID(id)
        if err != nil {
                return nil, err
        }
        testRun.Intervals = intervals

        // Calculate total execution time
        totalTime, err := s.intervalRepo.CalculateTotalExecutionTime(id)
        if err != nil {
                return nil, err
        }
        testRun.TotalExecutionTime = &totalTime

        return testRun, nil
}

