package handlers

import (
        "encoding/json"
        "html/template"
        "log"
        "net/http"
        "path/filepath"
        "strconv"
        "strings"

        "github.com/galex-do/test-machine/internal/service"
)

// Handler holds all the dependencies for HTTP handlers
type Handler struct {
        projectService   *service.ProjectService
        testSuiteService *service.TestSuiteService
        testCaseService  *service.TestCaseService
        testRunService   *service.TestRunService
}

// NewHandler creates a new handler
func NewHandler(projectService *service.ProjectService, testSuiteService *service.TestSuiteService, testCaseService *service.TestCaseService, testRunService *service.TestRunService) *Handler {
        return &Handler{
                projectService:   projectService,
                testSuiteService: testSuiteService,
                testCaseService:  testCaseService,
                testRunService:   testRunService,
        }
}

// SetupRoutes sets up all the routes
func (h *Handler) SetupRoutes() *http.ServeMux {
        mux := http.NewServeMux()

        // Serve static files
        mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

        // Frontend routes
        mux.HandleFunc("/", h.indexHandler)
        mux.HandleFunc("/project/", h.projectHandler)
        mux.HandleFunc("/reports", h.reportsHandler)

        // API routes
        mux.HandleFunc("/api/projects", h.projectsAPIHandler)
        mux.HandleFunc("/api/projects/", h.projectAPIHandler)
        mux.HandleFunc("/api/test-suites", h.testSuitesAPIHandler)
        mux.HandleFunc("/api/test-suites/", h.testSuiteAPIHandler)
        mux.HandleFunc("/api/test-cases", h.testCasesAPIHandler)
        mux.HandleFunc("/api/test-cases/", h.testCaseAPIHandler)
        mux.HandleFunc("/api/test-runs", h.testRunsAPIHandler)
        mux.HandleFunc("/api/test-runs/", h.testRunAPIHandler)

        return mux
}

// parseTemplates parses HTML templates
func (h *Handler) parseTemplates(templateName string) *template.Template {
        tmpl, err := template.ParseFiles(filepath.Join("templates", templateName))
        if err != nil {
                log.Printf("Error parsing template %s: %v", templateName, err)
                return nil
        }
        return tmpl
}

// writeJSONError writes a JSON error response
func (h *Handler) writeJSONError(w http.ResponseWriter, message string, statusCode int) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeJSONResponse writes a JSON response
func (h *Handler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(data)
}

// Frontend handlers
func (h *Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
        tmpl := h.parseTemplates("index.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }

        data := map[string]interface{}{
                "title": "Dashboard",
        }
        tmpl.Execute(w, data)
}

func (h *Handler) projectHandler(w http.ResponseWriter, r *http.Request) {
        path := strings.TrimPrefix(r.URL.Path, "/project/")
        parts := strings.Split(path, "/")
        
        if len(parts) == 0 || parts[0] == "" {
                http.Error(w, "Invalid project URL", http.StatusBadRequest)
                return
        }
        
        projectID, err := strconv.Atoi(parts[0])
        if err != nil {
                http.Error(w, "Invalid project ID", http.StatusBadRequest)
                return
        }
        
        // Handle different hierarchical routes
        if len(parts) == 1 {
                // /project/{id} - Project details
                h.renderProjectDetails(w, projectID)
        } else if len(parts) >= 3 && parts[1] == "test-suite" {
                // /project/{pid}/test-suite/{sid} - Test suite details
                testSuiteID, err := strconv.Atoi(parts[2])
                if err != nil {
                        http.Error(w, "Invalid test suite ID", http.StatusBadRequest)
                        return
                }
                
                if len(parts) >= 5 && parts[3] == "test-case" {
                        // /project/{pid}/test-suite/{sid}/test-case/{cid} - Test case details
                        testCaseID, err := strconv.Atoi(parts[4])
                        if err != nil {
                                http.Error(w, "Invalid test case ID", http.StatusBadRequest)
                                return
                        }
                        h.renderTestCaseDetails(w, projectID, testSuiteID, testCaseID)
                } else {
                        h.renderTestSuiteDetails(w, projectID, testSuiteID)
                }
        } else {
                http.Error(w, "Invalid URL path", http.StatusBadRequest)
        }
}

func (h *Handler) renderProjectDetails(w http.ResponseWriter, projectID int) {
        tmpl := h.parseTemplates("project.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }

        data := map[string]interface{}{
                "title":     "Project Details",
                "projectID": projectID,
        }
        tmpl.Execute(w, data)
}

func (h *Handler) renderTestSuiteDetails(w http.ResponseWriter, projectID, testSuiteID int) {
        tmpl := h.parseTemplates("test-suite.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }

        data := map[string]interface{}{
                "title":       "Test Suite Details",
                "projectID":   projectID,
                "testSuiteID": testSuiteID,
        }
        tmpl.Execute(w, data)
}

func (h *Handler) renderTestCaseDetails(w http.ResponseWriter, projectID, testSuiteID, testCaseID int) {
        tmpl := h.parseTemplates("test-case.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }

        data := map[string]interface{}{
                "title":       "Test Case Details",
                "projectID":   projectID,
                "testSuiteID": testSuiteID,
                "testCaseID":  testCaseID,
        }
        tmpl.Execute(w, data)
}



func (h *Handler) reportsHandler(w http.ResponseWriter, r *http.Request) {
        tmpl := h.parseTemplates("reports.html")
        if tmpl == nil {
                http.Error(w, "Template error", http.StatusInternalServerError)
                return
        }

        data := map[string]interface{}{
                "title": "Reports",
        }
        tmpl.Execute(w, data)
}