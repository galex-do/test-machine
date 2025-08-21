package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"test-management/database"
	"test-management/handlers"
)

func main() {
	// Initialize database
	database.InitDatabase()

	// Create Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Serve static files
	r.Static("/static", "./static")

	// Frontend routes
	r.GET("/", handlers.IndexHandler)
	r.GET("/project/:id", handlers.ProjectHandler)
	r.GET("/test-suite/:id", handlers.TestSuiteHandler)
	r.GET("/test-case/:id", handlers.TestCaseHandler)
	r.GET("/test-run/:id", handlers.TestRunHandler)
	r.GET("/reports", handlers.ReportsHandler)

	// API routes
	api := r.Group("/api")
	{
		// Projects
		api.GET("/projects", handlers.GetProjects)
		api.POST("/projects", handlers.CreateProject)
		api.GET("/projects/:id", handlers.GetProject)
		api.PUT("/projects/:id", handlers.UpdateProject)
		api.DELETE("/projects/:id", handlers.DeleteProject)

		// Test Suites
		api.GET("/test-suites", handlers.GetTestSuites)
		api.POST("/test-suites", handlers.CreateTestSuite)
		api.GET("/test-suites/:id", handlers.GetTestSuite)
		api.PUT("/test-suites/:id", handlers.UpdateTestSuite)
		api.DELETE("/test-suites/:id", handlers.DeleteTestSuite)

		// Test Cases
		api.GET("/test-cases", handlers.GetTestCases)
		api.POST("/test-cases", handlers.CreateTestCase)
		api.GET("/test-cases/:id", handlers.GetTestCase)
		api.PUT("/test-cases/:id", handlers.UpdateTestCase)
		api.DELETE("/test-cases/:id", handlers.DeleteTestCase)
		api.GET("/test-cases/search", handlers.SearchTestCases)

		// Test Steps
		api.GET("/test-cases/:id/steps", handlers.GetTestSteps)
		api.POST("/test-cases/:id/steps", handlers.CreateTestStep)
		api.PUT("/test-steps/:id", handlers.UpdateTestStep)
		api.DELETE("/test-steps/:id", handlers.DeleteTestStep)

		// Test Runs
		api.GET("/test-runs", handlers.GetTestRuns)
		api.POST("/test-runs", handlers.CreateTestRun)
		api.GET("/test-runs/:id", handlers.GetTestRun)
		api.PUT("/test-runs/:id", handlers.UpdateTestRun)
		api.DELETE("/test-runs/:id", handlers.DeleteTestRun)

		// Test Executions
		api.GET("/test-runs/:id/executions", handlers.GetTestExecutions)
		api.POST("/test-runs/:id/executions", handlers.CreateTestExecution)
		api.PUT("/test-executions/:id", handlers.UpdateTestExecution)

		// Reports
		api.GET("/reports/summary", handlers.GetReportSummary)
		api.GET("/reports/export/:format", handlers.ExportReport)
	}

	// Start server on port 5000
	log.Println("Server starting on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", r))
}
