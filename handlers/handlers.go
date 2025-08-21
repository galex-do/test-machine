package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"test-management/database"
	"test-management/models"
)

// Frontend handlers

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Test Management Platform",
	})
}

func ProjectHandler(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusOK, "project.html", gin.H{
		"title":     "Project Details",
		"projectID": id,
	})
}

func TestSuiteHandler(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusOK, "test-suite.html", gin.H{
		"title":       "Test Suite Details",
		"testSuiteID": id,
	})
}

func TestCaseHandler(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusOK, "test-case.html", gin.H{
		"title":      "Test Case Details",
		"testCaseID": id,
	})
}

func TestRunHandler(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusOK, "test-run.html", gin.H{
		"title":     "Test Run Details",
		"testRunID": id,
	})
}

func ReportsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "reports.html", gin.H{
		"title": "Reports",
	})
}

// API handlers

// Project handlers
func GetProjects(c *gin.Context) {
	var projects []models.Project
	db := database.GetDB()
	
	if err := db.Preload("TestSuites").Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	
	c.JSON(http.StatusOK, projects)
}

func CreateProject(c *gin.Context) {
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := models.Project{
		Name:        req.Name,
		Description: req.Description,
	}

	db := database.GetDB()
	if err := db.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func GetProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	db := database.GetDB()
	
	if err := db.Preload("TestSuites.TestCases").First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var project models.Project
	
	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	project.Name = req.Name
	project.Description = req.Description

	if err := db.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	db := database.GetDB()
	
	if err := db.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// Test Suite handlers
func GetTestSuites(c *gin.Context) {
	var testSuites []models.TestSuite
	db := database.GetDB()
	
	query := db.Preload("Project").Preload("TestCases")
	
	if projectID := c.Query("project_id"); projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	
	if err := query.Find(&testSuites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test suites"})
		return
	}
	
	c.JSON(http.StatusOK, testSuites)
}

func CreateTestSuite(c *gin.Context) {
	var req models.CreateTestSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testSuite := models.TestSuite{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
	}

	db := database.GetDB()
	if err := db.Create(&testSuite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test suite"})
		return
	}

	c.JSON(http.StatusCreated, testSuite)
}

func GetTestSuite(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test suite ID"})
		return
	}

	var testSuite models.TestSuite
	db := database.GetDB()
	
	if err := db.Preload("Project").Preload("TestCases.TestSteps").First(&testSuite, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test suite not found"})
		return
	}

	c.JSON(http.StatusOK, testSuite)
}

func UpdateTestSuite(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test suite ID"})
		return
	}

	var req models.CreateTestSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var testSuite models.TestSuite
	
	if err := db.First(&testSuite, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test suite not found"})
		return
	}

	testSuite.Name = req.Name
	testSuite.Description = req.Description
	testSuite.ProjectID = req.ProjectID

	if err := db.Save(&testSuite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test suite"})
		return
	}

	c.JSON(http.StatusOK, testSuite)
}

func DeleteTestSuite(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test suite ID"})
		return
	}

	db := database.GetDB()
	
	if err := db.Delete(&models.TestSuite{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test suite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test suite deleted successfully"})
}

// Test Case handlers
func GetTestCases(c *gin.Context) {
	var testCases []models.TestCase
	db := database.GetDB()
	
	query := db.Preload("TestSuite.Project").Preload("TestSteps")
	
	if testSuiteID := c.Query("test_suite_id"); testSuiteID != "" {
		query = query.Where("test_suite_id = ?", testSuiteID)
	}
	
	if err := query.Find(&testCases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test cases"})
		return
	}
	
	c.JSON(http.StatusOK, testCases)
}

func CreateTestCase(c *gin.Context) {
	var req models.CreateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Priority == "" {
		req.Priority = "Medium"
	}

	testCase := models.TestCase{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		TestSuiteID: req.TestSuiteID,
		Status:      "Active",
	}

	db := database.GetDB()
	if err := db.Create(&testCase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test case"})
		return
	}

	c.JSON(http.StatusCreated, testCase)
}

func GetTestCase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test case ID"})
		return
	}

	var testCase models.TestCase
	db := database.GetDB()
	
	if err := db.Preload("TestSuite.Project").Preload("TestSteps").First(&testCase, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test case not found"})
		return
	}

	c.JSON(http.StatusOK, testCase)
}

func UpdateTestCase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test case ID"})
		return
	}

	var req models.CreateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var testCase models.TestCase
	
	if err := db.First(&testCase, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test case not found"})
		return
	}

	testCase.Title = req.Title
	testCase.Description = req.Description
	testCase.Priority = req.Priority
	testCase.TestSuiteID = req.TestSuiteID

	if err := db.Save(&testCase).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test case"})
		return
	}

	c.JSON(http.StatusOK, testCase)
}

func DeleteTestCase(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test case ID"})
		return
	}

	db := database.GetDB()
	
	if err := db.Delete(&models.TestCase{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test case"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test case deleted successfully"})
}

func SearchTestCases(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var testCases []models.TestCase
	db := database.GetDB()
	
	searchQuery := "%" + strings.ToLower(query) + "%"
	
	if err := db.Preload("TestSuite.Project").Preload("TestSteps").
		Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchQuery, searchQuery).
		Find(&testCases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search test cases"})
		return
	}

	c.JSON(http.StatusOK, testCases)
}

// Test Step handlers
func GetTestSteps(c *gin.Context) {
	testCaseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test case ID"})
		return
	}

	var testSteps []models.TestStep
	db := database.GetDB()
	
	if err := db.Where("test_case_id = ?", testCaseID).Order("step_number").Find(&testSteps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test steps"})
		return
	}

	c.JSON(http.StatusOK, testSteps)
}

func CreateTestStep(c *gin.Context) {
	testCaseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test case ID"})
		return
	}

	var req models.CreateTestStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testStep := models.TestStep{
		StepNumber:     req.StepNumber,
		Description:    req.Description,
		ExpectedResult: req.ExpectedResult,
		TestCaseID:     uint(testCaseID),
	}

	db := database.GetDB()
	if err := db.Create(&testStep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test step"})
		return
	}

	c.JSON(http.StatusCreated, testStep)
}

func UpdateTestStep(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test step ID"})
		return
	}

	var req models.CreateTestStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var testStep models.TestStep
	
	if err := db.First(&testStep, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test step not found"})
		return
	}

	testStep.StepNumber = req.StepNumber
	testStep.Description = req.Description
	testStep.ExpectedResult = req.ExpectedResult

	if err := db.Save(&testStep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test step"})
		return
	}

	c.JSON(http.StatusOK, testStep)
}

func DeleteTestStep(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test step ID"})
		return
	}

	db := database.GetDB()
	
	if err := db.Delete(&models.TestStep{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test step"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test step deleted successfully"})
}

// Test Run handlers
func GetTestRuns(c *gin.Context) {
	var testRuns []models.TestRun
	db := database.GetDB()
	
	query := db.Preload("TestSuite.Project").Preload("TestExecutions.TestCase")
	
	if testSuiteID := c.Query("test_suite_id"); testSuiteID != "" {
		query = query.Where("test_suite_id = ?", testSuiteID)
	}
	
	if err := query.Order("created_at DESC").Find(&testRuns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test runs"})
		return
	}
	
	c.JSON(http.StatusOK, testRuns)
}

func CreateTestRun(c *gin.Context) {
	var req models.CreateTestRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	testRun := models.TestRun{
		Name:        req.Name,
		Description: req.Description,
		TestSuiteID: req.TestSuiteID,
		Status:      "In Progress",
		StartedAt:   &now,
	}

	db := database.GetDB()
	
	// Start transaction
	tx := db.Begin()
	
	if err := tx.Create(&testRun).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test run"})
		return
	}

	// Create test executions for all test cases in the suite
	var testCases []models.TestCase
	if err := tx.Where("test_suite_id = ?", req.TestSuiteID).Find(&testCases).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test cases"})
		return
	}

	for _, testCase := range testCases {
		testExecution := models.TestExecution{
			TestRunID:  testRun.ID,
			TestCaseID: testCase.ID,
			Status:     "Not Executed",
		}
		if err := tx.Create(&testExecution).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test executions"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, testRun)
}

func GetTestRun(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test run ID"})
		return
	}

	var testRun models.TestRun
	db := database.GetDB()
	
	if err := db.Preload("TestSuite.Project").
		Preload("TestExecutions.TestCase.TestSteps").
		First(&testRun, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test run not found"})
		return
	}

	c.JSON(http.StatusOK, testRun)
}

func UpdateTestRun(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test run ID"})
		return
	}

	var req models.CreateTestRunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var testRun models.TestRun
	
	if err := db.First(&testRun, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test run not found"})
		return
	}

	testRun.Name = req.Name
	testRun.Description = req.Description

	if err := db.Save(&testRun).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test run"})
		return
	}

	c.JSON(http.StatusOK, testRun)
}

func DeleteTestRun(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test run ID"})
		return
	}

	db := database.GetDB()
	
	if err := db.Delete(&models.TestRun{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete test run"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test run deleted successfully"})
}

// Test Execution handlers
func GetTestExecutions(c *gin.Context) {
	testRunID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test run ID"})
		return
	}

	var testExecutions []models.TestExecution
	db := database.GetDB()
	
	if err := db.Preload("TestCase.TestSteps").
		Where("test_run_id = ?", testRunID).
		Find(&testExecutions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch test executions"})
		return
	}

	c.JSON(http.StatusOK, testExecutions)
}

func CreateTestExecution(c *gin.Context) {
	testRunID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test run ID"})
		return
	}

	var req struct {
		TestCaseID uint `json:"test_case_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	testExecution := models.TestExecution{
		TestRunID:  uint(testRunID),
		TestCaseID: req.TestCaseID,
		Status:     "Not Executed",
	}

	db := database.GetDB()
	if err := db.Create(&testExecution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create test execution"})
		return
	}

	c.JSON(http.StatusCreated, testExecution)
}

func UpdateTestExecution(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test execution ID"})
		return
	}

	var req models.UpdateTestExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var testExecution models.TestExecution
	
	if err := db.First(&testExecution, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Test execution not found"})
		return
	}

	testExecution.Status = req.Status
	testExecution.Notes = req.Notes
	testExecution.ExecutedBy = req.ExecutedBy
	
	if req.Status != "Not Executed" {
		now := time.Now()
		testExecution.ExecutedAt = &now
	}

	if err := db.Save(&testExecution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update test execution"})
		return
	}

	// Check if test run should be marked as completed
	var pendingExecutions int64
	db.Model(&models.TestExecution{}).
		Where("test_run_id = ? AND status = ?", testExecution.TestRunID, "Not Executed").
		Count(&pendingExecutions)

	if pendingExecutions == 0 {
		var testRun models.TestRun
		if err := db.First(&testRun, testExecution.TestRunID).Error; err == nil {
			now := time.Now()
			testRun.Status = "Completed"
			testRun.CompletedAt = &now
			db.Save(&testRun)
		}
	}

	c.JSON(http.StatusOK, testExecution)
}

// Report handlers
func GetReportSummary(c *gin.Context) {
	db := database.GetDB()
	var summary models.ReportSummary

	// Count totals
	db.Model(&models.Project{}).Count(&summary.TotalProjects)
	db.Model(&models.TestSuite{}).Count(&summary.TotalTestSuites)
	db.Model(&models.TestCase{}).Count(&summary.TotalTestCases)
	db.Model(&models.TestRun{}).Count(&summary.TotalTestRuns)

	// Execution summary
	db.Model(&models.TestExecution{}).Count(&summary.ExecutionSummary.TotalExecutions)
	db.Model(&models.TestExecution{}).Where("status = ?", "Pass").Count(&summary.ExecutionSummary.Passed)
	db.Model(&models.TestExecution{}).Where("status = ?", "Fail").Count(&summary.ExecutionSummary.Failed)
	db.Model(&models.TestExecution{}).Where("status = ?", "Blocked").Count(&summary.ExecutionSummary.Blocked)
	db.Model(&models.TestExecution{}).Where("status = ?", "Skip").Count(&summary.ExecutionSummary.Skipped)
	db.Model(&models.TestExecution{}).Where("status = ?", "Not Executed").Count(&summary.ExecutionSummary.NotExecuted)

	if summary.ExecutionSummary.TotalExecutions > 0 {
		summary.ExecutionSummary.PassRate = float64(summary.ExecutionSummary.Passed) / float64(summary.ExecutionSummary.TotalExecutions) * 100
	}

	// Recent test runs
	db.Preload("TestSuite.Project").
		Order("created_at DESC").
		Limit(10).
		Find(&summary.RecentRuns)

	// Project summaries
	var projects []models.Project
	db.Find(&projects)

	for _, project := range projects {
		var projectSummary models.ProjectExecutionSummary
		projectSummary.ProjectID = project.ID
		projectSummary.ProjectName = project.Name

		// Count executions for this project
		db.Table("test_executions").
			Joins("JOIN test_cases ON test_executions.test_case_id = test_cases.id").
			Joins("JOIN test_suites ON test_cases.test_suite_id = test_suites.id").
			Where("test_suites.project_id = ?", project.ID).
			Count(&projectSummary.Summary.TotalExecutions)

		db.Table("test_executions").
			Joins("JOIN test_cases ON test_executions.test_case_id = test_cases.id").
			Joins("JOIN test_suites ON test_cases.test_suite_id = test_suites.id").
			Where("test_suites.project_id = ? AND test_executions.status = ?", project.ID, "Pass").
			Count(&projectSummary.Summary.Passed)

		db.Table("test_executions").
			Joins("JOIN test_cases ON test_executions.test_case_id = test_cases.id").
			Joins("JOIN test_suites ON test_cases.test_suite_id = test_suites.id").
			Where("test_suites.project_id = ? AND test_executions.status = ?", project.ID, "Fail").
			Count(&projectSummary.Summary.Failed)

		db.Table("test_executions").
			Joins("JOIN test_cases ON test_executions.test_case_id = test_cases.id").
			Joins("JOIN test_suites ON test_cases.test_suite_id = test_suites.id").
			Where("test_suites.project_id = ? AND test_executions.status = ?", project.ID, "Blocked").
			Count(&projectSummary.Summary.Blocked)

		db.Table("test_executions").
			Joins("JOIN test_cases ON test_executions.test_case_id = test_cases.id").
			Joins("JOIN test_suites ON test_cases.test_suite_id = test_suites.id").
			Where("test_suites.project_id = ? AND test_executions.status = ?", project.ID, "Skip").
			Count(&projectSummary.Summary.Skipped)

		db.Table("test_executions").
			Joins("JOIN test_cases ON test_executions.test_case_id = test_cases.id").
			Joins("JOIN test_suites ON test_cases.test_suite_id = test_suites.id").
			Where("test_suites.project_id = ? AND test_executions.status = ?", project.ID, "Not Executed").
			Count(&projectSummary.Summary.NotExecuted)

		if projectSummary.Summary.TotalExecutions > 0 {
			projectSummary.Summary.PassRate = float64(projectSummary.Summary.Passed) / float64(projectSummary.Summary.TotalExecutions) * 100
		}

		summary.ProjectSummaries = append(summary.ProjectSummaries, projectSummary)
	}

	c.JSON(http.StatusOK, summary)
}

func ExportReport(c *gin.Context) {
	format := c.Param("format")
	
	if format != "csv" && format != "json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported format. Use 'csv' or 'json'"})
		return
	}

	db := database.GetDB()
	var executions []models.TestExecution

	db.Preload("TestRun.TestSuite.Project").
		Preload("TestCase").
		Find(&executions)

	if format == "json" {
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=test_report_%s.json", time.Now().Format("2006-01-02")))
		c.JSON(http.StatusOK, executions)
		return
	}

	// CSV format
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=test_report_%s.csv", time.Now().Format("2006-01-02")))

	csvData := "Project,Test Suite,Test Case,Status,Executed By,Executed At,Notes\n"
	for _, execution := range executions {
		projectName := ""
		testSuiteName := ""
		if execution.TestRun.TestSuite.Project.Name != "" {
			projectName = execution.TestRun.TestSuite.Project.Name
		}
		if execution.TestRun.TestSuite.Name != "" {
			testSuiteName = execution.TestRun.TestSuite.Name
		}
		
		executedAt := ""
		if execution.ExecutedAt != nil {
			executedAt = execution.ExecutedAt.Format("2006-01-02 15:04:05")
		}
		
		csvData += fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
			projectName,
			testSuiteName,
			execution.TestCase.Title,
			execution.Status,
			execution.ExecutedBy,
			executedAt,
			strings.ReplaceAll(execution.Notes, ",", ";"),
		)
	}

	c.String(http.StatusOK, csvData)
}
