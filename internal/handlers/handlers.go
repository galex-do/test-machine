package handlers

import (
        "encoding/json"
        "fmt"
        "net/http"
        "strconv"

        "github.com/galex-do/test-machine/internal/service"
)

// Handler holds all the dependencies for HTTP handlers
type Handler struct {
        projectService   *service.ProjectService
        testSuiteService *service.TestSuiteService
        testCaseService  *service.TestCaseService
        testRunService   *service.TestRunService
        keyService       *service.KeyService
        gitService       *service.GitService
}

// NewHandler creates a new handler
func NewHandler(projectService *service.ProjectService, testSuiteService *service.TestSuiteService, testCaseService *service.TestCaseService, testRunService *service.TestRunService, keyService *service.KeyService, gitService *service.GitService) *Handler {
        return &Handler{
                projectService:   projectService,
                testSuiteService: testSuiteService,
                testCaseService:  testCaseService,
                testRunService:   testRunService,
                keyService:       keyService,
                gitService:       gitService,
        }
}

// SetupRoutes sets up all the routes
func (h *Handler) SetupRoutes() http.Handler {
        mux := http.NewServeMux()

        // API routes only
        mux.HandleFunc("/api/projects", h.projectsAPIHandler)
        mux.HandleFunc("/api/projects/", h.projectAPIHandler)
        mux.HandleFunc("/api/test-suites", h.testSuitesAPIHandler)
        mux.HandleFunc("/api/test-suites/", h.testSuiteAPIHandler)
        mux.HandleFunc("/api/test-cases", h.testCasesAPIHandler)
        mux.HandleFunc("/api/test-cases/", h.testCaseAPIHandler)
        mux.HandleFunc("/api/test-runs", h.testRunsAPIHandler)
        mux.HandleFunc("/api/test-runs/", h.testRunAPIHandler)
        mux.HandleFunc("/api/test-steps/", h.testStepAPIHandler)
        mux.HandleFunc("/api/keys", h.keyAPIHandler)
        mux.HandleFunc("/api/keys/", h.keyByIDAPIHandler)
        mux.HandleFunc("/api/sync/", h.syncAPIHandler)
        mux.HandleFunc("/api/stats", h.statsAPIHandler)

        // Add CORS middleware
        return h.corsMiddleware(mux)
}

// corsMiddleware adds CORS headers for frontend integration
func (h *Handler) corsMiddleware(next *http.ServeMux) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Set CORS headers
                w.Header().Set("Access-Control-Allow-Origin", "*")
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

                // Handle preflight requests
                if r.Method == "OPTIONS" {
                        w.WriteHeader(http.StatusOK)
                        return
                }

                next.ServeHTTP(w, r)
        })
}

// writeJSONError writes a JSON error response
func (h *Handler) writeJSONError(w http.ResponseWriter, message string, statusCode int) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(statusCode)
        json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeJSONResponse writes a JSON response
func (h *Handler) writeJSONResponse(w http.ResponseWriter, data interface{}) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(data)
}

// statsAPIHandler handles /api/stats requests
func (h *Handler) statsAPIHandler(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Get all entities to calculate stats
        projects, err := h.projectService.GetAll()
        if err != nil {
                h.writeJSONError(w, "Error fetching projects", http.StatusInternalServerError)
                return
        }

        testSuites, err := h.testSuiteService.GetAll(nil)
        if err != nil {
                h.writeJSONError(w, "Error fetching test suites", http.StatusInternalServerError)
                return
        }

        testCases, err := h.testCaseService.GetAll(nil)
        if err != nil {
                h.writeJSONError(w, "Error fetching test cases", http.StatusInternalServerError)
                return
        }

        testRuns, err := h.testRunService.GetAll(nil)
        if err != nil {
                h.writeJSONError(w, "Error fetching test runs", http.StatusInternalServerError)
                return
        }

        stats := map[string]interface{}{
                "totalProjects":   len(projects),
                "totalTestSuites": len(testSuites),
                "totalTestCases":  len(testCases),
                "totalTestRuns":   len(testRuns),
        }

        h.writeJSONResponse(w, stats)
}

// syncAPIHandler handles sync-related API requests
func (h *Handler) syncAPIHandler(w http.ResponseWriter, r *http.Request) {
        // Parse the URL to extract project ID and action
        path := r.URL.Path
        
        // Handle /api/sync/projects/{id}/sync
        if r.Method == "POST" {
                var projectIDStr string
                if n, _ := fmt.Sscanf(path, "/api/sync/projects/%s/sync", &projectIDStr); n == 1 {
                        projectID, err := strconv.Atoi(projectIDStr)
                        if err != nil {
                                h.writeJSONError(w, "Invalid project ID", http.StatusBadRequest)
                                return
                        }

                        // Perform sync
                        response, err := h.gitService.SyncProjectRepository(projectID)
                        if err != nil {
                                h.writeJSONError(w, fmt.Sprintf("Sync failed: %v", err), http.StatusInternalServerError)
                                return
                        }

                        h.writeJSONResponse(w, response)
                        return
                }
        }

        h.writeJSONError(w, "Not found", http.StatusNotFound)
}