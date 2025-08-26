package handlers

import (
        "encoding/json"
        "fmt"
        "net/http"
        "strconv"

        "github.com/galex-do/test-machine/internal/service"
)

type SyncHandler struct {
        gitService *service.GitService
}

func NewSyncHandler(gitService *service.GitService) *SyncHandler {
        return &SyncHandler{
                gitService: gitService,
        }
}

// SyncProject handles POST /api/projects/{id}/sync
func (h *SyncHandler) SyncProject(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Extract project ID from URL path
        path := r.URL.Path
        // Expected path: /api/projects/{id}/sync
        var projectIDStr string
        if _, err := fmt.Sscanf(path, "/api/projects/%s/sync", &projectIDStr); err != nil {
                http.Error(w, "Invalid project ID", http.StatusBadRequest)
                return
        }

        projectID, err := strconv.Atoi(projectIDStr)
        if err != nil {
                http.Error(w, "Invalid project ID", http.StatusBadRequest)
                return
        }

        // Perform sync
        response, err := h.gitService.SyncProjectRepository(projectID)
        if err != nil {
                http.Error(w, fmt.Sprintf("Sync failed: %v", err), http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(response); err != nil {
                http.Error(w, "Failed to encode response", http.StatusInternalServerError)
                return
        }
}

// GetProjectRepository handles GET /api/projects/{id}/repository
func (h *SyncHandler) GetProjectRepository(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
                http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
                return
        }

        // Extract project ID from URL path
        path := r.URL.Path
        // Expected path: /api/projects/{id}/repository
        var projectIDStr string
        if _, err := fmt.Sscanf(path, "/api/projects/%s/repository", &projectIDStr); err != nil {
                http.Error(w, "Invalid project ID", http.StatusBadRequest)
                return
        }

        _, err := strconv.Atoi(projectIDStr)
        if err != nil {
                http.Error(w, "Invalid project ID", http.StatusBadRequest)
                return
        }

        // Get repository data (this would need a repository service method)
        // For now, return empty response as this is not the main focus
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message": "Repository data endpoint - to be implemented"}`))
}