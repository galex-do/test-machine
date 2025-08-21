package main

import (
        "encoding/json"
        "fmt"
        "html/template"
        "log"
        "net/http"
        "path/filepath"
        "strconv"
        "strings"
)

type Project struct {
        ID          int    `json:"id"`
        Name        string `json:"name"`
        Description string `json:"description"`
        CreatedAt   string `json:"created_at"`
}

type TestSuite struct {
        ID          int     `json:"id"`
        Name        string  `json:"name"`
        Description string  `json:"description"`
        ProjectID   int     `json:"project_id"`
        CreatedAt   string  `json:"created_at"`
        Project     Project `json:"project,omitempty"`
}

type TestCase struct {
        ID          int       `json:"id"`
        Title       string    `json:"title"`
        Description string    `json:"description"`
        Priority    string    `json:"priority"`
        Status      string    `json:"status"`
        TestSuiteID int       `json:"test_suite_id"`
        CreatedAt   string    `json:"created_at"`
        TestSuite   TestSuite `json:"test_suite,omitempty"`
}

// In-memory storage
var projects []Project
var testSuites []TestSuite
var testCases []TestCase
var nextProjectID = 1
var nextTestSuiteID = 1
var nextTestCaseID = 1

func main() {
        // Initialize sample data
        initSampleData()

        // Serve static files
        http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

        // Frontend routes
        http.HandleFunc("/", indexHandler)
        http.HandleFunc("/project/", projectHandler)
        http.HandleFunc("/test-suite/", testSuiteHandler)
        http.HandleFunc("/test-case/", testCaseHandler)
        http.HandleFunc("/test-run/", testRunHandler)
        http.HandleFunc("/reports", reportsHandler)

        // API routes
        http.HandleFunc("/api/projects", projectsAPIHandler)
        http.HandleFunc("/api/projects/", projectAPIHandler)
        http.HandleFunc("/api/test-suites", testSuitesAPIHandler)
        http.HandleFunc("/api/test-suites/", testSuiteAPIHandler)
        http.HandleFunc("/api/test-cases", testCasesAPIHandler)
        http.HandleFunc("/api/test-cases/", testCaseAPIHandler)

        log.Println("Server starting on port 5000...")
        log.Fatal(http.ListenAndServe(":5000", nil))
}

func initSampleData() {
        // Add sample project
        projects = append(projects, Project{
                ID:          nextProjectID,
                Name:        "Web Application Testing",
                Description: "Testing suite for the main web application",
                CreatedAt:   "2025-01-20T10:00:00Z",
        })
        nextProjectID++

        projects = append(projects, Project{
                ID:          nextProjectID,
                Name:        "Mobile App Testing",
                Description: "Testing suite for the mobile application",
                CreatedAt:   "2025-01-21T09:00:00Z",
        })
        nextProjectID++

        // Add sample test suite
        testSuites = append(testSuites, TestSuite{
                ID:          nextTestSuiteID,
                Name:        "User Authentication",
                Description: "Test cases for user login, registration, and password reset",
                ProjectID:   1,
                CreatedAt:   "2025-01-20T11:00:00Z",
        })
        nextTestSuiteID++

        testSuites = append(testSuites, TestSuite{
                ID:          nextTestSuiteID,
                Name:        "E-commerce Checkout",
                Description: "Test cases for shopping cart and checkout process",
                ProjectID:   1,
                CreatedAt:   "2025-01-20T12:00:00Z",
        })
        nextTestSuiteID++

        // Add sample test cases
        testCases = append(testCases, TestCase{
                ID:          nextTestCaseID,
                Title:       "Valid User Login",
                Description: "Test successful login with valid credentials",
                Priority:    "High",
                Status:      "Active",
                TestSuiteID: 1,
                CreatedAt:   "2025-01-20T11:30:00Z",
        })
        nextTestCaseID++

        testCases = append(testCases, TestCase{
                ID:          nextTestCaseID,
                Title:       "Invalid Password Login",
                Description: "Test login failure with invalid password",
                Priority:    "High",
                Status:      "Active",
                TestSuiteID: 1,
                CreatedAt:   "2025-01-20T11:45:00Z",
        })
        nextTestCaseID++

        testCases = append(testCases, TestCase{
                ID:          nextTestCaseID,
                Title:       "Add Item to Cart",
                Description: "Test adding products to shopping cart",
                Priority:    "Medium",
                Status:      "Active",
                TestSuiteID: 2,
                CreatedAt:   "2025-01-20T12:30:00Z",
        })
        nextTestCaseID++
}

// Template helper functions
func parseTemplates(templateName string) *template.Template {
        tmpl, err := template.ParseFiles(filepath.Join("templates", templateName))
        if err != nil {
                log.Printf("Error parsing template %s: %v", templateName, err)
                return nil
        }
        return tmpl
}

// Frontend handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {
        tmpl := parseTemplates("index.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }
        
        data := map[string]interface{}{
                "title": "Dashboard",
        }
        tmpl.Execute(w, data)
}

func projectHandler(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/project/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Invalid project ID", http.StatusBadRequest)
                return
        }

        tmpl := parseTemplates("project.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }
        
        data := map[string]interface{}{
                "title":     "Project Details",
                "projectID": id,
        }
        tmpl.Execute(w, data)
}

func testSuiteHandler(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/test-suite/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Invalid test suite ID", http.StatusBadRequest)
                return
        }

        tmpl := parseTemplates("test-suite.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }
        
        data := map[string]interface{}{
                "title":       "Test Suite Details",
                "testSuiteID": id,
        }
        tmpl.Execute(w, data)
}

func testCaseHandler(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/test-case/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Invalid test case ID", http.StatusBadRequest)
                return
        }

        tmpl := parseTemplates("test-case.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }
        
        data := map[string]interface{}{
                "title":      "Test Case Details",
                "testCaseID": id,
        }
        tmpl.Execute(w, data)
}

func testRunHandler(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/test-run/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Invalid test run ID", http.StatusBadRequest)
                return
        }

        tmpl := parseTemplates("test-run.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }
        
        data := map[string]interface{}{
                "title":     "Test Run Details",
                "testRunID": id,
        }
        tmpl.Execute(w, data)
}

func reportsHandler(w http.ResponseWriter, r *http.Request) {
        tmpl := parseTemplates("reports.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }
        
        data := map[string]interface{}{
                "title": "Reports",
        }
        tmpl.Execute(w, data)
}

// API handlers
func projectsAPIHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        switch r.Method {
        case "GET":
                // Return all projects with their test suites
                var enrichedProjects []Project
                for _, project := range projects {
                        enrichedProject := project
                        // Find test suites for this project
                        var projectTestSuites []TestSuite
                        for _, testSuite := range testSuites {
                                if testSuite.ProjectID == project.ID {
                                        projectTestSuites = append(projectTestSuites, testSuite)
                                }
                        }
                        enrichedProjects = append(enrichedProjects, enrichedProject)
                }
                json.NewEncoder(w).Encode(enrichedProjects)
                
        case "POST":
                var req struct {
                        Name        string `json:"name"`
                        Description string `json:"description"`
                }
                
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                        http.Error(w, "Invalid JSON", http.StatusBadRequest)
                        return
                }
                
                if req.Name == "" {
                        http.Error(w, "Name is required", http.StatusBadRequest)
                        return
                }
                
                project := Project{
                        ID:          nextProjectID,
                        Name:        req.Name,
                        Description: req.Description,
                        CreatedAt:   "2025-01-21T" + fmt.Sprintf("%02d:00:00Z", len(projects)+10),
                }
                nextProjectID++
                projects = append(projects, project)
                
                w.WriteHeader(http.StatusCreated)
                json.NewEncoder(w).Encode(project)
        }
}

func projectAPIHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Invalid project ID", http.StatusBadRequest)
                return
        }
        
        var project *Project
        for i := range projects {
                if projects[i].ID == id {
                        project = &projects[i]
                        break
                }
        }
        
        if project == nil {
                http.Error(w, "Project not found", http.StatusNotFound)
                return
        }
        
        switch r.Method {
        case "GET":
                json.NewEncoder(w).Encode(project)
        case "PUT":
                var req struct {
                        Name        string `json:"name"`
                        Description string `json:"description"`
                }
                
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                        http.Error(w, "Invalid JSON", http.StatusBadRequest)
                        return
                }
                
                if req.Name == "" {
                        http.Error(w, "Name is required", http.StatusBadRequest)
                        return
                }
                
                project.Name = req.Name
                project.Description = req.Description
                
                json.NewEncoder(w).Encode(project)
        case "DELETE":
                // Remove project
                for i := range projects {
                        if projects[i].ID == id {
                                projects = append(projects[:i], projects[i+1:]...)
                                break
                        }
                }
                
                w.WriteHeader(http.StatusOK)
                json.NewEncoder(w).Encode(map[string]string{"message": "Project deleted successfully"})
        }
}

func testSuitesAPIHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        switch r.Method {
        case "GET":
                projectIDStr := r.URL.Query().Get("project_id")
                var filteredTestSuites []TestSuite
                
                for _, testSuite := range testSuites {
                        if projectIDStr == "" || testSuite.ProjectID == mustParseInt(projectIDStr) {
                                // Add project info
                                for _, project := range projects {
                                        if project.ID == testSuite.ProjectID {
                                                testSuite.Project = project
                                                break
                                        }
                                }
                                filteredTestSuites = append(filteredTestSuites, testSuite)
                        }
                }
                
                json.NewEncoder(w).Encode(filteredTestSuites)
                
        case "POST":
                var req struct {
                        Name        string `json:"name"`
                        Description string `json:"description"`
                        ProjectID   int    `json:"project_id"`
                }
                
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                        http.Error(w, "Invalid JSON", http.StatusBadRequest)
                        return
                }
                
                if req.Name == "" || req.ProjectID == 0 {
                        http.Error(w, "Name and project_id are required", http.StatusBadRequest)
                        return
                }
                
                testSuite := TestSuite{
                        ID:          nextTestSuiteID,
                        Name:        req.Name,
                        Description: req.Description,
                        ProjectID:   req.ProjectID,
                        CreatedAt:   "2025-01-21T" + fmt.Sprintf("%02d:00:00Z", len(testSuites)+10),
                }
                nextTestSuiteID++
                testSuites = append(testSuites, testSuite)
                
                w.WriteHeader(http.StatusCreated)
                json.NewEncoder(w).Encode(testSuite)
        }
}

func testSuiteAPIHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        idStr := strings.TrimPrefix(r.URL.Path, "/api/test-suites/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Invalid test suite ID", http.StatusBadRequest)
                return
        }
        
        var testSuite *TestSuite
        for i := range testSuites {
                if testSuites[i].ID == id {
                        testSuite = &testSuites[i]
                        break
                }
        }
        
        if testSuite == nil {
                http.Error(w, "Test suite not found", http.StatusNotFound)
                return
        }
        
        switch r.Method {
        case "GET":
                // Add project info
                for _, project := range projects {
                        if project.ID == testSuite.ProjectID {
                                testSuite.Project = project
                                break
                        }
                }
                json.NewEncoder(w).Encode(testSuite)
        }
}

func testCasesAPIHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        switch r.Method {
        case "GET":
                testSuiteIDStr := r.URL.Query().Get("test_suite_id")
                var filteredTestCases []TestCase
                
                for _, testCase := range testCases {
                        if testSuiteIDStr == "" || testCase.TestSuiteID == mustParseInt(testSuiteIDStr) {
                                // Add test suite info
                                for _, testSuite := range testSuites {
                                        if testSuite.ID == testCase.TestSuiteID {
                                                testCase.TestSuite = testSuite
                                                // Add project info to test suite
                                                for _, project := range projects {
                                                        if project.ID == testSuite.ProjectID {
                                                                testCase.TestSuite.Project = project
                                                                break
                                                        }
                                                }
                                                break
                                        }
                                }
                                filteredTestCases = append(filteredTestCases, testCase)
                        }
                }
                
                json.NewEncoder(w).Encode(filteredTestCases)
                
        case "POST":
                var req struct {
                        Title       string `json:"title"`
                        Description string `json:"description"`
                        Priority    string `json:"priority"`
                        TestSuiteID int    `json:"test_suite_id"`
                }
                
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                        http.Error(w, "Invalid JSON", http.StatusBadRequest)
                        return
                }
                
                if req.Title == "" || req.TestSuiteID == 0 {
                        http.Error(w, "Title and test_suite_id are required", http.StatusBadRequest)
                        return
                }
                
                if req.Priority == "" {
                        req.Priority = "Medium"
                }
                
                testCase := TestCase{
                        ID:          nextTestCaseID,
                        Title:       req.Title,
                        Description: req.Description,
                        Priority:    req.Priority,
                        Status:      "Active",
                        TestSuiteID: req.TestSuiteID,
                        CreatedAt:   "2025-01-21T" + fmt.Sprintf("%02d:00:00Z", len(testCases)+10),
                }
                nextTestCaseID++
                testCases = append(testCases, testCase)
                
                w.WriteHeader(http.StatusCreated)
                json.NewEncoder(w).Encode(testCase)
        }
}

func testCaseAPIHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        idStr := strings.TrimPrefix(r.URL.Path, "/api/test-cases/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                http.Error(w, "Invalid test case ID", http.StatusBadRequest)
                return
        }
        
        var testCase *TestCase
        for i := range testCases {
                if testCases[i].ID == id {
                        testCase = &testCases[i]
                        break
                }
        }
        
        if testCase == nil {
                http.Error(w, "Test case not found", http.StatusNotFound)
                return
        }
        
        switch r.Method {
        case "GET":
                // Add test suite and project info
                for _, testSuite := range testSuites {
                        if testSuite.ID == testCase.TestSuiteID {
                                testCase.TestSuite = testSuite
                                for _, project := range projects {
                                        if project.ID == testSuite.ProjectID {
                                                testCase.TestSuite.Project = project
                                                break
                                        }
                                }
                                break
                        }
                }
                json.NewEncoder(w).Encode(testCase)
        }
}

func mustParseInt(s string) int {
        i, _ := strconv.Atoi(s)
        return i
}