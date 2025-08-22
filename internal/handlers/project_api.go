package handlers

import (
        "database/sql"
        "encoding/json"
        "net/http"
        "strconv"
        "strings"

        "github.com/galex-do/test-machine/internal/models"
)

// projectsAPIHandler handles API requests for projects collection
func (h *Handler) projectsAPIHandler(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
                h.getAllProjects(w, r)
        case "POST":
                h.createProject(w, r)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

// projectAPIHandler handles API requests for individual projects
func (h *Handler) projectAPIHandler(w http.ResponseWriter, r *http.Request) {
        idStr := strings.TrimPrefix(r.URL.Path, "/api/projects/")
        id, err := strconv.Atoi(idStr)
        if err != nil {
                h.writeJSONError(w, "Invalid project ID", http.StatusBadRequest)
                return
        }

        switch r.Method {
        case "GET":
                h.getProject(w, r, id)
        case "PUT":
                h.updateProject(w, r, id)
        case "DELETE":
                h.deleteProject(w, r, id)
        default:
                h.writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
}

func (h *Handler) getAllProjects(w http.ResponseWriter, r *http.Request) {
        projects, err := h.projectService.GetAll()
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        h.writeJSONResponse(w, projects)
}

func (h *Handler) getProject(w http.ResponseWriter, r *http.Request, id int) {
        project, err := h.projectService.GetByID(id)
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        if project == nil {
                h.writeJSONError(w, "Project not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, project)
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
        var req models.CreateProjectRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        project, err := h.projectService.Create(&req)
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        h.writeJSONResponse(w, project)
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request, id int) {
        var req models.UpdateProjectRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
                h.writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
                return
        }

        project, err := h.projectService.Update(id, &req)
        if err == sql.ErrNoRows {
                h.writeJSONError(w, "Project not found", http.StatusNotFound)
                return
        }
        if err != nil {
                h.writeJSONError(w, err.Error(), http.StatusBadRequest)
                return
        }

        if project == nil {
                h.writeJSONError(w, "Project not found", http.StatusNotFound)
                return
        }

        h.writeJSONResponse(w, project)
}

func (h *Handler) deleteProject(w http.ResponseWriter, r *http.Request, id int) {
        err := h.projectService.Delete(id)
        if err == sql.ErrNoRows {
                h.writeJSONError(w, "Project not found", http.StatusNotFound)
                return
        }
        if err != nil {
                h.writeJSONError(w, "Database error", http.StatusInternalServerError)
                return
        }

        h.writeJSONResponse(w, map[string]string{"message": "Project deleted successfully"})
}