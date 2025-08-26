package handlers

import (
        "encoding/json"
        "net/http"
        "strconv"

        "github.com/galex-do/test-machine/internal/models"
        "github.com/galex-do/test-machine/internal/repository"
        "github.com/galex-do/test-machine/internal/service"
)

type RepositoryAPIHandler struct {
        repositoryRepo *repository.RepositoryRepository
        projectRepo    *repository.ProjectRepository
        gitService     *service.GitService
}

func NewRepositoryAPIHandler(repositoryRepo *repository.RepositoryRepository, projectRepo *repository.ProjectRepository, gitService *service.GitService) *RepositoryAPIHandler {
        return &RepositoryAPIHandler{
                repositoryRepo: repositoryRepo,
                projectRepo:    projectRepo,
                gitService:     gitService,
        }
}

// GetRepositories handles GET /api/repositories - returns all repositories
func (h *RepositoryAPIHandler) GetRepositories(w http.ResponseWriter, r *http.Request) {
        repositories, err := h.repositoryRepo.GetAll()
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(repositories)
}

// GetRepository handles GET /api/repositories/{id} - returns a specific repository
func (h *RepositoryAPIHandler) GetRepository(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
                http.Error(w, "Invalid repository ID", http.StatusBadRequest)
                return
        }

        repository, err := h.repositoryRepo.GetByID(id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        if repository == nil {
                http.Error(w, "Repository not found", http.StatusNotFound)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(repository)
}

// GetRepositoryWithBranches handles GET /api/repositories/{id}/branches-tags - returns repository with branches and tags
func (h *RepositoryAPIHandler) GetRepositoryWithBranches(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
                http.Error(w, "Invalid repository ID", http.StatusBadRequest)
                return
        }

        repository, err := h.repositoryRepo.GetWithBranchesAndTags(id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        if repository == nil {
                http.Error(w, "Repository not found", http.StatusNotFound)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(repository)
}

// CreateRepository handles POST /api/repositories - creates a new repository
func (h *RepositoryAPIHandler) CreateRepository(w http.ResponseWriter, r *http.Request) {
        var req models.CreateRepositoryRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        if req.Name == "" {
                http.Error(w, "Repository name is required", http.StatusBadRequest)
                return
        }

        if req.RemoteURL == "" {
                http.Error(w, "Repository URL is required", http.StatusBadRequest)
                return
        }

        repository, err := h.repositoryRepo.Create(&req)
        if err != nil {
                if err.Error() == "pq: duplicate key value violates unique constraint \"unique_repository_url\"" {
                        http.Error(w, "A repository with this URL already exists", http.StatusConflict)
                        return
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(repository)
}

// UpdateRepository handles PUT /api/repositories/{id} - updates an existing repository
func (h *RepositoryAPIHandler) UpdateRepository(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
                http.Error(w, "Invalid repository ID", http.StatusBadRequest)
                return
        }

        var req models.UpdateRepositoryRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        if req.Name == "" {
                http.Error(w, "Repository name is required", http.StatusBadRequest)
                return
        }

        repository, err := h.repositoryRepo.Update(id, &req)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        if repository == nil {
                http.Error(w, "Repository not found", http.StatusNotFound)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(repository)
}

// DeleteRepository handles DELETE /api/repositories/{id} - deletes a repository
func (h *RepositoryAPIHandler) DeleteRepository(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
                http.Error(w, "Invalid repository ID", http.StatusBadRequest)
                return
        }

        // Check if repository is being used by any projects
        projectCount, err := h.projectRepo.CountProjectsByRepositoryID(id)
        if err != nil {
                http.Error(w, "Error checking repository usage", http.StatusInternalServerError)
                return
        }

        if projectCount > 0 {
                http.Error(w, "Cannot delete repository: it is linked to one or more projects", http.StatusConflict)
                return
        }

        err = h.repositoryRepo.Delete(id)
        if err != nil {
                if err.Error() == "repository not found" {
                        http.Error(w, "Repository not found", http.StatusNotFound)
                        return
                }
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}

// SyncRepository handles POST /api/repositories/{id}/sync - syncs a repository
func (h *RepositoryAPIHandler) SyncRepository(w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.PathValue("id"))
        if err != nil {
                http.Error(w, "Invalid repository ID", http.StatusBadRequest)
                return
        }

        response, err := h.gitService.SyncRepository(id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}