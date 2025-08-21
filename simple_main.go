package main

import (
        "database/sql"
        "encoding/json"
        "html/template"
        "log"
        "net/http"
        "os"
        "path/filepath"
        "strconv"
        "strings"
        "time"

        _ "github.com/lib/pq"
)

type Project struct {
        ID          int       `json:"id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
}

type TestSuite struct {
        ID          int       `json:"id"`
        Name        string    `json:"name"`
        Description string    `json:"description"`
        ProjectID   int       `json:"project_id"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
        Project     *Project  `json:"project,omitempty"`
}

type TestCase struct {
        ID          int       `json:"id"`
        Title       string    `json:"title"`
        Description string    `json:"description"`
        Priority    string    `json:"priority"`
        Status      string    `json:"status"`
        TestSuiteID int       `json:"test_suite_id"`
        CreatedAt   time.Time `json:"created_at"`
        UpdatedAt   time.Time `json:"updated_at"`
        TestSuite   *TestSuite `json:"test_suite,omitempty"`
}

var db *sql.DB

func main() {
        // Initialize database connection
        var err error
        db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
        if err != nil {
                log.Fatal("Failed to connect to database:", err)
        }
        defer db.Close()

        // Test database connection
        if err := db.Ping(); err != nil {
                log.Fatal("Failed to ping database:", err)
        }
        log.Println("Connected to PostgreSQL database")

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
                rows, err := db.Query("SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC")
                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error querying projects: %v", err)
                        return
                }
                defer rows.Close()

                var projects []Project
                for rows.Next() {
                        var p Project
                        err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
                        if err != nil {
                                http.Error(w, "Scan error", http.StatusInternalServerError)
                                log.Printf("Error scanning project: %v", err)
                                return
                        }
                        projects = append(projects, p)
                }

                json.NewEncoder(w).Encode(projects)

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

                var project Project
                err := db.QueryRow(
                        "INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at",
                        req.Name, req.Description,
                ).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error creating project: %v", err)
                        return
                }

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

        switch r.Method {
        case "GET":
                var project Project
                err := db.QueryRow(
                        "SELECT id, name, description, created_at, updated_at FROM projects WHERE id = $1",
                        id,
                ).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

                if err == sql.ErrNoRows {
                        http.Error(w, "Project not found", http.StatusNotFound)
                        return
                } else if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error querying project: %v", err)
                        return
                }

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

                var project Project
                err := db.QueryRow(
                        "UPDATE projects SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, name, description, created_at, updated_at",
                        req.Name, req.Description, id,
                ).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)

                if err == sql.ErrNoRows {
                        http.Error(w, "Project not found", http.StatusNotFound)
                        return
                } else if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error updating project: %v", err)
                        return
                }

                json.NewEncoder(w).Encode(project)

        case "DELETE":
                result, err := db.Exec("DELETE FROM projects WHERE id = $1", id)
                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error deleting project: %v", err)
                        return
                }

                rowsAffected, err := result.RowsAffected()
                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        return
                }

                if rowsAffected == 0 {
                        http.Error(w, "Project not found", http.StatusNotFound)
                        return
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
                
                var query string
                var args []interface{}
                
                if projectIDStr != "" {
                        query = `
                                SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                       p.id, p.name, p.description, p.created_at, p.updated_at
                                FROM test_suites ts
                                JOIN projects p ON ts.project_id = p.id
                                WHERE ts.project_id = $1
                                ORDER BY ts.created_at DESC
                        `
                        projectID, err := strconv.Atoi(projectIDStr)
                        if err != nil {
                                http.Error(w, "Invalid project_id", http.StatusBadRequest)
                                return
                        }
                        args = []interface{}{projectID}
                } else {
                        query = `
                                SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                       p.id, p.name, p.description, p.created_at, p.updated_at
                                FROM test_suites ts
                                JOIN projects p ON ts.project_id = p.id
                                ORDER BY ts.created_at DESC
                        `
                }

                rows, err := db.Query(query, args...)
                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error querying test suites: %v", err)
                        return
                }
                defer rows.Close()

                var testSuites []TestSuite
                for rows.Next() {
                        var ts TestSuite
                        var p Project
                        err := rows.Scan(
                                &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                                &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                        )
                        if err != nil {
                                http.Error(w, "Scan error", http.StatusInternalServerError)
                                log.Printf("Error scanning test suite: %v", err)
                                return
                        }
                        ts.Project = &p
                        testSuites = append(testSuites, ts)
                }

                json.NewEncoder(w).Encode(testSuites)

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

                var testSuite TestSuite
                err := db.QueryRow(
                        "INSERT INTO test_suites (name, description, project_id) VALUES ($1, $2, $3) RETURNING id, name, description, project_id, created_at, updated_at",
                        req.Name, req.Description, req.ProjectID,
                ).Scan(&testSuite.ID, &testSuite.Name, &testSuite.Description, &testSuite.ProjectID, &testSuite.CreatedAt, &testSuite.UpdatedAt)

                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error creating test suite: %v", err)
                        return
                }

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

        switch r.Method {
        case "GET":
                var ts TestSuite
                var p Project
                err := db.QueryRow(`
                        SELECT ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at
                        FROM test_suites ts
                        JOIN projects p ON ts.project_id = p.id
                        WHERE ts.id = $1
                `, id).Scan(
                        &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                        &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                )

                if err == sql.ErrNoRows {
                        http.Error(w, "Test suite not found", http.StatusNotFound)
                        return
                } else if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error querying test suite: %v", err)
                        return
                }

                ts.Project = &p
                json.NewEncoder(w).Encode(ts)
        }
}

func testCasesAPIHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        switch r.Method {
        case "GET":
                testSuiteIDStr := r.URL.Query().Get("test_suite_id")
                
                var query string
                var args []interface{}
                
                if testSuiteIDStr != "" {
                        query = `
                                SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                                       ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                       p.id, p.name, p.description, p.created_at, p.updated_at
                                FROM test_cases tc
                                JOIN test_suites ts ON tc.test_suite_id = ts.id
                                JOIN projects p ON ts.project_id = p.id
                                WHERE tc.test_suite_id = $1
                                ORDER BY tc.created_at DESC
                        `
                        testSuiteID, err := strconv.Atoi(testSuiteIDStr)
                        if err != nil {
                                http.Error(w, "Invalid test_suite_id", http.StatusBadRequest)
                                return
                        }
                        args = []interface{}{testSuiteID}
                } else {
                        query = `
                                SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                                       ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                                       p.id, p.name, p.description, p.created_at, p.updated_at
                                FROM test_cases tc
                                JOIN test_suites ts ON tc.test_suite_id = ts.id
                                JOIN projects p ON ts.project_id = p.id
                                ORDER BY tc.created_at DESC
                        `
                }

                rows, err := db.Query(query, args...)
                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error querying test cases: %v", err)
                        return
                }
                defer rows.Close()

                var testCases []TestCase
                for rows.Next() {
                        var tc TestCase
                        var ts TestSuite
                        var p Project
                        err := rows.Scan(
                                &tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt,
                                &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                                &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                        )
                        if err != nil {
                                http.Error(w, "Scan error", http.StatusInternalServerError)
                                log.Printf("Error scanning test case: %v", err)
                                return
                        }
                        ts.Project = &p
                        tc.TestSuite = &ts
                        testCases = append(testCases, tc)
                }

                json.NewEncoder(w).Encode(testCases)

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

                var testCase TestCase
                err := db.QueryRow(
                        "INSERT INTO test_cases (title, description, priority, test_suite_id) VALUES ($1, $2, $3, $4) RETURNING id, title, description, priority, status, test_suite_id, created_at, updated_at",
                        req.Title, req.Description, req.Priority, req.TestSuiteID,
                ).Scan(&testCase.ID, &testCase.Title, &testCase.Description, &testCase.Priority, &testCase.Status, &testCase.TestSuiteID, &testCase.CreatedAt, &testCase.UpdatedAt)

                if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error creating test case: %v", err)
                        return
                }

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

        switch r.Method {
        case "GET":
                var tc TestCase
                var ts TestSuite
                var p Project
                err := db.QueryRow(`
                        SELECT tc.id, tc.title, tc.description, tc.priority, tc.status, tc.test_suite_id, tc.created_at, tc.updated_at,
                               ts.id, ts.name, ts.description, ts.project_id, ts.created_at, ts.updated_at,
                               p.id, p.name, p.description, p.created_at, p.updated_at
                        FROM test_cases tc
                        JOIN test_suites ts ON tc.test_suite_id = ts.id
                        JOIN projects p ON ts.project_id = p.id
                        WHERE tc.id = $1
                `, id).Scan(
                        &tc.ID, &tc.Title, &tc.Description, &tc.Priority, &tc.Status, &tc.TestSuiteID, &tc.CreatedAt, &tc.UpdatedAt,
                        &ts.ID, &ts.Name, &ts.Description, &ts.ProjectID, &ts.CreatedAt, &ts.UpdatedAt,
                        &p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt,
                )

                if err == sql.ErrNoRows {
                        http.Error(w, "Test case not found", http.StatusNotFound)
                        return
                } else if err != nil {
                        http.Error(w, "Database error", http.StatusInternalServerError)
                        log.Printf("Error querying test case: %v", err)
                        return
                }

                ts.Project = &p
                tc.TestSuite = &ts
                json.NewEncoder(w).Encode(tc)
        }
}