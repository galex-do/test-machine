package handlers

import (
        "encoding/json"
        "fmt"
        "net/http"
        "strconv"

        "github.com/galex-do/test-machine/internal/models"
        "github.com/galex-do/test-machine/internal/repository"
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
        repositoryRepo   *repository.RepositoryRepository
}

// NewHandler creates a new handler
func NewHandler(projectService *service.ProjectService, testSuiteService *service.TestSuiteService, testCaseService *service.TestCaseService, testRunService *service.TestRunService, keyService *service.KeyService, gitService *service.GitService, repositoryRepo *repository.RepositoryRepository) *Handler {
        return &Handler{
                projectService:   projectService,
                testSuiteService: testSuiteService,
                testCaseService:  testCaseService,
                testRunService:   testRunService,
                keyService:       keyService,
                gitService:       gitService,
                repositoryRepo:   repositoryRepo,
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
        mux.HandleFunc("/api/repositories", h.repositoriesAPIHandler)
        mux.HandleFunc("/api/repositories/", h.repositoryAPIHandler)
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

        // Handle /api/sync/repositories/{id}/sync
        var repositoryIDStr string
        if n, _ := fmt.Sscanf(path, "/api/sync/repositories/%s/sync", &repositoryIDStr); n == 1 {
                repositoryID, err := strconv.Atoi(repositoryIDStr)
                if err != nil {
                        h.writeJSONError(w, "Invalid repository ID", http.StatusBadRequest)
                        return
                }

                // Perform repository sync
                response, err := h.gitService.SyncRepository(repositoryID)
                if err != nil {
                        h.writeJSONError(w, fmt.Sprintf("Sync failed: %v", err), http.StatusInternalServerError)
                        return
                }

                h.writeJSONResponse(w, response)
                return
        }

        h.writeJSONError(w, "Not found", http.StatusNotFound)
}

// repositoriesAPIHandler handles /api/repositories requests
func (h *Handler) repositoriesAPIHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                repositories, err := h.repositoryRepo.GetAll()
                if err != nil {
                        h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                h.writeJSONResponse(w, repositories)

        case "POST":
                var req models.CreateRepositoryRequest
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                        h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                        return
                }

                if req.Name == "" {
                        h.writeJSONError(w, "Repository name is required", http.StatusBadRequest)
                        return
                }

                if req.RemoteURL == "" {
                        h.writeJSONError(w, "Repository URL is required", http.StatusBadRequest)
                        return
                }

                repository, err := h.repositoryRepo.Create(&req)
                if err != nil {
                        if fmt.Sprintf("%v", err) == "pq: duplicate key value violates unique constraint \"unique_repository_url\"" {
                                h.writeJSONError(w, "A repository with this URL already exists", http.StatusConflict)
                                return
                        }
                        h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                w.WriteHeader(http.StatusCreated)
                h.writeJSONResponse(w, repository)

        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// repositoryAPIHandler handles /api/repositories/{id} requests  
func (h *Handler) repositoryAPIHandler(w http.ResponseWriter, r *http.Request) {
        // Extract ID from URL
        path := r.URL.Path
        var idStr string
        if n, _ := fmt.Sscanf(path, "/api/repositories/%s", &idStr); n != 1 {
                h.writeJSONError(w, "Invalid URL format", http.StatusBadRequest)
                return
        }

        // Handle sync endpoint separately
        if r.Method == "POST" && fmt.Sprintf("/api/repositories/%s/sync", idStr) == path {
                id, err := strconv.Atoi(idStr)
                if err != nil {
                        h.writeJSONError(w, "Invalid repository ID", http.StatusBadRequest)
                        return
                }

                response, err := h.gitService.SyncRepository(id)
                if err != nil {
                        h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                h.writeJSONResponse(w, response)
                return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
                h.writeJSONError(w, "Invalid repository ID", http.StatusBadRequest)
                return
        }

        switch r.Method {
        case "GET":
                repository, err := h.repositoryRepo.GetByID(id)
                if err != nil {
                        h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                if repository == nil {
                        h.writeJSONError(w, "Repository not found", http.StatusNotFound)
                        return
                }

                h.writeJSONResponse(w, repository)

        case "PUT":
                var req models.UpdateRepositoryRequest
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                        h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                        return
                }

                if req.Name == "" {
                        h.writeJSONError(w, "Repository name is required", http.StatusBadRequest)
                        return
                }

                repository, err := h.repositoryRepo.Update(id, &req)
                if err != nil {
                        h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                if repository == nil {
                        h.writeJSONError(w, "Repository not found", http.StatusNotFound)
                        return
                }

                h.writeJSONResponse(w, repository)

        case "DELETE":
                err := h.repositoryRepo.Delete(id)
                if err != nil {
                        if err.Error() == "repository not found" {
                                h.writeJSONError(w, "Repository not found", http.StatusNotFound)
                                return
                        }
                        h.writeJSONError(w, err.Error(), http.StatusInternalServerError)
                        return
                }

                w.WriteHeader(http.StatusNoContent)

        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}